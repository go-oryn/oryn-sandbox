package config

import "embed"

// Option is a functional option for configuring Config
type Option func(*Config) error

// WithEmbedFS sets the embedded filesystem and the base path for config files.
// The path parameter specifies the folder name in the embedded filesystem (e.g., "configs", "config").
// Example: WithEmbedFS(configFS, "configs") will load files from the "configs" folder in the embedded FS.
func WithEmbedFS(fs embed.FS, path string) Option {
	return func(c *Config) error {
		c.embedFS = fs
		c.embedFSPath = path

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
