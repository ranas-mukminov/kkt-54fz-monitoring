package collector

import (
	"context"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
)

// Collector is the interface for data collectors
type Collector interface {
	// Start begins collecting data
	Start(ctx context.Context) error

	// Stop stops the collector
	Stop() error

	// Name returns the collector name
	Name() string

	// Metrics returns the current metrics channel
	Metrics() <-chan domain.Metrics

	// Errors returns the errors channel
	Errors() <-chan domain.KKTError
}
