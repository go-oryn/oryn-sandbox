# Configuration Guide

This project uses Viper to manage configuration from multiple YAML files with environment variable overrides.

## Directory Structure

```
config/
├── database.yaml  # Database configuration
├── http.yaml      # HTTP server configuration
├── logs.yaml      # Logging configuration
└── README.md      # This file
```

## Usage

### Basic Usage

```go
import "github.com/go-oryn/oryn-sandbox/pkg/config"

// Load all config files from ./config directory
cfg, err := config.LoadConfig("./config", "APP")
if err != nil {
    log.Fatal(err)
}

// Access values using dot notation
port := cfg.GetInt("http.port")
logLevel := cfg.GetString("logs.level")
```

### Alternative: Panic on Error

```go
// Use MustLoadConfig to panic if config loading fails
cfg := config.MustLoadConfig("./config", "APP")
```

## Environment Variable Overrides

All configuration values can be overridden using environment variables with the following rules:

1. **Prefix**: Use the prefix specified in `LoadConfig()` (e.g., `APP_`)
2. **Separator**: Replace dots (`.`) with underscores (`_`)
3. **Case**: Use UPPERCASE

### Examples

| Config Key | Config Value | Env Var | Override Value |
|------------|--------------|---------|----------------|
| `http.port` | `8080` | `APP_HTTP_PORT=9090` | `9090` |
| `logs.level` | `info` | `APP_LOGS_LEVEL=debug` | `debug` |
| `database.host` | `localhost` | `APP_DATABASE_HOST=prod-db` | `prod-db` |
| `http.tls.enabled` | `false` | `APP_HTTP_TLS_ENABLED=true` | `true` |

### Running with Overrides

```bash
# Override single value
APP_HTTP_PORT=9090 go run main.go

# Override multiple values
APP_HTTP_PORT=9090 APP_LOGS_LEVEL=debug go run main.go

# In production with export
export APP_DATABASE_HOST=prod-db.example.com
export APP_DATABASE_PASSWORD=secure_password
go run main.go
```

## Adding New Config Files

Simply create a new `.yaml` or `.yml` file in the config directory:

```yaml
# config/cache.yaml
cache:
  provider: redis
  host: localhost
  port: 6379
  ttl: 3600
```

The new file will automatically be loaded and merged with existing configuration.

## Config Access Methods

Viper provides type-safe accessors:

```go
cfg.GetString("key")           // string
cfg.GetInt("key")              // int
cfg.GetBool("key")             // bool
cfg.GetFloat64("key")          // float64
cfg.GetDuration("key")         // time.Duration
cfg.GetStringSlice("key")      // []string
cfg.GetStringMap("key")        // map[string]interface{}
```

## Best Practices

1. **Default Values**: Always provide sensible defaults in YAML files
2. **Sensitive Data**: Never commit secrets; use env vars for passwords/tokens
3. **Documentation**: Document each config option in the YAML file
4. **Validation**: Add validation logic after loading config
5. **Environment-Specific**: Use env vars for environment-specific values (dev/staging/prod)
