// Package web implements a web server.
package web

import (
	"embed"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

//go:embed static/*
var static embed.FS

//go:embed html/*
var templates embed.FS

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	node *peer.Node
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(node *peer.Node) error {
	server := &Server{node}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", server.Index)
	router.Handler(http.MethodGet, "/:peer_id", View(server.Author))
	router.Handler(http.MethodGet, "/:peer_id/repositories/:name/:refs/:head/tree", View(server.Tree))
	router.Handler(http.MethodGet, "/:peer_id/repositories/:name/:refs/:head/tree/*file", View(server.Tree))

	http.Handle("/", router)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	return http.ListenAndServe(BindAddr, nil)
}

// Index redirects to the current author page.
func (s *Server) Index(w http.ResponseWriter, req *http.Request) {
	peerID := s.node.PeerID().String()
	url := path.Join("/", peerID)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
