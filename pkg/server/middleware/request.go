package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

//SetRequestID set request id header
func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewV4()
		setHeader(&c.Request.Header, requestIDXHeader, id.String())
		c.Next()
	}
}

//SetRequestName set request name header
func SetRequestName(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		setHeader(&c.Request.Header, requestNameXHeader, name)
		c.Next()
	}
}
