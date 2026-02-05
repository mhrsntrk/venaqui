package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ValidateURL checks if a string is a valid URL
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http:// or https://)")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// ValidatePath checks if a path is valid and writable
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check if it's an absolute path
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("path must be absolute: %s", cleanPath)
	}

	// Check if the path itself exists and is a directory
	info, err := os.Stat(cleanPath)
	if err == nil {
		// Path exists, check if it's a directory and writable
		if !info.IsDir() {
			return fmt.Errorf("path exists but is not a directory: %s", cleanPath)
		}
		if info.Mode().Perm()&0200 == 0 {
			return fmt.Errorf("directory exists but is not writable: %s", cleanPath)
		}
		return nil
	}

	// Path doesn't exist, check parent directory
	if !os.IsNotExist(err) {
		return fmt.Errorf("cannot access path: %w", err)
	}

	// Check if parent directory exists and is writable
	parentDir := filepath.Dir(cleanPath)
	parentInfo, err := os.Stat(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("parent directory does not exist: %s", parentDir)
		}
		return fmt.Errorf("cannot access parent directory: %w", err)
	}

	if !parentInfo.IsDir() {
		return fmt.Errorf("parent path is not a directory: %s", parentDir)
	}

	// Check if parent directory is writable
	if parentInfo.Mode().Perm()&0200 == 0 {
		return fmt.Errorf("parent directory is not writable: %s", parentDir)
	}

	return nil
}

// NormalizePath normalizes a path string
func NormalizePath(path string) string {
	// Expand user home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = strings.Replace(path, "~", homeDir, 1)
		}
	}

	// Clean the path
	return filepath.Clean(path)
}
