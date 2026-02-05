package utils

import (
	"os"
	"path/filepath"
)

// GetDownloadFolder returns the default download folder for the current OS
func GetDownloadFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Cross-platform downloads folder
	return filepath.Join(homeDir, "Downloads"), nil
}

// EnsureDirExists creates a directory if it doesn't exist
func EnsureDirExists(dir string) error {
	return os.MkdirAll(dir, 0755)
}
