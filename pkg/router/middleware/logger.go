package middleware

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		//Write log after
		latency := time.Now().Sub(start)
		log.WithFields(log.Fields{
			"Method":  r.Method,
			"Path":    r.RequestURI,
			"Latency": fmt.Sprintf("%v", latency),
		}).Info("Request")
	})
}
