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

	authenticator := options.NewAuthenticator()

	if authenticator.authOptions != options {
		t.Errorf("Actual AuthOptions differ from expected: %v", authenticator.authOptions)
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

	authenticator := NewAuthenticator(options)

	if authenticator.authOptions != options {
		t.Errorf("Actual AuthOptions differ from expected: %v", authenticator.authOptions)
	}

}

func TestAuthenticator_Authenticate(t *testing.T) {
	authenticator := AuthenticatorFromEnv()

	_, _, err := authenticator.Authenticate()
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

	return options.NewAuthenticator()
}

func TestGenerateSessionId(t *testing.T) {
	sessionId, err := generateSessionId()

	if err != nil {
		t.Errorf("Error generating session id")
	}

	if len(sessionId) != 32 {
		t.Errorf("Length of session id not equal to 32: %d", len(sessionId))
	}

	log.Println(sessionId)
}

func TestGenerateRequestId(t *testing.T) {
	generateRequestId()
}
