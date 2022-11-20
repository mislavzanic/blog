package posts

import (
	"fmt"
	"log"
	"time"
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
		bfchroma.Style("solarized-dark"),
	)
}

func AfterEpoch(t time.Time) bool {
    return t.After(time.Time{})
}

func ToMarkdown(args ...interface{}) template.HTML {
	content := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)), blackfriday.WithRenderer(renderer()), blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return template.HTML(content)
}

func GetUrl(uri string) func(string) string {
	return func(path string) string {
		return fmt.Sprintf("/%s/%s", uri, strings.Split(strings.Split(path, "/")[2], ".")[0])
	}
}

func RenderFromTemplate(w http.ResponseWriter, templateName string, templatePath string, funcMap map[string]any, data any) {
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
