package comdirect

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/j-sattler/comdirect-golang/internal/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: Should the Authenticator struct be "stateless" (access token, request info, session etc. not part of struct).
// TODO: Not sure how and where to encapsulate state during the authentication flow. There is state and it should be clear how to handle it.
// TODO: Currently it is not possible to chose between TAN types (Photo TAN, Push TAN etc.). Only testing with P_TAN_PUSH at the moment.
// TODO: Let the user pass in a preconfigured http.Client to follow dependency injection principle.
// TODO: Think about accessibility of functions, structs, etc.

// Authenticator is responsible for authenticating against the comdirect REST API.
// It uses the given AuthOptions for authentication and returns an AccessToken in case
// the authentication flow was successful. Authenticator is using golang's default http.Client.
type Authenticator struct {
	authOptions *AuthOptions
	http        *http.Client
	accessToken *AccessToken
}

// authState encapsulates the state that is required for the comdirect authentication flow.
// The authState is passed between requests during the authentication flow.
type authState struct {
	accessToken  *AccessToken
	session      *session
	requestInfo  *requestInfo
	onceAuthInfo *onceAuthenticationInfo
}

// AuthOptions encapsulates the information for authentication against the comdirect REST API.
// Username also used for signing in via web frontend (Zugangsnummer).
// Password also used for signing in via web frontend (Online Banking PIN).
// ClientId needs to be generated over the comdirect web frontend.
// ClientSecret needs to be generated over the comdirect web frontend.
// AutoRefresh indicates whether or not the Authenticator should automatically refresh the AccessToken.
type AuthOptions struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	AutoRefresh  bool
}

// AccessToken represents an OAuth2 token that is returned after authenticating with the comdirect REST API.
type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	CustomerId   string `json:"kdnr"`
	BpId         int    `json:"bpid"`
	ContactId    int    `json:"kontaktId"`
}

type requestInfo struct {
	ClientRequestId clientRequestId `json:"clientRequestId"`
}

type clientRequestId struct {
	SessionId string `json:"sessionId"`
	RequestId string `json:"requestId"`
}

type session struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}

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

type OAuther interface {
	Authenticate()
	Refresh()
	Revoke()
}

// Create a new Authenticator from AuthOptions with a sessionId passed to the function.
// Default http.Client with a Timeout of 30 seconds will be used.
func (a *AuthOptions) NewAuthenticator() *Authenticator {
	return &Authenticator{
		authOptions: a,
		http: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}
}

// Create a new Authenticator by passing in AuthOptions and sessionId.
// Default http.Client with a Timeout of 30 seconds will be used.
func NewAuthenticator(options *AuthOptions) *Authenticator {
	return &Authenticator{
		authOptions: options,
		http: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}
}

// Authenticate against the comdirect REST API with given AuthOptions.
// Calling this function initializes a new authState using the session passed to the Authenticator.
func (a *Authenticator) Authenticate() (*AccessToken, error) {

	state, err := a.initializeAuthState()

	if err != nil {
		return nil, err
	}

	// Step 2.1: OAuth2 Resource Owner Password Credentials Grant
	if err = a.passwordGrant(state); err != nil {
		return nil, err
	}

	// Step 2.2: Request the session status
	if err = a.fetchSessionStatus(state); err != nil {
		return nil, err
	}

	// Step 2.3: Validate session TAN
	if err = a.validateSessionTan(state); err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 30) // give the user 30 seconds to solve TAN challenge

	// Step 2.4: Activate session TAN
	if err = a.activateSessionTan(state); err != nil {
		return nil, err
	}

	// Step 2.5: OAuth2 Comdirect Secondary Flow to extend scopes
	if err = a.secondaryFlow(state); err != nil {
		return nil, err
	}

	return state.accessToken, nil
}

func IsAuthenticated() bool {
	// TODO: decide based on AccessToken information
	return false
}

func (*Authenticator) Refresh() {
	// TODO
}

func (*Authenticator) Revoke() {
	// TODO
}

// Step 2.1
// TODO: outsource state validation
func (a *Authenticator) passwordGrant(state *authState) error {

	if state.requestInfo.ClientRequestId.SessionId == "" {
		return errors.New("sessionId of authState cannot be empty")
	}

	urlEncoded := url.Values{
		"username":      {a.authOptions.Username},
		"password":      {a.authOptions.Password},
		"grant_type":    {PasswordGrantType},
		"client_id":     {a.authOptions.ClientId},
		"client_secret": {a.authOptions.ClientSecret},
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

	state.accessToken = &AccessToken{}
	_, err := a.request(req, state.accessToken)

	if err != nil {
		return err
	}

	return nil
}

// Step 2.2
// TODO: outsource state validation
func (a *Authenticator) fetchSessionStatus(state *authState) error {

	if state.requestInfo == nil {
		return errors.New("requestInfo of authState cannot be nil")
	}
	if state.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	requestInfoJson, err := json.Marshal(state.requestInfo)
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    comdirectUrl("/api/session/clients/user/v1/sessions"),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJson)},
		},
	}

	var sessions []session
	if _, err = a.request(req, &sessions); err != nil {
		return err
	}

	if len(sessions) == 0 {
		return errors.New("length of returned session array is zero; expected at least one")
	}

	state.session = &sessions[0]

	return nil
}

