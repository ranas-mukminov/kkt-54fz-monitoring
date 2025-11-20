package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/config"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/pkg/logger"
)

// FileLogCollector collects data from file logs
type FileLogCollector struct {
	cfg          config.FileLogConfig
	log          *logger.Logger
	metricsChan  chan domain.Metrics
	errorsChan   chan domain.KKTError
	stopChan     chan struct{}
}

// NewFileLogCollector creates a new file log collector
func NewFileLogCollector(cfg config.FileLogConfig, log *logger.Logger) *FileLogCollector {
	return &FileLogCollector{
		cfg:         cfg,
		log:         log,
		metricsChan: make(chan domain.Metrics, 100),
		errorsChan:  make(chan domain.KKTError, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start begins collecting data
func (c *FileLogCollector) Start(ctx context.Context) error {
	c.log.Info("Starting file log collector", "path", c.cfg.Path)

	go c.collect(ctx)

	return nil
}

// Stop stops the collector
func (c *FileLogCollector) Stop() error {
	c.log.Info("Stopping file log collector")
	close(c.stopChan)
	return nil
}

// Name returns the collector name
func (c *FileLogCollector) Name() string {
	return "file_log"
}

// Metrics returns the metrics channel
func (c *FileLogCollector) Metrics() <-chan domain.Metrics {
	return c.metricsChan
}

// Errors returns the errors channel
func (c *FileLogCollector) Errors() <-chan domain.KKTError {
	return c.errorsChan
}

// collect is the main collection loop
func (c *FileLogCollector) collect(ctx context.Context) {
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
				c.log.Error("Failed to collect from file logs", "error", err)
			}
		}
	}
}

// collectOnce performs one collection cycle
func (c *FileLogCollector) collectOnce() error {
	// TODO: Implement actual file log parsing
	// This is a stub implementation
	c.log.Debug("Collecting from file logs", "path", c.cfg.Path)

	// Stub: generate sample metrics
	metrics := domain.Metrics{
		KKTID:            "stub-kkt-001",
		Timestamp:        time.Now(),
		Status:           domain.KKTStatusRunning,
		DocumentsTotal:   1000,
		ErrorsByType:     make(map[domain.ErrorType]int64),
		OFDSyncStatus:    domain.OFDSyncStatusSynced,
		ShiftStatus:      domain.ShiftStatusOpen,
		LastDocumentTime: time.Now().Add(-5 * time.Minute),
		FDMemoryUsage:    45.5,
		DocumentsPerHour: 50.0,
		AverageSyncTime:  2.5,
	}

	select {
	case c.metricsChan <- metrics:
	default:
		return fmt.Errorf("metrics channel full")
	}

	return nil
}
