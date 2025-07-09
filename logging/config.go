// Package logging provides structured logging with OpenTelemetry tracing integration.
package logging

import (
	"log/slog"
	"os"
	"strings"
)

// LoggingConfig holds configuration for structured logging
type LoggingConfig struct {
	Level   slog.Level
	AsJSON  bool
	Tracing bool
}

// FromEnv creates LoggingConfig from environment variables
func FromEnv() LoggingConfig {
	return LoggingConfig{
		Level:   parseLogLevel(getEnvOrDefault("LOG_LEVEL", "INFO")),
		AsJSON:  parseBool(getEnvOrDefault("LOG_AS_JSON", "false")),
		Tracing: parseBool(getEnvOrDefault("LOG_TRACING", "false")),
	}
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseBool(value string) bool {
	return strings.ToLower(value) == "true" || strings.ToLower(value) == "t"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
