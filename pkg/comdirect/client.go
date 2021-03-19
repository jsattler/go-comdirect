package comdirect

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// TODO: remove duplicated code fragments
// TODO: error handling
// TODO: check regarding separation of concerns

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

type AccountBalance struct {
	Account                Account     `json:"account"`
	AccountId              string      `json:"accountId"`
	Balance                AmountValue `json:"balance"`
	BalanceEUR             AmountValue `json:"balanceEUR"`
	AvailableCashAmount    AmountValue `json:"availableCashAmount"`
	AvailableCashAmountEUR AmountValue `json:"availableCashAmountEUR"`
}

type Account struct {
	AccountId        string      `json:"accountId"`
	AccountDisplayId string      `json:"accountDisplayId"`
	Currency         string      `json:"currency"`
	ClientId         string      `json:"clientId"`
	AccountType      AccountType `json:"accountType"`
	Iban             string      `json:"iban"`
	CreditLimit      AmountValue `json:"creditLimit"`
}

type AccountType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type AmountValue struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

type AccountBalances struct {
	Values []AccountBalance `json:"values"`
}

type TransactionResponse struct {
	Values []Transaction `json:"values"`
}

type Transaction struct {
	Reference             string          `json:"reference"`
	BookingStatus         string          `json:"bookingStatus"`
	BookingDate           string          `json:"bookingDate"`
	Amount                AmountValue     `json:"amount"`
	Remitter              Remitter        `json:"remitter"`
	Deptor                string          `json:"deptor"`
	Creditor              Creditor        `json:"creditor"`
	ValutaDate            string          `json:"valutaDate"`
	DirectDebitCreditorId string          `json:"directDebitCreditorId"`
	DirectDebitMandateId  string          `json:"directDebitMandateId"`
	EndToEndReference     string          `json:"endToEndReference"`
	NewTransaction        bool            `json:"newTransaction"`
	RemittanceInfo        string          `json:"remittanceInfo"`
	TransactionType       TransactionType `json:"transactionType"`
}

type TransactionType struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

type Remitter struct {
	HolderName string `json:"holderName"`
}

type Creditor struct {
	HolderName string `json:"holderName"`
	Iban       string `json:"iban"`
	Bic        string `json:"bic"`
}

func NewWithAuthenticator(authenticator *Authenticator) *Client {
	return &Client{
		authenticator: authenticator,
		http:          &http.Client{Timeout: DefaultHttpTimeout},
	}
}

func NewWithAuthOptions(options *AuthOptions) *Client {
	return NewWithAuthenticator(options.NewAuthenticator())
}

// TODO: Refactor
func (c *Client) Balances() ([]AccountBalance, error) {
	var authState *AuthState
	var err error
	if !c.authenticator.IsAuthenticated() { // not authenticated
		authState, err = c.authenticator.Authenticate()
		if err != nil {
			return nil, err
		}
	} else { // already authenticated
		authState = c.authenticator.authState
	}

	riJson, err := json.Marshal(authState.requestInfo)
	if err != nil {
		log.Fatal(err)
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiUrl("/banking/clients/user/v2/accounts/balances"),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + authState.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(riJson)},
		},
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	abr := &AccountBalances{}
	err = json.NewDecoder(res.Body).Decode(abr)
	if err != nil {
		return nil, err
	}

	return abr.Values, nil
}

// TODO: Refactor
func (c *Client) Balance(accountId string) (*AccountBalance, error) {
	var authState *AuthState
	var err error
	if !c.authenticator.IsAuthenticated() { // not authenticated
		authState, err = c.authenticator.Authenticate()
		if err != nil {
			return nil, err
		}
	} else { // already authenticated
		authState = c.authenticator.authState
	}

	riJson, err := json.Marshal(authState.requestInfo)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("/banking/v2/accounts/%s/balances", accountId)
	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiUrl(path),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + authState.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(riJson)},
		},
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	abr := &AccountBalance{}
	err = json.NewDecoder(res.Body).Decode(abr)
	if err != nil {
		return nil, err
	}

	return abr, nil

}

// TODO: Refactor
func (c *Client) Transactions(accountId string) ([]Transaction, error) {
	var authState *AuthState
	var err error
	if !c.authenticator.IsAuthenticated() { // not authenticated
		authState, err = c.authenticator.Authenticate()
		if err != nil {
			return nil, err
		}
	} else { // already authenticated
		authState = c.authenticator.authState
	}

	riJson, err := json.Marshal(authState.requestInfo)
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("/banking/v1/accounts/%s/transactions", accountId)
	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiUrl(path),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + authState.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(riJson)},
		},
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	abr := &TransactionResponse{}
	err = json.NewDecoder(res.Body).Decode(abr)
	if err != nil {
		return nil, err
	}

	return abr.Values, nil

}

func apiUrl(path string) *url.URL {
	return &url.URL{
		Host:   Host,
		Scheme: HttpsScheme,
		Path:   ApiPath + path,
	}
}
