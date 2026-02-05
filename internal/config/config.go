package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	RealDebridAPIToken string
	Aria2RPCUrl        string
	Aria2Secret        string
	DefaultDownloadDir string
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".venaqui")
	configFile := filepath.Join(configPath, "config.yaml")

	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Set defaults
	viper.SetDefault("aria2.rpc_url", "http://localhost:6800/jsonrpc")
	viper.SetDefault("aria2.secret", "")
	viper.SetDefault("download.default_dir", "")

	// Read config file (ignore error if file doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is okay, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Get default download directory
	defaultDir, err := GetDefaultDownloadDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get default download directory: %w", err)
	}

	// Override with config if specified
	if viper.GetString("download.default_dir") != "" {
		defaultDir = viper.GetString("download.default_dir")
	}

	// Get Real-Debrid API token (required)
	apiToken := viper.GetString("realdebrid.api_token")
	if apiToken == "" {
		return nil, fmt.Errorf("realdebrid.api_token is required in %s", configFile)
	}

	cfg := &Config{
		RealDebridAPIToken: apiToken,
		Aria2RPCUrl:        viper.GetString("aria2.rpc_url"),
		Aria2Secret:        viper.GetString("aria2.secret"),
		DefaultDownloadDir: defaultDir,
	}

	return cfg, nil
}
