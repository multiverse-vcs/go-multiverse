package main

import (
	"errors"
	"os"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var importCommand = &cli.Command{
	Action:    importAction,
	Name:      "import",
	Usage:     "Import a repo",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "Repo type",
			Value:   "git",
		},
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "Repo url",
		},
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{"d"},
			Usage:   "Repo directory",
		},
	},
}

func importAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	name := c.Args().Get(0)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := FindConfig(cwd); err == nil {
		return errors.New("repo already exists")
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.ImportArgs{
		Name: name,
		Type: c.String("type"),
		URL:  c.String("url"),
		Dir:  c.String("dir"),
	}

	var reply rpc.ImportReply
	if err := client.Call("Service.Import", &args, &reply); err != nil {
		return err
	}

	config := NewConfig(cwd)
	config.Name = name
	return config.Save()
}
