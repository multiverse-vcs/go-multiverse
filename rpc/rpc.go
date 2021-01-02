// Package rpc implements a remote procedure call server.
package rpc

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// Service implements an RPC service.
type Service struct {
	dag    ipld.DAGService
	dstore datastore.Batching
}

// SockAddr returns the unix sock file path.
func SockAddr() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".multiverse", "rpc.sock"), nil
}

// NewClient returns a new RPC client.
func NewClient() (*rpc.Client, error) {
	sock, err := SockAddr()
	if err != nil {
		return nil, err
	}

	return rpc.DialHTTP("unix", sock)
}

// ListenAndServer starts an RPC listener.
func ListenAndServe(dag ipld.DAGService, dstore datastore.Batching) error {
	sock, err := SockAddr()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(sock); err != nil {
		return err
	}

	service := Service{
		dag:    dag,
		dstore: namespace.Wrap(dstore, data.Prefix),
	}

	rpc.Register(&service)
	rpc.HandleHTTP()

	listener, err := net.Listen("unix", sock)
	if err != nil {
		return err
	}

	return http.Serve(listener, nil)
}
