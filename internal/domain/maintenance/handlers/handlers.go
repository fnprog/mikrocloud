package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/domain/maintenance"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type MaintenanceHandler struct {
	deps *deps.Dependencies
}

func NewMaintenanceHandler(deps *deps.Dependencies) *MaintenanceHandler {
	return &MaintenanceHandler{deps: deps}
}

func (h *MaintenanceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	resp := maintenance.HealthCheckResponse{
		Status:    "ok",
		Service:   "mikrocloud",
		Version:   "0.1.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	utils.SendJSON(w, http.StatusOK, resp)
}

func (h *MaintenanceHandler) SystemStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	status := "healthy"
	components := struct {
		Database  string `json:"database"`
		Storage   string `json:"storage"`
		Container string `json:"container"`
	}{}

	if err := h.deps.DB.DB().PingContext(ctx); err != nil {
		components.Database = "error"
		status = "degraded"
	} else {
		components.Database = "ok"
	}

	components.Storage = "ok"

	if _, err := h.deps.ContainerService.ListContainers(ctx); err != nil {
		components.Container = "error"
		status = "degraded"
	} else {
		components.Container = "ok"
	}

	resp := map[string]interface{}{
		"status":     status,
		"components": components,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}
	utils.SendJSON(w, http.StatusOK, resp)
}

func (h *MaintenanceHandler) GetResources(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := map[string]int{
		"projects":     0,
		"applications": 0,
		"databases":    0,
		"services":     0,
	}

	orgID := middleware.GetOrgID(r)
	orgDomainID, err := users.OrganizationIDFromString(orgID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_org_id", "Invalid organization ID")
		return
	}

	projects, err := h.deps.ProjectService.ListProjects(ctx, orgDomainID)
	if err == nil {
		resp["projects"] = len(projects)
	}

	applications, err := h.deps.ApplicationService.ListApplications(ctx)
	if err == nil {
		resp["applications"] = len(applications)
	}

	// databases, err := h.deps.DatabaseService. ListAllWithContainers()
	// if err == nil {
	// 	resp["databases"] = len(databases)
	// }

	templates, err := h.deps.TemplateService.ListTemplates(ctx)
	if err == nil {
		resp["services"] = len(templates)
	}

	utils.SendJSON(w, http.StatusOK, resp)
}

func (h *MaintenanceHandler) SystemInfo(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"version":     "0.1.0",
		"platform":    runtime.GOOS + "/" + runtime.GOARCH,
		"go_version":  runtime.Version(),
		"environment": "production",
	}
	utils.SendJSON(w, http.StatusOK, resp)
}

func (h *MaintenanceHandler) ListDomains(w http.ResponseWriter, r *http.Request) {
	resp := maintenance.DomainListResponse{
		Domains: []maintenance.DomainInfo{},
		Total:   0,
	}
	utils.SendJSON(w, http.StatusOK, resp)
}

func (h *MaintenanceHandler) AddDomain(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, http.StatusNotImplemented, "not_implemented", "Domain management infrastructure not yet available. Database schema and repositories need to be created first.")
}

func (h *MaintenanceHandler) RemoveDomain(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, http.StatusNotImplemented, "not_implemented", "Domain management infrastructure not yet available. Database schema and repositories need to be created first.")
}

func (h *MaintenanceHandler) EnableSSL(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, http.StatusNotImplemented, "not_implemented", "SSL management infrastructure not yet available. Domain tables and certificate management need to be implemented first.")
}
