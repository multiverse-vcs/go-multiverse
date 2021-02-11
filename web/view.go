package web

import (
	"bytes"
	"html/template"
	"net/http"
)

// TODO uncomment
//var layout = template.Must(template.New("index.html").Funcs(funcs).ParseFS(views, "views/*"))

// View is an http handler that renders a view.
type View func(http.ResponseWriter, *http.Request) (*ViewModel, error)

// ViewModel contains the data for the template.
type ViewModel struct {
	Name string
	Data interface{}
}

// ServeHTTP handles http requests to a route.
func (v View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO remove
	layout := template.Must(template.New("index.html").Funcs(funcs).ParseGlob("web/views/*"))

	model, err := v(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if model == nil {
		return
	}

	var page bytes.Buffer
	if err := layout.ExecuteTemplate(&page, model.Name, model.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := template.HTML(page.String())
	if err := layout.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
