package main

import (
	"os"

	"github.com/yondero/multiverse/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
