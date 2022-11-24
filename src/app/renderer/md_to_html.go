package renderer

import (
	"fmt"
	"time"
	"strings"
	"html/template"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/russross/blackfriday/v2"
	"github.com/Depado/bfchroma"
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
		bfchroma.Style("xcode-dark"),
	)
}

func AfterEpoch(t time.Time) bool {
    return t.After(time.Time{})
}

func ToMarkdown(args ...interface{}) template.HTML {
	content := blackfriday.Run(
		[]byte(
			fmt.Sprintf("%s", args...),
		),
		blackfriday.WithRenderer(renderer()),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)
	return template.HTML(content)
}

func GetUrl(uri string) func(string) string {
	return func(path string) string {
		return fmt.Sprintf("/%s/%s", uri, strings.Split(strings.Split(path, "/")[2], ".")[0])
	}
}

