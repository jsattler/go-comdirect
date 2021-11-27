package comdirect

import (
	"errors"
	"fmt"
	"net/http"
)

type Depot struct {
	DepotId                    string   `json:"depotId"`
	DepotDisplayId             string   `json:"depotDisplayId"`
	ClientId                   string   `json:"clientId"`
	DefaultSettlementAccountId string   `json:"defaultSettlementAccountId"`
	SettlementAccountIds       []string `json:"settlementAccountIds"`
}

type DepotPosition struct {
	DepotId                  string      `json:"depotId"`
	PositionId               string      `json:"positionId"`
	Wkn                      string      `json:"wkn"`
	CustodyType              string      `json:"custodyType"`
	Quantity                 AmountValue `json:"quantity"`
	AvailableQuantity        AmountValue `json:"availableQuantity"`
	CurrentPrice             Price       `json:"currentPrice"`
	PrevDayPrice             Price       `json:"prevDayPrice"`
	CurrentValue             AmountValue `json:"currentValue"`
	PurchaseValue            AmountValue `json:"purchaseValue"`
	ProfitLossPurchaseAbs    AmountValue `json:"profitLossPurchaseAbs"`
	ProfitLossPurchaseRel    string      `json:"profitLossPurchaseRel"`
	ProfitLossPrevDayAbs     AmountValue `json:"profitLossPrevDayAbs"`
	ProfitLossPrevDayRel     string      `json:"profitLossPrevDayRel"`
	AvailableQuantityToHedge AmountValue `json:"availableQuantityToHedge"`
}

type Price struct {
	Price         AmountValue `json:"price"`
	PriceDateTime string      `json:"priceDateTime"`
}

type DepotTransaction struct {
	TransactionID        string      `json:"transactionId"`
	Instrument           Instrument  `json:"instrument"`
	ExecutionPrice       AmountValue `json:"executionPrice"`
	TransactionValue     AmountValue `json:"transactionValue"`
	TransactionDirection string      `json:"transactionDirection"`
	TransactionType      string      `json:"transactionType"`
	FXRate               string      `json:"fxRate"`
}

type DepotTransactions struct {
	Values []DepotTransaction `json:"values"`
}

type Positions struct {
	Values []DepotPosition `json:"values"`
}

type Depots struct {
	Values []Depot `json:"values"`
}

// Depots retrieves all depots for the current Authentication.
func (c *Client) Depots() ([]Depot, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/brokerage/clients/user/v3/depots"),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	depots := &Depots{}
	_, err = c.http.exchange(req, depots)
	return depots.Values, err
}

// DepotPositions retrieves all positions for a specific depot ID.
func (c *Client) DepotPositions(depotID string) ([]DepotPosition, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/brokerage/v3/depots/%s/positions", depotID)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	depots := &Positions{}
	_, err = c.http.exchange(req, depots)
	return depots.Values, err
}

// DepotPosition retrieves a position by its ID from the depot specified by its ID.
func (c *Client) DepotPosition(depotID string, positionID string) (*DepotPosition, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/brokerage/v3/depots/%s/positions/%s", depotID, positionID)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	position := &DepotPosition{}
	_, err = c.http.exchange(req, position)
	return position, err
}

// DepotTransactions retrieves all transactions for a depot specified by its ID.
func (c *Client) DepotTransactions(depotID string) ([]DepotTransaction, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/brokerage/v3/depots/%s/transactions", depotID)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	depotTransactions := &DepotTransactions{}
	_, err = c.http.exchange(req, depotTransactions)
	return depotTransactions.Values, err
}
