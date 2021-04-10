package comdirect

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) Balances() ([]AccountBalance, error) {

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
	if c.authentication.accessToken == nil || c.authentication.IsExpired() {
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
