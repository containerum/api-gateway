package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/gatewayErrors"
	"git.containerum.net/ch/api-gateway/pkg/model"
	"git.containerum.net/ch/api-gateway/pkg/server/middleware"
	"git.containerum.net/ch/api-gateway/pkg/utils/wslimiter"
	"github.com/containerum/cherry/adaptors/gonic"
	h "github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type proxyTransport struct {
	header http.Header
}

const (
	wsReadBuffer          = 1024
	wsWriteBuffer         = 1024
	httpTimeout           = 15 * time.Second
	httpKeepAlive         = 30 * time.Second
	maxWSConnectionsPerIP = 100
)

var (
	headersToDelete = []string{"Connection", "Sec-Websocket-Key", "Sec-Websocket-Version", "Sec-Websocket-Extensions", "Upgrade"}
	httpDialer      = &net.Dialer{
		Timeout:   httpTimeout,
		KeepAlive: httpKeepAlive,
	}
	httpTransport = http.Transport{DialContext: httpDialer.DialContext}
	wsUpgrader    = wslimiter.NewPerIPLimiter(maxWSConnectionsPerIP, &websocket.Upgrader{
		ReadBufferSize:  wsReadBuffer,
		WriteBufferSize: wsWriteBuffer,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}, func(r *http.Request) string {
		return r.Header.Get(h.UserIPXHeader)
	})
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
			err := proxyWS(c, request)
			switch err {
			case nil:
				// pass
			case wslimiter.ErrLimitReached:
				gonic.Gonic(gatewayErrors.ErrTooManyRequests().AddDetailF("websocket connections limit reached"), c)
			}
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
	conn, connBackend, err := makeWSconnections(c, backend)
	if err != nil {
		return err
	}
	errc := make(chan error, 2)
	var once sync.Once

	replicateConn := func(dst, src *wslimiter.Conn, dstName, srcName string) {
		defer once.Do(func() {
			dst.Close()
			src.Close()
		})

		var buf [1024]byte
		_, err := io.CopyBuffer(dst.UnderlyingConn(), src.UnderlyingConn(), buf[:])
		switch {
		case err == nil:
			// pass
		case isClose(err),
			isBrokenPipe(err),
			isNetTimeout(err):
			// pass
		default:
			log.WithError(err).Errorf("websocket: replicate bytes from %s to %s failed", srcName, dstName)
		}
		errc <- err
	}
	go replicateConn(conn, connBackend, "client", "backend")
	go replicateConn(connBackend, conn, "backend", "client")
	<-errc
	return nil
}

func makeWSconnections(c *gin.Context, backend *url.URL) (conn, connBackend *wslimiter.Conn, err error) {
	if conn, err = wsUpgrader.Upgrade(c.Writer, c.Request, nil); err != nil {
		log.WithError(err).Error("Unable to upgrade to WebSocket")
		return
	}
	deleteHeaders(&c.Request.Header, headersToDelete...)
	_connBackend, _, err := websocket.DefaultDialer.Dial(backend.String(), c.Request.Header)
	if err != nil {
		log.WithError(err).Error("Unable to dial to WebSocket")
		return
	}
	connBackend = &wslimiter.Conn{Conn: _connBackend}
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

//
// Websocket-specific connection errors
//

func isNetTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func isBrokenPipe(err error) bool {
	opErr, isOpErr := err.(*net.OpError)
	if !isOpErr {
		return false
	}
	syscallErr, ok := opErr.Err.(*os.SyscallError)
	return ok && syscallErr.Err == syscall.EPIPE
}

func isClose(err error) bool {
	_, isClose := err.(*websocket.CloseError)
	if isClose {
		return true
	}
	opErr, isOpErr := err.(*net.OpError)
	if !isOpErr {
		return false
	}
	return opErr.Err.Error() == "use of closed network connection"
}
