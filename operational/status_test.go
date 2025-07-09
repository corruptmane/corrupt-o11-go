package operational

import (
	"testing"
)

func TestNewStatus(t *testing.T) {
	status := NewStatus()

	if !status.IsAlive() {
		t.Error("Expected new status to be alive")
	}

	if status.IsReady() {
		t.Error("Expected new status to not be ready")
	}
}

func TestSetReady(t *testing.T) {
	status := NewStatus()

	status.SetReady(true)
	if !status.IsReady() {
		t.Error("Expected status to be ready after SetReady(true)")
	}

	status.SetReady(false)
	if status.IsReady() {
		t.Error("Expected status to not be ready after SetReady(false)")
	}
}

func TestSetAlive(t *testing.T) {
	status := NewStatus()

	status.SetAlive(false)
	if status.IsAlive() {
		t.Error("Expected status to not be alive after SetAlive(false)")
	}

	status.SetAlive(true)
	if !status.IsAlive() {
		t.Error("Expected status to be alive after SetAlive(true)")
	}
}
