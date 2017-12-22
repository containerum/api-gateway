package middleware

import (
	"net"
	"net/http"

	rate "git.containerum.net/ch/ratelimiter"

	log "github.com/Sirupsen/logrus"
)

//Rate limitate reuests per second
func Rate(limiter *rate.PerIPLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter != nil {
				host, _, _ := net.SplitHostPort(r.RemoteAddr)
				ok, _ := limiter.Limit(host)
				log.WithFields(log.Fields{
					"Host":   host,
					"Status": ok,
				}).Debug("Rate limit")
				if !ok {
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
