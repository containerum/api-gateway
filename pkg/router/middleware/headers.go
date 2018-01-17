package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	userClientHeaderName = "User-Client"

	userClientXHeaderName = "X-User-Client"
	userIPXHeaderName     = "X-Client-IP"
	userAgentXHeaderName  = "X-User-Agent"
	гыук
)

var (
	xHeaderRegexp, _           = regexp.Compile("^X-[a-zA-Z0-9]+")
	xHeaderUserClientRegexp, _ = regexp.Compile("^[a-f0-9]{32}$")
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
			log.WithField("Header", k).WithField("Value", v).Debug("Header deleted from response")
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

func TranslateUserXHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String()
		w.Header().Add(userIPXHeaderName, ip) //Add X-Client-IP
		log.WithField("Name", userIPXHeaderName).WithField("Value", ip).Debug("Add X-Header")
		w.Header().Add(userAgentXHeaderName, r.UserAgent()) //Add X-User-Agent
		log.WithField("Name", userAgentXHeaderName).WithField("Value", r.UserAgent()).Debug("Add X-Header")
		//Translate user header to X header
		for k, v := range r.Header {
			if k == userClientHeaderName && len(v) > 0 { //Add X-User-Client
				if xHeaderUserClientRegexp.MatchString(v[0]) {
					w.Header()[userClientXHeaderName] = v
					log.WithField("Name", userClientXHeaderName).WithField("Value", v[0]).Debug("Add X-Header")
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
