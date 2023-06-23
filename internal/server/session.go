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

func (sh *ServerHandler) createUserSession(r *http.Request, user *database.User) error {
	sessionToken := generateSecureSessionToken()
	sh.DB.Create(&database.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	})
	r.AddCookie(&http.Cookie{
		Name:    "session",
		Value:   sessionToken,
		Expires: time.Now().Add(time.Hour * 24 * 30),
	})
	return nil
}

func (sh *ServerHandler) validateUserSession(r *http.Request) error {
	sessionToken := sh.getSessionToken(r)
	if sessionToken == "" {
		return fmt.Errorf("no session token provided")
	}

	var session database.Session
	sh.DB.Where("token = ?", sessionToken).First(&session)
	if session.ID == 0 {
		return fmt.Errorf("invalid session token")
	} else if session.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("session expired")
	}

	return nil
}

func (sh *ServerHandler) invalidateUserSession(r *http.Request) error {
	sessionToken := sh.getSessionToken(r)
	if sessionToken == "" {
		return fmt.Errorf("no session token provided")
	}

	sh.DB.Where("token = ?", sessionToken).Delete(&database.Session{})
	r.AddCookie(&http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Now().Add(time.Hour * -1),
	})
	return nil
}
