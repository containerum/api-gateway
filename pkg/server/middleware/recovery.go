package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				path := c.Request.URL.Path
				raw := c.Request.URL.RawQuery
				if raw != "" {
					path = path + "?" + raw
				}
				log.WithFields(log.Fields{
					"Time":   time.Now().Format(timeFormat),
					"Path":   path,
					"Method": c.Request.Method,
					"IP":     c.ClientIP(),
					"Code":   500,
				}).Error("Recovery")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
