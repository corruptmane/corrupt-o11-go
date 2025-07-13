// Package internal provides internal utilities for the observability library.
package internal

import (
	"fmt"
	"os"
	"strings"
)

// EnvBool parses a boolean environment variable with flexible value support.
// Supports: true/false, t/f, 1/0, yes/no, y/n, on/off (case insensitive).
func EnvBool(envVar, defaultValue string) (bool, error) {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(envVar)))
	if value == "" {
		value = strings.ToLower(strings.TrimSpace(defaultValue))
	}

	switch value {
	case "true", "t", "1", "yes", "y", "on":
		return true, nil
	case "false", "f", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value for %s: '%s'. Use true/false, 1/0, yes/no, on/off", envVar, value)
	}
}

// MustEnvBool is like EnvBool but panics on error.
func MustEnvBool(envVar, defaultValue string) bool {
	result, err := EnvBool(envVar, defaultValue)
	if err != nil {
		panic(err)
	}
	return result
}
