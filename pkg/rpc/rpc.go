package rpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// SockFile is the name of the unix socket file.
const SockFile = "rpc.sock"

// DialErrMsg is an error message for failed RPC connections.
var DialErrMsg = `
Could not connect to local RPC server.
Make sure the Multiverse daemon is up.
See 'multi help daemon' for more info.
`

// NewClient returns a new RPC client.
func NewClient() (*rpc.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	socket := filepath.Join(home, remote.DotDir, SockFile)
	return rpc.DialHTTP("unix", socket)
}

// ListenAndServe starts the RPC listener.
func ListenAndServe(server *remote.Server) error {
	socket := filepath.Join(server.Root, SockFile)
	if err := os.RemoveAll(socket); err != nil {
		return err
	}

	rpc.RegisterName("Author", &author.Service{server})
	rpc.RegisterName("Repo", &repo.Service{server})
	rpc.HandleHTTP()

	listener, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	return http.Serve(listener, nil)
}
