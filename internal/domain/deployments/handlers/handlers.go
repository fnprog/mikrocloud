package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments/logs"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	"github.com/mikrocloud/mikrocloud/internal/utils"
	"github.com/mikrocloud/mikrocloud/pkg/containers"
)

type DeploymentHandler struct {
	deploymentService  *service.DeploymentService
	applicationService service.ApplicationService
}

func NewDeploymentHandler(deploymentService *service.DeploymentService, appService service.ApplicationService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService:  deploymentService,
		applicationService: appService,
	}
}

type CreateDeploymentRequest struct {
	ApplicationID    string                  `json:"application_id" validate:"required"`
	IsProduction     bool                    `json:"is_production"`
	TriggerType      deployments.TriggerType `json:"trigger_type" validate:"required"`
	ImageTag         string                  `json:"image_tag"`
	GitCommitHash    string                  `json:"git_commit_hash,omitempty"`
	GitCommitMessage string                  `json:"git_commit_message,omitempty"`
	GitBranch        string                  `json:"git_branch,omitempty"`
	GitAuthorName    string                  `json:"git_author_name,omitempty"`
}

type DeploymentResponse struct {
	ID                    string                       `json:"id"`
	ApplicationID         string                       `json:"application_id"`
	DeploymentNumber      int                          `json:"deployment_number"`
	Status                deployments.DeploymentStatus `json:"status"`
	IsProduction          bool                         `json:"is_production"`
	TriggeredBy           *string                      `json:"triggered_by,omitempty"`
	TriggeredByUsername   *string                      `json:"triggered_by_username,omitempty"`
	TriggerType           deployments.TriggerType      `json:"trigger_type"`
	ImageTag              string                       `json:"image_tag"`
	ImageDigest           string                       `json:"image_digest,omitempty"`
	ContainerID           string                       `json:"container_id,omitempty"`
	GitCommitHash         string                       `json:"git_commit_hash,omitempty"`
	GitCommitMessage      string                       `json:"git_commit_message,omitempty"`
	GitBranch             string                       `json:"git_branch,omitempty"`
	GitAuthorName         string                       `json:"git_author_name,omitempty"`
	BuildStartedAt        *string                      `json:"build_started_at,omitempty"`
	BuildCompletedAt      *string                      `json:"build_completed_at,omitempty"`
	BuildDurationSeconds  *int                         `json:"build_duration_seconds,omitempty"`
	DeployStartedAt       *string                      `json:"deploy_started_at,omitempty"`
	DeployCompletedAt     *string                      `json:"deploy_completed_at,omitempty"`
	DeployDurationSeconds *int                         `json:"deploy_duration_seconds,omitempty"`
	BuildLogs             string                       `json:"build_logs,omitempty"`
	DeployLogs            string                       `json:"deploy_logs,omitempty"`
	ErrorMessage          string                       `json:"error_message,omitempty"`
	CreatedAt             string                       `json:"created_at"`
	UpdatedAt             string                       `json:"updated_at"`
}

func (h *DeploymentHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context and convert to UserID
	userIDStr := middleware.GetUserID(r)
	if userIDStr == "" {
		utils.SendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
		return
	}

	userID, err := users.UserIDFromString(userIDStr)
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "invalid_user", "Invalid user ID")
		return
	}

	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	var req CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request_body", "Invalid request body")
		return
	}

	// Validate that the application ID in the request matches the URL
	if req.ApplicationID != applicationID.String() {
		utils.SendError(w, http.StatusBadRequest, "application_id_mismatch", "Application ID mismatch")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	// Create deployment command
	cmd := service.CreateDeploymentCommand{
		ApplicationID:    applicationID,
		IsProduction:     req.IsProduction,
		TriggeredBy:      &userID,
		TriggerType:      req.TriggerType,
		ImageTag:         req.ImageTag,
		GitCommitHash:    req.GitCommitHash,
		GitCommitMessage: req.GitCommitMessage,
		GitBranch:        req.GitBranch,
		GitAuthorName:    req.GitAuthorName,
	}

	// If no image tag provided, generate one
	if cmd.ImageTag == "" {
		cmd.ImageTag = generateImageTag(app, cmd.GitCommitHash)
	}

	// Create and execute deployment with build integration
	deployment, err := h.deploymentService.CreateAndExecuteDeployment(r.Context(), cmd, h.applicationService)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "deployment_creation_failed", "Failed to create deployment: "+err.Error())
		return
	}

	response := h.mapDeploymentToResponse(deployment)
	utils.SendJSON(w, http.StatusCreated, response)
}

