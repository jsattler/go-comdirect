package cmd

import (
	"bufio"
	"context"
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

	if username == "" {
		fmt.Print("Enter Username: ")
		scanner.Scan()
		username = scanner.Text()
	}

	if password == "" {
		fmt.Print("Enter Password: ")
		bytePassword, _ := term.ReadPassword(syscall.Stdin)
		password = string(bytePassword)
		fmt.Println()
	}

	if clientID == "" {
		fmt.Print("Enter Client ID: ")
		scanner.Scan()
		clientID = scanner.Text()
	}

	if clientSecret == "" {
		fmt.Print("Enter Client Secret: ")
		byteClientSecret, _ := term.ReadPassword(syscall.Stdin)
		clientSecret = string(byteClientSecret)
		fmt.Println()
	}

	options := &comdirect.AuthOptions{
		Username:     username,
		Password:     password,
		ClientId:     clientID,
		ClientSecret: clientSecret,
	}

	if err := keychain.StoreAuthOptions(options); err != nil {
		log.Fatal(err)
	}

	authenticator := comdirect.NewAuthenticator(options)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println("Open your comdirect photoTAN app to complete the login")

	authentication, err := authenticator.Authenticate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := keychain.StoreAuthentication(authentication); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully logged in - the session will expire at %s\n",
		authentication.ExpiryTime().
			Add(time.Duration(authentication.AccessToken().ExpiresIn)*time.Second).
			Format(time.RFC3339))
}
