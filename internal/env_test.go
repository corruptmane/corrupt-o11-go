package internal

import (
	"testing"
)

func TestEnvBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		hasError bool
	}{
		// True values
		{"true", true, false},
		{"TRUE", true, false},
		{"t", true, false},
		{"T", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"YES", true, false},
		{"y", true, false},
		{"Y", true, false},
		{"on", true, false},
		{"ON", true, false},

		// False values
		{"false", false, false},
		{"FALSE", false, false},
		{"f", false, false},
		{"F", false, false},
		{"0", false, false},
		{"no", false, false},
		{"NO", false, false},
		{"n", false, false},
		{"N", false, false},
		{"off", false, false},
		{"OFF", false, false},

		// Invalid values
		{"invalid", false, true},
		{"maybe", false, true},
		{"2", false, true},
	}

	for _, test := range tests {
		result, err := EnvBool("TEST_VAR", test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("EnvBool with value '%s' should return error but didn't", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("EnvBool with value '%s' should not return error but got: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("EnvBool with value '%s' = %v, expected %v", test.input, result, test.expected)
			}
		}
	}
}

func TestMustEnvBool(t *testing.T) {
	// Test that valid values don't panic
	result := MustEnvBool("TEST_VAR", "true")
	if result != true {
		t.Errorf("MustEnvBool('true') = %v, expected true", result)
	}

	// Test that invalid values panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustEnvBool with invalid value should panic")
		}
	}()
	MustEnvBool("TEST_VAR", "invalid")
}
