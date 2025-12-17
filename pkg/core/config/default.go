package config

import "time"

// GetStringOrDefault returns the value associated with the key as a string.
// If the key is not found or the value is empty, it returns the provided default value.
func (c *Config) GetStringOrDefault(key string, defaultValue string) string {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetString(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// GetStringSliceOrDefault returns the value associated with the key as a slice of strings.
// If the key is not found or the value is nil/empty, it returns the provided default value.
func (c *Config) GetStringSliceOrDefault(key string, defaultValue []string) []string {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetStringSlice(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// GetStringMapOrDefault returns the value associated with the key as a map of interfaces.
// If the key is not found or the value is nil/empty, it returns the provided default value.
func (c *Config) GetStringMapOrDefault(key string, defaultValue map[string]interface{}) map[string]interface{} {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetStringMap(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// GetStringMapStringOrDefault returns the value associated with the key as a map of strings.
// If the key is not found or the value is nil/empty, it returns the provided default value.
func (c *Config) GetStringMapStringOrDefault(key string, defaultValue map[string]string) map[string]string {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetStringMapString(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// GetStringMapStringSliceOrDefault returns the value associated with the key as a map of string slices.
// If the key is not found or the value is nil/empty, it returns the provided default value.
func (c *Config) GetStringMapStringSliceOrDefault(key string, defaultValue map[string][]string) map[string][]string {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetStringMapStringSlice(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// GetBoolOrDefault returns the value associated with the key as a boolean.
// If the key is not found, it returns the provided default value.
func (c *Config) GetBoolOrDefault(key string, defaultValue bool) bool {
	if !c.IsSet(key) {
		return defaultValue
	}

	return c.GetBool(key)
}

// GetFloat64OrDefault returns the value associated with the key as a float64.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetFloat64OrDefault(key string, defaultValue float64) float64 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetFloat64(key)
	if value == 0.0 {
		return defaultValue
	}

	return value
}

// GetIntOrDefault returns the value associated with the key as an integer.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetIntOrDefault(key string, defaultValue int) int {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetInt(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetInt32OrDefault returns the value associated with the key as an int32.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetInt32OrDefault(key string, defaultValue int32) int32 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetInt32(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetInt64OrDefault returns the value associated with the key as an int64.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetInt64OrDefault(key string, defaultValue int64) int64 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetInt64(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetIntSliceOrDefault returns the value associated with the key as a slice of ints.
// If the key is not found or the value is nil/empty, it returns the provided default value.
func (c *Config) GetIntSliceOrDefault(key string, defaultValue []int) []int {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetIntSlice(key)
	if len(value) == 0 {
		return defaultValue
	}

	return value
}

// GetUintOrDefault returns the value associated with the key as an uint.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetUintOrDefault(key string, defaultValue uint) uint {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetUint(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetUint16OrDefault returns the value associated with the key as an uint16.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetUint16OrDefault(key string, defaultValue uint16) uint16 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetUint16(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetUint32OrDefault returns the value associated with the key as an uint32.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetUint32OrDefault(key string, defaultValue uint32) uint32 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetUint32(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetUint64OrDefault returns the value associated with the key as an uint64.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetUint64OrDefault(key string, defaultValue uint64) uint64 {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetUint64(key)
	if value == 0 {
		return defaultValue
	}

	return value
}

// GetTimeOrDefault returns the value associated with the key as time.Time.
// If the key is not found, it returns the provided default value.
func (c *Config) GetTimeOrDefault(key string, defaultValue time.Time) time.Time {
	if !c.IsSet(key) {
		return defaultValue
	}

	return c.GetTime(key)
}

// GetDurationOrDefault returns the value associated with the key as time.Duration.
// If the key is not found or the value is zero, it returns the provided default value.
func (c *Config) GetDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if !c.IsSet(key) {
		return defaultValue
	}

	value := c.GetDuration(key)
	if value == 0 {
		return defaultValue
	}

	return value
}
