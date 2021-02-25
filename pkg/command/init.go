package command

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/pkg/command/context"
	"github.com/urfave/cli/v2"
)

// NewInitCommand returns a new cli command.
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a repository",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if _, err := context.Root(cwd); err == nil {
				return errors.New("repo already exists")
			}

			root := filepath.Join(cwd, context.DotDir)
			if err := os.Mkdir(root, 0755); err != nil {
				return err
			}

			config := context.NewConfig(root)
			return config.Write()
		},
	}
}
