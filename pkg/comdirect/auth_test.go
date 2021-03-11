package comdirect

import (
	"log"
	"os"
	"testing"
)

func TestNewAuthenticator(t *testing.T) {
	options := &AuthOptions{
		Username:     "",
		Password:     "",
		ClientId:     "",
		ClientSecret: "",
		AutoRefresh:  true,
	}
	sessionId, err := GenerateSessionId()

	if err != nil {
		t.Errorf("Failed to generate SessionId")
		return
	}

	authenticator := options.NewAuthenticator(sessionId)

	if authenticator.AuthOptions != options {
		t.Errorf("Actual AuthOptions differ from expected: %v", authenticator.AuthOptions)
	}
}

func TestNewAuthenticator2(t *testing.T) {
	options := &AuthOptions{
		Username:     "",
		Password:     "",
		ClientId:     "",
		ClientSecret: "",
		AutoRefresh:  true,
	}

	sessionId, err := GenerateSessionId()

	if err != nil {
		t.Errorf("Failed to generate SessionId")
		return
	}

	authenticator := NewAuthenticator(options, sessionId)

	if authenticator.AuthOptions != options {
		t.Errorf("Actual AuthOptions differ from expected: %v", authenticator.AuthOptions)
	}

}

func TestAuthenticator_Authenticate(t *testing.T) {
	authenticator := AuthenticatorFromEnv()

	_, err := authenticator.Authenticate()
	if err != nil {
		t.Errorf("Token should not be nil")
		return
	}
}

// set env variables locally
func AuthenticatorFromEnv() *Authenticator {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
		AutoRefresh:  true,
	}
	sessionId, err := GenerateSessionId()

	if err != nil {
		log.Fatal("Failed to generate session id")
	}

	return options.NewAuthenticator(sessionId)
}
