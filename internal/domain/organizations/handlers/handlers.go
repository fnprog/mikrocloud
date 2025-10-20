package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/domain/organizations/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

type OrganizationResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	OwnerID      string `json:"owner_id"`
	BillingEmail string `json:"billing_email"`
	Plan         string `json:"plan"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (h *OrganizationHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	organizations, err := h.orgService.ListOrganizations(r.Context())
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "list_failed", "Failed to list organizations: "+err.Error())
		return
	}

	response := make([]OrganizationResponse, 0, len(organizations))
	for _, org := range organizations {
		response = append(response, OrganizationResponse{
			ID:           org.ID().String(),
			Name:         org.Name(),
			Slug:         org.Slug(),
			Description:  org.Description(),
			OwnerID:      org.OwnerID().String(),
			BillingEmail: org.BillingEmail(),
			Plan:         string(org.Plan()),
			Status:       string(org.Status()),
			CreatedAt:    org.CreatedAt().Format("2006-01-02T15:04:05Z"),
			UpdatedAt:    org.UpdatedAt().Format("2006-01-02T15:04:05Z"),
		})
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *OrganizationHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	org, err := h.orgService.GetOrganization(r.Context(), orgID)
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "not_found", "Organization not found: "+err.Error())
		return
	}

	response := OrganizationResponse{
		ID:           org.ID().String(),
		Name:         org.Name(),
		Slug:         org.Slug(),
		Description:  org.Description(),
		OwnerID:      org.OwnerID().String(),
		BillingEmail: org.BillingEmail(),
		Plan:         string(org.Plan()),
		Status:       string(org.Status()),
		CreatedAt:    org.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    org.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

type CreateOrganizationRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	BillingEmail string `json:"billing_email"`
}

func (h *OrganizationHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var req CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "Invalid request body: "+err.Error())
		return
	}

	if req.Name == "" {
		utils.SendError(w, http.StatusBadRequest, "validation_error", "Organization name is required")
		return
	}

	if req.Slug == "" {
		utils.SendError(w, http.StatusBadRequest, "validation_error", "Organization slug is required")
		return
	}

	userID := r.Context().Value("user_id").(string)

	org, err := h.orgService.CreateOrganization(r.Context(), req.Name, req.Slug, req.Description, req.BillingEmail, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "creation_failed", "Failed to create organization: "+err.Error())
		return
	}

	response := OrganizationResponse{
		ID:           org.ID().String(),
		Name:         org.Name(),
		Slug:         org.Slug(),
		Description:  org.Description(),
		OwnerID:      org.OwnerID().String(),
		BillingEmail: org.BillingEmail(),
		Plan:         string(org.Plan()),
		Status:       string(org.Status()),
		CreatedAt:    org.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    org.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	utils.SendJSON(w, http.StatusCreated, response)
}

type UpdateOrganizationRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	BillingEmail string `json:"billing_email"`
}

func (h *OrganizationHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	var req UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "Invalid request body: "+err.Error())
		return
	}

	org, err := h.orgService.UpdateOrganization(r.Context(), orgID, req.Name, req.Description, req.BillingEmail)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "update_failed", "Failed to update organization: "+err.Error())
		return
	}

	response := OrganizationResponse{
		ID:           org.ID().String(),
		Name:         org.Name(),
		Slug:         org.Slug(),
		Description:  org.Description(),
		OwnerID:      org.OwnerID().String(),
		BillingEmail: org.BillingEmail(),
		Plan:         string(org.Plan()),
		Status:       string(org.Status()),
		CreatedAt:    org.CreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    org.UpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *OrganizationHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	if err := h.orgService.DeleteOrganization(r.Context(), orgID); err != nil {
		if err.Error() == "cannot delete the last organization" || err.Error() == "cannot delete the default organization" {
			utils.SendError(w, http.StatusBadRequest, "deletion_not_allowed", err.Error())
			return
		}
		utils.SendError(w, http.StatusInternalServerError, "deletion_failed", "Failed to delete organization: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type OrganizationMemberResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	UserEmail      string `json:"user_email"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	InvitedBy      string `json:"invited_by,omitempty"`
	InvitedAt      string `json:"invited_at,omitempty"`
	JoinedAt       string `json:"joined_at,omitempty"`
}

func (h *OrganizationHandler) ListOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	membersWithUsers, err := h.orgService.ListOrganizationMembers(r.Context(), orgID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "list_failed", "Failed to list organization members: "+err.Error())
		return
	}

	response := make([]OrganizationMemberResponse, 0, len(membersWithUsers))
	for _, mwu := range membersWithUsers {
		memberResp := OrganizationMemberResponse{
			ID:             mwu.Member.ID().String(),
			OrganizationID: mwu.Member.OrganizationID().String(),
			UserID:         mwu.Member.UserID().String(),
			UserName:       mwu.User.Name(),
			UserEmail:      mwu.User.Email().String(),
			Role:           string(mwu.Member.Role()),
			Status:         string(mwu.Member.Status()),
		}

		if mwu.Member.InvitedBy() != nil {
			memberResp.InvitedBy = mwu.Member.InvitedBy().String()
		}
		if mwu.Member.InvitedAt() != nil {
			memberResp.InvitedAt = mwu.Member.InvitedAt().Format("2006-01-02T15:04:05Z")
		}
		if mwu.Member.JoinedAt() != nil {
			memberResp.JoinedAt = mwu.Member.JoinedAt().Format("2006-01-02T15:04:05Z")
		}

		response = append(response, memberResp)
	}

	utils.SendJSON(w, http.StatusOK, response)
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *OrganizationHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	var req InviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "Invalid request body: "+err.Error())
		return
	}

	if req.Email == "" {
		utils.SendError(w, http.StatusBadRequest, "validation_error", "Email is required")
		return
	}

	if req.Role == "" {
		req.Role = "member"
	}

	userID := r.Context().Value("user_id").(string)

	err := h.orgService.InviteMember(r.Context(), orgID, req.Email, req.Role, userID)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "invite_failed", "Failed to invite member: "+err.Error())
		return
	}

	utils.SendJSON(w, http.StatusCreated, map[string]string{"message": "Member invited successfully"})
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role"`
}

func (h *OrganizationHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	memberID := chi.URLParam(r, "member_id")

	var req UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "Invalid request body: "+err.Error())
		return
	}

	if req.Role == "" {
		utils.SendError(w, http.StatusBadRequest, "validation_error", "Role is required")
		return
	}

	err := h.orgService.UpdateMemberRole(r.Context(), memberID, req.Role)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "update_failed", "Failed to update member role: "+err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Member role updated successfully"})
}

func (h *OrganizationHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	memberID := chi.URLParam(r, "member_id")

	if err := h.orgService.RemoveMember(r.Context(), memberID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "removal_failed", "Failed to remove member: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type TransferOwnershipRequest struct {
	NewOwnerID string `json:"new_owner_id"`
}

func (h *OrganizationHandler) TransferOwnership(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "organization_id")

	var req TransferOwnershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "Invalid request body: "+err.Error())
		return
	}

	if req.NewOwnerID == "" {
		utils.SendError(w, http.StatusBadRequest, "validation_error", "New owner ID is required")
		return
	}

	if err := h.orgService.TransferOwnership(r.Context(), orgID, req.NewOwnerID); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "transfer_failed", "Failed to transfer ownership: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
