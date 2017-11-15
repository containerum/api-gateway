package middleware

import (
	"github.com/rs/cors"
)

// TODO: Write own ServeHTTP
// TODO: Add Requred Headers

//Cors make return Allowed CORS for server
func Cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST"},
		Debug:            true,
	})
}