func (h *DeploymentHandler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get deployment ID from URL path
	deploymentID, err := deployments.DeploymentIDFromString(chi.URLParam(r, "deployment_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_deployment_id", "Invalid deployment ID")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	deploymentWithMeta, err := h.deploymentService.GetDeploymentWithMetadata(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "deployment_not_found", "Deployment not found")
		return
	}

	if deploymentWithMeta.Deployment.ApplicationID() != applicationID {
		utils.SendError(w, http.StatusForbidden, "deployment_forbidden", "Deployment does not belong to this application")
		return
	}

	response := h.mapDeploymentWithMetadataToResponse(deploymentWithMeta)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *DeploymentHandler) ListDeployments(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	// Parse query parameters for pagination and filtering
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	status := r.URL.Query().Get("status")

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var deploymentsList []*repository.DeploymentWithMetadata

	if status != "" {
		utils.SendError(w, http.StatusNotImplemented, "status_filter_not_supported", "Status filtering not yet supported with metadata")
		return
	} else {
		deploymentsList, err = h.deploymentService.ListDeploymentsByApplicationWithMetadata(r.Context(), applicationID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "list_deployments_failed", "Failed to list deployments")
			return
		}
	}

	total := len(deploymentsList)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedDeployments := deploymentsList[start:end]
	responses := make([]DeploymentResponse, len(paginatedDeployments))
	for i, deploymentWithMeta := range paginatedDeployments {
		responses[i] = h.mapDeploymentWithMetadataToResponse(deploymentWithMeta)
	}

	result := map[string]interface{}{
		"deployments": responses,
		"pagination": map[string]interface{}{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	}

	utils.SendJSON(w, http.StatusOK, result)
}

func (h *DeploymentHandler) StopDeployment(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get deployment ID from URL path
	deploymentID, err := deployments.DeploymentIDFromString(chi.URLParam(r, "deployment_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_deployment_id", "Invalid deployment ID")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	// Verify deployment exists and belongs to the application
	deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "deployment_not_found", "Deployment not found")
		return
	}

	if deployment.ApplicationID() != applicationID {
		utils.SendError(w, http.StatusForbidden, "deployment_forbidden", "Deployment does not belong to this application")
		return
	}

	// Stop the deployment
	if err := h.deploymentService.StopDeployment(r.Context(), deploymentID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "stop_deployment_failed", "Failed to stop deployment: "+err.Error())
		return
	}

	// Return updated deployment
	updatedDeployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "get_deployment_failed", "Failed to get updated deployment")
		return
	}

	response := h.mapDeploymentToResponse(updatedDeployment)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *DeploymentHandler) CancelDeployment(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get deployment ID from URL path
	deploymentID, err := deployments.DeploymentIDFromString(chi.URLParam(r, "deployment_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_deployment_id", "Invalid deployment ID")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	// Verify deployment exists and belongs to the application
	deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "deployment_not_found", "Deployment not found")
		return
	}

	if deployment.ApplicationID() != applicationID {
		utils.SendError(w, http.StatusForbidden, "deployment_forbidden", "Deployment does not belong to this application")
		return
	}

	// Cancel the deployment
	if err := h.deploymentService.CancelDeployment(r.Context(), deploymentID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "cancel_deployment_failed", "Failed to cancel deployment: "+err.Error())
		return
	}

	// Return updated deployment
	updatedDeployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "get_deployment_failed", "Failed to get updated deployment")
		return
	}

	response := h.mapDeploymentToResponse(updatedDeployment)
	utils.SendJSON(w, http.StatusOK, response)
}

