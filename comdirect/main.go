package main

import (
	"github.com/jsattler/go-comdirect/comdirect/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
