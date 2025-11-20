package ai

import (
	"context"
	"fmt"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
)

// MockProvider is a mock AI provider for testing and development
type MockProvider struct{}

// NewMockProvider creates a new mock AI provider
func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

// ClusterErrors clusters similar errors together (mock implementation)
func (m *MockProvider) ClusterErrors(ctx context.Context, errors []domain.KKTError) ([]ErrorCluster, error) {
	if len(errors) == 0 {
		return []ErrorCluster{}, nil
	}

	// Simple mock clustering: group by error type
	clusters := make(map[domain.ErrorType][]domain.KKTError)
	for _, err := range errors {
		clusters[err.ErrorType] = append(clusters[err.ErrorType], err)
	}

	result := make([]ErrorCluster, 0, len(clusters))
	clusterID := 1
	for errType, errs := range clusters {
		if len(errs) == 0 {
			continue
		}

		cluster := ErrorCluster{
			ID:         fmt.Sprintf("cluster-%d", clusterID),
			Errors:     errs,
			Pattern:    fmt.Sprintf("Errors of type: %v", errType),
			Severity:   errs[0].Severity,
			Count:      len(errs),
			FirstSeen:  errs[0].Timestamp.Format("2006-01-02 15:04:05"),
			LastSeen:   errs[len(errs)-1].Timestamp.Format("2006-01-02 15:04:05"),
			Suggestion: generateMockSuggestion(errType),
		}
		result = append(result, cluster)
		clusterID++
	}

	return result, nil
}

// GenerateAlertRecommendations generates alert recommendations (mock implementation)
func (m *MockProvider) GenerateAlertRecommendations(ctx context.Context, metrics []domain.Metrics) ([]AlertRecommendation, error) {
	recommendations := []AlertRecommendation{
		{
			ID:          "rec-1",
			Type:        "kkt_unavailable",
			Condition:   "kkt_status == 0",
			Threshold:   5.0,
			Severity:    "critical",
			Description: "Alert when KKT is unavailable for more than 5 minutes",
			Rationale:   "Based on historical data, KKT unavailability beyond 5 minutes typically indicates a serious issue requiring immediate attention.",
		},
		{
			ID:          "rec-2",
			Type:        "fd_memory_high",
			Condition:   "kkt_fd_memory_usage_percent > 80",
			Threshold:   80.0,
			Severity:    "warning",
			Description: "Alert when fiscal drive memory usage exceeds 80%",
			Rationale:   "High memory usage may lead to fiscal drive overflow. Early warning allows for proactive replacement.",
		},
		{
			ID:          "rec-3",
			Type:        "ofd_sync_failure",
			Condition:   "kkt_ofd_sync_status == 3",
			Threshold:   0.0,
			Severity:    "high",
			Description: "Alert on OFD synchronization failures",
			Rationale:   "OFD sync failures can lead to compliance violations. Immediate action required.",
		},
		{
			ID:          "rec-4",
			Type:        "low_document_rate",
			Condition:   "kkt_documents_per_hour < 10",
			Threshold:   10.0,
			Severity:    "warning",
			Description: "Alert when document processing rate is unusually low",
			Rationale:   "Low document rate may indicate KKT malfunction or business operation issues.",
		},
	}

	return recommendations, nil
}

// Name returns the provider name
func (m *MockProvider) Name() string {
	return "mock"
}

// generateMockSuggestion generates a mock suggestion based on error type
func generateMockSuggestion(errType domain.ErrorType) string {
	switch errType {
	case domain.ErrorTypeNetwork:
		return "Check network connectivity and firewall settings"
	case domain.ErrorTypeFiscalDrive:
		return "Verify fiscal drive status and consider replacement if near capacity"
	case domain.ErrorTypeOFD:
		return "Check OFD service status and credentials"
	case domain.ErrorTypePrinter:
		return "Check printer connection and paper supply"
	case domain.ErrorTypeHardware:
		return "Inspect hardware components and connections"
	case domain.ErrorTypeSoftware:
		return "Review software logs and consider updating to latest version"
	case domain.ErrorTypeConfiguration:
		return "Review configuration settings and compare with documentation"
	default:
		return "Contact technical support for assistance"
	}
}
