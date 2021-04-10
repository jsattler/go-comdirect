package comdirect

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_AccountBalances(t *testing.T) {
	client := getClientFromEnv()
	auth, err := client.authenticator.Authenticate()
	if err != nil {
		t.Errorf("failed to authenticate: %s", err)
	}

	fmt.Printf("%+v\n", auth)
	client.authentication = *auth
	balances, err := client.Balances()
	if err != nil {
		t.Errorf("failed to exchange account balances %s", err)
	}

	fmt.Printf("%+v\n", balances)
}

func TestClient_Balance(t *testing.T) {
	client := getClientFromEnv()
	accountId := os.Getenv("COMDIRECT_ACCOUNT_ID")
	_, err := client.Balance(accountId)

	if err != nil {
		t.Errorf("failed to exchange account balance %s", err)
		return
	}
}

func TestClient_Transactions(t *testing.T) {
	client := getClientFromEnv()
	accountId := os.Getenv("COMDIRECT_ACCOUNT_ID")
	transactions, err := client.Transactions(accountId)

	if err != nil {
		t.Errorf("failed to exchange account balance %s", err)
		return
	}

	for _, t := range transactions {
		fmt.Printf("%v\n", t)
	}
}

func getClientFromEnv() *Client {

	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}

	return NewWithAuthOptions(options)
}
