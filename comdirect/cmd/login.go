package cmd

import (
	"bufio"
	"fmt"
	"github.com/jsattler/go-comdirect/comdirect/keychain"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
	"time"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "log in to comdirect",
		Run:   login,
	}
)

func login(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	if usernameFlag == "" {
		fmt.Print("Username: ")
		scanner.Scan()
		usernameFlag = scanner.Text()
	}

	if passwordFlag == "" {
		fmt.Print("Password: ")
		bytePassword, _ := term.ReadPassword(syscall.Stdin)
		passwordFlag = string(bytePassword)
		fmt.Println()
	}

	if clientIDFlag == "" {
		fmt.Print("Client ID: ")
		scanner.Scan()
		clientIDFlag = scanner.Text()
	}

	if clientSecretFlag == "" {
		fmt.Print("Client Secret: ")
		byteClientSecret, _ := term.ReadPassword(syscall.Stdin)
		clientSecretFlag = string(byteClientSecret)
		fmt.Println()
	}

	options := &comdirect.AuthOptions{
		Username:     usernameFlag,
		Password:     passwordFlag,
		ClientId:     clientIDFlag,
		ClientSecret: clientSecretFlag,
	}

	if err := keychain.StoreAuthOptions(options); err != nil {
		log.Fatal(err)
	}

	authenticator := comdirect.NewAuthenticator(options)
	ctx, cancel := contextWithTimeout()
	defer cancel()

	fmt.Println("Open your comdirect photoTAN app to complete the login")

	authentication, err := authenticator.Authenticate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := keychain.StoreAuthentication(authentication); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully logged in - the session will expire in 10 minutes (%s)\n",
		authentication.ExpiryTime().
			Add(time.Duration(authentication.AccessToken().ExpiresIn)*time.Second).
			Format(time.RFC3339))
}
