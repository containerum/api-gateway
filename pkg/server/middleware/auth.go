package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"git.containerum.net/ch/kube-client/pkg/cherry/auth"

	"git.containerum.net/ch/auth/proto"
	"git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	"git.containerum.net/ch/kube-client/pkg/cherry/api-gateway"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//CheckAuth check user token and roles
func CheckAuth(roles []string, authClient *authProto.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}
		if authClient == nil {
			log.Errorf("auth client not set up")
			gonic.Gonic(gatewayErrors.ErrInternal(), c)
			return
		}
		accessToken := c.Request.Header.Get(authorizationHeader)
		userIP := c.ClientIP()
		log.WithFields(log.Fields{
			"AccessToken": accessToken,
			"UserAgent":   c.GetHeader(userAgentXHeader),
			"FingerPrint": c.GetHeader(userClientXHeader),
			"UserIp":      userIP,
		}).Debug("Check token")
		token, err := (*authClient).CheckToken(context.Background(), &authProto.CheckTokenRequest{
			AccessToken: accessToken,
			UserAgent:   c.GetHeader(userAgentXHeader),
			FingerPrint: c.GetHeader(userClientXHeader),
			UserIp:      userIP,
		})
		switch err := err.(type) {
		case nil:
			// pass
		case *cherry.Err:
			log.WithError(err).Warnf("CheckToken() error")
			switch {
			case cherry.In(autherr.ErrInvalidToken(),
				autherr.ErrTokenNotOwnedBySender()):
				gonic.Gonic(autherr.ErrInvalidToken(), c)
			default:
				gonic.Gonic(err, c)
			}
			return
		default:
			log.WithError(err).Errorf("internal error while token checking")
			gonic.Gonic(gatewayErrors.ErrInternal(), c)
			return
		}
		if ok := checkUserRole(token.UserRole, roles); !ok {
			log.WithError(gatewayErrors.ErrUserPermissionDenied()).Warnf("user permission denied")
			gonic.Gonic(gatewayErrors.ErrUserPermissionDenied(), c)
			return
		}
		ns, vol, err := encodeAccessToBase64(token.GetAccess())
		if err != nil {
			log.Error(err)
			gonic.Gonic(gatewayErrors.ErrInternal(), c)
			return
		}
		setHeader(&c.Request.Header, tokenIDXHeader, token.TokenId)
		setHeader(&c.Request.Header, userIDXHeader, token.UserId)
		setHeader(&c.Request.Header, userRoleXHeader, token.UserRole)
		setHeader(&c.Request.Header, userNamespacesXHeader, ns)
		setHeader(&c.Request.Header, userVolumesXHeader, vol)
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
