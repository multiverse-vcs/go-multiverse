package main

import (
	"os"

	"github.com/multiverse-vcs/go-multiverse/cmd"
)

func main() {
	cmd.NewApp().Run(os.Args)
}
