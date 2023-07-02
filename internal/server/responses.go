package server

import "net/http"

func Unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	if _, err := w.Write([]byte("Unauthorized.\n")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
