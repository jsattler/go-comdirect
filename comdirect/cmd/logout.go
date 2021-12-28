package cmd

import (
	"fmt"
	"github.com/jsattler/go-comdirect/comdirect/keychain"
	"github.com/spf13/cobra"
)

var (
	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "log out of comdirect",
		Run:   logout,
	}
)

func logout(cmd *cobra.Command, args []string) {
	keychain.DeleteAuthentication()
	keychain.DeleteAuthOptions()
	fmt.Println("Successfully logged out")
}
