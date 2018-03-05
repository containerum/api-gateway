package middleware

import (
	"fmt"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	timeFormat = "02.01.2006 15:04:05"
)

// Logger will write the request logs and save it's in clickhouse
func Logger(m *model.Metrics) gin.HandlerFunc {
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
		status := c.Writer.Status()
		ip := c.ClientIP()
		method := c.Request.Method
		log.WithFields(log.Fields{
			"Time":    time.Now().Format(timeFormat),
			"Latency": latency,
			"IP":      ip,
			"Method":  method,
			"Code":    status,
			"Path":    path,
		}).Info("Request")

		reqCount := m.RTotal.(*prometheus.CounterVec)
		reqCount.WithLabelValues(method, getStatus(status)).Inc()
		reqCountIP := m.RUserIP.(*prometheus.CounterVec)
		reqCountIP.WithLabelValues(method, getStatus(status), ip).Inc()
	}
}

func getStatus(status int) string {
	return fmt.Sprintf("%d", status)
}
