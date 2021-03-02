package http

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/author"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc/repo"
)

// ApiFile is the name of the api file.
const ApiFile = "api"

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

// Route is an http handler that returns an error
type Route func(http.ResponseWriter, *http.Request) error

// ServeHTTP serves the request and handles any errors.
func (r Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if err := r(w, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// Server is an HTTP gateway.
type Server struct {
	*remote.Server
}

// ApiAddress returns the address for the api.
func ApiAddress(home string) (string, error) {
	path := filepath.Join(home, remote.DotDir, ApiFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ListenAndServe starts the HTTP listener.
func ListenAndServe(s *remote.Server) error {
	path := filepath.Join(s.Root, ApiFile)
	data := []byte(s.Config.HttpAddress)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	server := &Server{s}
	router := httprouter.New()
	router.Handler(http.MethodGet, "/:peer/:repo/:branch/*file", Route(server.File))

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	rpc.RegisterName("Author", &author.Service{s})
	rpc.RegisterName("Repo", &repo.Service{s})

	http.HandleFunc("/_jsonRPC_", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)

		if req.Method == http.MethodPost {
			jsonrpc.ServeConn(&HttpConn{req.Body, w})
		}
	})

	http.Handle("/", router)
	return http.ListenAndServe(s.Config.HttpAddress, nil)
}
