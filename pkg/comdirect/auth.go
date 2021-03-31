package comdirect

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jsattler/comdirect-golang/internal/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TODO: Should the Authenticator struct be "stateless" (access token, request info, session etc. not part of struct).
// TODO: Currently the Authenticator is responsible for the AuthState and is responsible for state transitions. Not sure if this is a good approach.
// TODO: Not sure how and where to encapsulate state during the authentication flow. There is state and it should be clear how to handle it.
// TODO: Currently it is not possible to chose between TAN types (Photo TAN, Push TAN etc.). Only testing with P_TAN_PUSH at the moment.
// TODO: Let the user pass in a preconfigured http.Client to follow dependency injection principle.
// TODO: Think about accessibility of functions, structs, etc.
// TODO: Provide an interface

// Authenticator is responsible for authenticating against the comdirect REST API.
// It uses the given AuthOptions for authentication and returns an AccessToken in case
// the authentication flow was successful. Authenticator is using golang's default http.Client.
type Authenticator struct {
	authOptions *AuthOptions
	http        *http.Client
	authState   *AuthState
}

// AuthState encapsulates the state that is required for the comdirect authentication flow.
// The AuthState is passed between requests during the authentication flow.
type AuthState struct {
	lock                         sync.Mutex
	accessToken                  *AccessToken
	requestInfo                  *RequestInfo
	lastSuccessfulAuthentication *time.Time
	session                      *session
	onceAuthInfo                 *onceAuthenticationInfo
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

type RequestInfo struct {
	ClientRequestId ClientRequestId `json:"clientRequestId"`
}

type ClientRequestId struct {
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

// Create a new Authenticator from AuthOptions with a sessionId passed to the function.
// Default http.Client with a Timeout of 30 seconds will be used.
func (a *AuthOptions) NewAuthenticator() *Authenticator {
	return &Authenticator{
		authOptions: a,
		authState: &AuthState{
			requestInfo: newRequestInfo(),
		},
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
		authState: &AuthState{
			requestInfo: newRequestInfo(),
		},
		http: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}
}

// Authenticate against the comdirect REST API with given AuthOptions.
// Calling this function initializes a new AuthState using the session passed to the Authenticator.
func (a *Authenticator) Authenticate() (*AuthState, error) {
	a.authState.lock.Lock()         // acquire the lock for the full authentication process
	defer a.authState.lock.Unlock() // release the lock after authentication

	// Step 2.1: OAuth2 Resource Owner Password Credentials Grant
	if err := a.passwordGrant(); err != nil {
		return nil, err
	}

	// Step 2.2: Request the session status
	if err := a.fetchSessionStatus(); err != nil {
		return nil, err
	}

	// Step 2.3: Validate session TAN
	if err := a.validateSessionTan(); err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 10) // give the user 10 seconds to solve TAN challenge

	// Step 2.4: Activate session TAN
	if err := a.activateSessionTan(); err != nil {
		return nil, err
	}

	// Step 2.5: OAuth2 Comdirect Secondary Flow to extend scopes
	if err := a.secondaryFlow(); err != nil {
		return nil, err
	}
	successFullAuthTime := time.Now()
	a.authState.lastSuccessfulAuthentication = &successFullAuthTime

	return a.authState, nil
}

// Returns whether or not the Authenticator is still authenticated.
// This is an offline check, which means, that no request is made to check a valid authentication.
// The check is based on the information contained in the AuthState and AccessToken.
func (a *Authenticator) IsAuthenticated() bool {
	if a.authState.accessToken == nil {
		return false
	}

	if a.authState.lastSuccessfulAuthentication == nil {
		return false
	}

	expiresIn := time.Second * time.Duration(a.authState.accessToken.ExpiresIn)
	lastSuccessfulAuthentication := a.authState.lastSuccessfulAuthentication
	expires := lastSuccessfulAuthentication.Add(expiresIn)

	return expires.After(time.Now())
}

func (a *Authenticator) Refresh() {
	// TODO
}

func (a *Authenticator) Revoke() {
	// TODO
}

// Step 2.1
// TODO: outsource state validation
func (a *Authenticator) passwordGrant() error {

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

	a.authState.accessToken = &AccessToken{}
	_, err := a.request(req, a.authState.accessToken)

	if err != nil {
		return err
	}

	return nil
}

// Step 2.2
// TODO: outsource state validation
func (a *Authenticator) fetchSessionStatus() error {

	if a.authState.requestInfo == nil {
		return errors.New("requestInfo of authState cannot be nil")
	}
	if a.authState.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	requestInfoJson, err := json.Marshal(a.authState.requestInfo)
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    comdirectUrl("/api/session/clients/user/v1/sessions"),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + a.authState.accessToken.AccessToken},
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

	a.authState.session = &sessions[0]

	return nil
}

