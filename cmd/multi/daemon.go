package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ipfs/go-ds-badger2"
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

	config, err := peer.OpenConfig(root)
	if err != nil {
		return err
	}

	node, err := peer.New(c.Context, dstore, config)
	if err != nil {
		return err
	}

	// ensure any changes made offline will be published
	if err := node.Authors().Publish(c.Context); err != nil {
		return err
	}

	go web.ListenAndServe(node)
	go rpc.ListenAndServe(node)

	fmt.Printf(daemonBanner)
	fmt.Printf("Peer ID: %s\n", node.ID().Pretty())
	fmt.Printf("Web URL: %s\n", web.BindAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	return nil
}
