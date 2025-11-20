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
	options
}

// options holds configuration options for Config
type options struct {
	embedFS     embed.FS
	embedFSPath string
	values      map[string]any
}

// NewConfig creates a new Config instance with optional configuration.
func NewConfig(opts ...Option) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Configure environment variable support
	v.SetEnvPrefix("ORYN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Initialize the config instance
	cfg := &Config{
		Viper: v,
		options: options{
			embedFSPath: "configs", // default
		},
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Load config files if embedFS is provided
	if cfg.embedFS != (embed.FS{}) {
		// Load all YAML files from the root config directory
		if err := loadConfigFiles(cfg.Viper, cfg.embedFS, cfg.embedFSPath); err != nil {
			return nil, fmt.Errorf("failed to load root config files: %w", err)
		}

		// Check if ORYN_ENV is set and load environment-specific config files
		if env := os.Getenv("ORYN_ENV"); env != "" {
			if err := loadConfigFiles(cfg.Viper, cfg.embedFS, filepath.Join(cfg.embedFSPath, env)); err != nil {
				return nil, fmt.Errorf("failed to load %s config files: %w", env, err)
			}
		}
	}

	// Merge programmatically provided values if any
	if cfg.values != nil {
		if err := cfg.Viper.MergeConfigMap(cfg.values); err != nil {
			return nil, fmt.Errorf("failed to merge provided values: %w", err)
		}
	}

	return cfg, nil
}

// loadConfigFiles loads all YAML files from the specified directory in the embedded filesystem
func loadConfigFiles(v *viper.Viper, configFS embed.FS, dir string) error {
	entries, err := fs.ReadDir(configFS, dir)
	if err != nil {
		// If the directory doesn't exist (e.g., prod folder), that's okay
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

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

		// Resolve environment variable placeholders matching for example ${ENV_VAR}
		envVarPattern, err := regexp.Compile(`\$\{([^}]+)\}`)
		if err != nil {
			return fmt.Errorf("failed to compile env var pattern: %w", err)
		}

		resolvedData := envVarPattern.ReplaceAllStringFunc(string(data), func(match string) string {
			// Extract the env var name from the placeholder
			envVar := envVarPattern.FindStringSubmatch(match)
			if len(envVar) < 2 {
				return ""
			}

			// Returns the environment variable value, if not set, return empty string
			return os.Getenv(envVar[1])
		})

		// Merge config from file
		if err := v.MergeConfig(strings.NewReader(resolvedData)); err != nil {
			return fmt.Errorf("failed to merge config from %s: %w", filePath, err)
		}
	}

	return nil
}
