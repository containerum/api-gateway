package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
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
)

var (
	xHeaderRegexp, _           = regexp.Compile("^X-[a-zA-Z0-9]+")
	xHeaderUserClientRegexp, _ = regexp.Compile("^[a-f0-9]{32}$")

	requiredXHeaders = []string{userClientXHeaderName}
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

func CheckRequiredXHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(w.Header())
		var errHeaders []error
		for _, header := range requiredXHeaders {
			headerExist := false
			for h := range w.Header() {
				log.WithField("H", h).WithField("Header", header).Debug("Compare")
				if h == header {
					headerExist = true
					continue
				}
			}
			if !headerExist {
				errHeaders = append(errHeaders, errors.New(strings.TrimPrefix(header, "X-")))
			}
		}
		if len(errHeaders) != 0 {
			answer := struct {
				Error string
			}{
				Error: fmt.Sprintf("required headers %v was not provided", errHeaders),
			}
			data, _ := json.Marshal(&answer)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// const (
// 	UserIDHeader      = "X-User-ID"
// 	SessionIDHeader   = "X-Session-ID"
// 	FingerprintHeader = "X-User-Client"
// 	UserAgentHeader   = "X-User-Agent"
// 	ClientIPHeader    = "X-Client-IP"
// 	TokenIDHeader     = "X-Token-ID"
// 	PartTokenIDHeader = "X-User-Part-Token-ID"
// 	UserRoleHeader    = "X-User-Role"
// )
