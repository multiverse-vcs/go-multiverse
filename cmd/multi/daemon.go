package main

import (
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ipfs/go-ds-badger2"
	"github.com/multiverse-vcs/go-multiverse/http"
	"github.com/multiverse-vcs/go-multiverse/node"
	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

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

	node, err := node.Init(c.Context, dstore)
	if err != nil {
		return err
	}

	go http.ListenAndServe(node)
	go rpc.ListenAndServe(node)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	return nil
}
