package comdirect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jsattler/go-comdirect/internal/mediatype"
)

// Authenticator is responsible for authenticating against the comdirect REST API.
// It uses the given AuthOptions for authentication and returns an AccessToken in case
// the authentication flow was successful. Authenticator is using golang's default http.Client.
type Authenticator struct {
	authOptions *AuthOptions
	http        *HTTPClient
}

// authContext encapsulates the state that is passed through the comdirect authentication flow.
type authContext struct {
	accessToken  AccessToken
	requestInfo  requestInfo
	session      session
	onceAuthInfo onceAuthenticationInfo
}

// Authentication represents an authentication object for the comdirect REST API.
type Authentication struct {
	accessToken AccessToken
	sessionID   string
	time        time.Time
}

// AuthOptions encapsulates the information for authentication against the comdirect REST API.
type AuthOptions struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
}

// AccessToken represents an OAuth2 token that is returned from the comdirect REST API.
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CustomerID   string `json:"kdnr"`
	BPID         int    `json:"bpid"`
	ContactID    int    `json:"kontaktId"`
}

// requestInfo represents an x-http-request-info header exchanged with the comdirect REST API.
type requestInfo struct {
	ClientRequestID clientRequestID `json:"clientRequestId"`
}

// clientRequestID represents a client request ID that is exchanged with the comdirect REST API.
type clientRequestID struct {
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
}

// session represents a TAN session that is exchanged with the comdirect REST API.
type session struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}

// onceAuthenticationInfo represents the x-once-authentication-info header exchanged with
// the comdirect REST API.
type onceAuthenticationInfo struct {
	Id             string   `json:"id"`
	Typ            string   `json:"typ"`
	AvailableTypes []string `json:"availableTypes"`
	Link           link     `json:"link"`
}

type link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type authStatus struct {
	AuthenticationId string `json:"authenticationId"`
	Status           string `json:"status"`
}

// NewAuthenticator creates a new Authenticator by passing AuthOptions
// and an http.Client with a timeout of DefaultHttpTimeout.
func NewAuthenticator(options *AuthOptions) *Authenticator {
	return &Authenticator{
		authOptions: options,
		http:        &HTTPClient{http.Client{Timeout: DefaultHttpTimeout}},
	}
}

func NewAuthentication(accessToken AccessToken, sessionID string, time time.Time) *Authentication {
	return &Authentication{
		accessToken: accessToken,
		sessionID:   sessionID,
		time:        time,
	}
}

func (a *Authentication) AccessToken() AccessToken {
	return a.accessToken
}

func (a *Authentication) SessionID() string {
	return a.sessionID
}

func (a *Authentication) ExpiryTime() time.Time {
	return a.time
}

// Authenticate authenticates against the comdirect REST API.
func (a *Authenticator) Authenticate(ctx context.Context) (*Authentication, error) {

	authCtx, err := a.passwordGrant(ctx, a.authOptions)
	if err != nil {
		return nil, err
	}

	authCtx, err = a.fetchSessionStatus(ctx, authCtx)
	if err != nil {
		return nil, err
	}

	authCtx, err = a.validateSessionTan(ctx, authCtx)
	if err != nil {
		return nil, err
	}

	authCtx, err = a.checkAuthenticationStatus(ctx, authCtx)
	if err != nil {
		return nil, err
	}

	authCtx, err = a.activateSessionTan(ctx, authCtx)
	if err != nil {
		return nil, err
	}

	authCtx, err = a.secondaryFlow(ctx, authCtx)
	if err != nil {
		return nil, err
	}

	return &Authentication{
		accessToken: authCtx.accessToken,
		sessionID:   authCtx.requestInfo.ClientRequestID.SessionID,
		time:        time.Now(),
	}, err
}

