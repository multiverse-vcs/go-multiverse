package cmd

import (
	"os"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/multiverse-vcs/go-multiverse/config"
	"github.com/multiverse-vcs/go-multiverse/storage"
	"github.com/urfave/cli/v2"
)

// DefaultKeyType is the type of key to use.
const DefaultKeyType = crypto.Ed25519

// NewInitCommand returns a new init command.
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "create a new repo",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if _, err := Root(cwd); err == nil {
				return cli.Exit("repo already exists", 1)
			}

			if err := os.Mkdir(storage.DotDir, 0755); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			store, err := storage.NewOsStore(cwd)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			priv, _, err := crypto.GenerateKeyPair(DefaultKeyType, -1)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.WriteConfig(config.Default()); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if err := store.WriteKey(priv); err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return nil
		},
	}
}
