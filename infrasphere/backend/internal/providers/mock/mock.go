package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/infrasphere/control-plane/backend/internal/models"
	"github.com/infrasphere/control-plane/backend/pkg/cloud"
)

type Provider struct {
	accountID string
	region    string
}

func New() *Provider {
	return &Provider{accountID: "mock-prod", region: "us-west-2"}
}

func (p *Provider) Name() string { return "mock" }

func (p *Provider) Connect(ctx context.Context, cfg cloud.ProviderConfig) error {
	if cfg.AccountID != "" {
		p.accountID = cfg.AccountID
	}
	if cfg.Region != "" {
		p.region = cfg.Region
	}
	return nil
}

func (p *Provider) ValidateCredentials(ctx context.Context) error { return nil }

func (p *Provider) ListInventory(ctx context.Context) (*models.Inventory, error) {
	now := time.Now().UTC()
	resources := []models.Resource{
		resource("res-vpc-1", "prod-vpc", "mock", p.accountID, p.region, "", "vpc", "Network", "available", "healthy", 121, map[string]string{"env": "prod"}, map[string]interface{}{"cidr": "10.40.0.0/16"}, now),
		resource("res-subnet-1", "payments-private-a", "mock", p.accountID, p.region, "us-west-2a", "subnet", "Network", "available", "healthy", 24, map[string]string{"app": "payments"}, map[string]interface{}{"cidr": "10.40.1.0/24"}, now),
		resource("res-eks-1", "payments-eks", "mock", p.accountID, p.region, "", "kubernetes_cluster", "Kubernetes", "running", "warning", 1780, map[string]string{"app": "payments", "env": "prod"}, map[string]interface{}{"version": "1.28", "nodes": 12}, now),
		resource("res-vm-1", "payments-api-7c9d", "mock", p.accountID, p.region, "us-west-2a", "compute_instance", "Compute", "running", "healthy", 312, map[string]string{"app": "payments"}, map[string]interface{}{"cpuUtilization": 18, "memoryUtilization": 41}, now),
		resource("res-db-1", "payments-postgres", "mock", p.accountID, p.region, "us-west-2a", "database", "Database", "running", "healthy", 950, map[string]string{"app": "payments", "tier": "data"}, map[string]interface{}{"engine": "postgres", "storageEncrypted": true}, now),
		resource("res-lb-1", "payments-public", "mock", p.accountID, p.region, "", "load_balancer", "Load Balancing", "active", "critical", 203, map[string]string{"internet-facing": "true"}, map[string]interface{}{"scheme": "internet-facing", "openPorts": []int{80, 443}}, now),
	}
	return &models.Inventory{
		Resources: resources,
		Relationships: []models.ResourceRelationship{
			{ID: "rel-1", SourceID: "res-vm-1", TargetID: "res-subnet-1", Relationship: "attached_to"},
			{ID: "rel-2", SourceID: "res-subnet-1", TargetID: "res-vpc-1", Relationship: "belongs_to"},
			{ID: "rel-3", SourceID: "res-eks-1", TargetID: "res-vm-1", Relationship: "runs_on"},
			{ID: "rel-4", SourceID: "res-lb-1", TargetID: "res-vm-1", Relationship: "routes_to"},
			{ID: "rel-5", SourceID: "res-vm-1", TargetID: "res-db-1", Relationship: "uses"},
		},
	}, nil
}

func (p *Provider) DeployWorkload(ctx context.Context, spec cloud.WorkloadSpec) (*cloud.DeploymentResult, error) {
	return &cloud.DeploymentResult{
		DeploymentID: "dep-" + time.Now().UTC().Format("20060102150405"),
		Status:       "WAITING_FOR_APPROVAL",
		Logs:         []string{"validated workload spec", "generated provider-neutral deployment plan", "approval required before apply"},
		CostEstimate: 482.25,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

func (p *Provider) GetMetrics(ctx context.Context, query cloud.MetricsQuery) (*cloud.MetricsResult, error) {
	now := time.Now().UTC()
	return &cloud.MetricsResult{Series: []cloud.MetricPoint{
		{Timestamp: now.Add(-4 * time.Hour), Value: 44},
		{Timestamp: now.Add(-3 * time.Hour), Value: 51},
		{Timestamp: now.Add(-2 * time.Hour), Value: 38},
		{Timestamp: now.Add(-1 * time.Hour), Value: 67},
		{Timestamp: now, Value: 48},
	}}, nil
}

func (p *Provider) GetCosts(ctx context.Context, query cloud.CostQuery) (*cloud.CostResult, error) {
	return &cloud.CostResult{Total: 3390, ByProvider: map[string]float64{"AWS": 1420, "GCP": 730, "Azure": 810, "VCF": 430}}, nil
}

func resource(id, name, provider, account, region, zone, typ, category, status, health string, cost float64, tags map[string]string, metadata map[string]interface{}, now time.Time) models.Resource {
	return models.Resource{ID: id, Name: name, Provider: provider, AccountID: account, Region: region, Zone: zone, ResourceType: typ, Category: category, Status: status, Health: health, CostMonthly: cost, Tags: tags, Metadata: metadata, CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now}
}

func ExampleQuestionAnswer(question string) string {
	return fmt.Sprintf("Based on inventory, the highest-risk item related to %q is the internet-facing payments load balancer. Recommended next step: verify target groups, tighten security policy, and confirm observability is enabled before changing production.", question)
}

