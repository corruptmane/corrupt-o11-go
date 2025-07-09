# corrupt-o11y-go

[![CI](https://github.com/corruptmane/corrupt-o11y-go/actions/workflows/ci.yml/badge.svg)](https://github.com/corruptmane/corrupt-o11y-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/corruptmane/corrupt-o11y-go.svg)](https://pkg.go.dev/github.com/corruptmane/corrupt-o11y-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/corruptmane/corrupt-o11y-go)](https://goreportcard.com/report/github.com/corruptmane/corrupt-o11y-go)

A comprehensive observability library for Go applications providing structured logging, metrics collection, distributed tracing, and operational endpoints.

## Features

- **Structured Logging**: Built on Go's standard `log/slog` with OpenTelemetry integration
- **Metrics Collection**: Prometheus metrics with service info patterns and built-in collectors
- **Distributed Tracing**: OpenTelemetry tracing with stdout, HTTP, and gRPC exporters
- **Operational Endpoints**: Health checks, readiness probes, metrics exposure, and service info
- **Service Metadata**: Centralized service information for consistent labeling
- **Environment Configuration**: All modules support environment-based configuration

## Quick Start

```go
package main

import (
    "context"
    "log/slog"

    "github.com/corruptmane/corrupt-o11y-go/logging"
    "github.com/corruptmane/corrupt-o11y-go/metadata"
    "github.com/corruptmane/corrupt-o11y-go/metrics"
    "github.com/corruptmane/corrupt-o11y-go/operational"
    "github.com/corruptmane/corrupt-o11y-go/tracing"
)

func main() {
    // Initialize service metadata
    serviceInfo := metadata.FromEnv()

    // Configure logging
    logConfig := logging.FromEnv()
    logging.ConfigureLogging(logConfig)
    logger := logging.GetLogger("main")

    // Setup metrics
    metricsCollector := metrics.NewMetricsCollector()
    metricsCollector.CreateServiceInfoMetricFromServiceInfo(serviceInfo)

    // Setup tracing
    tracingConfig, _ := tracing.FromEnv()
    tracing.ConfigureTracing(context.Background(), tracingConfig, serviceInfo.Name, serviceInfo.Version)

    // Setup operational server
    opConfig := operational.FromEnv()
    status := operational.NewStatus()
    status.SetReady(true)

    opServer := operational.NewOperationalServer(opConfig, serviceInfo, status, metricsCollector)
    opServer.Start(context.Background())

    logger.Info("Service started", slog.String("url", opServer.ServerURL()))

    // Your application logic here...
}
```

## Configuration

All modules support environment-based configuration:

### Service Metadata
- `SERVICE_NAME` - Service name (default: "unknown-dev")
- `SERVICE_VERSION` - Service version (default: "dev")
- `INSTANCE_ID` - Instance identifier (default: "unknown-dev")
- `COMMIT_SHA` - Git commit SHA (default: "unknown-dev")
- `BUILD_TIME` - Build timestamp (default: "unknown-dev")

### Logging
- `LOG_LEVEL` - Log level: DEBUG, INFO, WARN, ERROR (default: "INFO")
- `LOG_AS_JSON` - Output as JSON: true/false (default: "false")
- `LOG_TRACING` - Include tracing info: true/false (default: "false")

**Note**: To include tracing information in logs, you must use the context-aware logging methods (`InfoContext`, `ErrorContext`, etc.) and pass the span context:

```go
tracer := tracing.GetTracer("my-service")
spanCtx, span := tracer.Start(context.Background(), "operation")
defer span.End()

// This will include span information in the log
logger.InfoContext(spanCtx, "Operation completed")

// This will NOT include span information
logger.Info("Operation completed")
```

### Tracing
- `TRACING_EXPORTER_TYPE` - Exporter type: stdout, http, grpc (default: "stdout")
- `TRACING_EXPORTER_ENDPOINT` - Exporter endpoint URL (required for http/grpc)

### Operational Server
- `OPERATIONAL_HOST` - Bind host (default: "0.0.0.0")
- `OPERATIONAL_PORT` - Bind port (default: 42069)

## Operational Endpoints

The operational server provides:

- `GET /health` - Health check (200 if alive, 503 if not)
- `GET /ready` - Readiness check (200 if ready, 503 if not)
- `GET /metrics` - Prometheus metrics
- `GET /info` - Service information as JSON

## Installation

```bash
go get github.com/corruptmane/corrupt-o11y-go
```

## Development

### Requirements
- Go 1.23+
- [just](https://github.com/casey/just) (task runner)
- [golangci-lint](https://golangci-lint.run/) (linting)

### Commands

```bash
# Run tests
just test

# Run tests with race detector
just test-race

# Run tests with coverage
just test-coverage

# Run linter
just lint

# Format code
just fmt

# Run all checks
just check
```

## License

MIT
