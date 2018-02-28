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
