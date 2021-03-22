package comdirect

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
