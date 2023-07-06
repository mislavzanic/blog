package md2html

import (
	"github.com/mislavzanic/blog/src/app/md2html/renderer"

	"github.com/yuin/goldmark"
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

func NewMd2HTML(data []byte) md2html {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRenderer(
			gr.NewRenderer(
				gr.WithNodeRenderers(
					util.Prioritized(renderer.NewCodeRenderer(), 900),
					util.Prioritized(html.NewRenderer(), 1000),
				),
			),
		),
		
	)

	return md2html{
		md: md,
		reader: text.NewReader(data),
	}
}
