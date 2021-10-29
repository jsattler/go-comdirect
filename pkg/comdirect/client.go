package comdirect

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	HttpRequestInfoHeaderKey        = "X-Http-Request-Info"
	OnceAuthenticationInfoHeaderKey = "X-Once-Authentication-Info"
	OnceAuthenticationHeaderKey     = "X-Once-Authentication"
	AuthorizationHeaderKey          = "Authorization"
	ContentTypeHeaderKey            = "Content-Type"
	AcceptHeaderKey                 = "Accept"

	Host           = "api.comdirect.de"
	ApiPath        = "/api"
	OAuthTokenPath = "/oauth/token"

	PasswordGrantType     = "password"
	SecondaryGrantType    = "cd_secondary"
	RefreshTokenGrantType = "refresh_token"

	DefaultHttpTimeout = time.Second * 30
	HttpsScheme        = "https"
	BearerPrefix       = "Bearer "
)

type Client struct {
	authenticator  *Authenticator
	http           *HTTPClient
	authentication *Authentication
}

type AmountValue struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

func NewWithAuthenticator(authenticator *Authenticator) *Client {
	return &Client{
		authenticator: authenticator,
		http:          &HTTPClient{http.Client{Timeout: DefaultHttpTimeout}},
	}
}

func NewWithAuthOptions(options *AuthOptions) *Client {
	return NewWithAuthenticator(options.NewAuthenticator())
}

// Authenticate uses the underlying Authenticator to authenticate against the comdirect REST API.
func (c *Client) Authenticate(ctx context.Context) (*Authentication, error) {
	if c.authenticator == nil {
		return nil, errors.New("authenticator cannot be nil")
	}
	authentication, err := c.authenticator.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	c.authentication = authentication
	return c.authentication, nil
}

// Revoke uses the underlying Authenticator to revoke an access token.
func (c *Client) Revoke() error {
	if c.authenticator == nil {
		return errors.New("authenticator cannot be nil")
	}
	err := c.authenticator.Revoke(*c.authentication)
	if err != nil {
		return err
	}
	c.authentication = nil
	return nil
}

// Refresh uses the underlying Authenticator to refresh an access token.
func (c *Client) Refresh() (*Authentication, error) {
	if c.authenticator == nil {
		return nil, errors.New("authenticator cannot be nil")
	}
	authentication, err := c.authenticator.Refresh(*c.authentication)
	if err != nil {
		return nil, err
	}
	c.authentication = &authentication
	return c.authentication, nil
}
