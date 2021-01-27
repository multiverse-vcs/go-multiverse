// Package web implements a web server.
package web

import (
	"embed"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var templates embed.FS

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	client *peer.Client
	store  *data.Store
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
func ListenAndServe(client *peer.Client, store *data.Store) error {
	server := &Server{client, store}

	router := httprouter.New()
	router.Handler(http.MethodGet, "/", View(server.Home))
	router.Handler(http.MethodGet, "/:name", View(server.Tree))
	router.Handler(http.MethodGet, "/:name/:ref/commits", View(server.Commits))
	router.Handler(http.MethodGet, "/:name/:ref/tree", View(server.Tree))
	router.Handler(http.MethodGet, "/:name/:ref/tree/*file", View(server.Tree))
	router.Handler(http.MethodGet, "/:name/:ref/blob/*file", View(server.Blob))

	http.Handle("/", router)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	return http.ListenAndServe(BindAddr, nil)
}