func (a *Authenticator) Refresh(auth Authentication) (Authentication, error) {
	encoded := url.Values{
		"grant_type":    {RefreshTokenGrantType},
		"client_id":     {a.authOptions.ClientId},
		"client_secret": {a.authOptions.ClientSecret},
		"refresh_token": {auth.accessToken.RefreshToken},
	}.Encode()

	body := ioutil.NopCloser(strings.NewReader(encoded))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Host: Host, Scheme: HttpsScheme, Path: OAuthTokenPath},
		Header: http.Header{
			http.CanonicalHeaderKey(AcceptHeaderKey):      {mediatype.ApplicationJson},
			http.CanonicalHeaderKey(ContentTypeHeaderKey): {mediatype.XWWWFormUrlEncoded},
		},
		Body: body,
	}

	var accessToken AccessToken
	_, err := a.http.exchange(req, &accessToken)

	auth.accessToken = accessToken
	auth.time = time.Now()
	return auth, err
}

func (a *Authenticator) Revoke(auth Authentication) error {

	req := &http.Request{
		Method: http.MethodDelete,
		URL:    &url.URL{Host: Host, Scheme: HttpsScheme, Path: OAuthTokenPath},
		Header: http.Header{
			http.CanonicalHeaderKey(AcceptHeaderKey):        {mediatype.ApplicationJson},
			http.CanonicalHeaderKey(ContentTypeHeaderKey):   {mediatype.XWWWFormUrlEncoded},
			http.CanonicalHeaderKey(AuthorizationHeaderKey): {BearerPrefix + auth.accessToken.AccessToken},
		},
	}
	response, err := a.http.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 204 {
		return errors.New("could not revoke access token")
	}
	return nil
}

func (a *Authentication) IsExpired() bool {
	expiresIn := time.Duration(a.accessToken.ExpiresIn)
	return a.time.Add(expiresIn * time.Second).Before(time.Now())
}

// Step 2.1
func (a *Authenticator) passwordGrant(ctx context.Context, options *AuthOptions) (authContext, error) {

	urlEncoded := url.Values{
		"username":      {options.Username},
		"password":      {options.Password},
		"grant_type":    {PasswordGrantType},
		"client_id":     {options.ClientId},
		"client_secret": {options.ClientSecret},
	}.Encode()

	body := ioutil.NopCloser(strings.NewReader(urlEncoded))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Host: Host, Scheme: HttpsScheme, Path: OAuthTokenPath},
		Header: http.Header{
			http.CanonicalHeaderKey(AcceptHeaderKey):      {mediatype.ApplicationJson},
			http.CanonicalHeaderKey(ContentTypeHeaderKey): {mediatype.XWWWFormUrlEncoded},
		},
		Body: body,
	}
	req = req.WithContext(ctx)

	authCtx := authContext{
		accessToken: AccessToken{},
	}
	_, err := a.http.exchange(req, &authCtx.accessToken)

	return authCtx, err
}

// Step 2.2
func (a *Authenticator) fetchSessionStatus(ctx context.Context, authCtx authContext) (authContext, error) {
	authCtx.requestInfo = requestInfo{
		ClientRequestID: clientRequestID{
			SessionID: generateSessionID(),
			RequestID: generateRequestID(),
		},
	}
	requestInfoJSON, err := json.Marshal(authCtx.requestInfo)
	if err != nil {
		return authCtx, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    comdirectURL("/api/session/clients/user/v1/sessions"),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + authCtx.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJSON)},
		},
	}
	req = req.WithContext(ctx)

	var sessions []session
	if _, err = a.http.exchange(req, &sessions); err != nil {
		return authCtx, err
	}

	if len(sessions) == 0 {
		return authCtx, errors.New("length of returned session array is zero; expected at least one")
	}

	authCtx.session = sessions[0]

	return authCtx, err
}

