package server

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/exler/nurli/internal"
	"github.com/exler/nurli/internal/core"
)

func (sh *ServerHandler) prepareTemplates() (err error) {
	funcMap := template.FuncMap{
		"stringIn": core.StringIn,
		"domain":   core.GetDomainFromURL,
	}

	// `templates/**/*.html` doesn't pick up the files in the `templates/` directory
	// so we have to add the top directory to the patterns list.
	sh.templates, err = template.New("").Funcs(funcMap).ParseFS(internal.TemplateFS, "templates/*.html", "templates/**/*.html")
	return err
}

func (sh *ServerHandler) renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// We have to clone the templates to ensure that the {{ content }} is actually from the desired template
	// and not from the first parsed.
	clone := template.Must(sh.templates.Clone())
	clone = template.Must(clone.ParseFS(internal.TemplateFS, "templates/"+tmpl+".html"))
	// Get only the filename from the path as the templates are added into the namespace
	// by their filename only.
	if slashIndex := strings.LastIndex(tmpl, "/"); slashIndex != -1 {
		tmpl = tmpl[slashIndex+1:]
	}
	err := clone.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
