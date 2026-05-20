package cloud

import (
	"context"
	"time"

	"github.com/infrasphere/control-plane/backend/internal/models"
)

type ProviderConfig struct {
	OrganizationID string                 `json:"organizationId"`
	AccountID      string                 `json:"accountId"`
	Region         string                 `json:"region"`
	Credentials    map[string]string      `json:"credentials,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type WorkloadSpec struct {
	Name         string                 `json:"name"`
	Provider     string                 `json:"provider"`
	Target       map[string]string      `json:"target"`
	WorkloadType string                 `json:"workloadType"`
	Config       map[string]interface{} `json:"config"`
}

type DeploymentResult struct {
	DeploymentID string    `json:"deploymentId"`
	Status       string    `json:"status"`
	Logs         []string  `json:"logs"`
	CostEstimate float64   `json:"costEstimate"`
	CreatedAt    time.Time `json:"createdAt"`
}

type MetricsQuery struct {
	ResourceID string    `json:"resourceId"`
	Metric     string    `json:"metric"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}

type MetricsResult struct {
	Series []MetricPoint `json:"series"`
}

type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type CostQuery struct {
	OrganizationID string    `json:"organizationId"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
}

type CostResult struct {
	Total      float64            `json:"total"`
	ByProvider map[string]float64 `json:"byProvider"`
}

type CloudProvider interface {
	Name() string
	Connect(ctx context.Context, config ProviderConfig) error
	ListInventory(ctx context.Context) (*models.Inventory, error)
	DeployWorkload(ctx context.Context, spec WorkloadSpec) (*DeploymentResult, error)
	GetMetrics(ctx context.Context, query MetricsQuery) (*MetricsResult, error)
	GetCosts(ctx context.Context, query CostQuery) (*CostResult, error)
	ValidateCredentials(ctx context.Context) error
}

