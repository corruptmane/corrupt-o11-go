package metadata

import (
	"os"
)

const (
	// defaultUnknownValue is the default value for unknown service metadata
	defaultUnknownValue = "unknown-dev"
)

// ServiceInfo holds metadata about the service instance
type ServiceInfo struct {
	Name       string
	Version    string
	InstanceID string
	CommitSHA  string
	BuildTime  string
}

// FromEnv creates ServiceInfo from environment variables
func FromEnv() ServiceInfo {
	return ServiceInfo{
		Name:       getEnvOrDefault("SERVICE_NAME"),
		Version:    getEnvOrDefault("SERVICE_VERSION"),
		InstanceID: getEnvOrDefault("INSTANCE_ID"),
		CommitSHA:  getEnvOrDefault("COMMIT_SHA"),
		BuildTime:  getEnvOrDefault("BUILD_TIME"),
	}
}

// AsMap returns service info as a map for consistent labeling
func (si ServiceInfo) AsMap() map[string]string {
	return map[string]string{
		"service_name": si.Name,
		"version":      si.Version,
		"instance_id":  si.InstanceID,
		"commit_sha":   si.CommitSHA,
		"build_time":   si.BuildTime,
	}
}

func getEnvOrDefault(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultUnknownValue
}
