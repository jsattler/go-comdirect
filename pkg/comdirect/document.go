package comdirect

import (
	"errors"
	"fmt"
	"net/http"
)

type Document struct {
	DocumentID       string           `json:"document_id"`
	Name             string           `json:"name"`
	DateCreation     string           `json:"dateCreation"`
	MimeType         string           `json:"mimeType"`
	Deletable        bool             `json:"deletable"`
	Advertisement    bool             `json:"advertisement"`
	DocumentMetaData DocumentMetaData `json:"documentMetaData"`
}

type Documents struct {
	Values []Document `json:"values"`
}

type DocumentMetaData struct {
	Archived          bool `json:"archived"`
	AlreadyRead       bool `json:"alreadyRead"`
	PreDocumentExists bool `json:"predocumentExists"`
}

func (c *Client) Documents() ([]Document, error) {
	if c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL("/messages/clients/user/v2/documents"),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	documents := &Documents{}
	_, err = c.http.exchange(req, documents)
	return documents.Values, err
}

func (c *Client) Document(documentID string) (*Document, error) {
	if c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/messages/clients/user/v2/documents/%s", documentID)),
		Header: http.Header{
			AcceptHeaderKey:          {"application/json"},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}

	document := &Document{}
	_, err = c.http.exchange(req, document)
	return document, err
}
