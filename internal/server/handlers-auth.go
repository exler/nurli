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
			err := sh.templates.ExecuteTemplate(w, "login.html", map[string]interface{}{"error": "Invalid username or password"})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		sh.createUserSession(w, &user)
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
	sh.invalidateUserSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}
