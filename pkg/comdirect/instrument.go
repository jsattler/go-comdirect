package comdirect

type Instrument struct {
	InstrumentID string     `json:"instrumentId"`
	WKN          string     `json:"wkn"`
	ISIN         string     `json:"isin"`
	Mnemonic     string     `json:"mnemonic"`
	Name         string     `json:"name"`
	ShortName    string     `json:"shortName"`
	StaticData   StaticData `json:"staticData"`
}

type StaticData struct {
	Notation               string `json:"notation"`
	Currency               string `json:"currency"`
	InstrumentType         string `json:"instrumentType"`
	PriipsRelevant         bool   `json:"priipsRelevant"`
	KidAvailable           bool   `json:"kidAvailable"`
	ShippingWaiverRequired bool   `json:"shippingWaiverRequired"`
	FundRedemptionLimited  bool   `json:"fundRedemptionLimited"`
}
