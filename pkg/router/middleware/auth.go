package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"git.containerum.net/ch/grpc-proto-files/auth"

	log "github.com/Sirupsen/logrus"
)

const (
	authorizationHeader = "User-Token"
)

var (
	//ErrUnableParseToken - error when cant' parse JWT token
	ErrUnableParseToken = errors.New("Unable to parse JWT token")
	//ErrInvalidToken - error when token is invalid
	ErrInvalidToken = errors.New("Invalid access token")
	//ErrInvalidRole- error when user has "bad" role
	ErrInvalidRole = errors.New("Invalid user role")
	//ErrAccessForbidden - error when access is forbidden
	ErrAccessForbidden = errors.New("Access forbidden")

	errorInvalidTokenMsg = struct {
		Error string
	}{
		Error: ErrInvalidToken.Error(),
	}
)

func CheckAuthToken(authClient *auth.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get(authorizationHeader)
			log.WithFields(log.Fields{
				"useragent":   w.Header().Get(userAgentXHeaderName),
				"finger":      w.Header().Get(userClientXHeaderName),
				"ip":          w.Header().Get(userIPXHeaderName),
				"accessToken": accessToken,
			}).Debug("Check token info")

			token, err := (*authClient).CheckToken(context.Background(), &auth.CheckTokenRequest{
				AccessToken: accessToken,
				UserAgent:   w.Header().Get(userAgentXHeaderName),
				FingerPrint: w.Header().Get(userClientXHeaderName),
				UserIp:      w.Header().Get(userIPXHeaderName),
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				b, _ := json.Marshal(errorInvalidTokenMsg)
				w.Write([]byte(b))
				log.WithError(err).Warn(ErrInvalidToken)
				return
			}
			ns, vol, err := encodeAccessToBase64(token.GetAccess())
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				b, _ := json.Marshal(errorInvalidTokenMsg)
				w.Write([]byte(b))
				log.WithError(err).Warn(ErrInvalidToken)
				return
			}
			w.Header().Add(userIDXHeaderName, token.UserId.Value)
			log.WithField("Name", userIDXHeaderName).WithField("Value", token.UserId.Value).Debug("Add X-Header")

			w.Header().Add(userRoleHeaderName, token.UserRole)
			log.WithField("Name", userRoleHeaderName).WithField("Value", token.UserRole).Debug("Add X-Header")

			w.Header().Add(tokenIDXHeaderName, token.TokenId.Value)
			log.WithField("Name", tokenIDXHeaderName).WithField("Value", token.TokenId.Value).Debug("Add X-Header")

			w.Header().Add(userNamespacesXHeaderName, ns)
			log.WithField("Name", userNamespacesXHeaderName).WithField("Value", ns).Debug("Add X-Header")

			w.Header().Add(userVolumesXHeaderName, vol)
			log.WithField("Name", userVolumesXHeaderName).WithField("Value", vol).Debug("Add X-Header")

			next.ServeHTTP(w, r)
		})
	}
}

func IsAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if w.Header().Get(userRoleHeaderName) == "admin" {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusForbidden)
			b, _ := json.Marshal(errorInvalidTokenMsg)
			log.Warn(ErrAccessForbidden)
			w.Write([]byte(b))
		})
	}
}

func CheckUserRole(targetRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := w.Header().Get(userRoleHeaderName)
			for _, role := range targetRoles {
				if role == userRole {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusForbidden)
			b, _ := json.Marshal(errorInvalidTokenMsg)
			log.WithField("Role", w.Header().Get(userRoleHeaderName)).Warn(ErrInvalidRole)
			w.Write([]byte(b))
		})
	}
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
