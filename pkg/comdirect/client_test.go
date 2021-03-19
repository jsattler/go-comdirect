package comdirect

import (
	"os"
	"testing"
)

func TestClient_AccountBalances(t *testing.T) {
	client := getClientFromEnv()

	_, err := client.Balances()

	if err != nil {
		t.Errorf("failed to request account balances %s", err)
	}
}

func TestClient_Balance(t *testing.T) {
	client := getClientFromEnv()
	accountId := os.Getenv("COMDIRECT_ACCOUNT_ID")
	_, err := client.Balance(accountId)

	if err != nil {
		t.Errorf("failed to request account balance %s", err)
		return
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
