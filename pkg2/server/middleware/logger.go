package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	timeFormat = "02.01.2006 15:04:05"
)

// Logger will write the request logs and save it's in clickhouse
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		log.WithFields(log.Fields{
			"Time":    time.Now().Format(timeFormat),
			"Latency": latency,
			"IP":      c.ClientIP(),
			"Method":  c.Request.Method,
			"Code":    c.Writer.Status(),
			"Path":    path,
		}).Info("Request")
	}
}
