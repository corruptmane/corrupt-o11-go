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

func TestCreateCounterVec(t *testing.T) {
	config := MetricsConfig{MetricPrefix: "myapp_"}
	collector := NewMetricsCollectorWithConfig(config)

	counter := collector.CreateCounterVec(
		prometheus.CounterOpts{
			Name: "requests",
			Help: "Total requests",
		},
		[]string{"method", "status"},
	)

	if counter == nil {
		t.Error("Expected counter to be created")
	}

	// Test that metric is registered with original name (internal tracking)
	if !collector.IsRegistered("requests") {
		t.Error("Expected counter to be registered with original name")
	}

	// Test counter functionality
	counter.WithLabelValues("GET", "200").Inc()
	counter.WithLabelValues("POST", "404").Add(5)
}

func TestCreateGaugeFunc(t *testing.T) {
	config := MetricsConfig{MetricPrefix: "test_"}
	collector := NewMetricsCollectorWithConfig(config)

	value := 42.0
	gauge := collector.CreateGaugeFunc(
		prometheus.GaugeOpts{
			Name: "temperature",
			Help: "Current temperature",
		},
		func() float64 { return value },
	)

	if gauge == nil {
		t.Error("Expected gauge to be created")
	}

	// Test that metric is registered with original name (internal tracking)
	if !collector.IsRegistered("temperature") {
		t.Error("Expected gauge to be registered with original name")
	}
}

func TestMultipleCollectorsWithSameMetric(t *testing.T) {
	// Test the shared source pattern with Func metrics
	var sharedCounter int64 = 0

	collector1 := NewMetricsCollector()
	collector2 := NewMetricsCollector()

	// Create identical metrics reading from the same source
	counter1 := collector1.CreateCounterFunc(
		prometheus.CounterOpts{Name: "shared_requests", Help: "Shared requests"},
		func() float64 { return float64(sharedCounter) },
	)

	counter2 := collector2.CreateCounterFunc(
		prometheus.CounterOpts{Name: "shared_requests", Help: "Shared requests"},
		func() float64 { return float64(sharedCounter) },
	)

	if counter1 == nil || counter2 == nil {
		t.Error("Expected both counters to be created successfully")
	}

	// Both should be registered in their respective collectors
	if !collector1.IsRegistered("shared_requests") {
		t.Error("Expected counter to be registered in collector1")
	}
	if !collector2.IsRegistered("shared_requests") {
		t.Error("Expected counter to be registered in collector2")
	}

	// Modify shared source
	sharedCounter = 100
	// Both metrics will now report 100 when scraped
}
