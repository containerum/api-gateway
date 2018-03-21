package middleware

import (
	"time"

	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	"git.containerum.net/ch/kube-client/pkg/cherry/api-gateway"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

//Limiter keeps tollbooth limiter for limiting http requests
type Limiter struct {
	*limiter.Limiter
	rate int
}

//CreateLimiter return rate limiter for http
func CreateLimiter(rate int) *Limiter {
	limit := tollbooth.NewLimiter(float64(rate), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limit.SetIPLookups([]string{"X-Client-IP", "X-Forwarded-For", "X-Real-IP"})
	return &Limiter{limit, rate}
}

//Limit middleware for limiting http requests
func (l *Limiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByKeys(l.Limiter, []string{c.ClientIP()})
		if httpError != nil {
			gonic.Gonic(gatewayErrors.ErrTooManyRequests().AddDetailF("max request count: %v", l.rate), c)
		} else {
			c.Next()
		}
	}
}
