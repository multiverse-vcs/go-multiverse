package http

import (
	"context"
	"html/template"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-unixfs"
)

// Highlight returns a syntax highlighted version of the given code.
func (s *Server) Highlight(name, code string) (template.HTML, error) {
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("pastie")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithLineNumbers(true),
		html.LineNumbersInTable(true),
		html.LinkableLineNumbers(true, ""),
	)

	lexer = chroma.Coalesce(lexer)
	token, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	if err := formatter.Format(&b, style, token); err != nil {
		return "", err
	}

	return template.HTML(b.String()), nil
}

// IsDir returns true if the node with the given CID is a directory.
func (s *Server) IsDir(id cid.Cid) (bool, error) {
	ctx := context.Background()

	node, err := s.dag.Get(ctx, id)
	if err != nil {
		return false, err
	}

	fsnode, err := unixfs.ExtractFSNode(node)
	if err != nil {
		return false, err
	}

	return fsnode.IsDir(), nil
}
