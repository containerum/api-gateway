package middleware

import (
	"net/http"
	"regexp"

	"git.containerum.net/ch/api-gateway/pkg/gatewayErrors"
	"github.com/containerum/cherry/adaptors/gonic"

	h "github.com/containerum/utils/httputil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	//XHeaderRegexp keeps regexp for detecting X-Headers
	XHeaderRegexp, _    = regexp.Compile("^X-[a-zA-Z0-9]+")
	userClientRegexp, _ = regexp.Compile("^[a-f0-9]{32}$")
)

//ClearXHeaders clear all request and response X-Headers
func ClearXHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		clearHeaders(&c.Request.Header, "request")
		c.Next()
	}
}

//SetMainUserXHeaders write X-User-IP, X-User-Client, X-User-Agent for next services
func SetMainUserXHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		setHeader(&c.Request.Header, h.UserIPXHeader, c.ClientIP())
		setHeader(&c.Request.Header, h.UserClientXHeader, c.GetHeader(h.UserClientHeader))
		setHeader(&c.Request.Header, h.UserAgentXHeader, c.Request.UserAgent())
	}
}

//SetHeaderFromQuery write X-User-IP, X-User-Client, X-User-Agent for next services
func SetHeaderFromQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClient, _ := c.GetQuery(h.UserClientHeader)
		userToken, _ := c.GetQuery(h.AuthorizationHeader)
		if c.GetHeader(h.UserClientHeader) == "" && userClient != "" {
			c.Request.Header.Add(h.UserClientHeader, userClient)
		}
		if c.GetHeader(h.AuthorizationHeader) == "" && userToken != "" {
			c.Request.Header.Add(h.AuthorizationHeader, userToken)
		}
	}
}

//CheckUserClientHeader validate User-Client header
func CheckUserClientHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(h.UserClientHeader) == "" {
			gonic.Gonic(gatewayErrors.ErrHeaderNotProvided().AddDetailsErr(h.ErrHeaderRequired(h.UserClientHeader)), c)
			return
		}
		if !userClientRegexp.MatchString(c.GetHeader(h.UserClientHeader)) {
			gonic.Gonic(gatewayErrors.ErrInvalidformat().AddDetailsErr(h.ErrInvalidFormat(h.UserClientHeader)), c)
			HeaderEntry(h.UserClientHeader, c.GetHeader(h.UserClientHeader)).WithError(h.ErrInvalidFormat(h.UserClientHeader)).Warn("Invalid header")
			return
		}
		c.Next()
	}
}

func setHeader(h *http.Header, key string, value string) {
	h.Add(key, value)
	HeaderEntry(key, value).Debug("Header added")
}

func clearHeaders(h *http.Header, source string) {
	for key, value := range *h {
		if XHeaderRegexp.MatchString(key) {
			h.Del(key)
			HeaderEntry(key, value).Debug("Header deleted from " + source)
		}
	}
}

//HeaderEntry return logrus Entry with Header and Value params
func HeaderEntry(header string, value interface{}) *log.Entry {
	return log.WithFields(log.Fields{
		"Header": header,
		"Value":  value,
	})
}
