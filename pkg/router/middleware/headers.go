package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

var (
	xHeaderRegexp, _ = regexp.Compile("^X-[a-zA-Z0-9]+")
)

type ModifierMiddleware struct {
	handler http.Handler
}

func (m *ModifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	// passing a ResponseRecorder instead of the original RW
	m.handler.ServeHTTP(rec, r)
	// we copy the original headers and remove X-headers from response
	for k, v := range rec.Header() {
		if xHeaderRegexp.MatchString(k) {
			log.WithField("Header", k).Debug("Header deleted from response")
			continue
		}
		w.Header()[k] = v
	}
	//Write Body
	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())
}

//ClearXHeaders remove all requested X-Headers and remove all responsed X-Headers
func ClearXHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Clear request headers
		for header := range r.Header {
			if xHeaderRegexp.MatchString(header) {
				r.Header.Del(header)
				log.WithField("Header", header).Debug("Header deleted from request")
			}
		}
		x := ModifierMiddleware{next}
		x.ServeHTTP(w, r)
	})
}
