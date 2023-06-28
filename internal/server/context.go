package server

import (
	"net/http"

	"github.com/exler/nurli/internal/database"
)

type contextKey int

const authenticatedUserKey contextKey = 0

func getUserFromRequest(r *http.Request) *database.User {
	return r.Context().Value(authenticatedUserKey).(*database.User)
}
