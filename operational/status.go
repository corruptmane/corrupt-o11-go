package operational

import (
	"sync/atomic"
)

// Status provides thread-safe service status tracking
type Status struct {
	ready atomic.Bool
	alive atomic.Bool
}

// NewStatus creates a new Status instance
func NewStatus() *Status {
	status := &Status{}
	status.alive.Store(true)
	return status
}

// IsReady returns whether the service is ready to accept requests
func (s *Status) IsReady() bool {
	return s.ready.Load()
}

// IsAlive returns whether the service is alive (for health checks)
func (s *Status) IsAlive() bool {
	return s.alive.Load()
}

// SetReady sets the service ready status
func (s *Status) SetReady(ready bool) {
	s.ready.Store(ready)
}

// SetAlive sets the service alive status
func (s *Status) SetAlive(alive bool) {
	s.alive.Store(alive)
}