// Step: 2.3
func (a *Authenticator) validateSessionTan(ctx context.Context, authCtx authContext) (authContext, error) {
	authCtx.session.SessionTanActive = true
	authCtx.session.Activated2FA = true
	jsonSession, err := json.Marshal(authCtx.session)
	if err != nil {
		return authCtx, err
	}

	authCtx.requestInfo.ClientRequestID.RequestID = generateRequestID()
	JSONRequestInfo, err := json.Marshal(authCtx.requestInfo)
	if err != nil {
		return authCtx, err
	}

	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s/validate", authCtx.session.Identifier)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    comdirectURL(path),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + authCtx.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(JSONRequestInfo)},
		},
		Body: body,
	}
	req = req.WithContext(ctx)

	res, err := a.http.exchange(req, &authCtx.session)
	if err != nil || res == nil {
		return authCtx, err
	}

	jsonOnceAuthInfo := res.Header.Get(OnceAuthenticationInfoHeaderKey)
	if jsonOnceAuthInfo == "" {
		return authCtx, errors.New("x-once-authentication-info header missing in response")
	}

	err = json.Unmarshal([]byte(jsonOnceAuthInfo), &authCtx.onceAuthInfo)

	if err != nil {
		return authCtx, err
	}

	return authCtx, res.Body.Close()
}

// Step 2.4
func (a *Authenticator) activateSessionTan(ctx context.Context, authCtx authContext) (authContext, error) {
	authCtx.requestInfo.ClientRequestID.RequestID = generateRequestID()
	requestInfoJSON, err := json.Marshal(authCtx.requestInfo)
	if err != nil {
		return authCtx, err
	}

	onceAuthInfoHeader := fmt.Sprintf(`{"id":"%s"}`, authCtx.onceAuthInfo.Id)
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s", authCtx.session.Identifier)
	JSONSession, err := json.Marshal(authCtx.session)

	if err != nil {
		return authContext{}, err
	}

	req := &http.Request{
		Method: http.MethodPatch,
		URL:    comdirectURL(path),
		Header: http.Header{
			AuthorizationHeaderKey:          {BearerPrefix + authCtx.accessToken.AccessToken},
			AcceptHeaderKey:                 {mediatype.ApplicationJson},
			ContentTypeHeaderKey:            {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey:        {string(requestInfoJSON)},
			OnceAuthenticationInfoHeaderKey: {onceAuthInfoHeader},
		},
		Body: ioutil.NopCloser(strings.NewReader(string(JSONSession))),
	}
	req = req.WithContext(ctx)

	_, err = a.http.exchange(req, &authCtx.session)

	if err != nil {
		return authCtx, err
	}

	return authCtx, nil
}

// Step 2.5
func (a *Authenticator) secondaryFlow(ctx context.Context, authCtx authContext) (authContext, error) {

	body := url.Values{
		"token":         {authCtx.accessToken.AccessToken},
		"grant_type":    {SecondaryGrantType},
		"client_id":     {a.authOptions.ClientId},
		"client_secret": {a.authOptions.ClientSecret},
	}.Encode()

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Host: Host, Scheme: HttpsScheme, Path: OAuthTokenPath},
		Header: http.Header{
			http.CanonicalHeaderKey(AcceptHeaderKey):      {mediatype.ApplicationJson},
			http.CanonicalHeaderKey(ContentTypeHeaderKey): {mediatype.XWWWFormUrlEncoded},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
	req = req.WithContext(ctx)

	res, err := a.http.Do(req)

	if err != nil {
		return authCtx, err
	}

	if err = json.NewDecoder(res.Body).Decode(&authCtx.accessToken); err != nil {
		return authCtx, err
	}

	return authCtx, res.Body.Close()
}

func (a *Authenticator) checkAuthenticationStatus(ctx context.Context, authCtx authContext) (authContext, error) {
	authCtx.requestInfo.ClientRequestID.RequestID = generateRequestID()
	requestInfoJson, err := json.Marshal(authCtx.requestInfo)
	if err != nil {
		return authCtx, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    comdirectURL(authCtx.onceAuthInfo.Link.Href),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + authCtx.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJson)},
		},
	}
	req = req.WithContext(ctx)

	for {
		select {
		// Poll authentication status every 3 seconds
		case <-time.After(3 * time.Second):
			response, err := a.http.Do(req)
			if err != nil {
				return authCtx, err
			}
			var authStatus authStatus
			if err = json.NewDecoder(response.Body).Decode(&authStatus); err != nil {
				return authCtx, err
			}

			if authStatus.Status == "AUTHENTICATED" {
				return authCtx, nil
			}
		// When timeout is reached return with error
		case <-ctx.Done():
			return authCtx, ctx.Err()
		}
	}

}
