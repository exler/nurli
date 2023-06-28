package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/exler/nurli/internal/database"
)

func generateSecureSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (sh *ServerHandler) getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (sh *ServerHandler) createUserSession(w http.ResponseWriter, user *database.User) error {
	sessionToken := generateSecureSessionToken()
	sh.DB.Create(&database.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	})
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionToken,
		Expires: time.Now().Add(time.Hour * 24 * 30),
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (sh *ServerHandler) validateUserSession(r *http.Request) (*database.User, error) {
	sessionToken := sh.getSessionToken(r)
	if sessionToken == "" {
		return nil, fmt.Errorf("no session token provided")
	}

	var session database.Session
	sh.DB.Where("token = ?", sessionToken).Preload("User").First(&session)
	if session.ID == 0 {
		return nil, fmt.Errorf("invalid session token")
	} else if session.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("session expired")
	}

	return &session.User, nil
}

func (sh *ServerHandler) invalidateUserSession(w http.ResponseWriter, r *http.Request) error {
	sessionToken := sh.getSessionToken(r)
	if sessionToken == "" {
		return fmt.Errorf("no session token provided")
	}

	sh.DB.Where("token = ?", sessionToken).Delete(&database.Session{})
	cookie := http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now().Add(time.Hour * -1),
	}
	http.SetCookie(w, &cookie)
	return nil
}
