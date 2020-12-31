package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:     "multi",
	HelpName: "multi",
	Usage:    "Multiverse command line interface",
	Description: `Multiverse is a decentralized version control system for peer-to-peer software development.`,
	Version: "0.0.1",
	Authors: []*cli.Author{
		{Name: "Keenan Nemetz", Email: "keenan.nemetz@pm.me"},
	},
	Commands: []*cli.Command{
		commitCommand,
		daemonCommand,
		initCommand,
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
