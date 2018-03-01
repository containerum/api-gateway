package middleware

import (
	"errors"
	"net/http"
	"regexp"

	errs "git.containerum.net/ch/api-gateway/pkg2/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	requestIDXHeader      = "X-Request-ID"
	userIDXHeader         = "X-User-ID"
	userClientXHeader     = "X-User-Client"
	userAgentXHeader      = "X-User-Agent"
	userIPXHeader         = "X-Client-IP"
	tokenIDXHeader        = "X-Token-ID"
	userRoleXHeader       = "X-User-Role"
	userNamespacesXHeader = "X-User-Namespace"
	userVolumesXHeader    = "X-User-Volume"
	userHideDataXHeader   = "X-User-Hide-Data"
)

const (
	userClientHeader    = "User-Client"
	authorizationHeader = "User-Token"
)

var (
	//XHeaderRegexp keeps regexp for detecting X-Headers
	XHeaderRegexp, _    = regexp.Compile("^X-[a-zA-Z0-9]+")
	userClientRegexp, _ = regexp.Compile("^[a-f0-9]{32}$")
)

var (
	//ErrInvalidUserClientHeader returns if User-Client header is empty or invalid
	ErrInvalidUserClientHeader = errors.New("Invalid User-Client header")
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
		setHeader(&c.Request.Header, userIPXHeader, c.ClientIP())
		setHeader(&c.Request.Header, userClientXHeader, c.GetHeader(userClientHeader))
		setHeader(&c.Request.Header, userAgentXHeader, c.Request.UserAgent())
	}
}

//CheckUserClientHeader validate User-Client header
func CheckUserClientHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(userClientHeader) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errs.New(ErrInvalidUserClientHeader.Error(), "Not provided"))
			return
		}
		if !userClientRegexp.MatchString(c.GetHeader(userClientHeader)) {
			log.WithFields(log.Fields{
				"Header": userClientHeader,
				"Value":  c.GetHeader(userClientHeader),
			}).Debug(ErrInvalidUserClientHeader)
			c.AbortWithStatusJSON(http.StatusInternalServerError, errs.New(ErrInvalidUserClientHeader.Error(), "Invalid format"))
			return
		}
		c.Next()
	}
}

func setHeader(h *http.Header, key string, value string) {
	h.Add(key, value)
	log.WithFields(log.Fields{
		"Header": key,
		"Value":  value,
	}).Debug("Header added")
}

func clearHeaders(h *http.Header, source string) {
	for k, v := range *h {
		if XHeaderRegexp.MatchString(k) {
			h.Del(k)
			log.WithFields(log.Fields{
				"Header": k,
				"Value":  v,
			}).Debug("Header deleted from " + source)
		}
	}
}
