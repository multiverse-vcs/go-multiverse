package rpc

import (
	"context"
	"net"
	"net/rpc"

	datastore "github.com/ipfs/go-datastore"
	"github.com/multiverse-vcs/go-multiverse/peer"
)

// makeNode returns a new peer node.
func makeNode(ctx context.Context) (*peer.Node, error) {
	dstore := datastore.NewMapDatastore()

	config, err := peer.NewConfig("")
	if err != nil {
		return nil, err
	}

	return peer.New(ctx, dstore, config)
}

// makeClient starts an rpc server and returns a connected client.
func makeClient(node *peer.Node) (*rpc.Client, error) {
	service := Service{
		node: node,
	}

	server := rpc.NewServer()
	if err := server.Register(&service); err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, err
	}
	go server.Accept(listener)

	return rpc.Dial("tcp", listener.Addr().String())
}
