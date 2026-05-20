package observability

import (
	"context"

	"github.com/infrasphere/control-plane/backend/internal/models"
	"github.com/infrasphere/control-plane/backend/pkg/cloud"
)

type LogsQuery struct{ ResourceID string }
type TracesQuery struct{ Service string }
type LogsResult struct{ Lines []string `json:"lines"` }
type TracesResult struct{ Spans []map[string]string `json:"spans"` }

type ObservabilityProvider interface {
	Name() string
	QueryMetrics(ctx context.Context, query cloud.MetricsQuery) (*cloud.MetricsResult, error)
	QueryLogs(ctx context.Context, query LogsQuery) (*LogsResult, error)
	QueryTraces(ctx context.Context, query TracesQuery) (*TracesResult, error)
	ListAlerts(ctx context.Context) ([]models.Alert, error)
}
