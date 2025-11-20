package ai

import (
	"context"
	"testing"
	"time"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
)

func TestMockProvider_ClusterErrors(t *testing.T) {
	provider := NewMockProvider()

	errors := []domain.KKTError{
		{
			ID:        "err-1",
			KKTID:     "kkt-001",
			ErrorType: domain.ErrorTypeNetwork,
			Severity:  domain.ErrorSeverityWarning,
			Message:   "Network timeout",
			Timestamp: time.Now(),
		},
		{
			ID:        "err-2",
			KKTID:     "kkt-001",
			ErrorType: domain.ErrorTypeNetwork,
			Severity:  domain.ErrorSeverityWarning,
			Message:   "Network connection lost",
			Timestamp: time.Now(),
		},
		{
			ID:        "err-3",
			KKTID:     "kkt-002",
			ErrorType: domain.ErrorTypeFiscalDrive,
			Severity:  domain.ErrorSeverityCritical,
			Message:   "Fiscal drive error",
			Timestamp: time.Now(),
		},
	}

	clusters, err := provider.ClusterErrors(context.Background(), errors)
	if err != nil {
		t.Fatalf("ClusterErrors failed: %v", err)
	}

	if len(clusters) != 2 {
		t.Errorf("Expected 2 clusters, got %d", len(clusters))
	}

	// Check that we have network and fiscal drive clusters
	hasNetwork := false
	hasFiscalDrive := false
	for _, cluster := range clusters {
		if cluster.Count == 2 {
			hasNetwork = true
		}
		if cluster.Count == 1 {
			hasFiscalDrive = true
		}
	}

	if !hasNetwork || !hasFiscalDrive {
		t.Error("Expected one network cluster (2 errors) and one fiscal drive cluster (1 error)")
	}
}

func TestMockProvider_GenerateAlertRecommendations(t *testing.T) {
	provider := NewMockProvider()

	metrics := []domain.Metrics{
		{
			KKTID:          "kkt-001",
			Status:         domain.KKTStatusRunning,
			DocumentsTotal: 1000,
		},
	}

	recommendations, err := provider.GenerateAlertRecommendations(context.Background(), metrics)
	if err != nil {
		t.Fatalf("GenerateAlertRecommendations failed: %v", err)
	}

	if len(recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}

	// Check that we have expected recommendation types
	hasKKTUnavailable := false
	hasFDMemory := false
	hasOFDSync := false

	for _, rec := range recommendations {
		switch rec.Type {
		case "kkt_unavailable":
			hasKKTUnavailable = true
		case "fd_memory_high":
			hasFDMemory = true
		case "ofd_sync_failure":
			hasOFDSync = true
		}
	}

	if !hasKKTUnavailable {
		t.Error("Expected kkt_unavailable recommendation")
	}
	if !hasFDMemory {
		t.Error("Expected fd_memory_high recommendation")
	}
	if !hasOFDSync {
		t.Error("Expected ofd_sync_failure recommendation")
	}
}

func TestMockProvider_Name(t *testing.T) {
	provider := NewMockProvider()
	if provider.Name() != "mock" {
		t.Errorf("Expected provider name 'mock', got '%s'", provider.Name())
	}
}