// Step: 2.3
// TODO: outsource state validation
func (a *Authenticator) validateSessionTan(state *authState) error {

	if state.session == nil {
		return errors.New("session of authState cannot be nil")
	}
	if state.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	state.session.Activated2FA = true
	state.session.SessionTanActive = true
	jsonSession, err := json.Marshal(state.session)
	if err != nil {
		return err
	}

	if err = state.updateRequestId(); err != nil {
		return err
	}
	jsonRequestInfo, err := json.Marshal(state.requestInfo)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s/validate", state.session.Identifier)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    comdirectUrl(path),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(jsonRequestInfo)},
		},
		Body: body,
	}

	res, err := a.request(req, state.session)

	if err != nil || res == nil {
		return err
	}

	if _, ok := res.Header[OnceAuthenticationInfoHeaderKey]; !ok {
		return errors.New("x-once-authentication-info header missing in response")
	}

	jsonOnceAuthInfo := res.Header.Get(OnceAuthenticationInfoHeaderKey)

	state.onceAuthInfo = &onceAuthenticationInfo{}
	err = json.Unmarshal([]byte(jsonOnceAuthInfo), state.onceAuthInfo)

	if err != nil {
		return err
	}

	return res.Body.Close()
}

func (a *Authenticator) isActivated(info *onceAuthenticationInfo, requestInfo string, token *AccessToken) (bool, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Host:   Host,
			Scheme: HttpsScheme,
			Path:   info.Link.Href,
		},
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + token.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {requestInfo},
		},
	}

	res, err := a.http.Do(req)

	if err != nil {
		return false, err
	}

	_, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Step 2.4
// TODO: outsource state validation
func (a *Authenticator) activateSessionTan(state *authState) error {
	if state.session == nil {
		return errors.New("session of authState cannot be nil")
	}

	if state.onceAuthInfo == nil {
		return errors.New("onceAuthInfo of authState cannot be nil")
	}

	if state.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	if err := state.updateRequestId(); err != nil {
		return err
	}
	requestInfoJson, err := json.Marshal(state.requestInfo)
	if err != nil {
		return err
	}

	onceAuthInfoHeader := fmt.Sprintf(`{"id":"%s"}`, state.onceAuthInfo.Id)
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s", state.session.Identifier)
	jsonSession, err := json.Marshal(state.session)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))

	req := &http.Request{
		Method: http.MethodPatch,
		URL:    comdirectUrl(path),
		Header: http.Header{
			AuthorizationHeaderKey:          {BearerPrefix + state.accessToken.AccessToken},
			AcceptHeaderKey:                 {mediatype.ApplicationJson},
			ContentTypeHeaderKey:            {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey:        {string(requestInfoJson)},
			OnceAuthenticationInfoHeaderKey: {onceAuthInfoHeader},
		},
		Body: body,
	}

	_, err = a.request(req, state.session)

	if err != nil {
		return err
	}

	return nil
}

// Step 2.5
// TODO: outsource state validation
func (a *Authenticator) secondaryFlow(state *authState) error {

	if state.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	body := url.Values{
		"token":         {state.accessToken.AccessToken},
		"grant_type":    {ComdirectSecondaryGrantType},
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
		return err
	}

	if err = json.NewDecoder(res.Body).Decode(state.accessToken); err != nil {
		return err
	}

	return res.Body.Close()
}

func (a *Authenticator) request(request *http.Request, target interface{}) (*http.Response, error) {
	res, err := a.http.Do(request)
	if err != nil {
		return res, err
	}

	if err = json.NewDecoder(res.Body).Decode(target); err != nil {
		return res, err
	}

	return res, res.Body.Close()
}

func (a *Authenticator) initializeAuthState() (*authState, error) {

	state := &authState{}
	if err := state.initializeRequestInfo(); err != nil {
		return nil, err
	}

	return state, nil
}

func (a *authState) updateRequestId() error {
	if a.requestInfo == nil {
		err := a.initializeRequestInfo()
		if err != nil {
			return err
		}
	}

	requestId := generateRequestId()
	a.requestInfo.ClientRequestId.RequestId = requestId
	return nil
}

func (a *authState) initializeRequestInfo() error {
	requestId := generateRequestId()
	sessionId, err := GenerateSessionId()

	if err != nil {
		return err
	}

	a.requestInfo = &requestInfo{
		ClientRequestId: clientRequestId{
			SessionId: sessionId,
			RequestId: requestId,
		},
	}
	return nil
}

// Construct a url.URL with Host 'api.comdirect.de' HttpsScheme and a given path.
func comdirectUrl(path string) *url.URL {
	return &url.URL{Host: Host, Scheme: HttpsScheme, Path: path}
}
