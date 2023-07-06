package md2html

import (
	"fmt"
	"strings"
	"time"
	"html/template"
	"log"
	"net/http"
	"bytes"
)

func GetUrl(uri string) func(string) string {
	return func(path string) string {
		return fmt.Sprintf("/%s/%s", uri, strings.Split(strings.Split(path, "/")[2], ".")[0])
	}
}

func AfterEpoch(t time.Time) bool {
    return t.After(time.Time{})
}

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

func ToMarkdown(md ...interface{}) template.HTML {
	var b bytes.Buffer

	data := []byte(fmt.Sprintf("%s", md...))
	blogProcessor := NewMd2HTML(data)
	node := blogProcessor.md.Parser().Parse(blogProcessor.reader)
	blogProcessor.md.Renderer().Render(&b, data, node)

	return template.HTML(b.Bytes())
}
