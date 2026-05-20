package models

import "time"

type Role string

const (
	RolePlatformAdmin   Role = "platform_admin"
	RoleCloudAdmin      Role = "cloud_admin"
	RoleSRE             Role = "sre"
	RoleDeveloper       Role = "developer"
	RoleFinOpsAnalyst   Role = "finops_analyst"
	RoleSecurityAnalyst Role = "security_analyst"
	RoleAuditor         Role = "read_only_auditor"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrganizationMember struct {
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
	Role           Role   `json:"role"`
}

type CloudAccount struct {
	ID             string            `json:"id"`
	OrganizationID string            `json:"organizationId"`
	Name           string            `json:"name"`
	Provider       string            `json:"provider"`
	AccountID      string            `json:"accountId"`
	DefaultRegion  string            `json:"defaultRegion"`
	Status         string            `json:"status"`
	Tags           map[string]string `json:"tags"`
	LastSyncAt     *time.Time        `json:"lastSyncAt,omitempty"`
	CreatedAt      time.Time         `json:"createdAt"`
}

type Resource struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Provider     string                 `json:"provider"`
	AccountID    string                 `json:"accountId"`
	Region       string                 `json:"region"`
	Zone         string                 `json:"zone"`
	ResourceType string                 `json:"resourceType"`
	Category     string                 `json:"category"`
	Status       string                 `json:"status"`
	Health       string                 `json:"health"`
	CostMonthly  float64                `json:"costMonthly"`
	Tags         map[string]string      `json:"tags"`
	Metadata     map[string]interface{} `json:"metadata"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type ResourceRelationship struct {
	ID           string `json:"id"`
	SourceID     string `json:"sourceId"`
	TargetID     string `json:"targetId"`
	Relationship string `json:"relationship"`
}

type Deployment struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Provider       string    `json:"provider"`
	Target         string    `json:"target"`
	WorkloadType   string    `json:"workloadType"`
	Status         string    `json:"status"`
	CostEstimate   float64   `json:"costEstimate"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Finding struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	ResourceID  string `json:"resourceId"`
	Provider    string `json:"provider"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Fix         string `json:"fix"`
}

type Alert struct {
	ID        string    `json:"id"`
	Severity  string    `json:"severity"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"startedAt"`
}

type CostSummary struct {
	CurrentMonth      float64            `json:"currentMonth"`
	ForecastMonthEnd  float64            `json:"forecastMonthEnd"`
	ByProvider        map[string]float64 `json:"byProvider"`
	TopCostDrivers    []Resource         `json:"topCostDrivers"`
	Recommendations   []string           `json:"recommendations"`
}

type Inventory struct {
	Resources     []Resource             `json:"resources"`
	Relationships []ResourceRelationship `json:"relationships"`
}

