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
	spanContext := span.SpanContext()

	// Add tracing information if we have a valid span context
	if spanContext.IsValid() {
		spanInfo := map[string]any{
			"span_id":  spanContext.SpanID().String(),
			"trace_id": spanContext.TraceID().String(),
		}
		record.Add("span", spanInfo)
	}

	return h.Handler.Handle(ctx, record)
}

func (h *tracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &tracingHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *tracingHandler) WithGroup(name string) slog.Handler {
	return &tracingHandler{Handler: h.Handler.WithGroup(name)}
}
