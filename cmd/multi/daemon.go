package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ipfs/go-ds-badger2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/multiverse-vcs/go-multiverse/p2p"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/multiverse-vcs/go-multiverse/web"
	"github.com/urfave/cli/v2"
)

const daemonBanner = `
  __  __       _ _   _                         
 |  \/  |_   _| | |_(_)_   _____ _ __ ___  ___ 
 | |\/| | | | | | __| \ \ / / _ \ '__/ __|/ _ \
 | |  | | |_| | | |_| |\ V /  __/ |  \__ \  __/
 |_|  |_|\__,_|_|\__|_| \_/ \___|_|  |___/\___|
                                               
`

var daemonCommand = &cli.Command{
	Action: daemonAction,
	Name:   "daemon",
	Usage:  "Starts a client",
}

func daemonAction(c *cli.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(home, ".multiverse", "datastore")
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	dstore, err := badger.NewDatastore(path, &badger.DefaultOptions)
	if err != nil {
		return err
	}

	key, err := p2p.GenerateKey()
	if err != nil {
		return err
	}

	peerID, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return err
	}

	node, err := node.Init(c.Context, dstore, key)
	if err != nil {
		return err
	}

	go web.ListenAndServe(node)
	go rpc.ListenAndServe(node)

	fmt.Printf(daemonBanner)
	fmt.Printf("Peer ID:    %s\n", peerID.Pretty())
	fmt.Printf("Web Server: %s\n", web.BindAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	return nil
}
