package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// APIResponse represents a generic API response
type APIResponse struct {
	StatusCode int         `json:"-"`
	Body       interface{} `json:"body,omitempty"`
	Error      error       `json:"error,omitempty"`
}

// FetchApi makes an HTTP request to the specified URL with the given method and data
func FetchApi(url string, method string, headers map[string]string, data interface{}) APIResponse {
	var payloadBytes []byte
	if data != nil {
		var err error
		payloadBytes, err = json.Marshal(data)
		if err != nil {
			return APIResponse{StatusCode: http.StatusInternalServerError, Error: err}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return APIResponse{StatusCode: http.StatusInternalServerError, Error: err}
	}
	req.Header.Set("Content-Type", "application/json")
	// Add headers to request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return APIResponse{StatusCode: http.StatusInternalServerError, Error: err}
	}
	defer resp.Body.Close()

	// Deserialize response body to JSON
	var body interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return APIResponse{StatusCode: resp.StatusCode, Error: err}
	}

	// Create APIResponse object
	apiResp := APIResponse{StatusCode: resp.StatusCode, Body: body}

	// Return APIResponse object
	return apiResp
}
