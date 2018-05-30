package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"git.containerum.net/ch/api-gateway/pkg/gatewayErrors"
	"git.containerum.net/ch/auth/proto"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	h "github.com/containerum/utils/httputil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	//ErrAuthClientNotSet returns if grpc client is nil
	ErrAuthClientNotSet = errors.New("Auth client not set up")
	//ErrUserPermissionDenied return if user don't have permissions
	ErrUserPermissionDenied = errors.New("User permission denied")
)

//CheckAuth check user token and roles
func CheckAuth(roles []string, authClient *authProto.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}
		if authClient == nil {
			log.Error(ErrAuthClientNotSet.Error())
			gonic.Gonic(gatewayErrors.ErrInternal(), c)
			return
		}
		accessToken := c.Request.Header.Get(h.AuthorizationHeader)
		userAgent, userFinger := c.GetHeader(h.UserAgentXHeader), c.GetHeader(h.UserClientXHeader)
		userIP := c.ClientIP()
		getTokenEntry(accessToken, userAgent, userFinger, userIP).Debug("Check token")
		token, err := (*authClient).CheckToken(context.Background(), &authProto.CheckTokenRequest{
			AccessToken: accessToken,
			UserAgent:   userAgent,
			FingerPrint: userFinger,
			UserIp:      userIP,
		})

		switch err := err.(type) {
		case nil:
			// pass
		case *cherry.Err:
			log.WithError(err).Warnf("CheckToken error")
			c.AbortWithStatusJSON(err.StatusHTTP, err)
			return
		default:
			log.WithError(err).Errorf("Something is wrong with Auth server")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if ok := checkUserRole(token.UserRole, roles); !ok {
			log.WithError(gatewayErrors.ErrUserPermissionDenied()).Warnf(ErrUserPermissionDenied.Error())
			gonic.Gonic(gatewayErrors.ErrUserPermissionDenied(), c)
			return
		}
		ns, vol, err := encodeAccessToBase64(token.GetAccess())
		if err != nil {
			log.Error(err)
			gonic.Gonic(gatewayErrors.ErrInternal(), c)
			return
		}
		setHeader(&c.Request.Header, h.TokenIDXHeader, token.TokenId)
		setHeader(&c.Request.Header, h.UserIDXHeader, token.UserId)
		setHeader(&c.Request.Header, h.UserRoleXHeader, token.UserRole)
		setHeader(&c.Request.Header, h.UserNamespacesXHeader, ns)
		setHeader(&c.Request.Header, h.UserVolumesXHeader, vol)
	}
}

func checkUserRole(userRole string, roles []string) bool {
	if len(roles) == 0 {
		return true
	}
	for _, role := range roles {
		switch role {
		case "*":
			return true
		case "user":
			if userRole == "user" {
				return true
			}
		case "admin":
			if userRole == "admin" {
				return true
			}
		}
	}
	return false
}

func encodeAccessToBase64(access *authProto.ResourcesAccess) (ns string, vol string, err error) {
	userNamespaces := access.GetNamespace()
	userVolumes := access.GetVolume()
	bNs, e := json.Marshal(userNamespaces)
	if e != nil {
		err = e
		return
	}
	ns = base64.StdEncoding.EncodeToString(bNs)
	bVol, e := json.Marshal(userVolumes)
	if e != nil {
		err = e
		return
	}
	vol = base64.StdEncoding.EncodeToString(bVol)
	return
}

func getTokenEntry(token, agent, finger, ip string) *log.Entry {
	return log.WithFields(log.Fields{
		"AccessToken": token,
		"UserAgent":   agent,
		"FingerPrint": finger,
		"UserIp":      ip,
	})
}
