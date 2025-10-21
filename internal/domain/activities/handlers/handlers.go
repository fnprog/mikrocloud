package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/activities/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type ActivitiesHandlers struct {
	service *service.ActivitiesService
}

func NewActivitiesHandlers(service *service.ActivitiesService) *ActivitiesHandlers {
	return &ActivitiesHandlers{service: service}
}

func (h *ActivitiesHandlers) GetRecentActivities(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusUnauthorized, "unauthorized", "Invalid token")
		return
	}

	orgIDStr, ok := claims["org_id"].(string)

	if !ok || orgIDStr == "" {
		userIDStr, _ := claims["user_id"].(string)
		orgIDStr = userIDStr
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid organization ID", err.Error())
		return
	}

	limit := 50

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	activities, err := h.service.GetRecentActivities(orgID, limit, offset)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch activities", err.Error())
		return
	}

	response := map[string]any{
		"activities": activities,
		"total":      len(activities),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *ActivitiesHandlers) GetResourceActivities(w http.ResponseWriter, r *http.Request) {
	resourceType := chi.URLParam(r, "resource_type")
	resourceIDStr := chi.URLParam(r, "resource_id")

	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid resource ID", err.Error())
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	activities, err := h.service.GetResourceActivities(resourceType, resourceID, limit)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to fetch resource activities", err.Error())
		return
	}

	response := map[string]any{
		"activities": activities,
		"total":      len(activities),
	}

	utils.SendJSON(w, http.StatusOK, response)
}
