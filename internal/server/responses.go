package server

import "net/http"

func Unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized.\n"))
}
