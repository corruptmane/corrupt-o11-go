package tracing

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Test with default values
	config, err := FromEnv()
	if err != nil {
		t.Errorf("Expected no error with default values, got %v", err)
	}

	if config.ExportType != ExportTypeStdout {
		t.Errorf("Expected ExportType to be 'stdout', got %s", config.ExportType)
	}
	if config.Endpoint != "" {
		t.Errorf("Expected Endpoint to be empty, got %s", config.Endpoint)
	}
}

func TestFromEnvWithValues(t *testing.T) {
	// Set environment variables
	os.Setenv("TRACING_EXPORTER_TYPE", "http")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "http://localhost:4318")
	defer func() {
		os.Unsetenv("TRACING_EXPORTER_TYPE")
		os.Unsetenv("TRACING_EXPORTER_ENDPOINT")
	}()

	config, err := FromEnv()
	if err != nil {
		t.Errorf("Expected no error with valid values, got %v", err)
	}

	if config.ExportType != ExportTypeHTTP {
		t.Errorf("Expected ExportType to be 'http', got %s", config.ExportType)
	}
	if config.Endpoint != "http://localhost:4318" {
		t.Errorf("Expected Endpoint to be 'http://localhost:4318', got %s", config.Endpoint)
	}
}

func TestFromEnvWithInvalidExportType(t *testing.T) {
	// Set invalid export type
	os.Setenv("TRACING_EXPORTER_TYPE", "invalid")
	defer func() {
		os.Unsetenv("TRACING_EXPORTER_TYPE")
	}()

	_, err := FromEnv()
	if err == nil {
		t.Error("Expected error with invalid export type")
	}
}

func TestParseExportType(t *testing.T) {
	tests := []struct {
		input     string
		expected  ExportType
		shouldErr bool
	}{
		{"stdout", ExportTypeStdout, false},
		{"http", ExportTypeHTTP, false},
		{"grpc", ExportTypeGRPC, false},
		{"invalid", "", true},
	}

	for _, test := range tests {
		result, err := parseExportType(test.input)
		if test.shouldErr {
			if err == nil {
				t.Errorf("parseExportType(%s) should return error", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseExportType(%s) should not return error, got %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("parseExportType(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	}
}
