package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	depotCmd = &cobra.Command{
		Use:   "depot",
		Short: "list basic depot information",
		Run:   depot,
	}
)

func depot(cmd *cobra.Command, args []string) {
	fmt.Println("Running depots..")
}
