package middleware

import (
	"net/http"
	"regexp"

	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	"git.containerum.net/ch/kube-client/pkg/cherry/api-gateway"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	requestIDXHeader      = "X-Request-ID"
	requestNameXHeader    = "X-Request-Name"
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

//ClearXHeaders clear all request and response X-Headers
func ClearXHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.WithField("Headers", c.Request.Header).Debug("Header")

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

//SetMainUserXHeaders write X-User-IP, X-User-Client, X-User-Agent for next services
func SetHeaderFromQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClient, _ := c.GetQuery(userClientHeader)
		userToken, _ := c.GetQuery(authorizationHeader)
		if c.GetHeader(userClientHeader) == "" && userClient != "" {
			c.Request.Header.Add(userClientHeader, userClient)
		}
		if c.GetHeader(authorizationHeader) == "" && userToken != "" {
			c.Request.Header.Add(authorizationHeader, userToken)
		}
	}
}

//CheckUserClientHeader validate User-Client header
func CheckUserClientHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(userClientHeader) == "" {
			gonic.Gonic(gatewayErrors.ErrHeaderNotProvided().AddDetailF("header %s required", userClientHeader), c)
			return
		}
		if !userClientRegexp.MatchString(c.GetHeader(userClientHeader)) {
			log.WithFields(log.Fields{
				"Header": userClientHeader,
				"Value":  c.GetHeader(userClientHeader),
			}).Warnf("invalid header format")
			gonic.Gonic(gatewayErrors.ErrInvalidformat().AddDetailF("header %s has invalid format", userClientHeader), c)
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
