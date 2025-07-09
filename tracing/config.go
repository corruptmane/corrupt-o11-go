package tracing

import (
	"fmt"
	"os"
)

// ExportType represents the type of OpenTelemetry exporter
type ExportType string

const (
	ExportTypeStdout ExportType = "stdout"
	ExportTypeHTTP   ExportType = "http"
	ExportTypeGRPC   ExportType = "grpc"
)

// TracingConfig holds configuration for OpenTelemetry tracing
type TracingConfig struct {
	ExportType ExportType
	Endpoint   string
}

// FromEnv creates TracingConfig from environment variables
func FromEnv() (TracingConfig, error) {
	exportTypeStr := getEnvOrDefault("TRACING_EXPORTER_TYPE", "stdout")
	exportType, err := parseExportType(exportTypeStr)
	if err != nil {
		return TracingConfig{}, err
	}

	return TracingConfig{
		ExportType: exportType,
		Endpoint:   getEnvOrDefault("TRACING_EXPORTER_ENDPOINT", ""),
	}, nil
}

func parseExportType(exportType string) (ExportType, error) {
	switch exportType {
	case "stdout":
		return ExportTypeStdout, nil
	case "http":
		return ExportTypeHTTP, nil
	case "grpc":
		return ExportTypeGRPC, nil
	default:
		return "", fmt.Errorf("invalid export type: %s", exportType)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
