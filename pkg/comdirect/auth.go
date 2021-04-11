package comdirect

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsattler/comdirect-golang/internal/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: Currently it is not possible to chose between TAN types (Photo TAN, Push TAN etc.). Only testing with P_TAN_PUSH at the moment.
// TODO: Let the user pass in a preconfigured http.Client to follow dependency injection principle.
// TODO: Think about accessibility of functions, structs, etc.
// TODO: Provide an interface

// Authenticator is responsible for authenticating against the comdirect REST API.
// It uses the given AuthOptions for authentication and returns an AccessToken in case
// the authentication flow was successful. Authenticator is using golang's default http.Client.
type Authenticator struct {
	authOptions *AuthOptions
	http        *HttpClient
}

// authState encapsulates the state that is passed through the comdirect authentication flow.
type authState struct {
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

// NewAuthenticator creates a new Authenticator from AuthOptions
// with an http.Client configured with a timeout of DefaultHttpTimeout.
func (a AuthOptions) NewAuthenticator() *Authenticator {
	return &Authenticator{
		authOptions: &a,
		http:        &HttpClient{http.Client{Timeout: DefaultHttpTimeout}},
	}
}

// NewAuthenticator creates a new Authenticator by passing AuthOptions
// and an http.Client with a timeout of DefaultHttpTimeout.
func NewAuthenticator(options *AuthOptions) *Authenticator {
	return &Authenticator{
		authOptions: options,
		http:        &HttpClient{http.Client{Timeout: DefaultHttpTimeout}},
	}
}

// Authenticate authenticates against the comdirect REST API.
func (a *Authenticator) Authenticate() (*Authentication, error) {

	state, err := a.passwordGrant(a.authOptions)
	if err != nil {
		return nil, err
	}

	state, err = a.fetchSessionStatus(state)
	if err != nil {
		return nil, err
	}

	state, err = a.validateSessionTan(state)
	if err != nil {
		return nil, err
	}

	// https://community.comdirect.de/t5/Website-Apps/REST-API-Schritt-2-4-Aktivierung-einer-Session-TAN/td-p/153737/page/2
	time.Sleep(10 * time.Second) // TODO: Workaround until we can fetch the authentication status

	state, err = a.activateSessionTan(state)
	if err != nil {
		return nil, err
	}

	// Step 2.5: OAuth2 Comdirect Secondary Flow to extend scopes
	state, err = a.secondaryFlow(state)
	if err != nil {
		return nil, err
	}

	return &Authentication{
		accessToken: state.accessToken,
		sessionID:   state.requestInfo.ClientRequestID.SessionID,
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
func (a *Authenticator) passwordGrant(options *AuthOptions) (authState, error) {

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

	state := authState{
		accessToken: AccessToken{},
	}
	_, err := a.http.exchange(req, &state.accessToken)

	return state, err
}

// Step 2.2
func (a *Authenticator) fetchSessionStatus(state authState) (authState, error) {
	state.requestInfo = requestInfo{
		ClientRequestID: clientRequestID{
			SessionID: generateSessionID(),
			RequestID: generateRequestID(),
		},
	}
	requestInfoJson, err := json.Marshal(state.requestInfo)
	if err != nil {
		return state, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    comdirectURL("/api/session/clients/user/v1/sessions"),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJson)},
		},
	}

	var sessions []session
	if _, err = a.http.exchange(req, &sessions); err != nil {
		return state, err
	}

	if len(sessions) == 0 {
		return state, errors.New("length of returned session array is zero; expected at least one")
	}

	state.session = sessions[0]

	return state, err
}

// Step: 2.3
func (a *Authenticator) validateSessionTan(state authState) (authState, error) {
	state.session.SessionTanActive = true
	state.session.Activated2FA = true
	jsonSession, err := json.Marshal(state.session)
	if err != nil {
		return state, err
	}

	state.requestInfo.ClientRequestID.RequestID = generateRequestID()
	JSONRequestInfo, err := json.Marshal(state.requestInfo)
	if err != nil {
		return state, err
	}

	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s/validate", state.session.Identifier)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    comdirectURL(path),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(JSONRequestInfo)},
		},
		Body: body,
	}

	res, err := a.http.exchange(req, &state.session)
	if err != nil || res == nil {
		return state, err
	}

	jsonOnceAuthInfo := res.Header.Get(OnceAuthenticationInfoHeaderKey)
	if jsonOnceAuthInfo == "" {
		return state, errors.New("x-once-authentication-info header missing in response")
	}

	err = json.Unmarshal([]byte(jsonOnceAuthInfo), &state.onceAuthInfo)

	if err != nil {
		return state, err
	}

	return state, res.Body.Close()
}

// Step 2.4
func (a *Authenticator) activateSessionTan(state authState) (authState, error) {
	state.requestInfo.ClientRequestID.RequestID = generateRequestID()
	requestInfoJSON, err := json.Marshal(state.requestInfo)
	if err != nil {
		return state, err
	}

	onceAuthInfoHeader := fmt.Sprintf(`{"id":"%s"}`, state.onceAuthInfo.Id)
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s", state.session.Identifier)
	JSONSession, err := json.Marshal(state.session)
	req := &http.Request{
		Method: http.MethodPatch,
		URL:    comdirectURL(path),
		Header: http.Header{
			AuthorizationHeaderKey:          {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:                 {mediatype.ApplicationJson},
			ContentTypeHeaderKey:            {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey:        {string(requestInfoJSON)},
			OnceAuthenticationInfoHeaderKey: {onceAuthInfoHeader},
		},
		Body: ioutil.NopCloser(strings.NewReader(string(JSONSession))),
	}

	_, err = a.http.exchange(req, &state.session)

	if err != nil {
		return state, err
	}

	return state, nil
}

// Step 2.5
func (a *Authenticator) secondaryFlow(state authState) (authState, error) {

	body := url.Values{
		"token":         {state.accessToken.AccessToken},
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
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}

	res, err := a.http.Do(req)

	if err != nil {
		return state, err
	}

	if err = json.NewDecoder(res.Body).Decode(&state.accessToken); err != nil {
		return state, err
	}

	return state, res.Body.Close()
}
