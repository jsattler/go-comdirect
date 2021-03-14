package comdirect

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

const (
	Host           = "api.comdirect.de"
	ApiPath        = "/api"
	OAuthTokenPath = "/oauth/token"

	HttpRequestInfoHeaderKey        = "X-Http-Request-Info"
	OnceAuthenticationInfoHeaderKey = "X-Once-Authentication-Info"
	OnceAuthenticationHeaderKey     = "X-Once-Authentication"
	AuthorizationHeaderKey          = "Authorization"
	ContentTypeHeaderKey            = "Content-Type"
	AcceptHeaderKey                 = "Accept"

	DefaultHttpTimeout          = time.Second * 30
	HttpsScheme                 = "https"
	BearerPrefix                = "Bearer "
	PasswordGrantType           = "password"
	ComdirectSecondaryGrantType = "cd_secondary"
)

type Client struct {
	authenticator *Authenticator
	http          *http.Client
}

func NewWithAuthenticator(authenticator *Authenticator) *Client {
	return nil
}

func NewWithAuthOptions(options *AuthOptions) *Client {
	return nil
}

func GenerateSessionId() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func generateRequestId() string {
	unix := time.Now().Unix()
	id := fmt.Sprintf("%09d", unix)
	return id[0:9]
}
