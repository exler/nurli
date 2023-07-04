package server

import (
	"encoding/json"
	"net/http"
)

func Unauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
	w.WriteHeader(http.StatusUnauthorized)
	if _, err := w.Write([]byte(http.StatusText(http.StatusUnauthorized))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteAsJSON(w http.ResponseWriter, obj interface{}) error {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(obj); err != nil {
		return err
	}

	return nil
}
