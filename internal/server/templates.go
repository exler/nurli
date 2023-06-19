package server

import (
	"html/template"
	"net/http"
)

func (sh *ServerHandler) prepareTemplates() (err error) {
	funcMap := template.FuncMap{}

	sh.templates, err = template.New("").Funcs(funcMap).ParseGlob("./internal/templates/*.html")
	return
}

func (sh *ServerHandler) renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := sh.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
