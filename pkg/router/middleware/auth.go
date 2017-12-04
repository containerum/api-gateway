package middleware

import (
	"fmt"
	"net/http"

	"git.containerum.net/ch/grpc-proto-files/auth"
)

func CheckAuthToken(authClient *auth.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Print("CHECK TOKEN")
		})
	}
}
