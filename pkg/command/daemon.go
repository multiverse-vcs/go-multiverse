package command

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/multiverse-vcs/go-multiverse/pkg/remote"
	"github.com/multiverse-vcs/go-multiverse/pkg/rpc"
	"github.com/nasdf/ulimit"
	"github.com/urfave/cli/v2"
)

const Banner = `
  __  __       _ _   _                         
 |  \/  |_   _| | |_(_)_   _____ _ __ ___  ___ 
 | |\/| | | | | | __| \ \ / / _ \ '__/ __|/ _ \
 | |  | | |_| | | |_| |\ V /  __/ |  \__ \  __/
 |_|  |_|\__,_|_|\__|_| \_/ \___|_|  |___/\___|
                                               
`

// NewDaemonCommand returns a new cli command.
func NewDaemonCommand() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Run a persistent peer",
		Action: func(c *cli.Context) error {
			if err := ulimit.SetRlimit(8096); err != nil {
				return err
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			server, err := remote.NewServer(c.Context, home)
			if err != nil {
				return err
			}

			go rpc.ListenAndServe(server)

			fmt.Printf(Banner)
			fmt.Printf("Peer ID: %s\n", server.Peer.Host.ID().Pretty())

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)

			<-quit
			return nil
		},
	}
}
