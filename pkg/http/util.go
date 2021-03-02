package http

import (
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// formatter outputs highlighted code.
var formatter = html.New(
	html.WithLineNumbers(true),
	html.LineNumbersInTable(true),
)

// Highlight returns a syntax highlighted version of the given code.
func Highlight(name, source, theme string, w io.Writer) error {
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(theme)
	if style == nil {
		style = styles.Fallback
	}

	token, err := chroma.Coalesce(lexer).Tokenise(nil, source)
	if err != nil {
		return err
	}

	return formatter.Format(w, style, token)
}
