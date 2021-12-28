package comdirect

import (
	"context"
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
	Paging Paging           `json:"paging"`
	Values []AccountBalance `json:"values"`
}

type AccountTransactions struct {
	Paging Paging               `json:"paging"`
	Values []AccountTransaction `json:"values"`
}

type AccountTransaction struct {
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

func (c *Client) Balances(ctx context.Context) (*AccountBalances, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/banking/clients/user/v2/accounts/balances"),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}
	req = req.WithContext(ctx)

	accountBalances := &AccountBalances{}
	_, err = c.http.exchange(req, accountBalances)
	return accountBalances, err
}

func (c *Client) Balance(ctx context.Context, accountId string) (*AccountBalance, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}

	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/banking/v2/accounts/%s/balances", accountId)),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}
	accountBalance := &AccountBalance{}
	_, err = c.http.exchange(req, accountBalance)

	return accountBalance, err
}

func (c *Client) Transactions(ctx context.Context, accountId string, options ...Options) (*AccountTransactions, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		log.Fatal(err)
	}
	url := apiURL(fmt.Sprintf("/banking/v1/accounts/%s/transactions", accountId))
	encodeOptions(url, options)
	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	tr := &AccountTransactions{}
	_, err = c.http.exchange(req, tr)

	return tr, err
}
