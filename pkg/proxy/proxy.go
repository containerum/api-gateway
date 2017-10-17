package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

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
		r.URL, _ = url.Parse(target.UpstreamURL)
	}
}
