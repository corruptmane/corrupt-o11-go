package operational

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/corruptmane/corrupt-o11y-go/metadata"
	"github.com/corruptmane/corrupt-o11y-go/metrics"
)

func TestNewOperationalServer(t *testing.T) {
	config := OperationalServerConfig{
		Host: "127.0.0.1",
		Port: 0, // Use random port
	}

	serviceInfo := metadata.ServiceInfo{
		Name:       "test-service",
		Version:    "1.0.0",
		InstanceID: "test-instance",
		CommitSHA:  "abc123",
		BuildTime:  "2023-01-01T00:00:00Z",
	}

	status := NewStatus()
	metricsCollector := metrics.NewMetricsCollector()

	server := NewOperationalServer(config, serviceInfo, status, metricsCollector)

	if server == nil {
		t.Error("Expected NewOperationalServer to return non-nil server")
	}
}

func TestOperationalServerEndpoints(t *testing.T) {
	config := OperationalServerConfig{
		Host: "127.0.0.1",
		Port: 0, // Use random port
	}

	serviceInfo := metadata.ServiceInfo{
		Name:       "test-service",
		Version:    "1.0.0",
		InstanceID: "test-instance",
		CommitSHA:  "abc123",
		BuildTime:  "2023-01-01T00:00:00Z",
	}

	status := NewStatus()
	metricsCollector := metrics.NewMetricsCollector()

	server := NewOperationalServer(config, serviceInfo, status, metricsCollector)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop(ctx)

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	baseURL := server.ServerURL()

	// Test health endpoint (should be 200 - alive by default)
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Errorf("Failed to get health endpoint: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected health endpoint to return 200, got %d", resp.StatusCode)
		}
	}

	// Test ready endpoint (should be 503 - not ready by default)
	resp, err = http.Get(baseURL + "/ready")
	if err != nil {
		t.Errorf("Failed to get ready endpoint: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("Expected ready endpoint to return 503, got %d", resp.StatusCode)
		}
	}

	// Set ready and test again
	status.SetReady(true)
	resp, err = http.Get(baseURL + "/ready")
	if err != nil {
		t.Errorf("Failed to get ready endpoint after setting ready: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected ready endpoint to return 200 after setting ready, got %d", resp.StatusCode)
		}
	}

	// Test info endpoint
	resp, err = http.Get(baseURL + "/info")
	if err != nil {
		t.Errorf("Failed to get info endpoint: %v", err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected info endpoint to return 200, got %d", resp.StatusCode)
		}

		var info map[string]string
		err = json.NewDecoder(resp.Body).Decode(&info)
		if err != nil {
			t.Errorf("Failed to decode info response: %v", err)
		}

		if info["service_name"] != "test-service" {
			t.Errorf("Expected service_name to be 'test-service', got %s", info["service_name"])
		}
	}

	// Test metrics endpoint
	resp, err = http.Get(baseURL + "/metrics")
	if err != nil {
		t.Errorf("Failed to get metrics endpoint: %v", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected metrics endpoint to return 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/plain") {
			t.Errorf("Expected metrics endpoint to return text/plain content type, got %s", contentType)
		}
	}
}

func TestOperationalServerStop(t *testing.T) {
	config := OperationalServerConfig{
		Host: "127.0.0.1",
		Port: 0,
	}

	serviceInfo := metadata.ServiceInfo{
		Name:       "test-service",
		Version:    "1.0.0",
		InstanceID: "test-instance",
		CommitSHA:  "abc123",
		BuildTime:  "2023-01-01T00:00:00Z",
	}

	status := NewStatus()
	metricsCollector := metrics.NewMetricsCollector()

	server := NewOperationalServer(config, serviceInfo, status, metricsCollector)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start server
	err := server.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Stop server
	err = server.Stop(ctx)
	if err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}

	// Try to stop again (should not error)
	err = server.Stop(ctx)
	if err != nil {
		t.Errorf("Failed to stop server second time: %v", err)
	}
}
