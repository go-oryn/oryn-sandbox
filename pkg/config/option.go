package config

import "embed"

// Option is a functional option for configuring Config
type Option func(*Config) error

// WithEmbedFS sets the embedded filesystem for the config files.
// Example: WithEmbedFS(configFS) will load files from the provided FS.
func WithEmbedFS(fs embed.FS) Option {
	return func(c *Config) error {
		c.embedFS = fs

		return nil
	}
}

// WithValues sets configuration values programmatically.
// These values are merged after loading config files but can be overridden by environment variables.
// Example: WithValues(map[string]any{"database.host": "localhost", "database.port": 5432})
func WithValues(values map[string]any) Option {
	return func(c *Config) error {
		c.values = values

		return nil
	}
}
