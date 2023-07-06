package renderer

import (
	"log"
	"bytes"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

var (
	style = "xcode-dark"
)

type CodeBlockRenderer struct {
	formatter *html.Formatter
	highlighter *chroma.Style
}

func (r CodeBlockRenderer) highlight(w util.BufWriter, source, lang string) error {
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
	return r.formatter.Format(w, r.highlighter, it)
}

func (r CodeBlockRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
	    return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := string(n.Language(source))
	l := n.Lines().Len()

	var b bytes.Buffer

	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		b.Write(line.Value(source))
	}

	r.highlight(w, b.String(), lang)

	return ast.WalkContinue, nil
}

func (r CodeBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderCodeBlock)
}

func NewCodeRenderer() renderer.NodeRenderer {
	formatter := html.New(html.WithLineNumbers(false))
	if formatter == nil {
		log.Fatal("couldn't create html formatter")
	}

	highlighter := styles.Get(style)
	if highlighter == nil {
		log.Fatalf("didn't find style '%s'", style)
	}

	return CodeBlockRenderer{
		formatter: formatter,
		highlighter: highlighter,
	}
}
