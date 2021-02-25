package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
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

// Service is an HTTP gateway.
type Service struct {
	*remote.Server
}

// ListenAndServe starts the HTTP listener.
func ListenAndServe(server *remote.Server, bindAddr string) error {
	service := &Service{server}

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(cors)
	router.Handler(http.MethodGet, "/:peer/:repo", Route(service.Fetch))
	router.Handler(http.MethodPost, "/:peer/:repo/:branch", Route(service.Push))
	router.Handler(http.MethodGet, "/:peer/:repo/:branch/*file", Route(service.File))

	return http.ListenAndServe(bindAddr, router)
}

// cors handles cross origin pre-flight requests.
func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
}
