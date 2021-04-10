package comdirect

import (
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

	PasswordGrantType  = "password"
	SecondaryGrantType = "cd_secondary"

	DefaultHttpTimeout = time.Second * 30
	HttpsScheme        = "https"
	BearerPrefix       = "Bearer "
)

type Client struct {
	authenticator  *Authenticator
	http           *HttpClient
	authentication Authentication
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
	AccountID        string      `json:"accountId"`
	AccountDisplayID string      `json:"accountDisplayId"`
	Currency         string      `json:"currency"`
	ClientID         string      `json:"clientId"`
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
	DirectDebitCreditorID string          `json:"directDebitCreditorId"`
	DirectDebitMandateID  string          `json:"directDebitMandateId"`
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
		http:          &HttpClient{http.Client{Timeout: DefaultHttpTimeout}},
	}
}

func NewWithAuthOptions(options *AuthOptions) *Client {
	return NewWithAuthenticator(options.NewAuthenticator())
}
