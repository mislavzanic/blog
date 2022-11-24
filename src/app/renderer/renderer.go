package renderer

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)


func RenderFromTemplate(w http.ResponseWriter, templateName string, templatePath string, funcMap map[string]any, data any) {
	html_template, err := ioutil.ReadFile(templatePath)

	if err != nil {
		log.Fatal(err)
	}

	t := template.New(templateName)

	if funcMap != nil {
		t.Funcs(funcMap)
	}

	tmpl := template.Must(t.Parse(string(html_template)))

    if err = tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		log.Fatal(err)
	}
}
