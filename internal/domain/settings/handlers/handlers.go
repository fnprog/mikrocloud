package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mikrocloud/mikrocloud/internal/domain/settings"
	"github.com/mikrocloud/mikrocloud/internal/domain/settings/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type SettingsHandler struct {
	service   *service.SettingsService
	validator *validator.Validate
}

func NewSettingsHandler(svc *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		service:   svc,
		validator: validator.New(),
	}
}

func (h *SettingsHandler) GetGeneralSettings(w http.ResponseWriter, r *http.Request) {
	generalSettings, err := h.service.GetGeneralSettings()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, generalSettings)
}

func (h *SettingsHandler) SaveGeneralSettings(w http.ResponseWriter, r *http.Request) {
	var input settings.GeneralSettings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveGeneralSettings(&input); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Settings saved successfully"})
}

func (h *SettingsHandler) GetAdvancedSettings(w http.ResponseWriter, r *http.Request) {
	advancedSettings, err := h.service.GetAdvancedSettings()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, advancedSettings)
}

func (h *SettingsHandler) SaveAdvancedSettings(w http.ResponseWriter, r *http.Request) {
	var input settings.AdvancedSettings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveAdvancedSettings(&input); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Settings saved successfully"})
}

func (h *SettingsHandler) GetUpdateSettings(w http.ResponseWriter, r *http.Request) {
	updateSettings, err := h.service.GetUpdateSettings()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, updateSettings)
}

func (h *SettingsHandler) SaveUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var input settings.UpdateSettings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveUpdateSettings(&input); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Settings saved successfully"})
}

func (h *SettingsHandler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, http.StatusNotImplemented, "Backup feature not implemented yet", "")
}

func (h *SettingsHandler) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, http.StatusNotImplemented, "Restore feature not implemented yet", "")
}

func (h *SettingsHandler) GetInstanceInfo(w http.ResponseWriter, r *http.Request) {
	instanceInfo, err := h.service.GetInstanceInfo()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve instance info", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, instanceInfo)
}

func (h *SettingsHandler) GetSMTPSettings(w http.ResponseWriter, r *http.Request) {
	smtpSettings, err := h.service.GetSMTPSettings()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to retrieve SMTP settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, smtpSettings)
}

func (h *SettingsHandler) SaveSMTPSettings(w http.ResponseWriter, r *http.Request) {
	var input settings.UpdateSMTPSettings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.service.SaveSMTPSettings(&input); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to save SMTP settings", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "SMTP settings saved successfully"})
}

func (h *SettingsHandler) DetectIPAddresses(w http.ResponseWriter, r *http.Request) {
	detectedIPs, err := h.service.DetectIPAddresses()
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to detect IP addresses", err.Error())
		return
	}

	utils.SendJSON(w, http.StatusOK, detectedIPs)
}
