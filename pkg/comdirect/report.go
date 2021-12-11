package comdirect

import (
	"errors"
	"net/http"
)

type Report struct {
	ProductID            string         `json:"productId"`
	ProductType          string         `json:"productType"`
	TargetClientID       string         `json:"targetClientId"`
	ClientConnectionType string         `json:"clientConnectionType"`
	Balance              AccountBalance `json:"balance"`
}

type Reports struct {
	Values []Report `json:"values"`
}

// Reports returns the balance for all available accounts.
func (c *Client) Reports() ([]Report, error) {
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

	reports := &Reports{}
	_, err = c.http.exchange(req, reports)
	return reports.Values, err
}
