// Package rpc implements a remote procedure call server.
package rpc

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/peer"
)

// Service implements an RPC service.
type Service struct {
	client *peer.Client
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
func ListenAndServe(client *peer.Client) error {
	sock, err := SockAddr()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(sock); err != nil {
		return err
	}

	rpc.Register(&Service{client})
	rpc.HandleHTTP()

	listener, err := net.Listen("unix", sock)
	if err != nil {
		return err
	}

	return http.Serve(listener, nil)
}
