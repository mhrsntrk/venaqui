package realdebrid

import (
	"fmt"
	"net/http"
)

// ValidateToken checks if the API token is valid by making a test request
func (c *Client) ValidateToken() error {
	endpoint := fmt.Sprintf("%s/user", c.baseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return fmt.Errorf("invalid API token: unauthorized")
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("token validation failed: status %d", resp.StatusCode)
	}

	return nil
}
