package http

import (
	"net/http"
)

// Route is an http handler that returns an error
type Route func(http.ResponseWriter, *http.Request) error

// ServeHTTP serves the request and handles any errors.
func (r Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if err := r(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// CORS handles cross origin pre-flight requests.
func CORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
}
