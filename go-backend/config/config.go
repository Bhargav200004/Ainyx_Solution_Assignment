package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	// Server
	ServerPort string

	// Database
	DatabaseURL string

	// Environment: "development" or "production"
	Env string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ainyx?sslmode=disable"),
		Env:         getEnv("APP_ENV", "development"),
	}
}

// DBConnString returns the formatted database connection string.
func (c *Config) DBConnString() string {
	return c.DatabaseURL
}

// ServerAddr returns the formatted server address.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf(":%s", c.ServerPort)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
