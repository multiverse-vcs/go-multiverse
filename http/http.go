// Package http implements an HTML template server.
package http

import (
	"net/http"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// BindAddr is the address the http server binds to.
const BindAddr = "127.0.0.1:8080"

// Server contains http services.
type Server struct {
	dag    ipld.DAGService
	dstore datastore.Batching
}

// ListenAndServe starts an HTTP listener.
func ListenAndServe(dag ipld.DAGService, dstore datastore.Batching) error {
	server := Server{
		dag:    dag,
		dstore: namespace.Wrap(dstore, data.Prefix),
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", server.home)
	router.HandlerFunc(http.MethodGet, "/:repo", server.tree)
	router.HandlerFunc(http.MethodGet, "/:repo/tree/*file", server.tree)
	router.HandlerFunc(http.MethodGet, "/:repo/blob/*file", server.blob)

	return http.ListenAndServe(BindAddr, router)
}
