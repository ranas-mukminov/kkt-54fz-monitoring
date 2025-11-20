package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/config"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/pkg/logger"
)

// HTTPOFDCollector collects data from OFD HTTP API
type HTTPOFDCollector struct {
	cfg          config.HTTPOFDConfig
	log          *logger.Logger
	metricsChan  chan domain.Metrics
	errorsChan   chan domain.KKTError
	stopChan     chan struct{}
}

// NewHTTPOFDCollector creates a new HTTP OFD collector
func NewHTTPOFDCollector(cfg config.HTTPOFDConfig, log *logger.Logger) *HTTPOFDCollector {
	return &HTTPOFDCollector{
		cfg:         cfg,
		log:         log,
		metricsChan: make(chan domain.Metrics, 100),
		errorsChan:  make(chan domain.KKTError, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start begins collecting data
func (c *HTTPOFDCollector) Start(ctx context.Context) error {
	c.log.Info("Starting HTTP OFD collector", "url", c.cfg.URL)

	go c.collect(ctx)

	return nil
}

// Stop stops the collector
func (c *HTTPOFDCollector) Stop() error {
	c.log.Info("Stopping HTTP OFD collector")
	close(c.stopChan)
	return nil
}

// Name returns the collector name
func (c *HTTPOFDCollector) Name() string {
	return "http_ofd"
}

// Metrics returns the metrics channel
func (c *HTTPOFDCollector) Metrics() <-chan domain.Metrics {
	return c.metricsChan
}

// Errors returns the errors channel
func (c *HTTPOFDCollector) Errors() <-chan domain.KKTError {
	return c.errorsChan
}

// collect is the main collection loop
func (c *HTTPOFDCollector) collect(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.stopChan:
			return
		case <-ticker.C:
			if err := c.collectOnce(); err != nil {
				c.log.Error("Failed to collect from HTTP OFD", "error", err)
			}
		}
	}
}

// collectOnce performs one collection cycle
func (c *HTTPOFDCollector) collectOnce() error {
	// TODO: Implement actual HTTP OFD API calls
	// This is a stub implementation
	c.log.Debug("Collecting from HTTP OFD", "url", c.cfg.URL)

	// Stub: generate sample metrics
	metrics := domain.Metrics{
		KKTID:            "stub-kkt-002",
		Timestamp:        time.Now(),
		Status:           domain.KKTStatusRunning,
		DocumentsTotal:   2500,
		ErrorsByType:     make(map[domain.ErrorType]int64),
		OFDSyncStatus:    domain.OFDSyncStatusSynced,
		ShiftStatus:      domain.ShiftStatusOpen,
		LastDocumentTime: time.Now().Add(-2 * time.Minute),
		FDMemoryUsage:    62.3,
		DocumentsPerHour: 75.0,
		AverageSyncTime:  1.8,
	}

	select {
	case c.metricsChan <- metrics:
	default:
		return fmt.Errorf("metrics channel full")
	}

	return nil
}
