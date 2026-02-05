package config

import (
	"os"
	"path/filepath"
)

// GetDefaultDownloadDir returns the default download directory for the current OS
func GetDefaultDownloadDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Cross-platform downloads folder
	// Windows: C:\Users\<user>\Downloads
	// macOS/Linux: /Users/<user>/Downloads or /home/<user>/Downloads
	return filepath.Join(homeDir, "Downloads"), nil
}
