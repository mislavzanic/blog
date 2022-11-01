package posts

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"html/template"

	"github.com/Depado/bfchroma"
	"github.com/russross/blackfriday/v2"
	"github.com/alecthomas/chroma/formatters/html"
)

func renderer() *bfchroma.Renderer {
	return bfchroma.NewRenderer(
		bfchroma.WithoutAutodetect(),
		bfchroma.ChromaOptions(
			html.WithLineNumbers(false),
		),
		bfchroma.Extend(
			blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
				Flags: blackfriday.CommonHTMLFlags,
			}),
		),
		bfchroma.Style("doom-one"),
	)
}

func markDowner(args ...interface{}) template.HTML {
	content := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)), blackfriday.WithRenderer(renderer()), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return template.HTML(content)
}

func getUrl(path string) string {
	return fmt.Sprintf("/blog/%s", strings.Split(strings.Split(path, "/")[1], ".")[0])
}

func renderFromTemplate(w http.ResponseWriter, templateName string, templatePath string, funcMap map[string]any, data any) {
	html_template, _ := ioutil.ReadFile(templatePath)
	t := template.New(templateName)

	if funcMap != nil {
		t.Funcs(funcMap)
	}

	tmpl := template.Must(t.Parse(string(html_template)))

    if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		log.Fatal(err)
	}
}
