package comdirect

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewAuthenticator(t *testing.T) {
	options := &AuthOptions{
		Username:     "",
		Password:     "",
		ClientId:     "",
		ClientSecret: "",
	}
	authenticator := NewAuthenticator(options)
	if authenticator.authOptions != options {
		t.Errorf("actual AuthOptions differ from expected: %v", authenticator.authOptions)
	}
}

func TestNewAuthenticator2(t *testing.T) {
	options := &AuthOptions{
		Username:     "",
		Password:     "",
		ClientId:     "",
		ClientSecret: "",
	}
	authenticator := NewAuthenticator(options)
	if authenticator.authOptions != options {
		t.Errorf("actual AuthOptions differ from expected: %v", authenticator.authOptions)
	}
}

func TestAuthenticator_Authenticate(t *testing.T) {
	authenticator := AuthenticatorFromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := authenticator.Authenticate(ctx)
	if err != nil {
		t.Errorf("authentication failed %s", err)
	}
}

func AuthenticatorFromEnv() *Authenticator {
	options := &AuthOptions{
		Username:     os.Getenv("COMDIRECT_USERNAME"),
		Password:     os.Getenv("COMDIRECT_PASSWORD"),
		ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
		ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
	}
	return NewAuthenticator(options)
}

func TestGenerateSessionId(t *testing.T) {
	sessionID := generateSessionID()
	if len(sessionID) != 32 {
		t.Errorf("length of session id not equal to 32: %d", len(sessionID))
	}
}

func TestGenerateRequestId(t *testing.T) {
	requestID := generateRequestID()
	if len(requestID) != 9 {
		t.Errorf("length of request ID is not equal to 9: %d", len(requestID))
	}
}
