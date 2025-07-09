package operational

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Test with default values
	config := FromEnv()

	if config.Host != "0.0.0.0" {
		t.Errorf("Expected Host to be '0.0.0.0', got %s", config.Host)
	}
	if config.Port != 42069 {
		t.Errorf("Expected Port to be 42069, got %d", config.Port)
	}
}

func TestFromEnvWithValues(t *testing.T) {
	// Set environment variables
	os.Setenv("OPERATIONAL_HOST", "127.0.0.1")
	os.Setenv("OPERATIONAL_PORT", "8080")
	defer func() {
		os.Unsetenv("OPERATIONAL_HOST")
		os.Unsetenv("OPERATIONAL_PORT")
	}()

	config := FromEnv()

	if config.Host != "127.0.0.1" {
		t.Errorf("Expected Host to be '127.0.0.1', got %s", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
}

func TestFromEnvWithInvalidPort(t *testing.T) {
	// Set invalid port
	os.Setenv("OPERATIONAL_PORT", "invalid")
	defer func() {
		os.Unsetenv("OPERATIONAL_PORT")
	}()

	config := FromEnv()

	// Should fall back to default
	if config.Port != 42069 {
		t.Errorf("Expected Port to be 42069 (default) for invalid port, got %d", config.Port)
	}
}