func parseLogsIfJSON(logText string) interface{} {
	if logText == "" {
		return ""
	}

	// Try to parse as JSON (structured logs)
	var structuredLogs []logs.StructuredLog
	if err := json.Unmarshal([]byte(logText), &structuredLogs); err == nil {
		return structuredLogs
	}

	// Return raw text for backward compatibility
	return logText
}

func mergeLogs(buildLogs, deployLogs interface{}) interface{} {
	buildStructured, buildIsStructured := buildLogs.([]logs.StructuredLog)
	deployStructured, deployIsStructured := deployLogs.([]logs.StructuredLog)

	// If both are structured, merge them with deduplication
	if buildIsStructured && deployIsStructured {
		seen := make(map[string]bool)
		var merged []logs.StructuredLog

		for _, log := range buildStructured {
			key := log.Message
			if !seen[key] {
				seen[key] = true
				merged = append(merged, log)
			}
		}

		for _, log := range deployStructured {
			key := log.Message
			if !seen[key] {
				seen[key] = true
				merged = append(merged, log)
			}
		}

		return merged
	}

	// If either is raw text, concatenate as strings
	buildStr, buildOk := buildLogs.(string)
	deployStr, deployOk := deployLogs.(string)

	if buildOk && deployOk {
		if buildStr != "" && deployStr != "" {
			return buildStr + "\n" + deployStr
		}
		if buildStr != "" {
			return buildStr
		}
		return deployStr
	}

	// Mixed types, return concatenated strings
	return ""
}

func (h *DeploymentHandler) GetDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	// Get project ID from URL path
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	// Get application ID from URL path
	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	// Get deployment ID from URL path
	deploymentID, err := deployments.DeploymentIDFromString(chi.URLParam(r, "deployment_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_deployment_id", "Invalid deployment ID")
		return
	}

	// Verify application exists and belongs to the project
	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "deployment_not_found", "Deployment not found")
		return
	}

	// Verify deployment belongs to the application
	if deployment.ApplicationID() != applicationID {
		utils.SendError(w, http.StatusForbidden, "deployment_forbidden", "Deployment does not belong to this application")
		return
	}

	// Get log type from query parameter
	logType := r.URL.Query().Get("type")
	if logType == "" {
		logType = "all"
	}

	var logs interface{}
	switch logType {
	case "build":
		logs = parseLogsIfJSON(deployment.BuildLogs())
	case "deploy":
		logs = parseLogsIfJSON(deployment.DeployLogs())
	case "all":
		buildLogs := parseLogsIfJSON(deployment.BuildLogs())
		deployLogs := parseLogsIfJSON(deployment.DeployLogs())
		logs = mergeLogs(buildLogs, deployLogs)
	default:
		utils.SendError(w, http.StatusBadRequest, "invalid_log_type", "Invalid log type. Use 'build', 'deploy', or 'all'")
		return
	}

	response := map[string]interface{}{
		"logs": logs,
		"type": logType,
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *DeploymentHandler) StreamDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
		return
	}

	applicationID, err := applications.ApplicationIDFromString(chi.URLParam(r, "application_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_application_id", "Invalid application ID")
		return
	}

	deploymentID, err := deployments.DeploymentIDFromString(chi.URLParam(r, "deployment_id"))
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_deployment_id", "Invalid deployment ID")
		return
	}

	app, err := h.applicationService.GetApplication(r.Context(), applicationID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "application_not_found", "Application not found")
		return
	}

	if app.ProjectID() != projectID {
		utils.SendError(w, http.StatusForbidden, "application_forbidden", "Application does not belong to this project")
		return
	}

	deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "deployment_not_found", "Deployment not found")
		return
	}

	if deployment.ApplicationID() != applicationID {
		utils.SendError(w, http.StatusForbidden, "deployment_forbidden", "Deployment does not belong to this application")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.SendError(w, http.StatusInternalServerError, "streaming_not_supported", "Streaming not supported")
		return
	}

	parser := logs.NewLogParser(deployment.StartedAt())
	lastLogLength := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
			if err != nil {
				return
			}

			currentLogs := deployment.BuildLogs()
			if len(currentLogs) > lastLogLength {
				newLogs := currentLogs[lastLogLength:]
				for _, line := range strings.Split(newLogs, "\n") {
					if line != "" {
						structuredLog := parser.ParseLine(line)
						if structuredLog != nil {
							logJSON, err := json.Marshal(structuredLog)
							if err == nil {
								_, _ = w.Write([]byte("data: " + string(logJSON) + "\n\n"))
								flusher.Flush()
							}
						}
					}
				}
				lastLogLength = len(currentLogs)
			}

			if deployment.Status() != deployments.DeploymentStatusBuilding &&
				deployment.Status() != deployments.DeploymentStatusPending &&
				deployment.Status() != deployments.DeploymentStatusQueued {
				_, _ = w.Write([]byte("event: done\ndata: {\"status\": \"" + string(deployment.Status()) + "\"}\n\n"))
				flusher.Flush()
				return
			}
		}
	}
}

