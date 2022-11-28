package renderer

import (
	"html/template"
	"log"
	"net/http"
)


func RenderFromTemplate(w http.ResponseWriter, templateName string, templates []string, funcMap map[string]any, data any) {
	t := template.New(templateName)

	if funcMap != nil {
		t.Funcs(funcMap)
	}

	tmpl := template.Must(t.ParseFiles(templates...))

    if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		log.Fatal(err)
	}
}
