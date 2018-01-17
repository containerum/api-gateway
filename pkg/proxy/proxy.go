package proxy

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

	log "github.com/Sirupsen/logrus"
)

func CreateProxy(target *model.Listener, headers http.Header) *httputil.ReverseProxy {
	direct := createDirector(target, &headers)
	return &httputil.ReverseProxy{
		Director: direct,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second, //TODO: Get it from ENV
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
	}
}

func createDirector(target *model.Listener, headers *http.Header) func(r *http.Request) {
	return func(r *http.Request) {
		targetURL, _ := url.Parse(target.UpstreamURL)
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.URL.Path = singleJoiningSlash(targetURL.Path, r.URL.Path)

		if target.StripPath {
			strPath := stripPath(target.ListenPath, r.URL.Path)
			r.URL.Path = singleJoiningSlash(targetURL.Path, strPath)
		}
		if headers != nil {
			r.Header = *headers
			log.WithField("Headers", *headers).Debug("Add headers to proxy director")
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

func stripPath(listenPath, path string) string {
	rSlash, _ := regexp.Compile("/")
	listenPath = rSlash.ReplaceAllString(listenPath, `\/`)

	re := fmt.Sprintf("^%v", listenPath)
	r, _ := regexp.Compile(re)

	return r.ReplaceAllString(path, "")
}
