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
	HolderName                 string   `json:"holderName"`
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

type DepotAggregated struct {
	Depot                 Depot       `json:"depot"`
	PrevDayValue          AmountValue `json:"prevDayValue"`
	CurrentValue          AmountValue `json:"currentValue"`
	PurchaseValue         AmountValue `json:"purchaseValue"`
	ProfitLossPurchaseAbs AmountValue `json:"ProfitLossPurchaseAbs"`
	ProfitLossPurchaseRel string      `json:"profitLossPurchaseRel"`
	ProfitLossPrevDayAbs  AmountValue `json:"profitLossPrevDayAbs"`
	ProfitLossPrevDayRel  string      `json:"profitLossPrevDayRel"`
}

type DepotTransactions struct {
	Paging Paging             `json:"paging"`
	Values []DepotTransaction `json:"values"`
}

type DepotPositions struct {
	Paging     Paging          `json:"paging"`
	Aggregated DepotAggregated `json:"aggregated"`
	Values     []DepotPosition `json:"values"`
}

type Depots struct {
	Paging Paging  `json:"paging"`
	Values []Depot `json:"values"`
}

// Depots retrieves all depots for the current Authentication.
func (c *Client) Depots() (*Depots, error) {
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
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	depots := &Depots{}
	_, err = c.http.exchange(req, depots)
	return depots, err
}

// DepotPositions retrieves all positions for a specific depot ID.
func (c *Client) DepotPositions(depotID string, options ...Options) (*DepotPositions, error) {
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
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	depots := &DepotPositions{}
	_, err = c.http.exchange(req, depots)
	return depots, err
}

// DepotPosition retrieves a position by its ID from the depot specified by its ID.
func (c *Client) DepotPosition(depotID string, positionID string, options ...Options) (*DepotPosition, error) {
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
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	positions := &DepotPosition{}
	_, err = c.http.exchange(req, positions)
	return positions, err
}

// DepotTransactions retrieves all transactions for a depot specified by its ID.
func (c *Client) DepotTransactions(depotID string, options ...Options) (*DepotTransactions, error) {
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
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	depotTransactions := &DepotTransactions{}
	_, err = c.http.exchange(req, depotTransactions)
	return depotTransactions, err
}
