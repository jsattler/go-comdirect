package comdirect

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestClient_Balances(t *testing.T) {
	client := clientFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	auth, err := client.Authenticate(ctx)
	if err != nil {
		t.Errorf("failed to authenticate: %s", err)
	}
	fmt.Printf("%+v\n", auth)
	balances, err := client.Balances(ctx)
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
	ctx, cancel := contextTimeout10Seconds()
	defer cancel()
	_, err := client.Balance(ctx, os.Getenv("COMDIRECT_ACCOUNT_ID"))

	if err != nil {
		t.Errorf("failed to exchange account balance %s", err)
		return
	}
}

func TestClient_Transactions(t *testing.T) {
	client := clientFromEnv()
	ctx, cancel := contextTimeout10Seconds()
	defer cancel()
	transactions, err := client.Transactions(ctx, os.Getenv("COMDIRECT_ACCOUNT_ID"))

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

func contextTimeout10Seconds() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*10)
}
