package httpstatus

import "net/http"

func Is2xx(response *http.Response) bool {
	return response != nil && response.StatusCode >= 200 && response.StatusCode < 300
}

func Is3xx(response *http.Response) bool {
	return response != nil && response.StatusCode >= 300 && response.StatusCode < 400
}

func Is4xx(response *http.Response) bool {
	return response != nil && response.StatusCode >= 400 && response.StatusCode < 500
}

func Is5xx(response *http.Response) bool {
	return response != nil && response.StatusCode >= 500 && response.StatusCode < 600
}
