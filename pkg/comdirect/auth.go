package comdirect

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/j-sattler/comdirect-golang/internal/httpstatus"
	"github.com/j-sattler/comdirect-golang/internal/mediatype"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: lots of refactoring required, first attempt for auth flow

type Authenticator struct {
	AuthOptions *AuthOptions
	http        *http.Client
	SessionId   string
}

type AuthOptions struct {
	Username     string
	Password     string
	ClientId     string
	ClientSecret string
	AutoRefresh  bool
}

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

type SessionObject struct {
	Identifier       string `json:"identifier"`
	SessionTanActive bool   `json:"sessionTanActive"`
	Activated2FA     bool   `json:"activated2FA"`
}

type OnceAuthenticationInfo struct {
	Id             string   `json:"id"`
	Typ            string   `json:"typ"`
	AvailableTypes []string `json:"availableTypes"`
	Link           Link     `json:"link"`
}

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

func (a *AuthOptions) NewAuthenticator(sessionId string) *Authenticator {
	return &Authenticator{
		AuthOptions: a,
		SessionId:   sessionId,
		http: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}
}

func NewAuthenticator(options *AuthOptions, sessionId string) *Authenticator {
	return &Authenticator{
		AuthOptions: options,
		SessionId:   sessionId,
		http: &http.Client{
			Timeout: DefaultHttpTimeout,
		},
	}
}

func (a *Authenticator) Authenticate() (*AccessToken, error) {
	accessToken, err := a.fetchToken()
	if err != nil {
		return nil, err
	}

	sessionObject, err := a.fetchSessionStatus(accessToken)

	if err != nil {
		return nil, err
	}

	_, err = a.validateSessionTan(sessionObject, accessToken)

	if err != nil {
		return nil, err
	}

	secondaryFlow, err := a.secondaryFlow(accessToken)

	if err != nil {
		return nil, err
	}

	return secondaryFlow, nil
}

func (*Authenticator) Refresh() {
	// TODO
}

func (*Authenticator) Revoke() {
	// TODO
}

// Step 2.1
func (a *Authenticator) fetchToken() (*AccessToken, error) {
	body := url.Values{
		"username":      {a.AuthOptions.Username},
		"password":      {a.AuthOptions.Password},
		"grant_type":    {PasswordGrantType},
		"client_id":     {a.AuthOptions.ClientId},
		"client_secret": {a.AuthOptions.ClientSecret},
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
		return nil, err
	}

	if httpstatus.Is4xx(res) {
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		log.Println(string(bytes))
	}

	token := &AccessToken{}
	if err = json.NewDecoder(res.Body).Decode(token); err != nil {
		return nil, err
	}

	log.Println(token.AccessToken)
	return token, res.Body.Close()
}

// Step 2.2
func (a *Authenticator) fetchSessionStatus(token *AccessToken) (*SessionObject, error) {
	requestInfo := RequestInfo{
		ClientRequestId: ClientRequestId{
			SessionId: a.SessionId,
			RequestId: GenerateRequestId(),
		},
	}

	requestInfoJson, err := json.Marshal(requestInfo)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Host:   Host,
			Scheme: HttpsScheme,
			Path:   "/api/session/clients/user/v1/sessions",
		},
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + token.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJson)},
		},
	}

	res, err := a.http.Do(req)
	if err != nil {
		return nil, err
	}

	var sessionObject []SessionObject
	if err = json.NewDecoder(res.Body).Decode(&sessionObject); err != nil {
		return nil, err
	}

	if len(sessionObject) == 0 {
		return nil, errors.New("length of session object array is zero; expected minimum one")
	}

	return &sessionObject[0], res.Body.Close()
}

