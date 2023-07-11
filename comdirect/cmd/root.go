package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jsattler/go-comdirect/comdirect/keychain"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
)

var (
	folderFlag       string
	excludeFlag      string
	timeoutFlag      int
	formatFlag       string
	indexFlag        string
	countFlag        string
	sinceFlag        string
	downloadFlag     bool
	usernameFlag     string
	passwordFlag     string
	clientIDFlag     string
	clientSecretFlag string

	rootCmd = &cobra.Command{
		Use:   "comdirect",
		Short: "comdirect is a CLI tool to interact with the comdirect REST API",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	loginCmd.Flags().StringVarP(&passwordFlag, "password", "p", "", "comdirect password (PIN)")
	loginCmd.Flags().StringVarP(&usernameFlag, "username", "u", "", "comdirect username")
	loginCmd.Flags().StringVarP(&clientSecretFlag, "secret", "s", "", "comdirect client secret")
	loginCmd.Flags().StringVarP(&clientIDFlag, "id", "i", "", "comdirect client ID")

	documentCmd.Flags().StringVar(&folderFlag, "folder", "", "folder to save downloads")
	documentCmd.Flags().BoolVar(&downloadFlag, "download", false, "whether to download documents")

	transactionCmd.PersistentFlags().StringVar(&sinceFlag, "since", time.Now().Add(time.Hour*-1*24*30).Format("2006-01-02"), "Date of the earliest transaction date to retrieve in the form YYYY-MM-DD")

	rootCmd.PersistentFlags().StringVar(&indexFlag, "index", "0", "page index")
	rootCmd.PersistentFlags().StringVar(&countFlag, "count", "20", "page count")
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "markdown", "output format (markdown, csv or json)")
	rootCmd.PersistentFlags().IntVarP(&timeoutFlag, "timeout", "t", 30, "timeout in seconds to validate session TAN (default 30sec)")
	rootCmd.PersistentFlags().StringVar(&excludeFlag, "exclude", "", "exclude field from response")

	rootCmd.AddCommand(documentCmd)
	rootCmd.AddCommand(depotCmd)
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(versionCmd)

	accountCmd.AddCommand(balanceCmd)
	accountCmd.AddCommand(transactionCmd)

	depotCmd.AddCommand(positionCmd)
}

func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeoutFlag)*time.Second)
}

func formatAmountValue(av comdirect.AmountValue) string {
	value, err := strconv.ParseFloat(av.Value, 64)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%+5.2f", value)
}

func initClient() *comdirect.Client {
	authentication, err := keychain.RetrieveAuthentication()
	if err != nil || authentication.IsExpired() {
		// The session is expired, and we need to create a new session TAN
		authOptions, err := keychain.RetrieveAuthOptions()
		if err != nil {
			fmt.Println("You're not logged in. Please use 'comdirect login' to log in")
			os.Exit(1)
		}
		fmt.Println("Your session expired. Please open the comdirect photoTAN app to validate a new session.")
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
