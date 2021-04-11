package comdirect

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

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

func (c *Client) Balances() ([]AccountBalance, error) {
	if c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/banking/clients/user/v2/accounts/balances"),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	accountBalances := &AccountBalances{}
	_, err = c.http.exchange(req, accountBalances)
	return accountBalances.Values, err
}

func (c *Client) Balance(accountId string) (*AccountBalance, error) {
	if c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}

	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/banking/v2/accounts/%s/balances", accountId)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}
	accountBalance := &AccountBalance{}
	_, err = c.http.exchange(req, accountBalance)

	return accountBalance, err
}

func (c *Client) Transactions(accountId string) ([]Transaction, error) {
	if c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		log.Fatal(err)
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/banking/v1/accounts/%s/transactions", accountId)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	tr := &TransactionResponse{}
	_, err = c.http.exchange(req, tr)

	return tr.Values, err
}
