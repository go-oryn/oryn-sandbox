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

type Config struct {
	*viper.Viper
	env string
}

func NewConfig(opts ...Option) (*Config, error) {
	v := viper.New()

	// Configure defaults
	v.SetDefault("app.name", "oryn")
	v.SetDefault("app.version", "0.0.1")
	v.SetDefault("app.debug", false)

	// Configure environment variables support
	v.SetEnvPrefix("ORYN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	// Default options
	options := &Options{
		env: strings.ToLower(os.Getenv("ORYN_ENV")),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, fmt.Errorf("failed to apply config option: %w", err)
		}
	}

	// Load config files if embedFS is provided
	if options.embedFS != nil {
		// Load all YAML files from the root config directory
		if err := loadConfigFiles(v, *options.embedFS, "."); err != nil {
			return nil, fmt.Errorf("failed to load embed root config files: %w", err)
		}

		// Check if ORYN_ENV is set and load environment-specific config files
		if options.env != "" {
			if err := loadConfigFiles(v, *options.embedFS, filepath.Join(".", options.env)); err != nil {
				return nil, fmt.Errorf("failed to load embed %s config files: %w", options.env, err)
			}
		}
	}

	// Merge programmatically provided values if any
	if options.values != nil {
		if err := v.MergeConfigMap(options.values); err != nil {
			return nil, fmt.Errorf("failed to merge values: %w", err)
		}
	}

	return &Config{
		Viper: v,
		env:   options.env,
	}, nil
}

func (c *Config) Env() string {
	return c.env
}

func (c *Config) TestEnv() bool {
	return c.env == "test"
}

func loadConfigFiles(v *viper.Viper, configFS embed.FS, configFSDir string) error {
	entries, err := fs.ReadDir(configFS, configFSDir)
	if err != nil {
		// If the directory doesn't exist (e.g., prod folder), returns without error
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read config directory %s: %w", configFSDir, err)
	}

	// Compile a regexp to match environment variable placeholders
	envVarPattern, err := regexp.Compile(`\$\{([^}]+)\}`)
	if err != nil {
		return fmt.Errorf("failed to compile config env var pattern: %w", err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Detect the file format based on extension
		var configType string
		if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
			configType = "yaml"
		} else if strings.HasSuffix(fileName, ".json") {
			configType = "json"
		} else {
			// Skip unsupported file types
			continue
		}

		// Read the file content
		filePath := filepath.Join(configFSDir, fileName)
		data, err := fs.ReadFile(configFS, filePath)
		if err != nil {
			return fmt.Errorf("failed to read config file %s: %w", filePath, err)
		}

		// Resolve environment variable placeholders matching for example ${ENV_VAR}
		resolvedData := envVarPattern.ReplaceAllStringFunc(string(data), func(match string) string {
			// Extract the env var name from the placeholder
			envVar := envVarPattern.FindStringSubmatch(match)
			if len(envVar) < 2 {
				return ""
			}

			// Returns the environment variable value, if not set, returns empty string
			return os.Getenv(envVar[1])
		})

		// Set the config type
		v.SetConfigType(configType)

		// Merge the config
		if err := v.MergeConfig(strings.NewReader(resolvedData)); err != nil {
			return fmt.Errorf("failed to merge config from %s: %w", filePath, err)
		}
	}

	return nil
}
