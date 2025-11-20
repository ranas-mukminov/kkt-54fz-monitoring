package exporter

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/pkg/logger"
)

// Exporter implements Prometheus exporter for KKT metrics
type Exporter struct {
	log *logger.Logger

	// Metrics
	kktStatus           *prometheus.GaugeVec
	kktDocumentsTotal   *prometheus.GaugeVec
	kktErrorsTotal      *prometheus.GaugeVec
	kktOFDSyncStatus    *prometheus.GaugeVec
	kktShiftStatus      *prometheus.GaugeVec
	kktLastDocumentTime *prometheus.GaugeVec
	kktFDMemoryUsage    *prometheus.GaugeVec
	kktDocumentsPerHour *prometheus.GaugeVec
	kktAvgSyncTime      *prometheus.GaugeVec

	mu sync.RWMutex
}

// New creates a new Prometheus exporter
func New(log *logger.Logger) *Exporter {
	e := &Exporter{
		log: log,
	}

	e.initMetrics()
	e.registerMetrics()

	return e
}

// initMetrics initializes Prometheus metrics
func (e *Exporter) initMetrics() {
	e.kktStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_status",
			Help: "KKT device status (0=unavailable, 1=running, 2=error)",
		},
		[]string{"kkt_id"},
	)

	e.kktDocumentsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_documents_total",
			Help: "Total number of fiscal documents",
		},
		[]string{"kkt_id"},
	)

	e.kktErrorsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"kkt_id", "error_type"},
	)

	e.kktOFDSyncStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_ofd_sync_status",
			Help: "OFD synchronization status (0=unknown, 1=synced, 2=pending, 3=error)",
		},
		[]string{"kkt_id"},
	)

	e.kktShiftStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_shift_status",
			Help: "Shift status (0=closed, 1=open)",
		},
		[]string{"kkt_id"},
	)

	e.kktLastDocumentTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_last_document_timestamp",
			Help: "Timestamp of last fiscal document (Unix time)",
		},
		[]string{"kkt_id"},
	)

	e.kktFDMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_fd_memory_usage_percent",
			Help: "Fiscal drive memory usage percentage",
		},
		[]string{"kkt_id"},
	)

	e.kktDocumentsPerHour = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_documents_per_hour",
			Help: "Average documents processed per hour",
		},
		[]string{"kkt_id"},
	)

	e.kktAvgSyncTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kkt_average_sync_time_seconds",
			Help: "Average OFD synchronization time in seconds",
		},
		[]string{"kkt_id"},
	)
}

// registerMetrics registers metrics with Prometheus
func (e *Exporter) registerMetrics() {
	prometheus.MustRegister(
		e.kktStatus,
		e.kktDocumentsTotal,
		e.kktErrorsTotal,
		e.kktOFDSyncStatus,
		e.kktShiftStatus,
		e.kktLastDocumentTime,
		e.kktFDMemoryUsage,
		e.kktDocumentsPerHour,
		e.kktAvgSyncTime,
	)
}

// UpdateMetrics updates metrics from domain.Metrics
func (e *Exporter) UpdateMetrics(metrics domain.Metrics) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.kktStatus.WithLabelValues(metrics.KKTID).Set(float64(metrics.Status))
	e.kktDocumentsTotal.WithLabelValues(metrics.KKTID).Set(float64(metrics.DocumentsTotal))
	e.kktOFDSyncStatus.WithLabelValues(metrics.KKTID).Set(float64(metrics.OFDSyncStatus))
	e.kktShiftStatus.WithLabelValues(metrics.KKTID).Set(float64(metrics.ShiftStatus))
	e.kktLastDocumentTime.WithLabelValues(metrics.KKTID).Set(float64(metrics.LastDocumentTime.Unix()))
	e.kktFDMemoryUsage.WithLabelValues(metrics.KKTID).Set(metrics.FDMemoryUsage)
	e.kktDocumentsPerHour.WithLabelValues(metrics.KKTID).Set(metrics.DocumentsPerHour)
	e.kktAvgSyncTime.WithLabelValues(metrics.KKTID).Set(metrics.AverageSyncTime)

	// Update error gauges with current counts
	for errorType, count := range metrics.ErrorsByType {
		e.kktErrorsTotal.WithLabelValues(metrics.KKTID, errorTypeName(errorType)).Set(float64(count))
	}
}

// Handler returns the HTTP handler for metrics endpoint
func (e *Exporter) Handler() http.Handler {
	return promhttp.Handler()
}

// Start starts the exporter
func (e *Exporter) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", e.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		e.log.Info("Shutting down metrics server")
		if err := server.Shutdown(context.Background()); err != nil {
			e.log.Error("Failed to shutdown metrics server", "error", err)
		}
	}()

	e.log.Info("Starting metrics server", "addr", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start metrics server: %w", err)
	}

	return nil
}

// errorTypeName converts ErrorType to string
func errorTypeName(et domain.ErrorType) string {
	switch et {
	case domain.ErrorTypeNetwork:
		return "network"
	case domain.ErrorTypeFiscalDrive:
		return "fiscal_drive"
	case domain.ErrorTypeOFD:
		return "ofd"
	case domain.ErrorTypePrinter:
		return "printer"
	case domain.ErrorTypeHardware:
		return "hardware"
	case domain.ErrorTypeSoftware:
		return "software"
	case domain.ErrorTypeConfiguration:
		return "configuration"
	default:
		return "unknown"
	}
}
