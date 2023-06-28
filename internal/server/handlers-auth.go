package server

import (
	"net/http"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user database.User
		sh.DB.Where("username = ?", username).First(&user)
		if user.ID == 0 || !core.CheckPasswordHash(password, user.Password) {
			sh.renderTemplate(w, "login", map[string]interface{}{"error": "Invalid username or password"})
			return
		}

		sh.createUserSession(w, &user)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		sh.renderTemplate(w, "login", nil)
	}
}

func (sh *ServerHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sh.invalidateUserSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}
