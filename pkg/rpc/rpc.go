package rpc

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/file"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// DialErrMsg is an error message for failed RPC connections.
var DialErrMsg = `
Could not connect to local RPC server.
Make sure the Multiverse daemon is up.
See 'multi help daemon' for more info.
`

// HttpConn wraps an HTTP request.
type HttpConn struct {
	r io.Reader
	w io.Writer
}

// Read reads bytes from the reader.
func (c *HttpConn) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

// Write writes bytes to the writer.
func (c *HttpConn) Write(d []byte) (n int, err error) {
	return c.w.Write(d)
}

// Close does nothing.
func (c *HttpConn) Close() error {
	return nil
}

// NewClient returns a new RPC client.
func NewClient() (*rpc.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	config := remote.NewConfig(filepath.Join(home, remote.DotDir))
	if err := config.Read(); err != nil {
		return nil, err
	}

	return rpc.DialHTTP("tcp", config.HttpAddress)
}

// ListenAndServe starts the RPC listener.
func ListenAndServe(server *remote.Server) error {
	rpc.RegisterName("Author", &author.Service{server})
	rpc.RegisterName("File", &file.Service{server})
	rpc.RegisterName("Repo", &repo.Service{server})
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", server.Config.HttpAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	http.HandleFunc("/_jsonRPC_", ServeHTTP)
	return http.Serve(listener, nil)
}

// ServeHTTP serves json rpc connections over http.
func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)

	if req.Method == http.MethodPost {
		jsonrpc.ServeConn(&HttpConn{req.Body, w})
	}
}
