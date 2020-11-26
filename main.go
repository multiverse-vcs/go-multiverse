package main

import (
	"os"

	"github.com/ipfs/go-log"
	"github.com/multiverse-vcs/go-multiverse/cmd"
)

func main() {
	log.SetLogLevel("autonat", "debug")
	log.SetLogLevel("nat", "debug")

	cmd.NewApp().Run(os.Args)
}
