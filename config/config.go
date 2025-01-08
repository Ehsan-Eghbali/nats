package config

import (
	"os"
	"strings"
)

// AppConfig holds application-wide configuration values.
// In a real project, you might extend this struct with more fields.
type AppConfig struct {
	NATSURL    string
	StreamName string
}

// LoadConfig loads the configuration from environment variables (or default values).
func LoadConfig() AppConfig {
	cfg := AppConfig{
		NATSURL:    getEnv("NATS_URL", "nats://localhost:4222"),
		StreamName: getEnv("NATS_STREAM", "MY_STREAM"),
	}
	return cfg
}

// getEnv is a helper function to retrieve env variables or default.
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && strings.TrimSpace(val) != "" {
		return val
	}
	return defaultVal
}
