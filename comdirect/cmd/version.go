package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version of the comdirect CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("comdirect CLI v0.4.0")
		},
	}
)
