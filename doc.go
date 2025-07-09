// Package corrupt_o11y provides observability utilities for Go applications.
//
// This library provides modular observability components including:
//   - Structured logging with OpenTelemetry integration
//   - Prometheus metrics collection
//   - OpenTelemetry distributed tracing
//   - Operational HTTP endpoints (health, metrics, info)
//   - Service metadata management
//
// Example usage:
//
//	import (
//		"github.com/corruptmane/corrupt-o11y-go/logging"
//		"github.com/corruptmane/corrupt-o11y-go/metadata"
//		"github.com/corruptmane/corrupt-o11y-go/metrics"
//		"github.com/corruptmane/corrupt-o11y-go/operational"
//		"github.com/corruptmane/corrupt-o11y-go/tracing"
//	)
//
//	// Configure observability
//	serviceInfo := metadata.FromEnv()
//
//	// Setup logging
//	logConfig := logging.FromEnv()
//	logging.ConfigureLogging(logConfig)
//
//	// Setup metrics
//	metricsCollector := metrics.NewMetricsCollector()
//	metricsCollector.CreateServiceInfoMetricFromServiceInfo(serviceInfo)
//
//	// Setup operational server
//	opConfig := operational.FromEnv()
//	status := operational.NewStatus()
//	opServer := operational.NewOperationalServer(opConfig, serviceInfo, status, metricsCollector)
package corrupt_o11y
