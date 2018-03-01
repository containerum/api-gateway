package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	errs "git.containerum.net/ch/api-gateway/pkg2/errors"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/gin-gonic/gin"
)

var (
	errCheckToken          = errors.New("Unable to check token")
	errInvalidToken        = errors.New("Invalid token")
	errUnableGetAccess     = errors.New("Unable to get user access")
	errUserPermisionDenied = errors.New("User permission denied")
)

//CheckAuth check user token and roles
func CheckAuth(roles []string, authClient *auth.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}
		if authClient == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		accessToken := c.Request.Header.Get(authorizationHeader)
		token, err := (*authClient).CheckToken(context.Background(), &auth.CheckTokenRequest{
			AccessToken: accessToken,
			UserAgent:   c.GetHeader(userAgentXHeader),
			FingerPrint: c.GetHeader(userClientXHeader),
			UserIp:      c.GetHeader(userIPXHeader),
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errs.New(errCheckToken.Error(), errInvalidToken.Error()))
			return
		}
		if ok := checkUserRole(token.UserRole, roles); !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, errs.New(errCheckToken.Error(), errUserPermisionDenied.Error()))
			return
		}
		ns, vol, err := encodeAccessToBase64(token.GetAccess())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errs.New(errCheckToken.Error(), errUnableGetAccess.Error()))
			return
		}
		setHeader(&c.Request.Header, tokenIDXHeader, token.TokenId.Value)
		setHeader(&c.Request.Header, userIDXHeader, token.UserId.Value)
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

func encodeAccessToBase64(access *auth.ResourcesAccess) (ns string, vol string, err error) {
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
