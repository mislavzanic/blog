package md2html

import (
	"github.com/mislavzanic/blog/src/app/md2html/renderer"
	fences "github.com/stefanfritsch/goldmark-fences"
	// highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gr "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type md2html struct {
	md goldmark.Markdown
	reader text.Reader
}

func NewParser() parser.Parser {
	return parser.NewParser(parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}

func NewMd2HTML(data []byte) md2html {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			&fences.Extender{},
			emoji.Emoji, 
			// highlighting.NewHighlighting(
			// 	highlighting.WithStyle("xcode-dark"),
			// ),
		),
		goldmark.WithParser(
			NewParser(),
		),
		goldmark.WithRenderer(
			gr.NewRenderer(
				gr.WithNodeRenderers(
					util.Prioritized(renderer.NewCodeRenderer(), 900),
					util.Prioritized(html.NewRenderer(), 1000),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		
	)

	return md2html{
		md: md,
		reader: text.NewReader(data),
	}
}
