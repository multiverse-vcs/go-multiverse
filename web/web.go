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

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	client *peer.Client
}

// View is an http handler that renders a view.
type View func(http.ResponseWriter, *http.Request) error

// ServeHTTP handles http requests to a route.
func (v View) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := v(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(client *peer.Client) error {
	server := &Server{client}

	router := httprouter.New()
	router.Handler(http.MethodGet, "/", View(server.Home))
	router.Handler(http.MethodGet, "/:id", View(server.Tree))
	router.Handler(http.MethodGet, "/:id/:ref/commits", View(server.Commits))
	router.Handler(http.MethodGet, "/:id/:ref/tree", View(server.Tree))
	router.Handler(http.MethodGet, "/:id/:ref/tree/*file", View(server.Tree))
	router.Handler(http.MethodGet, "/:id/:ref/blob/*file", View(server.Blob))

	http.Handle("/", router)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	return http.ListenAndServe(BindAddr, nil)
}
