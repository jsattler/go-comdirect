package comdirect

import (
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
	ctx, cancel := contextTimeout10Seconds()
	defer cancel()
	if _, err := client.Authenticate(ctx); err != nil {
		t.Errorf("authentication failed: %s", err)
		return
	}

	_, err := client.Reports(ctx)
	if err != nil {
		t.Errorf("failed to retrieve instruments: %s", err)
	}
}
