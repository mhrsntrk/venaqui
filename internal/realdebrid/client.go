package realdebrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseURL = "https://api.real-debrid.com/rest/1.0"

// Client handles communication with the Real-Debrid API
type Client struct {
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Real-Debrid API client
func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		baseURL:  BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithBaseURL creates a new Real-Debrid API client with a custom base URL (for testing)
func NewClientWithBaseURL(apiToken, baseURL string) *Client {
	return &Client{
		apiToken: apiToken,
		baseURL:  baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UnrestrictedLink represents the response from the unrestrict endpoint
type UnrestrictedLink struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
	Filesize int64  `json:"filesize"`
	Link     string `json:"link"` // Direct download link
	Host     string `json:"host"`
	Chunks   int    `json:"chunks"`
	Download string `json:"download"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error"`
	ErrorCode int  `json:"error_code,omitempty"`
}

// UnrestrictLink converts a hoster link to an unrestricted download link
func (c *Client) UnrestrictLink(link string) (*UnrestrictedLink, error) {
	endpoint := fmt.Sprintf("%s/unrestrict/link", c.baseURL)

	payload := map[string]string{
		"link": link,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			// Handle specific error cases
			switch resp.StatusCode {
			case 401:
				return nil, fmt.Errorf("unauthorized: invalid API token")
			case 403:
				return nil, fmt.Errorf("forbidden: %s", errorResp.Error)
			case 429:
				return nil, fmt.Errorf("rate limited: too many requests, please wait")
			case 503:
				return nil, fmt.Errorf("service unavailable: Real-Debrid is temporarily unavailable")
			default:
				return nil, fmt.Errorf("RD API error (%d): %s", resp.StatusCode, errorResp.Error)
			}
		}
		return nil, fmt.Errorf("RD API error: %d - %s", resp.StatusCode, string(body))
	}

	var result UnrestrictedLink
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
