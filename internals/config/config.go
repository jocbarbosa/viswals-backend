package config

import (
	"os"
)

// Config holds all the configuration values needed across the application
type Config struct {
	EncryptionKey string
}

// NewConfig creates a Config that reads from environment variables
func NewConfig() Config {
	return Config{
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
