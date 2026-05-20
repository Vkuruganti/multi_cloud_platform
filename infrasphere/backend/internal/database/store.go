package database

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/infrasphere/control-plane/backend/internal/models"
	"github.com/infrasphere/control-plane/backend/pkg/cloud"
)

type Store struct {
	mu            sync.RWMutex
	users         []models.User
	orgs          []models.Organization
	members       []models.OrganizationMember
	accounts      []models.CloudAccount
	deployments   []models.Deployment
	auditLogs     []map[string]interface{}
	inventory     *models.Inventory
	provider      cloud.CloudProvider
}

func NewSeededStore(provider cloud.CloudProvider) *Store {
	now := time.Now().UTC()
	inv, _ := provider.ListInventory(context.Background())
	return &Store{
		provider:  provider,
		inventory: inv,
		users: []models.User{{ID: "user-admin", Email: "admin@infrasphere.local", Name: "InfraSphere Admin", Password: "ChangeMe123!", CreatedAt: now}},
		orgs: []models.Organization{{ID: "org-demo", Name: "Acme Platform Engineering", Slug: "acme-platform", CreatedAt: now}},
		members: []models.OrganizationMember{{UserID: "user-admin", OrganizationID: "org-demo", Role: models.RolePlatformAdmin}},
		accounts: []models.CloudAccount{
			{ID: "acct-mock", OrganizationID: "org-demo", Name: "Mock Production", Provider: "mock", AccountID: "mock-prod", DefaultRegion: "us-west-2", Status: "connected", Tags: map[string]string{"env": "prod"}, CreatedAt: now},
			{ID: "acct-aws", OrganizationID: "org-demo", Name: "AWS Starter", Provider: "aws", AccountID: "123456789012", DefaultRegion: "us-east-1", Status: "needs_credentials", Tags: map[string]string{"owner": "platform"}, CreatedAt: now},
		},
	}
}

func (s *Store) Authenticate(email, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if strings.EqualFold(u.Email, email) && u.Password == password {
			cp := u
			return &cp, nil
		}
	}
	return nil, errors.New("invalid credentials")
}

func (s *Store) UserByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			cp := u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *Store) OrganizationsForUser(userID string) []models.Organization {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Organization(nil), s.orgs...)
}

func (s *Store) Accounts() []models.CloudAccount {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.CloudAccount(nil), s.accounts...)
}

func (s *Store) AddAccount(a models.CloudAccount) models.CloudAccount {
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = "acct-" + time.Now().UTC().Format("150405")
	a.OrganizationID = "org-demo"
	a.Status = "connected"
	a.CreatedAt = time.Now().UTC()
	s.accounts = append(s.accounts, a)
	return a
}

func (s *Store) Inventory() *models.Inventory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := *s.inventory
	return &cp
}

func (s *Store) Resource(id string) (*models.Resource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.inventory.Resources {
		if r.ID == id {
			cp := r
			return &cp, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (s *Store) Sync(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	inv, err := s.provider.ListInventory(ctx)
	if err != nil {
		return err
	}
	s.inventory = inv
	now := time.Now().UTC()
	for i := range s.accounts {
		if s.accounts[i].Provider == "mock" {
			s.accounts[i].LastSyncAt = &now
		}
	}
	s.auditLogs = append(s.auditLogs, map[string]interface{}{"action": "inventory.sync", "actor": "system", "createdAt": now})
	return nil
}

func (s *Store) AddDeployment(d models.Deployment) models.Deployment {
	s.mu.Lock()
	defer s.mu.Unlock()
	d.ID = "dep-" + time.Now().UTC().Format("20060102150405")
	d.OrganizationID = "org-demo"
	d.Status = "DRAFT"
	d.CreatedAt = time.Now().UTC()
	s.deployments = append(s.deployments, d)
	return d
}

func (s *Store) Deployments() []models.Deployment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Deployment(nil), s.deployments...)
}

func (s *Store) CostSummary() models.CostSummary {
	resources := s.Inventory().Resources
	byProvider := map[string]float64{}
	total := 0.0
	for _, r := range resources {
		byProvider[r.Provider] += r.CostMonthly
		total += r.CostMonthly
	}
	return models.CostSummary{
		CurrentMonth:     total,
		ForecastMonthEnd: total * 1.18,
		ByProvider:       byProvider,
		TopCostDrivers:   resources,
		Recommendations:  []string{"Rightsize payments-api compute; average CPU is below 25%.", "Review internet-facing load balancer policy.", "Evaluate committed-use discounts for steady database spend."},
	}
}

func (s *Store) SecurityFindings() []models.Finding {
	return []models.Finding{
		{ID: "finding-1", Severity: "critical", ResourceID: "res-lb-1", Provider: "mock", Title: "Internet-facing load balancer", Explanation: "Public ingress is enabled for a production workload.", Fix: "Restrict allowed CIDRs or front with WAF and verified TLS policy."},
		{ID: "finding-2", Severity: "medium", ResourceID: "res-eks-1", Provider: "mock", Title: "Kubernetes version aging", Explanation: "Cluster is behind the target enterprise baseline.", Fix: "Plan upgrade after checking add-on compatibility."},
	}
}

func (s *Store) Alerts() []models.Alert {
	now := time.Now().UTC()
	return []models.Alert{
		{ID: "alert-1", Severity: "critical", Title: "Elevated 5xx on payments-api", Source: "prometheus", Status: "firing", StartedAt: now.Add(-42 * time.Minute)},
		{ID: "alert-2", Severity: "warning", Title: "Database storage forecast exceeds 80%", Source: "cloudwatch", Status: "firing", StartedAt: now.Add(-3 * time.Hour)},
	}
}

func (s *Store) AuditLogs() []map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]map[string]interface{}(nil), s.auditLogs...)
}

