package rpc

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// SocketAddr is RPC socket address.
const SocketAddr = "localhost:9001"

// ErrDialRPC is an error message for failed RPC connections.
var ErrDialRPC = errors.New(`
Could not connect to local RPC server.
Make sure the Multiverse daemon is up.
See 'multi help daemon' for more info.
`)

// Service wraps a remote and provides RPC.
type Service struct {
	*remote.Server
}

// NewClient returns a new RPC client.
func NewClient() (*rpc.Client, error) {
	return jsonrpc.Dial("tcp", SocketAddr)
}

// ListenAndServe starts the RPC listener.
func ListenAndServe(server *remote.Server) error {
	rpc.RegisterName("Remote", &Service{server})
	rpc.RegisterName("Repo", &repo.Service{server})

	listen, err := net.Listen("tcp", SocketAddr)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go jsonrpc.ServeConn(conn)
	}
	
	return nil
}
