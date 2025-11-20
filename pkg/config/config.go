package config

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

// Config wraps viper.Viper for configuration management
type Config struct {
	*viper.Viper
}

// NewConfig creates a new Config instance by loading configuration files from an embedded filesystem.
// It loads all YAML files from the root of the configs folder and merges them.
// If ORYN_ENV environment variable is set, it also loads and merges files from the configs/{ORYN_ENV} folder.
// For example: ORYN_ENV=prod loads from configs/prod, ORYN_ENV=test loads from configs/test.
// Environment variables with ORYN_ prefix can override any config value.
// Example: ORYN_LOGS_LEVEL=info will override logs.level
func NewConfig(configFS embed.FS) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Configure environment variable support
	// ORYN_LOGS_LEVEL will map to logs.level
	v.SetEnvPrefix("ORYN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Load all YAML files from root configs directory
	if err := loadConfigFiles(v, configFS, "configs", false); err != nil {
		return nil, fmt.Errorf("failed to load root configs files: %w", err)
	}

	// Check if ORYN_ENV is set and load environment-specific config files
	if env := os.Getenv("ORYN_ENV"); env != "" {
		envConfigPath := filepath.Join("configs", env)
		// Load and merge environment-specific config files
		if err := loadConfigFiles(v, configFS, envConfigPath, true); err != nil {
			return nil, fmt.Errorf("failed to load %s config files: %w", env, err)
		}
	}

	return &Config{Viper: v}, nil
}

// loadConfigFiles loads all YAML files from the specified directory in the embedded filesystem
func loadConfigFiles(v *viper.Viper, configFS embed.FS, dir string, merge bool) error {
	entries, err := fs.ReadDir(configFS, dir)
	if err != nil {
		// If the directory doesn't exist (e.g., prod folder), that's okay
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	firstFile := !merge
	for _, entry := range entries {
		// Skip directories and non-YAML files
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".yml") {
			continue
		}

		// Read the file content
		filePath := filepath.Join(dir, fileName)
		data, err := fs.ReadFile(configFS, filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Resolve environment variable placeholders
		resolvedData := resolvePlaceholders(string(data))

		if firstFile {
			// Use ReadConfig for the very first file
			if err := v.ReadConfig(strings.NewReader(resolvedData)); err != nil {
				return fmt.Errorf("failed to load config from %s: %w", filePath, err)
			}
			firstFile = false
		} else {
			// Use MergeConfig for all subsequent files
			if err := v.MergeConfig(strings.NewReader(resolvedData)); err != nil {
				return fmt.Errorf("failed to merge config from %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// envVarPattern matches ${ENV_VAR} placeholders
var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// resolvePlaceholders replaces ${ENV_VAR} placeholders with actual environment variable values
// If an environment variable is not set, the placeholder is replaced with an empty string
func resolvePlaceholders(data string) string {
	return envVarPattern.ReplaceAllStringFunc(data, func(match string) string {
		// Extract the env var name from ${ENV_VAR}
		envVar := envVarPattern.FindStringSubmatch(match)
		if len(envVar) < 2 {
			return ""
		}

		// Get the environment variable value
		// If not set, return empty string
		return os.Getenv(envVar[1])
	})
}
