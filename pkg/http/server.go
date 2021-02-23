package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
)

// BindAddr is the address the http server binds to.
const BindAddr = "localhost:2020"

// Service is an HTTP gateway.
type Service struct {
	*remote.Server
}

// ListenAndServe starts the HTTP listener.
func ListenAndServe(server *remote.Server) error {
	service := &Service{server}

	router := httprouter.New()
	router.GlobalOPTIONS = http.HandlerFunc(CORS)
	router.Handler(http.MethodGet, "/:peer/:repo", Route(service.Fetch))
	router.Handler(http.MethodPost, "/:peer/:repo/:branch", Route(service.Push))
	router.Handler(http.MethodGet, "/:peer/:repo/:branch/*file", Route(service.File))

	return http.ListenAndServe(BindAddr, router)
}
