package html

import (
	"context"
	"html/template"
	"path"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-unixfs"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/yuin/goldmark"
)

// util implements template utilities.
type util struct {
	node *node.Node
}

// formatter outputs highlighted code.
var formatter = html.New(
	html.WithLineNumbers(true),
	html.LineNumbersInTable(true),
	html.LinkableLineNumbers(true, "line-"),
)

// Breadcrumbs returns a list of ascending urls.
func (u *util) Breadcrumbs(url string) []string {
	var crumbs []string

	parts := strings.Split(strings.Trim(url, "/"), "/")
	for i := range parts {
		crumbs = append(crumbs, path.Join(parts[:i+1]...))
	}

	return crumbs
}

// Highlight returns a syntax highlighted version of the given code.
func (u *util) Highlight(name, source string) (template.HTML, error) {
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	token, err := chroma.Coalesce(lexer).Tokenise(nil, source)
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
func (u *util) IsDir(id cid.Cid) (bool, error) {
	f, err := u.node.Get(context.Background(), id)
	if err != nil {
		return false, err
	}

	fsnode, err := unixfs.ExtractFSNode(f)
	if err != nil {
		return false, err
	}

	return fsnode.IsDir(), nil
}

// Markdown renders the given markdown source into html.
func (u *util) Markdown(source string) (template.HTML, error) {
	var b strings.Builder
	if err := goldmark.Convert([]byte(source), &b); err != nil {
		return "", err
	}

	return template.HTML(b.String()), nil
}
