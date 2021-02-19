package command

import (
	"errors"
	"os"
	"path/filepath"

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

			if _, err := Root(cwd); err == nil {
				return errors.New("repo already exists")
			}

			root := filepath.Join(cwd, DotDir)
			if err := os.Mkdir(root, 0755); err != nil {
				return err
			}

			config := NewConfig(root)
			return config.Write()
		},
	}
}
