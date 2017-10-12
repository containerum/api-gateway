package middleware

import (
	"github.com/rs/cors"
)

//Cors make return Allowed CORS for server
func Cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST"},
	})
}
