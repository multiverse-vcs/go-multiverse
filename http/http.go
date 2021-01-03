// Package http implements a web server.
package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/html"
	"github.com/multiverse-vcs/go-multiverse/node"
)

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	node *node.Node
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(node *node.Node) error {
	server := Server{node}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", server.home)
	router.HandlerFunc(http.MethodGet, "/:repo", server.file)
	router.HandlerFunc(http.MethodGet, "/:repo/*file", server.file)

	return http.ListenAndServe(BindAddr, router)
}

func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	if err := html.Home(w, req, s.node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) file(w http.ResponseWriter, req *http.Request) {
	if err := html.File(w, req, s.node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
