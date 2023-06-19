package server

import "net/http"

func (sh *ServerHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
