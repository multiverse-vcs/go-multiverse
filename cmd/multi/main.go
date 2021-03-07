package main

import (
	"fmt"
	"os"

	"github.com/multiverse-vcs/go-multiverse/pkg/command"
)

// version is set by goreleaser
var version = "dev"

func main() {
	app := command.NewApp()
	app.Version = version

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
