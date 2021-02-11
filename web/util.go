package web

import (
	"fmt"
	"html/template"
	"path"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
)

// formatter outputs highlighted code.
var formatter = html.New(
	html.WithLineNumbers(true),
	html.LineNumbersInTable(true),
	html.LinkableLineNumbers(true, ""),
)

// funcs contains template helpers.
var funcs = template.FuncMap{
	"join": func(parts ...string) string {
		return path.Join(parts...)
	},
	"split": func(url string) []string {
		return strings.Split(url, "/")
	},
	"dir": func(url string) string {
		return path.Dir(url)
	},
	"base": func(url string) string {
		return path.Base(url)
	},
	"timestamp": func(t time.Time) string {
		return t.Format(time.RFC822)
	},
	"shortcid": func(c string) string {
		return c[len(c)-10:]
	},
	"shortpeer": func(c string) string {
		return fmt.Sprintf("%s..%s", c[:2], c[len(c)-6:])
	},
	"breadcrumbs": breadcrumbs,
	"highlight":   highlight,
	"markdown":    markdown,
}

// breadcrumbs returns a list of ascending urls.
func breadcrumbs(url string) []string {
	var crumbs []string
	for p := url; p != "/"; p = path.Dir(p) {
		crumbs = append([]string{p}, crumbs...)
	}

	return crumbs
}

// highlight returns a syntax highlighted version of the given code.
func highlight(name, source string) (template.HTML, error) {
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

// markdown renders the given markdown source into html.
func markdown(source string) (template.HTML, error) {
	var b strings.Builder
	if err := goldmark.Convert([]byte(source), &b); err != nil {
		return "", err
	}

	return template.HTML(b.String()), nil
}
