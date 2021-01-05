// Package view implements methods for rendering views.
package view

import (
	"net/http"
)

// View is an http route that renders a view.
type View func(http.ResponseWriter, *http.Request) error

// ServeHTTP handles http requests to a route.
func (v View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := v(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
