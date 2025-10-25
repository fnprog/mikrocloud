package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type TunnelHandler struct {
	tunnelService *service.TunnelService
	validator     *validator.Validate
}

func NewTunnelHandler(tunnelService *service.TunnelService) *TunnelHandler {
	return &TunnelHandler{
		tunnelService: tunnelService,
		validator:     validator.New(),
	}
}

type TunnelResponse struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	ProjectID      *string              `json:"project_id,omitempty"`
	OrganizationID string               `json:"organization_id"`
	Status         tunnels.TunnelStatus `json:"status"`
	HealthStatus   tunnels.HealthStatus `json:"health_status"`
	ContainerID    string               `json:"container_id"`
	LastError      string               `json:"last_error,omitempty"`
	CreatedAt      string               `json:"created_at"`
	UpdatedAt      string               `json:"updated_at"`
}

type CreateTunnelRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	ProjectID   *string `json:"project_id,omitempty" validate:"omitempty,uuid"`
	TunnelToken string  `json:"tunnel_token" validate:"required"`
}

type ListTunnelsResponse struct {
	Tunnels []TunnelResponse `json:"tunnels"`
}

func (h *TunnelHandler) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	var req CreateTunnelRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON format")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

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

	orgIDStr := middleware.GetOrgID(r)
	if orgIDStr == "" {
		utils.SendError(w, http.StatusBadRequest, "missing_organization", "Organization ID required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_organization_id", "Invalid organization ID")
		return
	}

	var projectID *uuid.UUID
	if req.ProjectID != nil {
		pid, err := uuid.Parse(*req.ProjectID)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
			return
		}
		projectID = &pid
	}

	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_user_id", "Invalid user ID format")
		return
	}

	cmd := service.CreateTunnelRequest{
		Name:           req.Name,
		ProjectID:      projectID,
		OrganizationID: orgID,
		TunnelToken:    req.TunnelToken,
		CreatedBy:      userUUID,
	}

	tunnel, err := h.tunnelService.CreateTunnel(r.Context(), cmd)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "create_failed", "Failed to create tunnel: "+err.Error())
		return
	}

	response := h.tunnelToResponse(tunnel)
	utils.SendJSON(w, http.StatusCreated, response)
}

func (h *TunnelHandler) GetTunnel(w http.ResponseWriter, r *http.Request) {
	tunnelIDStr := chi.URLParam(r, "tunnel_id")
	tunnelID, err := tunnels.TunnelIDFromString(tunnelIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_tunnel_id", "Invalid tunnel ID")
		return
	}

	tunnel, err := h.tunnelService.GetTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "tunnel_not_found", "Tunnel not found")
		return
	}

	response := h.tunnelToResponse(tunnel)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *TunnelHandler) ListTunnels(w http.ResponseWriter, r *http.Request) {
	orgIDStr := middleware.GetOrgID(r)
	if orgIDStr == "" {
		utils.SendError(w, http.StatusBadRequest, "missing_organization", "Organization ID required")
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_organization_id", "Invalid organization ID")
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	var tunnelList []*tunnels.CloudflareTunnel

	if projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "invalid_project_id", "Invalid project ID")
			return
		}
		tunnelList, err = h.tunnelService.ListTunnelsByProject(r.Context(), projectID)
	} else {
		tunnelList, err = h.tunnelService.ListTunnelsByOrganization(r.Context(), orgID)
	}

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "list_failed", "Failed to list tunnels: "+err.Error())
		return
	}

	responses := make([]TunnelResponse, len(tunnelList))
	for i, tun := range tunnelList {
		responses[i] = h.tunnelToResponse(tun)
	}

	result := ListTunnelsResponse{
		Tunnels: responses,
	}

	utils.SendJSON(w, http.StatusOK, result)
}

func (h *TunnelHandler) StartTunnel(w http.ResponseWriter, r *http.Request) {
	tunnelIDStr := chi.URLParam(r, "tunnel_id")
	tunnelID, err := tunnels.TunnelIDFromString(tunnelIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_tunnel_id", "Invalid tunnel ID")
		return
	}

	err = h.tunnelService.StartTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "start_failed", "Failed to start tunnel: "+err.Error())
		return
	}

	tunnel, err := h.tunnelService.GetTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "fetch_failed", "Failed to fetch tunnel")
		return
	}

	response := h.tunnelToResponse(tunnel)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *TunnelHandler) StopTunnel(w http.ResponseWriter, r *http.Request) {
	tunnelIDStr := chi.URLParam(r, "tunnel_id")
	tunnelID, err := tunnels.TunnelIDFromString(tunnelIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_tunnel_id", "Invalid tunnel ID")
		return
	}

	err = h.tunnelService.StopTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "stop_failed", "Failed to stop tunnel: "+err.Error())
		return
	}

	tunnel, err := h.tunnelService.GetTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "fetch_failed", "Failed to fetch tunnel")
		return
	}

	response := h.tunnelToResponse(tunnel)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *TunnelHandler) RestartTunnel(w http.ResponseWriter, r *http.Request) {
	tunnelIDStr := chi.URLParam(r, "tunnel_id")
	tunnelID, err := tunnels.TunnelIDFromString(tunnelIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_tunnel_id", "Invalid tunnel ID")
		return
	}

	err = h.tunnelService.RestartTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "restart_failed", "Failed to restart tunnel: "+err.Error())
		return
	}

	tunnel, err := h.tunnelService.GetTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "fetch_failed", "Failed to fetch tunnel")
		return
	}

	response := h.tunnelToResponse(tunnel)
	utils.SendJSON(w, http.StatusOK, response)
}

func (h *TunnelHandler) DeleteTunnel(w http.ResponseWriter, r *http.Request) {
	tunnelIDStr := chi.URLParam(r, "tunnel_id")
	tunnelID, err := tunnels.TunnelIDFromString(tunnelIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_tunnel_id", "Invalid tunnel ID")
		return
	}

	err = h.tunnelService.DeleteTunnel(r.Context(), tunnelID)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "delete_failed", "Failed to delete tunnel: "+err.Error())
		return
	}

	response := map[string]string{
		"message": "Tunnel deleted successfully",
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *TunnelHandler) tunnelToResponse(tun *tunnels.CloudflareTunnel) TunnelResponse {
	var projectID *string
	if tun.ProjectID() != nil {
		pid := tun.ProjectID().String()
		projectID = &pid
	}

	return TunnelResponse{
		ID:             tun.ID().String(),
		Name:           tun.Name().String(),
		ProjectID:      projectID,
		OrganizationID: tun.OrganizationID().String(),
		Status:         tun.Status(),
		HealthStatus:   tun.HealthStatus(),
		ContainerID:    tun.ContainerID(),
		LastError:      tun.ErrorMessage(),
		CreatedAt:      tun.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      tun.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}
