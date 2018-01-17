package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"git.containerum.net/ch/grpc-proto-files/common"

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

			res, err := (*authClient).CreateToken(context.Background(), &auth.CreateTokenRequest{
				UserAgent:   r.Header.Get(userAgentXHeaderName),
				Fingerprint: r.Header.Get(userClientXHeaderName),
				UserId:      &common.UUID{Value: "396b8880-02f5-4ca0-832e-90c2b7bb543c"},
				UserIp:      r.Header.Get(userIPXHeaderName),
				UserRole:    auth.Role_USER,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.WithField("Token", res.AccessToken).Info("Create token success")

			token, err := (*authClient).CheckToken(context.Background(), &auth.CheckTokenRequest{
				AccessToken: accessToken,
				UserAgent:   r.Header.Get(userAgentXHeaderName),
				FingerPrint: r.Header.Get(userClientXHeaderName),
				UserIp:      r.Header.Get(userIPXHeaderName),
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				b, _ := json.Marshal(errorInvalidTokenMsg)
				w.Write([]byte(b))
				log.WithError(err).Warn(ErrInvalidToken)
				return
			}
			log.WithField("Token", token).Debug("Valid token")
			next.ServeHTTP(w, r)
		})
	}
}

func IsAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

//User Id Header
//User Role Header
