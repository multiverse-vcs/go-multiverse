package remote

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

// SocketFile is the name of the unix socket file.
const SocketFile = "rpc.sock"

// SocketAddr returns the RPC socket address.
func SocketAddr(home string) string {
	return filepath.Join(home, DotDir, SocketFile)
}

// NewClient returns a new RPC client.
func NewClient(home string) (*rpc.Client, error) {
	return rpc.DialHTTP("unix", SocketAddr(home))
}

// ListenAndServe starts the RPC listener.
func ListenAndServe(home string, server *Server) error {
	socket := SocketAddr(home)
	if err := os.RemoveAll(socket); err != nil {
		return err
	}

	rpc.RegisterName("Remote", server)
	rpc.HandleHTTP()

	unix, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	return http.Serve(unix, nil)
}
