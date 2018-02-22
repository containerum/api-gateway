package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

//EnableCors make return Allowed CORS for server
func EnableCors(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATH", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "User-Client", "User-Token", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return c.Handler(next)
}
