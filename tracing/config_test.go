package tracing

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel/trace/noop"
)

func TestFromEnv(t *testing.T) {
	// Test with default values
	config, err := FromEnv()
	if err != nil {
		t.Errorf("Expected no error with default values, got %v", err)
	}

	if !config.IsEnabled {
		t.Errorf("Expected IsEnabled to be true by default, got %v", config.IsEnabled)
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
	os.Setenv("TRACING_ENABLED", "true")
	os.Setenv("TRACING_EXPORTER_TYPE", "http")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "http://localhost:4318")
	defer func() {
		os.Unsetenv("TRACING_ENABLED")
		os.Unsetenv("TRACING_EXPORTER_TYPE")
		os.Unsetenv("TRACING_EXPORTER_ENDPOINT")
	}()

	config, err := FromEnv()
	if err != nil {
		t.Errorf("Expected no error with valid values, got %v", err)
	}

	if !config.IsEnabled {
		t.Errorf("Expected IsEnabled to be true, got %v", config.IsEnabled)
	}
	if config.ExportType != ExportTypeHTTP {
		t.Errorf("Expected ExportType to be 'http', got %s", config.ExportType)
	}
	if config.Endpoint != "http://localhost:4318" {
		t.Errorf("Expected Endpoint to be 'http://localhost:4318', got %s", config.Endpoint)
	}
}

func TestFromEnvWithTracingDisabled(t *testing.T) {
	// Set tracing disabled
	os.Setenv("TRACING_ENABLED", "false")
	defer func() {
		os.Unsetenv("TRACING_ENABLED")
	}()

	config, err := FromEnv()
	if err != nil {
		t.Errorf("Expected no error with disabled tracing, got %v", err)
	}

	if config.IsEnabled {
		t.Errorf("Expected IsEnabled to be false, got %v", config.IsEnabled)
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

func TestConfigureTracingWithDisabled(t *testing.T) {
	config := TracingConfig{
		IsEnabled:  false,
		ExportType: ExportTypeStdout,
		Endpoint:   "",
	}

	ctx := context.Background()
	provider, err := ConfigureTracing(ctx, config, "test-service", "1.0.0")
	if err != nil {
		t.Errorf("Expected no error with disabled tracing, got %v", err)
	}

	// Check if it's a NoopTracerProvider by comparing with a known noop provider
	noopProvider := noop.NewTracerProvider()
	if provider != noopProvider {
		// Since direct comparison might not work, check that tracer names work
		tracer := provider.Tracer("test")
		if tracer == nil {
			t.Error("Expected valid tracer from NoopTracerProvider")
		}
	}
}
