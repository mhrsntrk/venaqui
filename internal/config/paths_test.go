package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultDownloadDir(t *testing.T) {
	dir, err := GetDefaultDownloadDir()
	if err != nil {
		t.Fatalf("GetDefaultDownloadDir() error = %v", err)
	}

	if dir == "" {
		t.Error("GetDefaultDownloadDir() returned empty string")
	}

	// Should be absolute path
	if !filepath.IsAbs(dir) {
		t.Errorf("GetDefaultDownloadDir() = %v, expected absolute path", dir)
	}

	// Should end with Downloads
	base := filepath.Base(dir)
	if base != "Downloads" {
		t.Errorf("GetDefaultDownloadDir() = %v, expected path ending with Downloads", dir)
	}

	// Should be in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	expectedDir := filepath.Join(homeDir, "Downloads")
	if dir != expectedDir {
		t.Errorf("GetDefaultDownloadDir() = %v, want %v", dir, expectedDir)
	}
}
