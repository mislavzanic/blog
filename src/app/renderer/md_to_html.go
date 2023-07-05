package renderer

import (
	"fmt"
	"time"
	"io"
	"strings"
	"html/template"

	md2html "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

func initCodeRenderer() {
	htmlFormatter = html.New(html.WithLineNumbers(false))
	if htmlFormatter == nil {
		panic("couldn't create html formatter")
	}
	styleName := "xcode-dark"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
}

func AfterEpoch(t time.Time) bool {
    return t.After(time.Time{})
}

func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}

func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}

	if leafNode, ok := node.(*Sidenote); ok {
	    renderSidenote(w, leafNode, entering)
	    return ast.GoToNext, true
	}

	return ast.GoToNext, false
}

func ToMarkdown(md ...interface{}) template.HTML {
	initCodeRenderer()
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	p.Opts.ParserHook = ParserHook
	p.Opts.Flags = parser.SkipFootnoteList
	doc := p.Parse([]byte(fmt.Sprintf("%s", md...)))

	htmlFlags := md2html.CommonFlags | md2html.HrefTargetBlank
	opts := md2html.RendererOptions{Flags: htmlFlags, RenderNodeHook: myRenderHook}
	renderer := md2html.NewRenderer(opts)

	return template.HTML(markdown.Render(doc, renderer))
}

func GetUrl(uri string) func(string) string {
	return func(path string) string {
		return fmt.Sprintf("/%s/%s", uri, strings.Split(strings.Split(path, "/")[2], ".")[0])
	}
}

