package rpc

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
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

// HttpConn wraps an HTTP request.
type HttpConn struct {
    r io.Reader
    w io.Writer
}

// Read reads bytes from the reader.
func (c *HttpConn) Read(p []byte) (n int, err error)  { 
	return c.r.Read(p) 
}

// Write writes bytes to the writer.
func (c *HttpConn) Write(d []byte) (n int, err error) { 
	return c.w.Write(d) 
}

// Close does nothing.
func (c *HttpConn) Close() error                      { 
	return nil 
}

// NewClient returns a new RPC client.
func NewClient() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", SocketAddr)
}

// ListenAndServe starts the RPC listener.
func ListenAndServe(server *remote.Server) error {
	rpc.RegisterName("Author", &author.Service{server})
	rpc.RegisterName("Repo", &repo.Service{server})
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", SocketAddr)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()
	
	http.HandleFunc("/", ServeHTTP)
	return http.Serve(listener, nil)
}

// ServeHTTP handles incoming RPC requests over HTTP.
func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(200)

	if req.Method != http.MethodOptions {
		jsonrpc.ServeConn(&HttpConn{req.Body, w})
	}
}