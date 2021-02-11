// Package web implements a web server.
package web

import (
	"embed"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

//go:embed static/*
var static embed.FS

//go:embed views/*
var views embed.FS

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	node peer.Peer
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(node peer.Peer) error {
	server := Server{node}
	author := Author(server)
	repository := Repository(server)

	router := httprouter.New()
	router.Handler(http.MethodGet, "/", View(author.Index))
	router.Handler(http.MethodGet, "/:peer_id", View(author.Index))
	router.Handler(http.MethodPost, "/:peer_id/follow", View(author.Follow))
	router.Handler(http.MethodPost, "/:peer_id/unfollow", View(author.Unfollow))
	router.Handler(http.MethodGet, "/:peer_id/repositories/:name/:refs/:head/tree", View(repository.Index))
	router.Handler(http.MethodGet, "/:peer_id/repositories/:name/:refs/:head/tree/*file", View(repository.Index))

	http.Handle("/", router)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	return http.ListenAndServe(BindAddr, nil)
}
