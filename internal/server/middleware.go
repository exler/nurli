package server

import (
	"context"
	"net/http"

	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *database.User
		var err error
		if user, err = sh.validateUserSession(r); err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), authenticatedUserKey, user)
		requestWithUser := r.WithContext(ctxWithUser)

		next.ServeHTTP(w, requestWithUser)
	})
}
