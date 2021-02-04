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

//go:embed html/*
var templates embed.FS

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	client *peer.Client
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(client *peer.Client) error {
	server := &Server{client}

	router := httprouter.New()
	router.Handler(http.MethodGet, "/", View(server.Home))
	router.Handler(http.MethodGet, "/:name/:refs/:head/tree", View(server.Tree))
	router.Handler(http.MethodGet, "/:name/:refs/:head/tree/*file", View(server.Tree))

	http.Handle("/", router)
	http.Handle("/static/", http.FileServer(http.FS(static)))
	return http.ListenAndServe(BindAddr, nil)
}
