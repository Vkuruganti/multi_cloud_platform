package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/infrasphere/control-plane/backend/internal/auth"
	"github.com/infrasphere/control-plane/backend/internal/config"
	"github.com/infrasphere/control-plane/backend/internal/database"
	"github.com/infrasphere/control-plane/backend/internal/middleware"
	"github.com/infrasphere/control-plane/backend/internal/models"
	"github.com/infrasphere/control-plane/backend/internal/providers/mock"
	"github.com/infrasphere/control-plane/backend/pkg/cloud"
)

type Server struct {
	cfg   config.Config
	store *database.Store
}

func New(cfg config.Config, store *database.Store) http.Handler {
	s := &Server{cfg: cfg, store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/api/auth/login", s.login)
	mux.HandleFunc("/api/auth/logout", s.logout)
	mux.HandleFunc("/api/auth/refresh", s.refresh)
	mux.Handle("/api/me", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.me)))
	mux.Handle("/api/orgs", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.orgs)))
	mux.Handle("/api/providers", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.providers)))
	mux.Handle("/api/cloud-accounts", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.cloudAccounts)))
	mux.Handle("/api/cloud-accounts/", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.cloudAccountAction)))
	mux.Handle("/api/inventory/resources", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.resources)))
	mux.Handle("/api/inventory/resources/", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.resource)))
	mux.Handle("/api/inventory/relationships", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.relationships)))
	mux.Handle("/api/inventory/topology", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.topology)))
	mux.Handle("/api/deployments", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.deployments)))
	mux.Handle("/api/deployments/", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.deploymentAction)))
	mux.Handle("/api/observability/metrics", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.metrics)))
	mux.Handle("/api/observability/logs", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.logs)))
	mux.Handle("/api/observability/traces", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.traces)))
	mux.Handle("/api/observability/alerts", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.alerts)))
	mux.Handle("/api/ai/chat", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.aiChat)))
	mux.Handle("/api/ai/recommend", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.aiRecommend)))
	mux.Handle("/api/ai/explain-resource", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.aiExplain)))
	mux.Handle("/api/ai/triage-incident", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.aiTriage)))
	mux.Handle("/api/cost/summary", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.costSummary)))
	mux.Handle("/api/cost/resources", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.costResources)))
	mux.Handle("/api/cost/anomalies", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.costAnomalies)))
	mux.Handle("/api/cost/recommendations", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.costRecommendations)))
	mux.Handle("/api/security/findings", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.securityFindings)))
	mux.Handle("/api/audit-logs", middleware.Auth(cfg.JWTSecret)(http.HandlerFunc(s.auditLogs)))
	return middleware.CORS(cfg.CORSAllowedOrigins)(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) { write(w, map[string]string{"status": "ok"}) }

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	u, err := s.store.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, _ := auth.Sign(s.cfg.JWTSecret, auth.Claims{Sub: u.ID, Email: u.Email, Roles: []string{string(models.RolePlatformAdmin)}})
	write(w, map[string]interface{}{"token": token, "user": u})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	write(w, map[string]string{"status": "ok"})
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		http.Error(w, "missing bearer token", http.StatusUnauthorized)
		return
	}
	claims, err := auth.Verify(s.cfg.JWTSecret, strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	token, _ := auth.Sign(s.cfg.JWTSecret, *claims)
	write(w, map[string]string{"token": token})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	u, _ := s.store.UserByID(claims.Sub)
	write(w, map[string]interface{}{"user": u, "roles": claims.Roles, "organizations": s.store.OrganizationsForUser(claims.Sub)})
}

func (s *Server) orgs(w http.ResponseWriter, r *http.Request) { write(w, s.store.OrganizationsForUser("user-admin")) }
func (s *Server) providers(w http.ResponseWriter, r *http.Request) {
	write(w, []map[string]string{{"name": "mock", "status": "ready"}, {"name": "aws", "status": "starter"}, {"name": "gcp", "status": "starter"}, {"name": "azure", "status": "starter"}, {"name": "vcf", "status": "starter"}})
}

func (s *Server) cloudAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var a models.CloudAccount
		_ = json.NewDecoder(r.Body).Decode(&a)
		write(w, s.store.AddAccount(a))
		return
	}
	write(w, s.store.Accounts())
}

