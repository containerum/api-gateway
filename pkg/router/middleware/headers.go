package middleware

import (
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

var (
	xHeaderRegexp, _ = regexp.Compile("^X-[a-zA-Z0-9]+")
)

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
		next.ServeHTTP(w, r)
		//Clear response headers
		for header := range w.Header() {
			if xHeaderRegexp.MatchString(header) {
				w.Header().Del(header)
				log.WithField("Header", header).Debug("Header deleted from response")
			}
		}
	})
}
