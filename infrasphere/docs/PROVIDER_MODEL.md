# Provider Model

Providers implement the `CloudProvider` interface in `backend/pkg/cloud`.

```go
type CloudProvider interface {
    Name() string
    Connect(ctx context.Context, config ProviderConfig) error
    ListInventory(ctx context.Context) (*Inventory, error)
    DeployWorkload(ctx context.Context, spec WorkloadSpec) (*DeploymentResult, error)
    GetMetrics(ctx context.Context, query MetricsQuery) (*MetricsResult, error)
    GetCosts(ctx context.Context, query CostQuery) (*CostResult, error)
    ValidateCredentials(ctx context.Context) error
}
```

## Adding a Provider

1. Create `backend/internal/providers/<name>`.
2. Implement `Name`, `Connect`, `ValidateCredentials`, `ListInventory`, `DeployWorkload`, `GetMetrics`, and `GetCosts`.
3. Normalize provider assets into `models.Resource`.
4. Preserve provider-specific details in `Metadata`.
5. Add credential validation and least-privilege docs.
6. Register the provider in the API and worker factory.

## Inventory Normalization

Common fields power cross-cloud filtering: provider, account, region, zone, resource type, category, status, tags, health, cost, and timestamps.

## Provider Metadata

Provider-specific fields such as AWS ARNs, GCP self links, Azure resource IDs, vCenter MoRefs, NSX segment IDs, Kubernetes versions, and service SKUs belong in `Metadata`.

## Deployment Abstraction

Deployment accepts a provider-neutral workload spec. Providers convert it into ECS/EKS/EC2/Lambda, GKE/Cloud Run/Compute Engine, AKS/Container Apps/VMs/Functions, or Tanzu/vSphere/Kubernetes actions.

