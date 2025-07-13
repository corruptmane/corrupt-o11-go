package metrics

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Test with environment variables set
	os.Setenv("METRICS_ENABLE_GO", "false")
	os.Setenv("METRICS_ENABLE_PROCESS", "true")
	os.Setenv("METRICS_PREFIX", "myapp_")
	defer func() {
		os.Unsetenv("METRICS_ENABLE_GO")
		os.Unsetenv("METRICS_ENABLE_PROCESS")
		os.Unsetenv("METRICS_PREFIX")
	}()

	config := FromEnv()

	if config.EnableGoCollector != false {
		t.Errorf("Expected EnableGoCollector to be false, got %v", config.EnableGoCollector)
	}
	if config.EnableProcessCollector != true {
		t.Errorf("Expected EnableProcessCollector to be true, got %v", config.EnableProcessCollector)
	}
	if config.MetricPrefix != "myapp_" {
		t.Errorf("Expected MetricPrefix to be 'myapp_', got %s", config.MetricPrefix)
	}
}

func TestFromEnvDefaults(t *testing.T) {
	// Test with no environment variables (defaults)
	config := FromEnv()

	if config.EnableGoCollector != true {
		t.Errorf("Expected EnableGoCollector default to be true, got %v", config.EnableGoCollector)
	}
	if config.EnableProcessCollector != true {
		t.Errorf("Expected EnableProcessCollector default to be true, got %v", config.EnableProcessCollector)
	}
	if config.MetricPrefix != "" {
		t.Errorf("Expected MetricPrefix default to be empty, got %s", config.MetricPrefix)
	}
}
