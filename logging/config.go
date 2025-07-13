// Package logging provides structured logging with OpenTelemetry tracing integration.
package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/corruptmane/corrupt-o11y-go/internal"
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
		AsJSON:  internal.MustEnvBool("LOG_AS_JSON", "false"),
		Tracing: internal.MustEnvBool("LOG_TRACING", "false"),
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

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
