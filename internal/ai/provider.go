package ai

import (
	"context"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/domain"
)

// AIProvider is the interface for AI providers
type AIProvider interface {
	// ClusterErrors clusters similar errors together
	ClusterErrors(ctx context.Context, errors []domain.KKTError) ([]ErrorCluster, error)

	// GenerateAlertRecommendations generates alert recommendations
	GenerateAlertRecommendations(ctx context.Context, metrics []domain.Metrics) ([]AlertRecommendation, error)

	// Name returns the provider name
	Name() string
}

// ErrorCluster represents a cluster of similar errors
type ErrorCluster struct {
	ID          string            `json:"id"`
	Errors      []domain.KKTError `json:"errors"`
	Pattern     string            `json:"pattern"`
	Severity    domain.ErrorSeverity `json:"severity"`
	Count       int               `json:"count"`
	FirstSeen   string            `json:"first_seen"`
	LastSeen    string            `json:"last_seen"`
	Suggestion  string            `json:"suggestion"`
}

// AlertRecommendation represents an alert recommendation
type AlertRecommendation struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Condition   string  `json:"condition"`
	Threshold   float64 `json:"threshold"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Rationale   string  `json:"rationale"`
}
