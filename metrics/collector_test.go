package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/corruptmane/corrupt-o11y-go/metadata"
)

func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	if collector == nil {
		t.Error("Expected NewMetricsCollector to return non-nil collector")
	}

	if collector.Registry() == nil {
		t.Error("Expected collector to have non-nil registry")
	}
}

func TestRegisterUnregister(t *testing.T) {
	collector := NewMetricsCollector()

	// Create a test metric
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	})

	// Register metric
	err := collector.Register("test_counter", counter)
	if err != nil {
		t.Errorf("Expected no error registering metric, got %v", err)
	}

	// Try to register the same metric again (should fail)
	err = collector.Register("test_counter", counter)
	if err == nil {
		t.Error("Expected error when registering duplicate metric")
	}

	// Unregister metric
	if !collector.Unregister("test_counter") {
		t.Error("Expected Unregister to return true for existing metric")
	}

	// Try to unregister non-existent metric
	if collector.Unregister("non_existent") {
		t.Error("Expected Unregister to return false for non-existent metric")
	}
}

func TestClear(t *testing.T) {
	collector := NewMetricsCollector()

	// Register a test metric
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "A test counter",
	})
	collector.Register("test_counter", counter)

	// Clear all metrics
	collector.Clear()

	// Try to unregister the metric (should return false as it's already cleared)
	if collector.Unregister("test_counter") {
		t.Error("Expected metric to be cleared")
	}
}

func TestCreateServiceInfoMetricFromServiceInfo(t *testing.T) {
	collector := NewMetricsCollector()

	serviceInfo := metadata.ServiceInfo{
		Name:       "test-service",
		Version:    "1.0.0",
		InstanceID: "test-instance",
		CommitSHA:  "abc123",
		BuildTime:  "2023-01-01T00:00:00Z",
	}

	metric := collector.CreateServiceInfoMetricFromServiceInfo(serviceInfo)

	if metric == nil {
		t.Error("Expected CreateServiceInfoMetricFromServiceInfo to return non-nil metric")
	}
}

func TestCreateServiceInfoMetric(t *testing.T) {
	metric := CreateServiceInfoMetric("test-service", "1.0.0", "test-instance", nil, nil)

	if metric == nil {
		t.Error("Expected CreateServiceInfoMetric to return non-nil metric")
	}

	// Test with all fields
	commitSHA := "abc123"
	buildTime := "2023-01-01T00:00:00Z"
	metric2 := CreateServiceInfoMetric("test-service", "1.0.0", "test-instance", &commitSHA, &buildTime)

	if metric2 == nil {
		t.Error("Expected CreateServiceInfoMetric with all fields to return non-nil metric")
	}
}