// Step: 2.3
func (a *Authenticator) validateSessionTan(sessionObject *SessionObject, token *AccessToken) (*SessionObject, error) {
	sessionObject.Activated2FA = true
	sessionObject.SessionTanActive = true
	sessionObjectJson, err := json.Marshal(sessionObject)

	if err != nil {
		return nil, err
	}

	requestInfo := &RequestInfo{
		ClientRequestId: ClientRequestId{
			SessionId: a.SessionId,
			RequestId: GenerateRequestId(),
		},
	}

	requestInfoJson, err := json.Marshal(requestInfo)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Host:   Host,
			Scheme: HttpsScheme,
			Path:   fmt.Sprintf("/api/session/clients/user/v1/sessions/%s/validate", sessionObject.Identifier),
		},
		Header: http.Header{
			AuthorizationHeaderKey:   {BearerPrefix + token.AccessToken},
			AcceptHeaderKey:          {mediatype.ApplicationJson},
			ContentTypeHeaderKey:     {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey: {string(requestInfoJson)},
		},
		Body: ioutil.NopCloser(strings.NewReader(string(sessionObjectJson))),
	}

	res, err := a.http.Do(req)
	if err != nil {
		return nil, err
	}

	newSessionObject := &SessionObject{}
	if err = json.NewDecoder(res.Body).Decode(newSessionObject); err != nil {
		return nil, err
	}

	if _, ok := res.Header[OnceAuthenticationInfoHeaderKey]; !ok {
		return nil, errors.New("missing once-authentication-info header")
	}

	onceAuthInfoJson := res.Header.Get(OnceAuthenticationInfoHeaderKey)

	onceAuthInfo := &OnceAuthenticationInfo{}
	err = json.Unmarshal([]byte(onceAuthInfoJson), onceAuthInfo)

	if err != nil {
		return nil, err
	}

	// comdirect broke their API. Now this is not working anymore...
	//err = a.isActivated(onceAuthInfo, string(requestInfoJson), token)

	time.Sleep(time.Second * 30) // give the user 30 seconds to activate Session TAN

	err = a.activateSessionTan(sessionObject, token, onceAuthInfo)

	return newSessionObject, res.Body.Close()
}

func (a *Authenticator) isActivated(info *OnceAuthenticationInfo, requestInfo string, token *AccessToken) (bool, error) {

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
func (a *Authenticator) activateSessionTan(object *SessionObject, token *AccessToken, info *OnceAuthenticationInfo) error {

	requestInfo := &RequestInfo{
		ClientRequestId: ClientRequestId{
			SessionId: a.SessionId,
			RequestId: GenerateRequestId(),
		},
	}

	requestInfoJson, err := json.Marshal(requestInfo)
	if err != nil {
		return err
	}

	soJson, err := json.Marshal(object)

	req := &http.Request{
		Method: http.MethodPatch,
		URL: &url.URL{
			Host:   Host,
			Scheme: HttpsScheme,
			Path:   fmt.Sprintf("/api/session/clients/user/v1/sessions/%s", object.Identifier),
		},
		Header: http.Header{
			AuthorizationHeaderKey:          {BearerPrefix + token.AccessToken},
			AcceptHeaderKey:                 {mediatype.ApplicationJson},
			ContentTypeHeaderKey:            {mediatype.ApplicationJson},
			HttpRequestInfoHeaderKey:        {string(requestInfoJson)},
			OnceAuthenticationInfoHeaderKey: {fmt.Sprintf(`{"id":"%s"}`, info.Id)},
		},
		Body: ioutil.NopCloser(strings.NewReader(string(soJson))),
	}

	res, err := a.http.Do(req)

	if err != nil {
		return err
	}

	log.Println(res.StatusCode)

	return nil
}

// Step 2.5
func (a *Authenticator) secondaryFlow(token *AccessToken) (*AccessToken, error) {

	body := url.Values{
		"token":         {token.AccessToken},
		"grant_type":    {ComdirectSecondaryGrantType},
		"client_id":     {a.AuthOptions.ClientId},
		"client_secret": {a.AuthOptions.ClientSecret},
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
		return nil, err
	}

	finalToken := &AccessToken{}
	if err = json.NewDecoder(res.Body).Decode(finalToken); err != nil {
		return nil, err
	}

	log.Println("Final token:", token.AccessToken)
	return finalToken, res.Body.Close()
}
