package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ipfs/go-ds-badger2"
	"github.com/multiverse-vcs/go-multiverse/key"
	"github.com/multiverse-vcs/go-multiverse/peer"
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

	root := filepath.Join(home, ".multiverse")
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	dpath := filepath.Join(root, "datastore")
	dstore, err := badger.NewDatastore(dpath, &badger.DefaultOptions)
	if err != nil {
		return err
	}

	kpath := filepath.Join(root, "keystore")
	kstore, err := key.NewKeystore(kpath)
	if err != nil {
		return err
	}

	key, err := kstore.DefaultKey()
	if err != nil {
		return err
	}

	client, err := peer.New(c.Context, dstore, key)
	if err != nil {
		return err
	}

	peerId, err := client.PeerID()
	if err != nil {
		return err
	}

	go web.ListenAndServe(client)
	go rpc.ListenAndServe(client)

	fmt.Printf(daemonBanner)
	fmt.Printf("Peer ID: %s\n", peerId.Pretty())
	fmt.Printf("Web URL: %s\n", web.BindAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	return nil
}
