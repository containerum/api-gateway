package middleware

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

//Logger write main logs
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := NewWrapResponseWriter(w)
		next.ServeHTTP(ww, r)
		//Write log after
		latency := time.Now().Sub(start)
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
