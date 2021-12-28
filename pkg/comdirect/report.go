package comdirect

import (
	"context"
	"errors"
	"net/http"
)

type Report struct {
	ProductID            string        `json:"productId"`
	ProductType          string        `json:"productType"`
	TargetClientID       string        `json:"targetClientId"`
	ClientConnectionType string        `json:"clientConnectionType"`
	Balance              ReportBalance `json:"balance"`
}

type ReportBalance struct {
	Account                Account     `json:"account"`
	AccountId              string      `json:"accountId"`
	Balance                AmountValue `json:"balance"`
	BalanceEUR             AmountValue `json:"balanceEUR"`
	AvailableCashAmount    AmountValue `json:"availableCashAmount"`
	AvailableCashAmountEUR AmountValue `json:"availableCashAmountEUR"`
	Depot                  Depot       `json:"depot"`
	DepotID                string      `json:"depotId"`
	DateLastUpdate         string      `json:"dateLastUpdate"`
	PrevDayValue           AmountValue `json:"prevDayValue"`
}

type ReportAggregated struct {
	BalanceEUR             AmountValue `json:"balanceEUR"`
	AvailableCashAmountEUR AmountValue `json:"availableCashAmountEUR"`
}

type Reports struct {
	Paging           Paging           `json:"paging"`
	ReportAggregated ReportAggregated `json:"Aggregated"`
	Values           []Report         `json:"values"`
}

// Reports returns the balance for all available accounts.
func (c *Client) Reports(ctx context.Context, options ...Options) (*Reports, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/reports/participants/user/v1/allbalances"),
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}
	req = req.WithContext(ctx)

	reports := &Reports{}
	_, err = c.http.exchange(req, reports)
	return reports, err
}
