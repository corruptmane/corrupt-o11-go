package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/corruptmane/corrupt-o11y-go/metadata"
)

// MetricsCollector provides a centralized registry for Prometheus metrics
type MetricsCollector struct {
	registry *prometheus.Registry
	metrics  map[string]prometheus.Collector
	mu       sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector with built-in metrics
func NewMetricsCollector() *MetricsCollector {
	registry := prometheus.NewRegistry()

	// Register built-in collectors
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &MetricsCollector{
		registry: registry,
		metrics:  make(map[string]prometheus.Collector),
	}
}

// Register registers a metric collector
func (mc *MetricsCollector) Register(name string, collector prometheus.Collector) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if err := mc.registry.Register(collector); err != nil {
		return err
	}

	mc.metrics[name] = collector
	return nil
}

// Unregister unregisters a metric collector
func (mc *MetricsCollector) Unregister(name string) bool {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if collector, exists := mc.metrics[name]; exists {
		mc.registry.Unregister(collector)
		delete(mc.metrics, name)
		return true
	}
	return false
}

// Registry returns the underlying Prometheus registry
func (mc *MetricsCollector) Registry() *prometheus.Registry {
	return mc.registry
}

// Clear removes all custom registered metrics
func (mc *MetricsCollector) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for name, collector := range mc.metrics {
		mc.registry.Unregister(collector)
		delete(mc.metrics, name)
	}
}

// CreateServiceInfoMetric creates a service info metric using this collector's registry
func (mc *MetricsCollector) CreateServiceInfoMetric(
	serviceName, serviceVersion, instanceID string,
	commitSHA, buildTime *string,
) *prometheus.GaugeVec {
	metric := CreateServiceInfoMetric(serviceName, serviceVersion, instanceID, commitSHA, buildTime)
	_ = mc.Register("service_info", metric)
	return metric
}

// CreateServiceInfoMetricFromServiceInfo creates a service info metric from ServiceInfo
func (mc *MetricsCollector) CreateServiceInfoMetricFromServiceInfo(serviceInfo metadata.ServiceInfo) *prometheus.GaugeVec {
	return mc.CreateServiceInfoMetric(
		serviceInfo.Name,
		serviceInfo.Version,
		serviceInfo.InstanceID,
		&serviceInfo.CommitSHA,
		&serviceInfo.BuildTime,
	)
}

// CreateServiceInfoMetric creates a service info metric following Prometheus best practices
func CreateServiceInfoMetric(
	serviceName, serviceVersion, instanceID string,
	commitSHA, buildTime *string,
) *prometheus.GaugeVec {
	labelNames := []string{"service", "version", "instance"}

	if commitSHA != nil {
		labelNames = append(labelNames, "commit")
	}
	if buildTime != nil {
		labelNames = append(labelNames, "build_time")
	}

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_info",
			Help: "Service information and build metadata",
		},
		labelNames,
	)

	labels := prometheus.Labels{
		"service":  serviceName,
		"version":  serviceVersion,
		"instance": instanceID,
	}

	if commitSHA != nil {
		labels["commit"] = *commitSHA
	}
	if buildTime != nil {
		labels["build_time"] = *buildTime
	}

	gauge.With(labels).Set(1)
	return gauge
}

// CreateServiceInfoMetricFromServiceInfo creates a service info metric from ServiceInfo instance
func CreateServiceInfoMetricFromServiceInfo(serviceInfo metadata.ServiceInfo) *prometheus.GaugeVec {
	return CreateServiceInfoMetric(
		serviceInfo.Name,
		serviceInfo.Version,
		serviceInfo.InstanceID,
		&serviceInfo.CommitSHA,
		&serviceInfo.BuildTime,
	)
}