func (s *Server) cloudAccountAction(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/sync") {
		_ = s.store.Sync(r.Context())
		write(w, map[string]string{"status": "sync_started"})
		return
	}
	if strings.HasSuffix(r.URL.Path, "/validate") {
		write(w, map[string]string{"status": "valid"})
		return
	}
	http.NotFound(w, r)
}

func (s *Server) resources(w http.ResponseWriter, r *http.Request) { write(w, s.store.Inventory().Resources) }

func (s *Server) resource(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/inventory/resources/")
	res, err := s.store.Resource(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	write(w, res)
}

func (s *Server) relationships(w http.ResponseWriter, r *http.Request) { write(w, s.store.Inventory().Relationships) }
func (s *Server) topology(w http.ResponseWriter, r *http.Request)      { write(w, s.store.Inventory()) }

func (s *Server) deployments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var d models.Deployment
		_ = json.NewDecoder(r.Body).Decode(&d)
		write(w, s.store.AddDeployment(d))
		return
	}
	write(w, s.store.Deployments())
}

func (s *Server) deploymentAction(w http.ResponseWriter, r *http.Request) {
	status := "PLANNING"
	if strings.HasSuffix(r.URL.Path, "/approve") {
		status = "WAITING_FOR_APPROVAL"
	}
	if strings.HasSuffix(r.URL.Path, "/apply") {
		status = "DEPLOYING"
	}
	if strings.HasSuffix(r.URL.Path, "/rollback") {
		status = "ROLLING_BACK"
	}
	write(w, map[string]interface{}{"status": status, "events": []string{"request accepted", "audit log written", "worker queued"}})
}

func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	p := mock.New()
	m, _ := p.GetMetrics(r.Context(), cloud.MetricsQuery{Metric: "cpu_utilization"})
	write(w, m)
}
func (s *Server) logs(w http.ResponseWriter, r *http.Request)   { write(w, []string{"deployment validated", "target cluster reachable", "no live log source configured"}) }
func (s *Server) traces(w http.ResponseWriter, r *http.Request) { write(w, []map[string]string{{"traceId": "trace-demo", "service": "payments-api", "status": "ok"}}) }
func (s *Server) alerts(w http.ResponseWriter, r *http.Request) { write(w, s.store.Alerts()) }

func (s *Server) aiChat(w http.ResponseWriter, r *http.Request) {
	var req struct{ Message string `json:"message"` }
	_ = json.NewDecoder(r.Body).Decode(&req)
	write(w, map[string]interface{}{"answer": mock.ExampleQuestionAnswer(req.Message), "citations": []string{"inventory.resources", "observability.alerts", "cost.summary"}, "requiresApproval": false})
}
func (s *Server) aiRecommend(w http.ResponseWriter, r *http.Request) { write(w, map[string]interface{}{"recommendation": "Use Kubernetes on the provider closest to data gravity; enable logs, traces, budget guardrails, and approval gates.", "risk": "medium"}) }
func (s *Server) aiExplain(w http.ResponseWriter, r *http.Request)   { write(w, map[string]string{"explanation": "This resource is part of the payments path and has public exposure. Check security posture before changes."}) }
func (s *Server) aiTriage(w http.ResponseWriter, r *http.Request)    { write(w, map[string]string{"summary": "Likely deployment or traffic-shift related incident. Correlate recent rollout events with elevated 5xx."}) }

func (s *Server) costSummary(w http.ResponseWriter, r *http.Request)         { write(w, s.store.CostSummary()) }
func (s *Server) costResources(w http.ResponseWriter, r *http.Request)       { write(w, s.store.Inventory().Resources) }
func (s *Server) costAnomalies(w http.ResponseWriter, r *http.Request)       { write(w, []string{"payments database spend is 14% above seven-day baseline"}) }
func (s *Server) costRecommendations(w http.ResponseWriter, r *http.Request) { write(w, s.store.CostSummary().Recommendations) }
func (s *Server) securityFindings(w http.ResponseWriter, r *http.Request)    { write(w, s.store.SecurityFindings()) }
func (s *Server) auditLogs(w http.ResponseWriter, r *http.Request)           { write(w, s.store.AuditLogs()) }

func write(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
