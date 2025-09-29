package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
	"github.com/mikrocloud/mikrocloud/internal/domain/applications/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type ApplicationHandler struct {
	appService *service.ApplicationService
	validator  *validator.Validate
}

func NewApplicationHandler(appService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
		validator:  validator.New(),
	}
}

// ApplicationResponse represents an application in API responses
type ApplicationResponse struct {
	ID               string                         `json:"id"`
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	ProjectID        string                         `json:"project_id"`
	EnvironmentID    string                         `json:"environment_id"`
	DeploymentSource applications.DeploymentSource  `json:"deployment_source"`
	Domain           string                         `json:"domain"`
	Buildpack        applications.BuildpackConfig   `json:"buildpack"`
	EnvVars          map[string]string              `json:"env_vars"`
	AutoDeploy       bool                           `json:"auto_deploy"`
	Status           applications.ApplicationStatus `json:"status"`
	CreatedAt        string                         `json:"created_at"`
	UpdatedAt        string                         `json:"updated_at"`
}

type CreateApplicationRequest struct {
	Name             string                        `json:"name" validate:"required,min=1,max=100"`
	Description      string                        `json:"description,omitempty"`
	EnvironmentID    string                        `json:"environment_id" validate:"required,uuid"`
	DeploymentSource applications.DeploymentSource `json:"deployment_source" validate:"required"`
	Buildpack        applications.BuildpackConfig  `json:"buildpack" validate:"required"`
	EnvVars          map[string]string             `json:"env_vars,omitempty"`
}

type UpdateApplicationRequest struct {
	Description      *string                        `json:"description,omitempty"`
	DeploymentSource *applications.DeploymentSource `json:"deployment_source,omitempty"`
	Domain           *string                        `json:"domain,omitempty"`
	Buildpack        *applications.BuildpackConfig  `json:"buildpack,omitempty"`
	EnvVars          map[string]string              `json:"env_vars,omitempty"`
	AutoDeploy       *bool                          `json:"auto_deploy,omitempty"`
}

type ApplicationListItem struct {
	ID            string                         `json:"id"`
	Name          string                         `json:"name"`
	Description   string                         `json:"description"`
	ProjectID     string                         `json:"project_id"`
	EnvironmentID string                         `json:"environment_id"`
	Domain        string                         `json:"domain"`
	Status        applications.ApplicationStatus `json:"status"`
	CreatedAt     string                         `json:"created_at"`
}

type ListApplicationsResponse struct {
	Applications []ApplicationListItem `json:"applications"`
}

type DeployApplicationRequest struct {
	Action string `json:"action" validate:"required,oneof=deploy stop"`
}

// CreateApplication creates a new application in a project
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var req CreateApplicationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Parse environment ID
	environmentID, err := uuid.Parse(req.EnvironmentID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_environment_id", "Invalid environment ID")
		return
	}

	cmd := service.CreateApplicationCommand{
		Name:             req.Name,
		Description:      req.Description,
		ProjectID:        projectID,
		EnvironmentID:    environmentID,
		DeploymentSource: req.DeploymentSource,
		BuildpackConfig:  req.Buildpack,
		EnvVars:          req.EnvVars,
	}

	app, err := h.appService.CreateApplication(r.Context(), cmd)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "create_failed", "Failed to create application: "+err.Error())
		return
	}

	response := ApplicationResponse{
		ID:               app.ID().String(),
		Name:             app.Name().String(),
		Description:      app.Description(),
		ProjectID:        app.ProjectID().String(),
		EnvironmentID:    app.EnvironmentID().String(),
		DeploymentSource: app.DeploymentSource(),
		Domain:           app.Domain(),
		Buildpack:        app.Buildpack(),
		EnvVars:          app.EnvVars(),
		AutoDeploy:       app.AutoDeploy(),
		Status:           app.Status(),
		CreatedAt:        app.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        app.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.SendJSON(w, http.StatusCreated, response)
}

