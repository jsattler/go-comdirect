package comdirect

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TODO: Think about where to put the stuff in here.
// TODO: Currently this is a pool for everything that does not fit somewhere else

type HTTPClient struct {
	http.Client
}

// api.comdirect.de/api/{path}
func apiURL(path string) *url.URL {
	return comdirectURL(ApiPath + path)
}

// api.comdirect.de/{path}
func comdirectURL(path string) *url.URL {
	return &url.URL{Host: Host, Scheme: HttpsScheme, Path: path}
}

func generateSessionID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%032d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

func generateRequestID() string {
	unix := time.Now().Unix()
	id := fmt.Sprintf("%09d", unix)
	return id[0:9]
}

func (h *HTTPClient) exchange(request *http.Request, target interface{}) (*http.Response, error) {
	res, err := h.Do(request)
	if err != nil {
		return res, err
	}

	if err = json.NewDecoder(res.Body).Decode(target); err != nil {
		return res, err
	}

	return res, res.Body.Close()
}

func requestInfoJSON(sessionID string) ([]byte, error) {
	info := &requestInfo{ClientRequestID: clientRequestID{
		SessionID: sessionID,
		RequestID: generateRequestID(),
	}}
	return json.Marshal(info)
}
