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

var (
	defaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (pt proxyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	transport := http.Transport{Dial: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	}
	resp, err := transport.RoundTrip(r)
	for k, v := range resp.Header {
		if middleware.XHeaderRegexp.MatchString(k) {
			log.WithFields(log.Fields{
				"Header": k,
				"Value":  v,
			}).Debug("Header deleted from response")
			continue
		}
		pt.header.Add(k, v[0])
	}
	resp.Header = pt.header
	return resp, err
}

func proxyHandler(route model.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithField("IsUpgrade", websocket.IsWebSocketUpgrade(c.Request)).WithField("Route", route).Debug("WS")
		if route.WS {
			u := c.Request.URL
			t, _ := url.Parse(route.Upstream)
			u.Scheme = "ws"
			u.Host = t.Host

			log.WithField("IsUpgrade", websocket.IsWebSocketUpgrade(c.Request)).WithField("U", u).Debug("WS")

			connPub, err := defaultUpgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				log.WithError(err).Error("Unable to upgrade to WebSocket")
				return
			}
			defer connPub.Close()

			c.Request.Header.Del("Connection")
			c.Request.Header.Del("Sec-Websocket-Key")
			c.Request.Header.Del("Sec-Websocket-Version")
			c.Request.Header.Del("Sec-Websocket-Extensions")
			c.Request.Header.Del("Upgrade")
			connBackend, _, err := websocket.DefaultDialer.Dial(u.String(), c.Request.Header)
			if err != nil {
				log.WithError(err).Error("Unable to dial to WebSocket")
				return
			}
			defer connBackend.Close()

			errc := make(chan error, 2)
			replicateWebsocketConn := func(dst, src *websocket.Conn, dstName, srcName string) {
				var err error
				for {
					msgType, msg, err := src.ReadMessage()
					if err != nil {
						log.Printf("websocketproxy: error when copying from %s to %s using ReadMessage: %v", srcName, dstName, err)
						break
					}
					err = dst.WriteMessage(msgType, msg)
					if err != nil {
						log.Printf("websocketproxy: error when copying from %s to %s using WriteMessage: %v", srcName, dstName, err)
						break
					}
				}
				errc <- err
			}
			go replicateWebsocketConn(connPub, connBackend, "client", "backend")
			go replicateWebsocketConn(connBackend, connPub, "backend", "client")
			<-errc
		} else {
			p := createProxy(&route)
			p.ServeHTTP(c.Writer, c.Request)
		}
	}
}

func createProxy(target *model.Route) *httputil.ReverseProxy {
	direct := createDirector(target)
	return &httputil.ReverseProxy{
		Director: direct,
		Transport: proxyTransport{
			header: http.Header{},
		},
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
