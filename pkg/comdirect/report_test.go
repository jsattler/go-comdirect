package comdirect

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Reports(t *testing.T) {
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

	reports, err := client.Reports()
	if err != nil {
		t.Errorf("failed to retrieve instruments: %s", err)
	}

	fmt.Printf("successfully retrieved instrument:\n%+v", reports[0])
}