// Step: 2.3
// TODO: outsource state validation
func (a *Authenticator) validateSessionTan() error {

	if a.authState.session == nil {
		return errors.New("session of authState cannot be nil")
	}
	if a.authState.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	a.authState.session.Activated2FA = true
	a.authState.session.SessionTanActive = true
	jsonSession, err := json.Marshal(a.authState.session)
	if err != nil {
		return err
	}

	a.authState.requestInfo.updateRequestId()

	jsonRequestInfo, err := json.Marshal(a.authState.requestInfo)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s/validate", a.authState.session.Identifier)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))
	req := &http.Request{
		Method: http.MethodPost,
		URL:    comdirectUrl(path),
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + a.authState.accessToken.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(jsonRequestInfo)},
		},
		Body: body,
	}

	res, err := a.request(req, a.authState.session)

	if err != nil || res == nil {
		return err
	}

	if _, ok := res.Header[OnceAuthenticationInfoHeaderKey]; !ok {
		return errors.New("x-once-authentication-info header missing in response")
	}

	jsonOnceAuthInfo := res.Header.Get(OnceAuthenticationInfoHeaderKey)

	a.authState.onceAuthInfo = &onceAuthenticationInfo{}
	err = json.Unmarshal([]byte(jsonOnceAuthInfo), a.authState.onceAuthInfo)

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
func (a *Authenticator) activateSessionTan() error {
	if a.authState.session == nil {
		return errors.New("session of authState cannot be nil")
	}

	if a.authState.onceAuthInfo == nil {
		return errors.New("onceAuthInfo of authState cannot be nil")
	}

	if a.authState.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	a.authState.requestInfo.updateRequestId()

	requestInfoJson, err := json.Marshal(a.authState.requestInfo)
	if err != nil {
		return err
	}

	onceAuthInfoHeader := fmt.Sprintf(`{"id":"%s"}`, a.authState.onceAuthInfo.Id)
	path := fmt.Sprintf("/api/session/clients/user/v1/sessions/%s", a.authState.session.Identifier)
	jsonSession, err := json.Marshal(a.authState.session)
	body := ioutil.NopCloser(strings.NewReader(string(jsonSession)))

	req := &http.Request{
		Method: http.MethodPatch,
		URL:    comdirectUrl(path),
		Header: http.Header{
			AuthorizationHeaderKey:          {BearerPrefix + a.authState.accessToken.AccessToken},
			AcceptHeaderKey:                 {mediatype.ApplicationJson},
			ContentTypeHeaderKey:            {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey:        {string(requestInfoJson)},
			OnceAuthenticationInfoHeaderKey: {onceAuthInfoHeader},
		},
		Body: body,
	}

	_, err = a.request(req, a.authState.session)

	if err != nil {
		return err
	}

	return nil
}

// Step 2.5
// TODO: outsource state validation
func (a *Authenticator) secondaryFlow() error {

	if a.authState.accessToken == nil {
		return errors.New("accessToken of authState cannot be nil")
	}

	body := url.Values{
		"token":         {a.authState.accessToken.AccessToken},
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

	if err = json.NewDecoder(res.Body).Decode(a.authState.accessToken); err != nil {
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

func (r *RequestInfo) updateRequestId() {
	requestId := generateRequestId()
	r.ClientRequestId.RequestId = requestId
}

func newRequestInfo() *RequestInfo {
	return &RequestInfo{
		ClientRequestId: ClientRequestId{
			SessionId: generateSessionId(),
			RequestId: generateRequestId(),
		},
	}
}

// Construct a url.URL with Host 'api.comdirect.de' HttpsScheme and a given path.
func comdirectUrl(path string) *url.URL {
	return &url.URL{Host: Host, Scheme: HttpsScheme, Path: path}
}

func generateSessionId() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%032d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func generateRequestId() string {
	unix := time.Now().Unix()
	id := fmt.Sprintf("%09d", unix)
	return id[0:9]
}
