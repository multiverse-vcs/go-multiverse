package main

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
