package comdirect

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	CLIENT_NOT_AUTHENTICATED = "client is not authenticated"
)

type Dimension struct {
	Venues []Venue `json:"venues"`
}

type Venue struct {
	Name          string     `json:"name"`
	VenueID       string     `json:"venueId"`
	Country       string     `json:"country"`
	Type          string     `json:"type"`
	Currencies    []string   `json:"currencies"`
	Sides         []string   `json:"sides"`
	ValidityTypes []string   `json:"validityTypes"`
	OrderTypes    OrderTypes `json:"orderTypes"`
}

type Dimensions struct {
	Paging Paging      `json:"paging"`
	Values []Dimension `json:"values"`
}

type OrderTypes struct {
	Quote              OrderType `json:"QUOTE"`
	Market             OrderType `json:"MARKET"`
	StopMarket         OrderType `json:"STOP_MARKET"`
	NextOrder          OrderType `json:"NEXT_ORDER"`
	OneCancelsOther    OrderType `json:"ONE_CANCELS_ORDER"`
	Limit              OrderType `json:"LIMIT"`
	TrailingStopMarket OrderType `json:"TRAILING_STOP_MARKET"`
}

type OrderType struct {
	LimitExtensions     []string `json:"limitExtensions"`
	TradingRestrictions []string `json:"tradingRestrictions"`
}

type OrderRequest struct {
	DepotID      string      `json:"depotId,omitempty"`
	OrderID      string      `json:"orderId,omitempty"`
	Side         string      `json:"side"`
	InstrumentID string      `json:"instrumentId"`
	OrderType    string      `json:"orderType"`
	Quantity     AmountValue `json:"quantity"`
	VenueID      string      `json:"venueId"`
	Limit        AmountValue `json:"limit"`
	ValidityType string      `json:"validityType"`
	Validity     string      `json:"validity"`
}

type Orders struct {
	Paging Paging  `json:"paging"`
	Values []Order `json:"values"`
}

type Order struct {
	DepotID             string      `json:"depotId"`
	SettlementAccountID string      `json:"settlementAccountId"`
	OrderID             string      `json:"orderID"`
	CreationTimestamp   string      `json:"creationTimestamp"`
	LegNumber           string      `json:"legNumber"`
	BestEx              bool        `json:"bestEx"`
	OrderType           string      `json:"orderType"`
	OrderStatus         string      `json:"orderStatus"`
	SubOrders           []Order     `json:"subOrders"`
	Side                string      `json:"side"`
	InstrumentID        string      `json:"instrumentId"`
	QuoteTicketID       string      `json:"quoteTicketId"`
	QuoteID             string      `json:"quoteID"`
	VenueID             string      `json:"venueID"`
	Quantity            AmountValue `json:"quantity"`
	LimitExtension      string      `json:"limitExtension"`
	TradingRestriction  string      `json:"tradingRestriction"`
	Limit               AmountValue `json:"limit"`
	TriggerLimit        AmountValue `json:"triggerLimit"`
	// TODO: AmountString
	TrailingLimitDistAbs string `json:"trailingLimitDistAbs"`
	// TODO: PercentageString
	TrailingLimitDistRel string      `json:"trailingLimitDistRel"`
	ValidityType         string      `json:"validityType"`
	Validity             string      `json:"validity"`
	OpenQuantity         AmountValue `json:"openQuantity"`
	CancelledQuantity    AmountValue `json:"cancelledQuantity"`
	ExecutedQuantity     AmountValue `json:"executedQuantity"`
	ExpectedValue        AmountValue `json:"expectedValue"`
	Executions           []Execution `json:"executions"`
}

type Execution struct {
	ExecutionID        string      `json:"executionID"`
	ExecutionNumber    int         `json:"executionNumber"`
	ExecutedQuantity   AmountValue `json:"executedQuantity"`
	ExecutionPrice     AmountValue `json:"executionPrice"`
	ExecutionTimestamp string      `json:"executionTimestamp"`
}

func (c *Client) Dimensions() ([]Dimension, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New(CLIENT_NOT_AUTHENTICATED)
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/brokerage/v3/orders/dimensions"),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	dimensions := &Dimensions{}
	_, err = c.http.exchange(req, dimensions)
	return dimensions.Values, err
}

func (c *Client) Orders(depotID string) ([]Order, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New(CLIENT_NOT_AUTHENTICATED)
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/brokerage/depots/%s/v3/orders", depotID)),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	orders := &Orders{}
	_, err = c.http.exchange(req, orders)
	return orders.Values, nil
}

func (c *Client) Order(orderID string) {

}

func (c *Client) CreateOrder(order OrderRequest, tan string) {
	// TODO
}

func (c *Client) UpdateOrder(orderID string, tan string) {
	// TODO
}

func (c *Client) DeleteOrder(orderID string, tan string) {
	// TODO
}

func (c *Client) PreValidateOrder() {
	// TODO
}

func (c *Client) ValidateOrder() {
	// TODO
}

func (c *Client) ExAnteOrder() {
	// TODO
}

func (c *Client) ValidateOrderUpdate(order OrderRequest) {
	// TODO
}

func (c *Client) ValidateOrderDeletion(orderID string) {
	// TODO
}
