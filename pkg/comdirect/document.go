package comdirect

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Document struct {
	DocumentID       string           `json:"documentId"`
	Name             string           `json:"name"`
	DateCreation     string           `json:"dateCreation"`
	MimeType         string           `json:"mimeType"`
	Deletable        bool             `json:"deletable"`
	Advertisement    bool             `json:"advertisement"`
	DocumentMetaData DocumentMetaData `json:"documentMetaData"`
}

type Documents struct {
	Paging Paging     `json:"paging"`
	Values []Document `json:"values"`
}

type DocumentMetaData struct {
	Archived          bool `json:"archived"`
	AlreadyRead       bool `json:"alreadyRead"`
	PreDocumentExists bool `json:"predocumentExists"`
}

func (c *Client) Documents(ctx context.Context, options ...Options) (*Documents, error) {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return nil, errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return nil, err
	}
	url := apiURL("/messages/clients/user/v2/documents")

	encodeOptions(url, options)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: defaultHeaders(c.authentication.accessToken.AccessToken, string(info)),
	}
	req = req.WithContext(ctx)
	documents := &Documents{}
	_, err = c.http.exchange(req, documents)
	return documents, err
}

func (c *Client) DownloadDocument(ctx context.Context, document *Document, folder string) error {
	if c.authentication == nil || c.authentication.accessToken.AccessToken == "" || c.authentication.IsExpired() {
		return errors.New("authentication is expired or not initialized")
	}
	info, err := requestInfoJSON(c.authentication.sessionID)
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    apiURL(fmt.Sprintf("/messages/v2/documents/%s", document.DocumentID)),
		Header: http.Header{
			AcceptHeaderKey:          {document.MimeType},
			ContentTypeHeaderKey:     {"application/json"},
			AuthorizationHeaderKey:   {BearerPrefix + c.authentication.accessToken.AccessToken},
			HttpRequestInfoHeaderKey: {string(info)},
		},
	}
	req = req.WithContext(ctx)
	res, err := c.http.Do(req)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	if folder == "" {
		folder, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	ext := strings.Split(document.MimeType, "/")
	fileName := strings.ReplaceAll(document.Name, " ", "_")
	file, err := os.Create(folder + "/" + document.DateCreation + "-" + fileName + "." + ext[1])
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = io.Copy(file, res.Body)

	return err
}
