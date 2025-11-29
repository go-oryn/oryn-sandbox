package config

import (
	"embed"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata
var testFS embed.FS

func TestWithEmbedFS(t *testing.T) {
	tests := []struct {
		name string
		fs   embed.FS
		path string
	}{
		{
			name: "set embedded filesystem with path",
			fs:   testFS,
			path: "configs",
		},
		{
			name: "set embedded filesystem with different path",
			fs:   testFS,
			path: "config",
		},
		{
			name: "set embedded filesystem with nested path",
			fs:   testFS,
			path: "configs/prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a config to apply the option to
			cfg := &Config{
				Viper: viper.New(),
			}

			// Apply the option
			opt := WithEmbedFS(tt.fs, tt.path)
			err := opt(cfg)

			// Verify no error
			require.NoError(t, err)

			// Verify the embedFS was set
			assert.NotEqual(t, embed.FS{}, cfg.embedFS, "embedFS should be set")

			// Verify the path was set correctly
			assert.Equal(t, tt.path, cfg.embedFSPath)
		})
	}
}

func TestWithValues(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]any
	}{
		{
			name: "set simple string value",
			values: map[string]any{
				"database.host": "localhost",
			},
		},
		{
			name: "set multiple values with different types",
			values: map[string]any{
				"database.host": "localhost",
				"database.port": 5432,
				"app.debug":     true,
			},
		},
		{
			name: "set nested map values",
			values: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
			},
		},
		{
			name: "set slice values",
			values: map[string]any{
				"servers": []string{"server1", "server2"},
			},
		},
		{
			name:   "set empty map",
			values: map[string]any{},
		},
		{
			name:   "set nil map",
			values: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a config to apply the option to
			cfg := &Config{
				Viper: viper.New(),
			}

			// Apply the option
			opt := WithValues(tt.values)
			err := opt(cfg)

			// Verify no error
			require.NoError(t, err)

			// Verify the values were set correctly
			if tt.values == nil {
				assert.Nil(t, cfg.values)
			} else {
				assert.Len(t, cfg.values, len(tt.values))

				// Check each value
				for key, expectedValue := range tt.values {
					assert.Contains(t, cfg.values, key)

					// For simple types, we can compare directly
					switch expectedValue.(type) {
					case string, int, bool:
						assert.Equal(t, expectedValue, cfg.values[key])
					}
				}
			}
		})
	}
}

func TestMultipleOptions(t *testing.T) {
	t.Run("apply both WithEmbedFS and WithValues", func(t *testing.T) {
		cfg := &Config{
			Viper: viper.New(),
		}

		// Apply multiple options
		opts := []Option{
			WithEmbedFS(testFS, "configs"),
			WithValues(map[string]any{
				"test.key": "test.value",
			}),
		}

		for _, opt := range opts {
			err := opt(cfg)
			require.NoError(t, err)
		}

		// Verify both options were applied
		assert.NotEqual(t, embed.FS{}, cfg.embedFS, "embedFS should be set")
		assert.Equal(t, "configs", cfg.embedFSPath)
		assert.NotNil(t, cfg.values)
		assert.Equal(t, "test.value", cfg.values["test.key"])
	})
}

func TestOptionOrder(t *testing.T) {
	t.Run("last option wins for WithValues", func(t *testing.T) {
		cfg := &Config{
			Viper: viper.New(),
		}

		// Apply WithValues twice with different values
		opt1 := WithValues(map[string]any{"key": "value1"})
		opt2 := WithValues(map[string]any{"key": "value2"})

		require.NoError(t, opt1(cfg))
		require.NoError(t, opt2(cfg))

		// The second option should overwrite the first
		assert.Equal(t, "value2", cfg.values["key"])
	})

	t.Run("last option wins for WithEmbedFS path", func(t *testing.T) {
		cfg := &Config{
			Viper: viper.New(),
		}

		// Apply WithEmbedFS twice with different paths
		opt1 := WithEmbedFS(testFS, "path1")
		opt2 := WithEmbedFS(testFS, "path2")

		require.NoError(t, opt1(cfg))
		require.NoError(t, opt2(cfg))

		// The second option should overwrite the first
		assert.Equal(t, "path2", cfg.embedFSPath)
	})
}
