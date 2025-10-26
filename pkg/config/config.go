package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig reads all YAML files from the specified config directory,
// merges them into a single Viper instance, and enables environment variable overrides.
//
// Environment variables can override config values using the following format:
// - Nested keys are separated by underscores
// - All uppercase
// - Optional prefix can be set
//
// Example:
//
//	Config: http.port = 8080
//	Env var: HTTP_PORT=9090 will override it
//
// Parameters:
//   - configDir: path to the directory containing config files
//   - envPrefix: optional prefix for environment variables (e.g., "APP")
//
// Returns:
//   - *viper.Viper: configured viper instance
//   - error: if any error occurs during loading
func LoadConfig(configDir string, envPrefix string) (*viper.Viper, error) {
	v := viper.New()

	// Configure environment variable support
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}

	// Replace dots and hyphens with underscores for env vars
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Enable automatic env var binding
	v.AutomaticEnv()

	// Check if config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory does not exist: %s", configDir)
	}

	// Read all YAML files from the config directory
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	// Also check for .yml extension
	ymlFiles, err := filepath.Glob(filepath.Join(configDir, "*.yml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}
	files = append(files, ymlFiles...)

	if len(files) == 0 {
		return nil, fmt.Errorf("no YAML config files found in %s", configDir)
	}

	// Load and merge all config files
	for _, file := range files {
		// Create a temporary viper instance for each file
		tempV := viper.New()
		tempV.SetConfigFile(file)

		if err := tempV.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", file, err)
		}

		// Merge the settings from this file into the main viper instance
		if err := v.MergeConfigMap(tempV.AllSettings()); err != nil {
			return nil, fmt.Errorf("failed to merge config from %s: %w", file, err)
		}

		fmt.Printf("Loaded config file: %s\n", filepath.Base(file))
	}

	return v, nil
}

// MustLoadConfig is a wrapper around LoadConfig that panics if an error occurs.
// Useful for initialization where you want to fail fast.
func MustLoadConfig(configDir string, envPrefix string) *viper.Viper {
	v, err := LoadConfig(configDir, envPrefix)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return v
}
