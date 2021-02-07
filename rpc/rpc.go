// Package rpc implements a remote procedure call server.
package rpc

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/peer"
)

// ErrConnect is the human friendly error message for failed connections.
var ErrConnect = errors.New(`Could not connect to local RPC server.
Make sure the Multiverse daemon is up.
See 'multi help daemon' for more info.`)

// Service implements an RPC service.
type Service struct {
	node *peer.Node
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

	client, err := rpc.DialHTTP("unix", sock)
	if err != nil {
		return nil, ErrConnect
	}

	return client, nil
}

// ListenAndServer starts an RPC listener.
func ListenAndServe(node *peer.Node) error {
	sock, err := SockAddr()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(sock); err != nil {
		return err
	}

	rpc.Register(&Service{node})
	rpc.HandleHTTP()

	listener, err := net.Listen("unix", sock)
	if err != nil {
		return err
	}

	return http.Serve(listener, nil)
}
