package comdirect

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Balances(t *testing.T) {
	client := clientFromEnv()
	auth, err := client.Authenticate()
	if err != nil {
		t.Errorf("failed to authenticate: %s", err)
	}
	fmt.Printf("%+v\n", auth)
	balances, err := client.Balances()
	if err != nil {
		t.Errorf("failed to exchange account balances %s", err)
	}

	fmt.Printf("%+v\n", balances)

	auth, err = client.Refresh()
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", auth)

	err = client.Revoke()
	if err != nil {
		return
	}
	fmt.Println("successfully revoked access token")
}

func TestClient_Balance(t *testing.T) {
	client := clientFromEnv()
	_, err := client.Balance(os.Getenv("COMDIRECT_ACCOUNT_ID"))

	if err != nil {
		t.Errorf("failed to exchange account balance %s", err)
		return
	}
}

func TestClient_Transactions(t *testing.T) {
	client := clientFromEnv()
	transactions, err := client.Transactions(os.Getenv("COMDIRECT_ACCOUNT_ID"))

	if err != nil {
		t.Errorf("failed to exchange account balance %s", err)
		return
	}

	for _, t := range transactions {
		fmt.Printf("%v\n", t)
	}
}

func clientFromEnv() *Client {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	return NewWithAuthOptions(options)
}
