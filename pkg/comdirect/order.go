package comdirect

type Dimension struct {
	Venues []Venue
}

type Venue struct {
	Name          string
	VenueID       string
	Country       string
	Type          string
	Currencies    []string
	Sides         []string
	ValidityTypes []string
	OrderTypes    OrderTypes
}

type OrderTypes struct {
	Quote              OrderType
	Market             OrderType
	StopMarket         OrderType
	NextOrder          OrderType
	OneCancelsOther    OrderType
	Limit              OrderType
	TrailingStopMarket OrderType
}

type OrderType struct {
	LimitExtensions     []string
	TradingRestrictions []string
}

type OrderRequest struct {
	DepotID      string
	OrderID		 string
	Side         string
	InstrumentID string
	OrderType    string
	Quantity     AmountValue
	VenueID      string
	Limit        AmountValue
	ValidityType string
	Validity     string
}

func (c *Client) Dimensions() {

}

func (c *Client) Orders(depotID string) {

}

func (c *Client) Order(orderID string) {

}

func (c *Client) CreateOrder(order OrderRequest, tan string){

}

func (c *Client) UpdateOrder(orderID string, tan string) {

}

func (c *Client) DeleteOrder(orderID string, tan string) {

}

func (c *Client) PreValidateOrder() {

}

func (c *Client) ValidateOrder() {

}

func (c *Client) ExAnteOrder() {

}

func (c *Client) ValidateOrderUpdate(order OrderRequest){

}

func (c *Client) ValidateOrderDeletion(orderID string){

}