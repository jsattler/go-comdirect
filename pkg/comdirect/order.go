package comdirect

import (
	"errors"
	"net/http"
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
	Values []Dimension
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

func (c *Client) Dimensions() ([]Dimension, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/brokerage/v3/orders/dimensions"),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	dimensions := &Dimensions{}
	_, err = c.http.exchange(req, dimensions)
	return dimensions.Values, err
}

func (c *Client) Orders(depotID string) {
	// TODO
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
