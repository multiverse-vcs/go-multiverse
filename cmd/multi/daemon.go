package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ipfs/go-ds-badger2"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/peer"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/multiverse-vcs/go-multiverse/web"
	"github.com/nasdf/ulimit"
	"github.com/urfave/cli/v2"
)

const daemonBanner = `
  __  __       _ _   _                         
 |  \/  |_   _| | |_(_)_   _____ _ __ ___  ___ 
 | |\/| | | | | | __| \ \ / / _ \ '__/ __|/ _ \
 | |  | | |_| | | |_| |\ V /  __/ |  \__ \  __/
 |_|  |_|\__,_|_|\__|_| \_/ \___|_|  |___/\___|
                                               
`

const daemonUlimit = 8096

var daemonCommand = &cli.Command{
	Action: daemonAction,
	Name:   "daemon",
	Usage:  "Starts a client",
}

func daemonAction(c *cli.Context) error {
	if err := ulimit.SetRlimit(daemonUlimit); err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	root := filepath.Join(home, ".multiverse")
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}

	dpath := filepath.Join(root, "datastore")
	dopts := badger.DefaultOptions

	dstore, err := badger.NewDatastore(dpath, &dopts)
	if err != nil {
		return err
	}

	config, err := peer.LoadConfig(root)
	if err != nil {
		return err
	}

	client, err := peer.New(c.Context, dstore, config)
	if err != nil {
		return err
	}

	store := data.NewStore(dstore)
	go web.ListenAndServe(client, store)
	go rpc.ListenAndServe(client, store)

	fmt.Printf(daemonBanner)
	fmt.Printf("Peer ID: %s\n", client.PeerID().Pretty())
	fmt.Printf("Web URL: %s\n", web.BindAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	return nil
}
