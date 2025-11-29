package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
)

//go:embed configs
var configFS embed.FS

//go:embed configs-json
var configJSONFS embed.FS

func main() {
	// Load configuration from embedded filesystem
	// Set ORYN_ENV=prod to load production overrides
	cfg, err := config.NewConfig(
		config.WithEmbedFS(configJSONFS, "configs-json"),
		config.WithValues(map[string]interface{}{"foo.bar.baz": 123}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration loaded successfully!\n")

	// Access config values from different files
	fmt.Println("=== HTTP Configuration ===")
	fmt.Printf("Server Address: %s\n", cfg.GetString("http.server.address"))

	fmt.Println("\n=== Logs Configuration ===")
	fmt.Printf("Level: %s\n", cfg.GetString("logs.level"))
	fmt.Printf("Format: %s\n", cfg.GetString("logs.format"))

	fmt.Println("\n=== Database Configuration ===")
	fmt.Printf("Host: %s\n", cfg.GetString("database.host"))
	fmt.Printf("Port: %d\n", cfg.GetInt("database.port"))
	fmt.Printf("Name: %s\n", cfg.GetString("database.name"))
	fmt.Printf("User: %s\n", cfg.GetString("database.user"))

	fmt.Println("\n=== Foo Configuration ===")
	fmt.Printf("Foo: %d\n", cfg.GetInt("foo.bar.baz"))

	fmt.Println("\n=== App Configuration (with placeholders) ===")
	fmt.Printf("Name: %s\n", cfg.GetString("app.name"))
	authors := cfg.GetStringSlice("app.authors")
	fmt.Printf("Authors: %v\n", authors)
	fmt.Printf("Version: %s\n", cfg.GetString("app.version"))
	fmt.Printf("Port: %d\n", cfg.GetInt("app.port"))
	fmt.Printf("Debug: %t\n", cfg.GetBool("app.debug"))
	fmt.Printf("Max Connections: %d\n", cfg.GetInt("app.maxConnections"))

	// Demonstrate environment variable override
	fmt.Println("\n=== Environment Variable Override ===")
	fmt.Println("Examples:")
	fmt.Println("  ORYN_ENV=prod go run main.go")
	fmt.Println("  ORYN_LOGS_LEVEL=error go run main.go")
	fmt.Println("  ORYN_DATABASE_HOST=custom-host go run main.go")
	fmt.Println("  ORYN_HTTP_SERVER_ADDRESS=:9000 go run main.go")
}
