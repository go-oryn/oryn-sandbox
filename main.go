package main

import (
	"fmt"
	"os"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
)

func main() {
	// Load configuration from the config directory
	// Environment variables with APP_ prefix can override any config value
	// Example: APP_HTTP_PORT=9090 will override http.port
	cfg, err := config.LoadConfig("./config", "APP")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Alternative: Use MustLoadConfig to panic on error
	// cfg := config.MustLoadConfig("./config", "APP")

	fmt.Println("Configuration loaded successfully!\n")

	// Access config values from different files
	fmt.Println("=== HTTP Configuration ===")
	fmt.Printf("Host: %s\n", cfg.GetString("http.host"))
	fmt.Printf("Port: %d\n", cfg.GetInt("http.port"))
	fmt.Printf("Timeout: %s\n", cfg.GetString("http.timeout"))
	fmt.Printf("TLS Enabled: %t\n", cfg.GetBool("http.tls.enabled"))

	fmt.Println("\n=== Logs Configuration ===")
	fmt.Printf("Level: %s\n", cfg.GetString("logs.level"))
	fmt.Printf("Format: %s\n", cfg.GetString("logs.format"))
	fmt.Printf("Output: %s\n", cfg.GetString("logs.output"))
	fmt.Printf("File Path: %s\n", cfg.GetString("logs.file.path"))

	fmt.Println("\n=== Database Configuration ===")
	fmt.Printf("Host: %s\n", cfg.GetString("database.host"))
	fmt.Printf("Port: %d\n", cfg.GetInt("database.port"))
	fmt.Printf("Name: %s\n", cfg.GetString("database.name"))
	fmt.Printf("User: %s\n", cfg.GetString("database.user"))
	fmt.Printf("Max Connections: %d\n", cfg.GetInt("database.max_connections"))

	// Demonstrate environment variable override
	fmt.Println("\n=== Environment Variable Override Example ===")
	fmt.Println("Try running with env vars to override:")
	fmt.Println("  APP_HTTP_PORT=9090 go run main.go")
	fmt.Println("  APP_DATABASE_HOST=prod-db.example.com go run main.go")
	fmt.Println("  APP_LOGS_LEVEL=debug go run main.go")
}
