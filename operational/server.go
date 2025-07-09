package operational

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/corruptmane/corrupt-o11y-go/logging"
	"github.com/corruptmane/corrupt-o11y-go/metadata"
	"github.com/corruptmane/corrupt-o11y-go/metrics"
)

// OperationalServer provides HTTP endpoints for operational concerns
type OperationalServer struct {
	config      OperationalServerConfig
	status      *Status
	metrics     *metrics.MetricsCollector
	serviceInfo metadata.ServiceInfo
	server      *http.Server
	serverURL   string
}

// NewOperationalServer creates a new operational server
func NewOperationalServer(
	config OperationalServerConfig,
	serviceInfo metadata.ServiceInfo,
	status *Status,
	metricsCollector *metrics.MetricsCollector,
) *OperationalServer {
	return &OperationalServer{
		config:      config,
		status:      status,
		metrics:     metricsCollector,
		serviceInfo: serviceInfo,
	}
}

// Start starts the operational server
func (s *OperationalServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	mux.HandleFunc("/info", s.handleInfo)
	mux.Handle("/metrics", promhttp.HandlerFor(s.metrics.Registry(), promhttp.HandlerOpts{}))

	// Create listener to get actual port
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Set server URL using actual assigned port
	host := s.config.Host
	if host == "0.0.0.0" {
		host = ""
	}
	actualPort := listener.Addr().(*net.TCPAddr).Port
	s.serverURL = fmt.Sprintf("http://%s:%d", host, actualPort)

	go func() {
		if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger := logging.GetLogger("operational")
			logger.Error("operational server error", slog.String("error", err.Error()))
		}
	}()

	return nil
}

// Stop gracefully stops the operational server
func (s *OperationalServer) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// ServerURL returns the server URL with actual assigned port
func (s *OperationalServer) ServerURL() string {
	return s.serverURL
}

func (s *OperationalServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if s.status.IsAlive() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (s *OperationalServer) handleReady(w http.ResponseWriter, r *http.Request) {
	if s.status.IsReady() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (s *OperationalServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(s.serviceInfo.AsMap())
}
