package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/multiverse-vcs/go-multiverse/p2p"
	"github.com/urfave/cli/v2"
)

// NewServeCommand returns a new serve command.
func NewServeCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "exchange data with peers",
		Action: func(c *cli.Context) error {
			store, err := Store()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.Online(c.Context); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Printf("bootstrapping network...\n")
			p2p.Bootstrap(c.Context, store.Host)

			fmt.Printf("Connected to network with peer id %s:\n", store.Host.ID().Pretty())
			fmt.Printf("  (listening on multiaddresses)\n")
			for _, a := range store.Host.Addrs() {
				fmt.Printf("\t%s%s/p2p/%s%s\n", ColorGreen, a, store.Host.ID().Pretty(), ColorReset)
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)

			<-interrupt
			return nil
		},
	}
}
