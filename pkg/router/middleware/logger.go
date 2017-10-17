package middleware

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/cactus/go-statsd-client/statsd"
)

//Statter connection with Statsd
var Statter *statsd.Statter

//Logger write main logs
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := NewWrapResponseWriter(w)

		next.ServeHTTP(ww, r)

		latency := time.Now().Sub(start)

		//Set status in Statsd
		if Statter != nil {
			statusCall := fmt.Sprintf("call.status.%v", ww.Status())
			methodCall := fmt.Sprintf("call.method.%v", r.Method)
			(*Statter).Inc("call.status.all", 1, 1.0)
			(*Statter).Inc(statusCall, 1, 1.0)
			(*Statter).Inc(methodCall, 1, 1.0)
		}

		//Write log after
		log.WithFields(log.Fields{
			"Method":       r.Method,
			"Path":         r.RequestURI,
			"Latency":      fmt.Sprintf("%v", latency),
			"Status":       ww.Status(),
			"RequestID":    w.Header().Get("X-Request-ID"),
			"ResponseSize": ww.BytesWritten(),
			"RequestSize":  r.ContentLength,
		}).Info("Request")
	})
}
