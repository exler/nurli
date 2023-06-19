package server

import "net/http"

func (sh *ServerHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	sh.renderTemplate(w, "index", nil)
}
