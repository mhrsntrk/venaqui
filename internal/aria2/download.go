package aria2

import (
	"fmt"

	"github.com/siku2/arigo"
)

// AddDownload adds a new download to aria2
func (c *Client) AddDownload(url, downloadDir string) (string, error) {
	options := &arigo.Options{
		Dir:                   downloadDir,
		MaxConnectionPerServer: 16,
		Split:                  16,
		MinSplitSize:          1048576, // 1M in bytes
	}

	gid, err := c.rpc.AddURI([]string{url}, options)
	if err != nil {
		return "", fmt.Errorf("failed to add download: %w", err)
	}

	return gid.GID, nil
}

// RemoveDownload removes a download from aria2
func (c *Client) RemoveDownload(gid string) error {
	return c.rpc.Remove(gid)
}

// PauseDownload pauses a download
func (c *Client) PauseDownload(gid string) error {
	return c.rpc.Pause(gid)
}

// ResumeDownload resumes a paused download
func (c *Client) ResumeDownload(gid string) error {
	return c.rpc.Unpause(gid)
}

// GetActiveDownloads returns all active downloads
func (c *Client) GetActiveDownloads() ([]string, error) {
	statuses, err := c.rpc.TellActive("gid")
	if err != nil {
		return nil, fmt.Errorf("failed to get active downloads: %w", err)
	}

	result := make([]string, len(statuses))
	for i, status := range statuses {
		result[i] = status.GID
	}

	return result, nil
}
