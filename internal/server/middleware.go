package server

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuthMiddleware(realm, username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			if len(auth) != 2 || auth[0] != "Basic" {
				Unauthorized(w, realm)
				return
			}

			payload, _ := base64.StdEncoding.DecodeString(auth[1])
			pair := strings.SplitN(string(payload), ":", 2)

			if len(pair) != 2 || (pair[0] != username || pair[1] != password) {
				Unauthorized(w, realm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
