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
		hashed_password, err := core.HashPassword(password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var user database.User
		sh.DB.Where("username = ? AND password = ?", username, hashed_password).First(&user)
		if user.ID == 0 {
			err := sh.templates.ExecuteTemplate(w, "login.html", map[string]interface{}{"error": "Invalid username or password"})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		sh.createUserSession(r, &user)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		err := sh.templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (sh *ServerHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sh.invalidateUserSession(r)
	http.Redirect(w, r, "/login", http.StatusFound)
}
