package comdirect

import (
	"fmt"
	"os"
	"testing"
)

func TestClient_Instrument(t *testing.T) {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	client := NewWithAuthOptions(options)
	ctx, cancel := contextTimeout10Seconds()
	defer cancel()
	if _, err := client.Authenticate(ctx); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	instruments, err := client.Instrument("865985")
	if err != nil {
		t.Errorf("failed to retrieve instruments: %s", err)
	}

	fmt.Printf("successfully retrieved instrument:\n%+v", instruments[0])
}
