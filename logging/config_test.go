package logging

import (
	"log/slog"
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Test with default values
	config := FromEnv()

	if config.Level != slog.LevelInfo {
		t.Errorf("Expected Level to be INFO, got %v", config.Level)
	}
	if config.AsJSON != false {
		t.Errorf("Expected AsJSON to be false, got %v", config.AsJSON)
	}
	if config.Tracing != false {
		t.Errorf("Expected Tracing to be false, got %v", config.Tracing)
	}
}

func TestFromEnvWithValues(t *testing.T) {
	// Set environment variables
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_AS_JSON", "true")
	os.Setenv("LOG_TRACING", "true")
	defer func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_AS_JSON")
		os.Unsetenv("LOG_TRACING")
	}()

	config := FromEnv()

	if config.Level != slog.LevelDebug {
		t.Errorf("Expected Level to be DEBUG, got %v", config.Level)
	}
	if config.AsJSON != true {
		t.Errorf("Expected AsJSON to be true, got %v", config.AsJSON)
	}
	if config.Tracing != true {
		t.Errorf("Expected Tracing to be true, got %v", config.Tracing)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"info", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"WARNING", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"error", slog.LevelError},
		{"INVALID", slog.LevelInfo}, // should default to INFO
	}

	for _, test := range tests {
		result := parseLogLevel(test.input)
		if result != test.expected {
			t.Errorf("parseLogLevel(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}
