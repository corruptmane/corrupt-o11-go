// Package metrics provides Prometheus metrics collection and management.
package metrics

import (
	"os"

	"github.com/corruptmane/corrupt-o11y-go/internal"
)

// MetricsConfig holds configuration for Prometheus metrics collection
type MetricsConfig struct {
	EnableGoCollector      bool   // Whether to enable Go runtime metrics
	EnableProcessCollector bool   // Whether to enable process metrics
	MetricPrefix           string // Optional prefix for all metrics
}

// FromEnv creates MetricsConfig from environment variables
func FromEnv() MetricsConfig {
	return MetricsConfig{
		EnableGoCollector:      internal.MustEnvBool("METRICS_ENABLE_GO", "true"),
		EnableProcessCollector: internal.MustEnvBool("METRICS_ENABLE_PROCESS", "true"),
		MetricPrefix:           os.Getenv("METRICS_PREFIX"),
	}
}