// GetApplication retrieves a specific application
func (h *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	appIDStr := chi.URLParam(r, "application_id")
	appID, err := applications.ApplicationIDFromString(appIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	app, err := h.appService.GetApplication(r.Context(), appID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	// Verify application belongs to the specified project
	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found in project")
		return
	}

	response := ApplicationResponse{
		ID:               app.ID().String(),
		Name:             app.Name().String(),
		Description:      app.Description(),
		ProjectID:        app.ProjectID().String(),
		EnvironmentID:    app.EnvironmentID().String(),
		DeploymentSource: app.DeploymentSource(),
		Domain:           app.Domain(),
		Buildpack:        app.Buildpack(),
		EnvVars:          app.EnvVars(),
		AutoDeploy:       app.AutoDeploy(),
		Status:           app.Status(),
		CreatedAt:        app.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        app.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

// ListApplications lists all applications for a project
func (h *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	apps, err := h.appService.ListApplicationsByProject(r.Context(), projectID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "list_failed", "Failed to list applications: "+err.Error())
		return
	}

	items := make([]ApplicationListItem, len(apps))
	for i, app := range apps {
		items[i] = ApplicationListItem{
			ID:            app.ID().String(),
			Name:          app.Name().String(),
			Description:   app.Description(),
			ProjectID:     app.ProjectID().String(),
			EnvironmentID: app.EnvironmentID().String(),
			Domain:        app.Domain(),
			Status:        app.Status(),
			CreatedAt:     app.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListApplicationsResponse{
		Applications: items,
	}

	utils.SendJSON(w, http.StatusOK, response)
}

// UpdateApplication updates an existing application
func (h *ApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	var req UpdateApplicationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Get application ID from URL
	appIDStr := chi.URLParam(r, "application_id")
	appID, err := applications.ApplicationIDFromString(appIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// First verify the application exists and belongs to the project
	app, err := h.appService.GetApplication(r.Context(), appID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found in project")
		return
	}

	cmd := service.UpdateApplicationCommand{
		ID:               appID,
		Description:      req.Description,
		DeploymentSource: req.DeploymentSource,
		Domain:           req.Domain,
		BuildpackConfig:  req.Buildpack,
		EnvVars:          req.EnvVars,
		AutoDeploy:       req.AutoDeploy,
	}

	updatedApp, err := h.appService.UpdateApplication(r.Context(), cmd)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "update_failed", "Failed to update application: "+err.Error())
		return
	}

	response := ApplicationResponse{
		ID:               updatedApp.ID().String(),
		Name:             updatedApp.Name().String(),
		Description:      updatedApp.Description(),
		ProjectID:        updatedApp.ProjectID().String(),
		EnvironmentID:    updatedApp.EnvironmentID().String(),
		DeploymentSource: updatedApp.DeploymentSource(),
		Domain:           updatedApp.Domain(),
		Buildpack:        updatedApp.Buildpack(),
		EnvVars:          updatedApp.EnvVars(),
		AutoDeploy:       updatedApp.AutoDeploy(),
		Status:           updatedApp.Status(),
		CreatedAt:        updatedApp.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        updatedApp.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

// DeleteApplication deletes an application
func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	// Get application ID from URL
	appIDStr := chi.URLParam(r, "application_id")
	appID, err := applications.ApplicationIDFromString(appIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// First verify the application exists and belongs to the project
	app, err := h.appService.GetApplication(r.Context(), appID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found in project")
		return
	}

	err = h.appService.DeleteApplication(r.Context(), appID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "delete_failed", "Failed to delete application: "+err.Error())
		return
	}

	response := map[string]string{
		"message": "Application deleted successfully",
	}

	utils.SendJSON(w, http.StatusOK, response)
}

// DeployApplication handles application deployment actions
func (h *ApplicationHandler) DeployApplication(w http.ResponseWriter, r *http.Request) {
	var req DeployApplicationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	// Get application ID from URL
	appIDStr := chi.URLParam(r, "application_id")
	appID, err := applications.ApplicationIDFromString(appIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get project ID from URL
	projectIDStr := chi.URLParam(r, "project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// First verify the application exists and belongs to the project
	app, err := h.appService.GetApplication(r.Context(), appID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found in project")
		return
	}

	switch req.Action {
	case "deploy":
		err = h.appService.StartDeployment(r.Context(), appID)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "deploy_failed", "Failed to start deployment: "+err.Error())
			return
		}
	case "stop":
		err = h.appService.StopApplication(r.Context(), appID)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "stop_failed", "Failed to stop application: "+err.Error())
			return
		}
	default:
		utils.SendError(w, http.StatusBadRequest, "invalid_action", "Invalid deployment action")
		return
	}

	response := map[string]string{
		"message": "Deployment action completed successfully",
		"action":  req.Action,
	}

	utils.SendJSON(w, http.StatusOK, response)
}
