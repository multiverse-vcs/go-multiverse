// Package html contains methods for rendering web pages.
package html

import (
	"html/template"
	"path"
	"strings"
)

// funcs contains template helpers.
var funcs = template.FuncMap{
	"join": func(parts ...string) string {
		return path.Join(parts...)
	},
	"split": func(url string) []string {
		return strings.Split(url, "/")
	},
	"base": func(url string) string {
		return path.Base(url)
	},
}

// compile is used to compile the templates for development.
func compile(page string) *template.Template {
	return template.Must(template.New("index.html").Funcs(funcs).ParseFiles("html/index.html", page))
}
