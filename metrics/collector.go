// Package metrics provides Prometheus metrics collection and management.
package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/corruptmane/corrupt-o11y-go/metadata"
)

// MetricsCollector provides a centralized registry for Prometheus metrics
type MetricsCollector struct {
	config   MetricsConfig
	registry *prometheus.Registry
	metrics  map[string]prometheus.Collector
	mu       sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector with default configuration
func NewMetricsCollector() *MetricsCollector {
	return NewMetricsCollectorWithConfig(FromEnv())
}

// NewMetricsCollectorWithConfig creates a new metrics collector with the given configuration
func NewMetricsCollectorWithConfig(config MetricsConfig) *MetricsCollector {
	registry := prometheus.NewRegistry()

	// Register built-in collectors based on configuration
	if config.EnableGoCollector {
		registry.MustRegister(collectors.NewGoCollector())
	}
	if config.EnableProcessCollector {
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	return &MetricsCollector{
		config:   config,
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

// IsRegistered checks if a metric with the given name is registered
func (mc *MetricsCollector) IsRegistered(name string) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	_, exists := mc.metrics[name]
	return exists
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

// applyPrefix applies the configured metric prefix to the opts Name field
func (mc *MetricsCollector) applyPrefix(name *string) {
	if mc.config.MetricPrefix != "" {
		*name = mc.config.MetricPrefix + *name
	}
}

// CreateCounterVec creates a CounterVec metric and registers it with this collector
func (mc *MetricsCollector) CreateCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	counter := prometheus.NewCounterVec(opts, labelNames)
	_ = mc.Register(originalName, counter)
	return counter
}

// CreateGaugeVec creates a GaugeVec metric and registers it with this collector
func (mc *MetricsCollector) CreateGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	gauge := prometheus.NewGaugeVec(opts, labelNames)
	_ = mc.Register(originalName, gauge)
	return gauge
}

// CreateHistogramVec creates a HistogramVec metric and registers it with this collector
func (mc *MetricsCollector) CreateHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	histogram := prometheus.NewHistogramVec(opts, labelNames)
	_ = mc.Register(originalName, histogram)
	return histogram
}

// CreateSummaryVec creates a SummaryVec metric and registers it with this collector
func (mc *MetricsCollector) CreateSummaryVec(opts prometheus.SummaryOpts, labelNames []string) *prometheus.SummaryVec {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	summary := prometheus.NewSummaryVec(opts, labelNames)
	_ = mc.Register(originalName, summary)
	return summary
}

// CreateCounterFunc creates a CounterFunc metric and registers it with this collector
func (mc *MetricsCollector) CreateCounterFunc(opts prometheus.CounterOpts, countFunc func() float64) prometheus.CounterFunc {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	counter := prometheus.NewCounterFunc(opts, countFunc)
	_ = mc.Register(originalName, counter)
	return counter
}

// CreateGaugeFunc creates a GaugeFunc metric and registers it with this collector
func (mc *MetricsCollector) CreateGaugeFunc(opts prometheus.GaugeOpts, gaugeFunc func() float64) prometheus.GaugeFunc {
	originalName := opts.Name
	mc.applyPrefix(&opts.Name)

	gauge := prometheus.NewGaugeFunc(opts, gaugeFunc)
	_ = mc.Register(originalName, gauge)
	return gauge
}
