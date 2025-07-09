// Package operational provides HTTP endpoints for health checks, metrics, and service info.
package operational

import (
	"os"
	"strconv"
)

// OperationalServerConfig holds configuration for operational HTTP server
type OperationalServerConfig struct {
	Host string
	Port int
}

// FromEnv creates OperationalServerConfig from environment variables
func FromEnv() OperationalServerConfig {
	port := 42069
	if portStr := os.Getenv("OPERATIONAL_PORT"); portStr != "" {
		if parsedPort, err := strconv.Atoi(portStr); err == nil {
			port = parsedPort
		}
	}

	return OperationalServerConfig{
		Host: getEnvOrDefault("OPERATIONAL_HOST", "0.0.0.0"),
		Port: port,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
