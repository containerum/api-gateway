package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Cors return cors rules
func Cors() gin.HandlerFunc {
	c := cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "PUT", "PATH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "User-Client", "User-Token", "Authorization", "Sec-Websocket-Version",
			"Connection", "Upgrade", "Sec-Websocket-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           2 * time.Hour,
	})
	return c
}
