package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// openDirectory opens a directory in the file manager
// If filePath is provided and exists, it will reveal the file in Finder (macOS) or open its directory
func openDirectory(dirPath string) error {
	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dirPath)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS - use 'open' to open directory in Finder
		cmd = exec.Command("open", dirPath)
	case "linux":
		cmd = exec.Command("xdg-open", dirPath)
	case "windows":
		cmd = exec.Command("explorer", dirPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// openFile opens a file with its default application
func openFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS - use 'open' to open file with default app
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		// On Windows, use 'start' command to open file with default application
		// The empty string after start is for the window title
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// revealFileInDirectory reveals a file in Finder/Explorer by opening its directory
// On macOS, this will highlight the file in Finder
func revealFileInDirectory(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	var cmd *exec.Cmd
	dir := filepath.Dir(filePath)

	switch runtime.GOOS {
	case "darwin": // macOS - use 'open -R' to reveal file in Finder
		cmd = exec.Command("open", "-R", filePath)
	case "linux":
		// On Linux, open the directory - some file managers support selecting the file
		cmd = exec.Command("xdg-open", dir)
	case "windows":
		// On Windows, use explorer with /select to highlight the file
		cmd = exec.Command("explorer", "/select,", filePath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}
