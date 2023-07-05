package renderer

import (
	"bytes"
	"io"

	"github.com/gomarkdown/markdown/ast"
)

type Sidenote struct {
	ast.Container
	heading []byte
}

var sidenote = []byte("<sidenote")

func ParserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parseSidenote(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func parseHeading(data []byte) ([]byte, int) {
	end := bytes.IndexByte(data, byte('>'))
	return data[:end+1], end+1
}

func parseSidenote(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, sidenote) {
		return nil, nil, 0
	}
	heading, i := parseHeading(data)
	end := bytes.Index(data[i:], []byte("</sidenote>"))
	if end < 0 {
		return nil, data, 0
	}
	end = end + i
	res := &Sidenote{heading: heading}
	return res, data[i+1:end], end + len([]byte("</sidenote>"))
}

func renderSidenote(w io.Writer, s *Sidenote, entering bool) {
	if entering {
		io.WriteString(w, string(s.heading))
	} else {
		io.WriteString(w, "</sidenote>")
	}
}
