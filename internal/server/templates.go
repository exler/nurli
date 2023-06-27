package server

import (
	"html/template"
	"net/http"

	"github.com/exler/nurli/internal"
	"github.com/exler/nurli/internal/core"
)

func (sh *ServerHandler) prepareTemplates() (err error) {
	funcMap := template.FuncMap{
		"domain": core.GetDomainFromURL,
	}

	// `templates/**/*.html` doesn't pick up the files in the `templates/` directory
	// so we have to add the top directory to the patterns list.
	sh.templates, err = template.New("").Funcs(funcMap).ParseFS(internal.TemplateFS, "templates/*.html", "templates/**/*.html")
	return err
}

func (sh *ServerHandler) renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := sh.templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
