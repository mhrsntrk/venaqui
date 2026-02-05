package realdebrid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TorrentInfo represents torrent information from Real-Debrid
type TorrentInfo struct {
	ID       string   `json:"id"`
	Filename string   `json:"filename"`
	Status   string   `json:"status"` // waiting_files_selection, queued, downloading, downloaded, error
	Files    []File   `json:"files"`
	Links    []string `json:"links"` // Download links after torrent is ready
	Progress float64  `json:"progress"`
}

// File represents a file in a torrent
type File struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Bytes    int64  `json:"bytes"`
	Selected int    `json:"selected"` // 0 = not selected, 1 = selected
}

// AddTorrentResponse represents the response from adding a torrent
type AddTorrentResponse struct {
	ID       string `json:"id"`
	URI      string `json:"uri"`
	Filename string `json:"filename"`
}

// AddTorrent adds a torrent file URL to Real-Debrid
func (c *Client) AddTorrent(torrentURL string) (*AddTorrentResponse, error) {
	endpoint := fmt.Sprintf("%s/torrents/addTorrent", c.baseURL)

	// For torrent URLs, we need to download the torrent file first
	// Real-Debrid API expects the torrent file content, not the URL
	resp, err := http.Get(torrentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download torrent file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download torrent file: status %d", resp.StatusCode)
	}

	torrentData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read torrent file: %w", err)
	}

	// Real-Debrid API expects the torrent file content directly in the PUT request body
	req, err := http.NewRequest("PUT", endpoint, bytes.NewReader(torrentData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/x-bittorrent")

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 201 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("RD API error (%d): %s", resp.StatusCode, errorResp.Error)
		}
		return nil, fmt.Errorf("RD API error: %d - %s", resp.StatusCode, string(body))
	}

	var result AddTorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AddMagnet adds a magnet link to Real-Debrid
func (c *Client) AddMagnet(magnetLink string) (*AddTorrentResponse, error) {
	endpoint := fmt.Sprintf("%s/torrents/addMagnet", c.baseURL)

	payload := map[string]string{
		"magnet": magnetLink,
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

	if resp.StatusCode != 201 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("RD API error (%d): %s", resp.StatusCode, errorResp.Error)
		}
		return nil, fmt.Errorf("RD API error: %d - %s", resp.StatusCode, string(body))
	}

	var result AddTorrentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// SelectFiles selects files in a torrent
func (c *Client) SelectFiles(torrentID string, fileIDs []int) error {
	endpoint := fmt.Sprintf("%s/torrents/selectFiles/%s", c.baseURL, torrentID)

	// Real-Debrid API expects form data with files as "all" or comma-separated IDs
	var filesParam string
	if len(fileIDs) == 0 {
		filesParam = "all"
	} else {
		// Convert to comma-separated string
		fileIDStrings := make([]string, len(fileIDs))
		for i, id := range fileIDs {
			fileIDStrings[i] = fmt.Sprintf("%d", id)
		}
		filesParam = strings.Join(fileIDStrings, ",")
	}

	// Use form data instead of JSON
	formData := fmt.Sprintf("files=%s", filesParam)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// According to API docs: returns 204 HTTP code, or 202 if action already done
	if resp.StatusCode != 204 && resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("RD API error (%d): %s", resp.StatusCode, errorResp.Error)
		}
		return fmt.Errorf("RD API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetTorrentInfo gets information about a torrent
func (c *Client) GetTorrentInfo(torrentID string) (*TorrentInfo, error) {
	endpoint := fmt.Sprintf("%s/torrents/info/%s", c.baseURL, torrentID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

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
			return nil, fmt.Errorf("RD API error (%d): %s", resp.StatusCode, errorResp.Error)
		}
		return nil, fmt.Errorf("RD API error: %d - %s", resp.StatusCode, string(body))
	}

	var result TorrentInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// WaitForTorrentReady waits for a torrent to be ready (downloaded status)
func (c *Client) WaitForTorrentReady(torrentID string, maxWait time.Duration) (*TorrentInfo, error) {
	deadline := time.Now().Add(maxWait)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		info, err := c.GetTorrentInfo(torrentID)
		if err != nil {
			return nil, err
		}

		switch info.Status {
		case "downloaded":
			return info, nil
		case "error":
			return nil, fmt.Errorf("torrent failed: %s", info.Filename)
		case "waiting_files_selection":
			// Need to select files first
			fileIDs := []int{}
			for _, file := range info.Files {
				fileIDs = append(fileIDs, file.ID)
			}
			if len(fileIDs) > 0 {
				if err := c.SelectFiles(torrentID, fileIDs); err != nil {
					return nil, fmt.Errorf("failed to select files: %w", err)
				}
			}
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for torrent to be ready")
		}

		<-ticker.C
	}
}
