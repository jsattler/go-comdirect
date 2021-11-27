package comdirect

import (
	"errors"
	"fmt"
	"net/http"
)

type Instrument struct {
	InstrumentID string     `json:"instrumentId"`
	WKN          string     `json:"wkn"`
	ISIN         string     `json:"isin"`
	Mnemonic     string     `json:"mnemonic"`
	Name         string     `json:"name"`
	ShortName    string     `json:"shortName"`
	StaticData   StaticData `json:"staticData"`
}

type Instruments struct {
	Values []Instrument `json:"values"`
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

// Instrument retrieves instrument information by WKN, ISIN or mnemonic
func (c *Client) Instrument(instrument string) ([]Instrument, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/brokerage/v1/instruments/%s", instrument)),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}

	instruments := &Instruments{}
	_, err = c.http.exchange(req, instruments)
	return instruments.Values, err
}
