// Package web implements a web server.
package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/multiverse-vcs/go-multiverse/web/view"
)

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:2020"

// Server contains http services.
type Server struct {
	node *node.Node
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(node *node.Node) error {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/", view.Home(node))
	router.Handler(http.MethodGet, "/:name", view.Repo(node))
	router.Handler(http.MethodGet, "/:name/:page", view.Repo(node))
	router.Handler(http.MethodGet, "/:name/:page/:ref", view.Repo(node))
	router.Handler(http.MethodGet, "/:name/:page/:ref/*file", view.Repo(node))

	var static http.Handler
	static = http.FileServer(http.Dir("web/static"))
	static = http.StripPrefix("/static", static)

	http.Handle("/", router)
	http.Handle("/static/", static)
	return http.ListenAndServe(BindAddr, nil)
}
