package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type proxyTransport struct {
	header http.Header
}

const (
	wsReadBuffer  = 1024
	wsWriteBuffer = 1024
	httpTimeout   = 15 * time.Second
	httpKeepAlive = 30 * time.Second
)

var (
	headersToDelete = []string{"Connection", "Sec-Websocket-Key", "Sec-Websocket-Version", "Sec-Websocket-Extensions", "Upgrade"}
	httpDialer      = &net.Dialer{
		Timeout:   httpTimeout,
		KeepAlive: httpKeepAlive,
	}
	httpTransport = http.Transport{Dial: httpDialer.Dial}
	wsUpgrader    = &websocket.Upgrader{
		ReadBufferSize:  wsReadBuffer,
		WriteBufferSize: wsWriteBuffer,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (pt proxyTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	if resp, err = httpTransport.RoundTrip(r); err != nil {
		return resp, err
	}
	for header, value := range resp.Header {
		if middleware.XHeaderRegexp.MatchString(header) {
			middleware.HeaderEntry(header, value).Debug("Header deleted from response")
			continue
		}
		pt.header.Add(header, value[0])
	}
	resp.Header = pt.header
	return resp, err
}

func proxyHandler(route model.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		if route.WS {
			request := c.Request.URL
			target, _ := url.Parse(route.Upstream)
			request.Scheme, request.Host = "ws", target.Host
			if err := proxyWS(c, request); err != nil {
				return
			}
			proxyWS(c, request)
		} else {
			p := proxyHTTP(&route)
			p.ServeHTTP(c.Writer, c.Request)
		}
	}
}

func proxyHTTP(target *model.Route) *httputil.ReverseProxy {
	direct := createDirector(target)
	return &httputil.ReverseProxy{
		Director: direct,
		Transport: proxyTransport{
			header: http.Header{},
		},
	}
}

func proxyWS(c *gin.Context, backend *url.URL) error {
	deleteHeaders(&c.Request.Header, headersToDelete...)
	conn, connBackend, err := makeWSconnections(c, backend)
	if err != nil {
		return err
	}
	errc := make(chan error, 2)
	replicateWebsocketConn := func(dst, src *websocket.Conn, dstName, srcName string) {
		var err error
		var msgType int
		var msg []byte
		for {
			msgType, msg, err = src.ReadMessage()
			if err != nil {
				log.Errorf("WebSocketProxy: Error when copying from %s to %s using ReadMessage: %v", srcName, dstName, err)
				break
			}
			err = dst.WriteMessage(msgType, msg)
			if err != nil {
				log.Errorf("WebSocketProxy: Error when copying from %s to %s using WriteMessage: %v", srcName, dstName, err)
				break
			}
		}
		errc <- err
	}
	go replicateWebsocketConn(conn, connBackend, "client", "backend")
	go replicateWebsocketConn(connBackend, conn, "backend", "client")
	<-errc
	return nil
}

func makeWSconnections(c *gin.Context, backend *url.URL) (conn, connBackend *websocket.Conn, err error) {
	if conn, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		log.WithError(err).Error("Unable to upgrade to WebSocket")
		return
	}
	if connBackend, _, err = websocket.DefaultDialer.Dial(backend.String(), c.Request.Header); err != nil {
		log.WithError(err).Error("Unable to dial to WebSocket")
		return
	}
	return
}

func deleteHeaders(header *http.Header, keys ...string) {
	for _, key := range keys {
		middleware.HeaderEntry(key, header.Get(key)).Debug("Header deleted")
		header.Del(key)
	}
}

func createDirector(target *model.Route) func(r *http.Request) {
	return func(r *http.Request) {
		targetURL, _ := url.Parse(target.Upstream)
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		if target.Strip {
			strPath := stripPath(r.URL.Path, target.Listen, targetURL.Path)
			r.URL.Path = singleJoiningSlash(buildHostURL(*targetURL), strPath)
		} else {
			r.URL.Path = singleJoiningSlash(targetURL.Path, r.URL.Path)
		}
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func stripPath(requestPath, listenPath, upstreamPath string) string {
	rSlash, _ := regexp.Compile("/")
	listenPath = rSlash.ReplaceAllString(listenPath, `\/`)

	re := fmt.Sprintf("^%v", listenPath)
	r, _ := regexp.Compile(re)

	diffPath := r.ReplaceAllString(requestPath, "")
	if diffPath == "" {
		return upstreamPath
	}
	return singleJoiningSlash(upstreamPath, diffPath)
}

func buildHostURL(u url.URL) string {
	return u.Scheme + "://" + u.Host
}
