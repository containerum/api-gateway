package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	errs "git.containerum.net/ch/api-gateway/pkg/errors"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
)

//Limiter keeps tollbooth limiter for limiting http requests
type Limiter struct {
	*limiter.Limiter
	rate int
}

var (
	errTooManyRequests = errors.New("Too many requests per second")
)

//CreateLimiter return rate limiter for http
func CreateLimiter(rate int) *Limiter {
	limiter := tollbooth.NewLimiter(float64(rate), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	limiter.SetIPLookups([]string{"X-Client-IP", "X-Forwarded-For", "X-Real-IP"})
	return &Limiter{limiter, rate}
}

//Limit middleware for limiting http requests
func (l *Limiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByKeys(l.Limiter, []string{c.ClientIP()})
		if httpError != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, errs.New(errTooManyRequests.Error(), fmt.Sprintf("Max request count: %v", l.rate)))
		} else {
			c.Next()
		}
	}
}