func (h *DeploymentHandler) mapDeploymentToResponse(deployment *deployments.Deployment) DeploymentResponse {
	response := DeploymentResponse{
		ID:               deployment.ID().String(),
		ApplicationID:    deployment.ApplicationID().String(),
		DeploymentNumber: deployment.DeploymentNumber(),
		Status:           deployment.Status(),
		IsProduction:     deployment.IsProduction(),
		TriggerType:      deployment.TriggerType(),
		ImageTag:         deployment.ImageTag(),
		ImageDigest:      deployment.ImageDigest(),
		ContainerID:      deployment.ContainerID(),
		GitCommitHash:    deployment.GitCommitHash(),
		GitCommitMessage: deployment.GitCommitMessage(),
		GitBranch:        deployment.GitBranch(),
		GitAuthorName:    deployment.GitAuthorName(),
		BuildLogs:        deployment.BuildLogs(),
		DeployLogs:       deployment.DeployLogs(),
		ErrorMessage:     deployment.ErrorMessage(),
		CreatedAt:        deployment.StartedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        deployment.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	if deployment.TriggeredBy() != nil {
		triggeredBy := deployment.TriggeredBy().String()
		response.TriggeredBy = &triggeredBy
	}

	if deployment.BuildStartedAt() != nil {
		buildStartedAt := deployment.BuildStartedAt().Format("2006-01-02T15:04:05Z07:00")
		response.BuildStartedAt = &buildStartedAt
	}

	if deployment.BuildCompletedAt() != nil {
		buildCompletedAt := deployment.BuildCompletedAt().Format("2006-01-02T15:04:05Z07:00")
		response.BuildCompletedAt = &buildCompletedAt
	}

	if deployment.DeployStartedAt() != nil {
		deployStartedAt := deployment.DeployStartedAt().Format("2006-01-02T15:04:05Z07:00")
		response.DeployStartedAt = &deployStartedAt
	}

	if deployment.DeployCompletedAt() != nil {
		deployCompletedAt := deployment.DeployCompletedAt().Format("2006-01-02T15:04:05Z07:00")
		response.DeployCompletedAt = &deployCompletedAt
	}

	response.BuildDurationSeconds = deployment.BuildDurationSeconds()
	response.DeployDurationSeconds = deployment.DeployDurationSeconds()

	return response
}

func (h *DeploymentHandler) mapDeploymentWithMetadataToResponse(deploymentWithMeta *repository.DeploymentWithMetadata) DeploymentResponse {
	response := h.mapDeploymentToResponse(deploymentWithMeta.Deployment)
	response.TriggeredByUsername = deploymentWithMeta.Username
	return response
}

func generateImageTag(app *applications.Application, gitCommitHash string) string {
	imageName := containers.SanitizeDockerName(app.Name().String())
	if gitCommitHash != "" {
		if len(gitCommitHash) > 7 {
			return imageName + ":" + gitCommitHash[:7]
		}
		return imageName + ":" + gitCommitHash
	}
	return imageName + ":latest"
}
