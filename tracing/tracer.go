package tracing

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ConfigureTracing configures OpenTelemetry tracing
func ConfigureTracing(ctx context.Context, config TracingConfig, serviceName, serviceVersion string) (*trace.TracerProvider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var exporter trace.SpanExporter

	switch config.ExportType {
	case ExportTypeStdout:
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
	case ExportTypeHTTP:
		if config.Endpoint == "" {
			return nil, errors.New("HTTP exporter requires an endpoint")
		}
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(config.Endpoint),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP exporter: %w", err)
		}
	case ExportTypeGRPC:
		if config.Endpoint == "" {
			return nil, errors.New("GRPC exporter requires an endpoint")
		}
		conn, err := grpc.NewClient(config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
		}
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithGRPCConn(conn),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create GRPC exporter: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported export type: %s", config.ExportType)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider, nil
}

// GetTracer returns an OpenTelemetry tracer
func GetTracer(name string) oteltrace.Tracer {
	return otel.Tracer(name)
}
