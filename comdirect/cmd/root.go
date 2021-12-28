package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jsattler/go-comdirect/comdirect/keychain"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	pageIndex string
	pageCount string

	username     string
	password     string
	clientID     string
	clientSecret string

	rootCmd = &cobra.Command{
		Use:   "comdirect",
		Short: "comdirect is a CLI tool to interact with the comdirect REST API",
	}

	boldGreen = color.New(color.FgGreen, color.Bold).SprintfFunc()
	boldRed   = color.New(color.FgRed, color.Bold).SprintfFunc()
	bold      = color.New(color.Bold).SprintFunc()
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	loginCmd.Flags().StringVarP(&password, "password", "p", "", "comdirect password (PIN)")
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "comdirect username")
	loginCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "comdirect client secret")
	loginCmd.Flags().StringVarP(&clientID, "id", "i", "", "comdirect client ID")

	documentCmd.Flags().StringVar(&folder, "folder", "", "folder to save downloads")

	rootCmd.PersistentFlags().StringVar(&pageIndex, "index", "0", "page index")
	rootCmd.PersistentFlags().StringVar(&pageCount, "count", "20", "page count")
	rootCmd.AddCommand(documentCmd)
	rootCmd.AddCommand(depotCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(versionCmd)
}

func InitClient() *comdirect.Client {
	authentication, err := keychain.RetrieveAuthentication()
	if err != nil || authentication.IsExpired() {
		// The session is expired, and we need to create a new session TAN
		authOptions, err := keychain.RetrieveAuthOptions()
		if err != nil {
			fmt.Println("You're not logged in - please use 'comdirect login' to log in")
			os.Exit(1)
		}
		fmt.Println("Session expired - please open the comdirect photoTAN app to validate a new session")
		client := comdirect.NewWithAuthOptions(authOptions)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		keychain.DeleteAuthentication()
		authentication, err = client.Authenticate(ctx)
		if err != nil {
			log.Fatal(err)
		}
		err = keychain.StoreAuthentication(authentication)
		if err != nil {
			log.Fatal(err)
		}
		return client
	}
	return comdirect.NewWithAuthentication(authentication)
}
