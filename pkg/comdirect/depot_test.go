package comdirect

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Depots(t *testing.T) {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	client := NewWithAuthOptions(options)
	if _, err := client.Authenticate(); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	depots, err := client.Depots()
	if err != nil {
		t.Errorf("failed to retrieve depots: %s", err)
	}

	fmt.Printf("successfully retrieved depots:\n%+v", depots)
}

func TestClient_DepotPositions(t *testing.T) {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	client := NewWithAuthOptions(options)
	if _, err := client.Authenticate(); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	depotPositions, err := client.DepotPositions(os.Getenv("COMDIRECT_DEPOT_ID"))
	if err != nil {
		t.Errorf("failed to retrieve depot positions: %s", err)
	}

	fmt.Printf("successfully retrieved depot positions:\n%+v", depotPositions)
}

func TestClient_DepotPosition(t *testing.T) {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	client := NewWithAuthOptions(options)
	if _, err := client.Authenticate(); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	depotPositions, err := client.DepotPosition(os.Getenv("COMDIRECT_DEPOT_ID"), os.Getenv("COMDIRECT_POSITION_ID"))
	if err != nil {
		t.Errorf("failed to retrieve depot position: %s", err)
	}

	fmt.Printf("successfully retrieved depot position:\n%+v", depotPositions)
}

func TestClient_DepotTransactions(t *testing.T) {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	client := NewWithAuthOptions(options)
	if _, err := client.Authenticate(); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	depotTransactions, err := client.DepotTransactions(os.Getenv("COMDIRECT_DEPOT_ID"))
	if err != nil {
		t.Errorf("failed to retrieve depot transactions: %s", err)
	}

	fmt.Printf("successfully retrieved depot transactions:\n%+v", depotTransactions)
}
