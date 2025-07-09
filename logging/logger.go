package logging

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// ConfigureLogging configures structured logging with the given configuration
func ConfigureLogging(config LoggingConfig) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: true,
	}

	if config.AsJSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	if config.Tracing {
		handler = &tracingHandler{Handler: handler}
	}

	slog.SetDefault(slog.New(handler))
}

// GetLogger returns a logger with the given name
func GetLogger(name string) *slog.Logger {
	return slog.With("logger", name)
}

// tracingHandler wraps slog.Handler to add OpenTelemetry tracing information
type tracingHandler struct {
	slog.Handler
}

func (h *tracingHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		spanContext := span.SpanContext()
		record.Add("span_id", spanContext.SpanID().String())
		record.Add("trace_id", spanContext.TraceID().String())
	}
	return h.Handler.Handle(ctx, record)
}
