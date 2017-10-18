package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"bitbucket.org/exonch/ch-gateway/pkg/model"
)

func CreateProxy(target *model.Router) *httputil.ReverseProxy {
	direct := createDirector(target)
	return &httputil.ReverseProxy{
		Director: direct,
	}
}

func createDirector(target *model.Router) func(r *http.Request) {
	return func(r *http.Request) {
		targetURL, _ := url.Parse(target.UpstreamURL)
		r.URL.Scheme = targetURL.Scheme
		r.URL.Host = targetURL.Host
		r.URL.Path = singleJoiningSlash(targetURL.Path, r.URL.Path)

		if target.StripPath {
			strPath := stripPath(target.ListenPath, r.URL.Path)
			r.URL.Path = singleJoiningSlash(targetURL.Path, strPath)
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
