package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// Config represents the application configuration
type Config struct {
	DatabasePath string `json:"database_path"`
	LogLevel     string `json:"log_level"`
	APITimeout   int    `json:"api_timeout"` // in seconds
}

// Load reads the configuration from a file and returns a Config struct
func Load(filename string) (*Config, error) {
	// Default configuration values
	config := &Config{
		DatabasePath: ":memory:", // SQLite in-memory database by default
		LogLevel:     "info",
		APITimeout:   30,
	}

	// Check if the configuration file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// If file doesn't exist, return the default configuration
		return config, nil
	}

	// Open the configuration file
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open configuration file")
	}
	defer file.Close() // Ensure file is closed even if an error occurs

	// Parse the JSON configuration file
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, errors.Wrap(err, "failed to parse configuration file")
	}

	// Validate the configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig ensures that the loaded configuration is valid
func validateConfig(config *Config) error {
	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
		"panic": true,
	}

	if _, valid := validLogLevels[config.LogLevel]; !valid {
		return errors.New("invalid log level: must be one of debug, info, warn, error, fatal, panic")
	}

	// Validate API timeout
	if config.APITimeout <= 0 {
		return errors.New("invalid API timeout: must be greater than 0")
	}

	return nil
}
