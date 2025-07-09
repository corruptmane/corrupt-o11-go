package metadata

import (
	"os"
	"testing"
)

func TestFromEnv(t *testing.T) {
	// Test with default values
	serviceInfo := FromEnv()

	if serviceInfo.Name != defaultUnknownValue {
		t.Errorf("Expected Name to be '%s', got %s", defaultUnknownValue, serviceInfo.Name)
	}
	if serviceInfo.Version != defaultUnknownValue {
		t.Errorf("Expected Version to be '%s', got %s", defaultUnknownValue, serviceInfo.Version)
	}
	if serviceInfo.InstanceID != defaultUnknownValue {
		t.Errorf("Expected InstanceID to be '%s', got %s", defaultUnknownValue, serviceInfo.InstanceID)
	}
	if serviceInfo.CommitSHA != defaultUnknownValue {
		t.Errorf("Expected CommitSHA to be '%s', got %s", defaultUnknownValue, serviceInfo.CommitSHA)
	}
	if serviceInfo.BuildTime != defaultUnknownValue {
		t.Errorf("Expected BuildTime to be '%s', got %s", defaultUnknownValue, serviceInfo.BuildTime)
	}
}

func TestFromEnvWithValues(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("SERVICE_VERSION", "1.0.0")
	os.Setenv("INSTANCE_ID", "test-instance")
	os.Setenv("COMMIT_SHA", "abc123")
	os.Setenv("BUILD_TIME", "2023-01-01T00:00:00Z")
	defer func() {
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("SERVICE_VERSION")
		os.Unsetenv("INSTANCE_ID")
		os.Unsetenv("COMMIT_SHA")
		os.Unsetenv("BUILD_TIME")
	}()

	serviceInfo := FromEnv()

	if serviceInfo.Name != "test-service" {
		t.Errorf("Expected Name to be 'test-service', got %s", serviceInfo.Name)
	}
	if serviceInfo.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got %s", serviceInfo.Version)
	}
	if serviceInfo.InstanceID != "test-instance" {
		t.Errorf("Expected InstanceID to be 'test-instance', got %s", serviceInfo.InstanceID)
	}
	if serviceInfo.CommitSHA != "abc123" {
		t.Errorf("Expected CommitSHA to be 'abc123', got %s", serviceInfo.CommitSHA)
	}
	if serviceInfo.BuildTime != "2023-01-01T00:00:00Z" {
		t.Errorf("Expected BuildTime to be '2023-01-01T00:00:00Z', got %s", serviceInfo.BuildTime)
	}
}

func TestAsMap(t *testing.T) {
	serviceInfo := ServiceInfo{
		Name:       "test-service",
		Version:    "1.0.0",
		InstanceID: "test-instance",
		CommitSHA:  "abc123",
		BuildTime:  "2023-01-01T00:00:00Z",
	}

	m := serviceInfo.AsMap()

	expectedKeys := []string{"service_name", "version", "instance_id", "commit_sha", "build_time"}
	for _, key := range expectedKeys {
		if _, exists := m[key]; !exists {
			t.Errorf("Expected key '%s' to exist in map", key)
		}
	}

	if m["service_name"] != "test-service" {
		t.Errorf("Expected service_name to be 'test-service', got %s", m["service_name"])
	}
	if m["version"] != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got %s", m["version"])
	}
}
