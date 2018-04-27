package middleware

import (
	"fmt"
	"net/url"
	"time"

	"git.containerum.net/ch/api-gateway/pkg/model"
	h "github.com/containerum/utils/httputil"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	timeFormat = "02.01.2006 15:04:05"
)

type req struct {
	start  time.Time
	finish time.Time
	path   string
	method string
	code   int
	ip     string
	agent  string
	route  string
}

type setFunc func(m *model.Metrics, r *req)

// Logger will write the request logs and save it's in clickhouse
func Logger(m *model.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := req{
			start:  time.Now(),
			path:   getPath(c.Request.URL),
			method: c.Request.Method,
			agent:  c.Request.UserAgent(),
			ip:     c.ClientIP(),
		}
		c.Next()
		r.finish = time.Now()
		r.route = getRoute(c)
		getReqEntry(&r).Info("Request")
		setReqFunc(m, &r, setReqCount, setReqCountIP, setReqCountRoute, setReqCountUserAgent)
	}
}

func setReqFunc(m *model.Metrics, r *req, fns ...setFunc) {
	for _, fn := range fns {
		fn(m, r)
	}
}

func setReqCount(m *model.Metrics, r *req) {
	reqCount := m.RTotal.(*prometheus.CounterVec)
	reqCount.WithLabelValues(r.method, getStatus(r.code)).Inc()
}

func setReqCountIP(m *model.Metrics, r *req) {
	reqCountIP := m.RUserIP.(*prometheus.CounterVec)
	reqCountIP.WithLabelValues(r.method, getStatus(r.code), r.ip).Inc()
}

func setReqCountRoute(m *model.Metrics, r *req) {
	reqCountRoute := m.RRoute.(*prometheus.CounterVec)
	reqCountRoute.WithLabelValues(r.method, getStatus(r.code), r.route).Inc()
}

func setReqCountUserAgent(m *model.Metrics, r *req) {
	reqCountUserAgent := m.RUserAgent.(*prometheus.CounterVec)
	reqCountUserAgent.WithLabelValues(r.method, getStatus(r.code), r.route, r.agent).Inc()
}

func getStatus(status int) string {
	return fmt.Sprintf("%d", status)
}

func getPath(u *url.URL) string {
	if u.RawQuery != "" {
		return u.Path + "?" + u.RawQuery
	}
	return u.Path
}

func getRoute(c *gin.Context) (route string) {
	if route = c.GetHeader(h.RequestNameXHeader); route == "" {
		return "unknow"
	}
	return
}

func getReqEntry(r *req) *log.Entry {
	return log.WithFields(log.Fields{
		"Time":    time.Now().Format(timeFormat),
		"Latency": r.finish.Sub(r.start),
		"IP":      r.ip,
		"Method":  r.method,
		"Code":    r.code,
		"Path":    r.path,
		"Route":   r.route,
	})
}
