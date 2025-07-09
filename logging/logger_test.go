package logging

import (
	"log/slog"
	"testing"
)

func TestConfigureLogging(t *testing.T) {
	// Test with default configuration
	config := LoggingConfig{
		Level:   slog.LevelInfo,
		AsJSON:  false,
		Tracing: false,
	}

	// This should not panic
	ConfigureLogging(config)

	// Test with JSON configuration
	config.AsJSON = true
	ConfigureLogging(config)

	// Test with tracing configuration
	config.Tracing = true
	ConfigureLogging(config)
}

func TestGetLogger(t *testing.T) {
	// Configure logging first
	config := LoggingConfig{
		Level:   slog.LevelInfo,
		AsJSON:  false,
		Tracing: false,
	}
	ConfigureLogging(config)

	logger := GetLogger("test-logger")

	if logger == nil {
		t.Error("Expected GetLogger to return non-nil logger")
	}

	// Test that we can create multiple loggers
	logger2 := GetLogger("test-logger-2")
	if logger2 == nil {
		t.Error("Expected GetLogger to return non-nil logger for second call")
	}
}
