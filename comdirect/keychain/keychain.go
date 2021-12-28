package keychain

import (
	"encoding/json"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/zalando/go-keyring"
	"log"
	"time"
)

const servicePrefix = "github.com.jsattler.comdirect-golang."
const user = "comdirect"

func StoreAuthOptions(options *comdirect.AuthOptions) error {
	if err := keyring.Set(servicePrefix+"username", user, options.Username); err != nil {
		return err
	}
	if err := keyring.Set(servicePrefix+"password", user, options.Password); err != nil {
		return err
	}
	if err := keyring.Set(servicePrefix+"clientID", user, options.ClientId); err != nil {
		return err
	}
	return keyring.Set(servicePrefix+"clientSecret", user, options.ClientSecret)
}

func StoreAuthentication(authentication *comdirect.Authentication) error {

	accessToken := authentication.AccessToken()
	loginTime := authentication.ExpiryTime()
	sessionID := authentication.SessionID()

	byteJson, err := json.Marshal(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	if err := keyring.Set(servicePrefix+"accessToken", user, string(byteJson)); err != nil {
		return err
	}

	if err := keyring.Set(servicePrefix+"loginTime", user, loginTime.Format(time.RFC3339)); err != nil {
		return err
	}

	return keyring.Set(servicePrefix+"sessionID", user, sessionID)
}

func RetrieveAuthOptions() (*comdirect.AuthOptions, error) {
	username, err := keyring.Get(servicePrefix+"username", user)
	if err != nil {
		return nil, err
	}
	password, err := keyring.Get(servicePrefix+"password", user)
	if err != nil {
		return nil, err
	}
	clientID, err := keyring.Get(servicePrefix+"clientID", user)
	if err != nil {
		return nil, err
	}
	clientSecret, err := keyring.Get(servicePrefix+"clientSecret", user)
	if err != nil {
		return nil, err
	}

	return &comdirect.AuthOptions{
		Username:     username,
		Password:     password,
		ClientId:     clientID,
		ClientSecret: clientSecret,
	}, nil

}

func RetrieveAuthentication() (*comdirect.Authentication, error) {
	loginTimeString, err := keyring.Get(servicePrefix+"loginTime", user)

	if err != nil {
		return nil, err
	}

	sessionID, err := keyring.Get(servicePrefix+"sessionID", user)
	if err != nil {
		return nil, err
	}

	accessTokenString, err := keyring.Get(servicePrefix+"accessToken", user)
	if err != nil {
		return nil, err
	}

	loginTime, err := time.Parse(time.RFC3339, loginTimeString)
	if err != nil {
		return nil, err
	}

	accessToken := comdirect.AccessToken{}

	err = json.Unmarshal([]byte(accessTokenString), &accessToken)

	if err != nil {
		return nil, err
	}

	return comdirect.NewAuthentication(accessToken, sessionID, loginTime), nil
}

func DeleteAuthentication() {
	_ = keyring.Delete(servicePrefix+"loginTime", user)
	_ = keyring.Delete(servicePrefix+"sessionID", user)
	_ = keyring.Delete(servicePrefix+"accessToken", user)
}

func DeleteAuthOptions() {
	_ = keyring.Delete(servicePrefix+"username", user)
	_ = keyring.Delete(servicePrefix+"password", user)
	_ = keyring.Delete(servicePrefix+"clientID", user)
	_ = keyring.Delete(servicePrefix+"clientSecret", user)
}
