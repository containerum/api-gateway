package middleware

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

//RequestID append Request ID to response
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewV4()
		log.WithField("ID", id.String()).Debug("Request ID created")
		w.Header().Set("X-Request-ID", id.String())
		next.ServeHTTP(w, r)
	})
}
