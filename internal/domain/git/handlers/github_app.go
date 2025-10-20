package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/domain/git"
	"github.com/mikrocloud/mikrocloud/internal/domain/git/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type GitHubAppManifest struct {
	Name          string                       `json:"name"`
	URL           string                       `json:"url"`
	RedirectURL   string                       `json:"redirect_url"`
	HookAttrs     GitHubAppManifestHookAttrs   `json:"hook_attributes"`
	Public        bool                         `json:"public"`
	DefaultPerms  GitHubAppManifestPermissions `json:"default_permissions"`
	DefaultEvents []string                     `json:"default_events"`
}

type GitHubAppManifestHookAttrs struct {
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type GitHubAppManifestPermissions struct {
	Contents     string `json:"contents"`
	Metadata     string `json:"metadata"`
	Emails       string `json:"emails"`
	PullRequests string `json:"pull_requests"`
}

type GitHubAppManifestConversionResponse struct {
	ID            int64  `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
}

type GitHubAppHandlers struct {
	service    *service.GitService
	cfg        *config.Config
	stateStore *OAuthStateStore
}

func NewGitHubAppHandlers(service *service.GitService, cfg *config.Config, stateStore *OAuthStateStore) *GitHubAppHandlers {
	return &GitHubAppHandlers{
		service:    service,
		cfg:        cfg,
		stateStore: stateStore,
	}
}

func (h *GitHubAppHandlers) GenerateManifest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	orgID := middleware.GetOrgID(r)
	if userID == "" {
		utils.SendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
		return
	}

	name := r.URL.Query().Get("name")
	webhookURL := r.URL.Query().Get("webhook_url")
	customURL := r.URL.Query().Get("custom_url")
	allowPreview := r.URL.Query().Get("allow_preview") == "true"

	if name == "" || webhookURL == "" {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "name and webhook_url are required")
		return
	}

	stateID, err := generateRandomString(32)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "state_generation_failed", "Failed to generate state")
		return
	}

	state := &OAuthState{
		UserID:                  userID,
		OrgID:                   orgID,
		Provider:                git.GitProviderGitHub,
		Name:                    name,
		AllowPreviewDeployments: allowPreview,
		CreatedAt:               time.Now(),
	}

	webhookURLParam := webhookURL
	state.WebhookURL = &webhookURLParam

	if customURL != "" {
		state.CustomURL = &customURL
	}

	h.stateStore.Set(stateID, state)

	baseURL := "https://github.com"
	if customURL != "" {
		baseURL = strings.TrimSuffix(customURL, "/")
	}

	instanceURL := h.cfg.GetPublicURL()

	manifest := GitHubAppManifest{
		Name:        name,
		URL:         instanceURL,
		RedirectURL: fmt.Sprintf("%s/api/auth/git/github-app/callback?state=%s", instanceURL, stateID),
		HookAttrs: GitHubAppManifestHookAttrs{
			URL:    webhookURL,
			Active: true,
		},
		Public: false,
		DefaultPerms: GitHubAppManifestPermissions{
			Contents:     "read",
			Metadata:     "read",
			Emails:       "read",
			PullRequests: "write",
		},
		DefaultEvents: []string{"push", "pull_request"},
	}

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "manifest_generation_failed", "Failed to generate manifest")
		return
	}

	manifestURL := fmt.Sprintf("%s/settings/apps/new?state=%s", baseURL, stateID)

	response := map[string]interface{}{
		"manifest_url": manifestURL,
		"manifest":     manifest,
		"form_data":    base64.StdEncoding.EncodeToString(manifestJSON),
	}

	utils.SendJSON(w, http.StatusOK, response)
}

func (h *GitHubAppHandlers) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	stateID := r.URL.Query().Get("state")

	if code == "" || stateID == "" {
		http.Redirect(w, r, "/dashboard/git?error=missing_parameters", http.StatusFound)
		return
	}

	state, exists := h.stateStore.Get(stateID)
	if !exists {
		http.Redirect(w, r, "/dashboard/git?error=invalid_state", http.StatusFound)
		return
	}

	defer h.stateStore.Delete(stateID)

	baseURL := "https://api.github.com"
	if state.CustomURL != nil && *state.CustomURL != "" {
		customURL := strings.TrimSuffix(*state.CustomURL, "/")
		if strings.Contains(customURL, "github.com") && !strings.Contains(customURL, "api.github.com") {
			baseURL = strings.Replace(customURL, "github.com", "api.github.com", 1)
		} else {
			baseURL = fmt.Sprintf("%s/api/v3", customURL)
		}
	}

	conversionURL := fmt.Sprintf("%s/app-manifests/%s/conversions", baseURL, code)

	req, err := http.NewRequest("POST", conversionURL, nil)
	if err != nil {
		slog.Error("Failed to create conversion request", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=conversion_failed", http.StatusFound)
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to convert manifest", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=conversion_failed", http.StatusFound)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Manifest conversion failed", "status", resp.StatusCode, "body", string(body))
		http.Redirect(w, r, "/dashboard/git?error=conversion_failed", http.StatusFound)
		return
	}

	var conversion GitHubAppManifestConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&conversion); err != nil {
		slog.Error("Failed to decode conversion response", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=conversion_decode_failed", http.StatusFound)
		return
	}

	appID := fmt.Sprintf("%d", conversion.ID)
	gitSource := &git.GitSource{
		ID:                      uuid.New().String(),
		OrgID:                   state.OrgID,
		UserID:                  state.UserID,
		Provider:                git.GitProviderGitHub,
		Name:                    state.Name,
		AccessToken:             "",
		WebhookURL:              state.WebhookURL,
		AllowPreviewDeployments: state.AllowPreviewDeployments,
		CustomURL:               state.CustomURL,
		IsGitHubApp:             true,
		GitHubAppID:             &appID,
		GitHubClientID:          &conversion.ClientID,
		GitHubClientSecret:      &conversion.ClientSecret,
		GitHubWebhookSecret:     &conversion.WebhookSecret,
		GitHubPrivateKey:        &conversion.PEM,
		GitHubAppSlug:           &conversion.Slug,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	if err := h.service.CreateGitSourceDirect(r.Context(), gitSource); err != nil {
		slog.Error("Failed to create git source", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=creation_failed", http.StatusFound)
		return
	}

	installURL := fmt.Sprintf("%s/apps/%s/installations/new",
		strings.Replace(baseURL, "/api/v3", "", 1),
		conversion.Slug)
	if !strings.Contains(baseURL, "/api") {
		baseWebURL := "https://github.com"
		if state.CustomURL != nil && *state.CustomURL != "" {
			baseWebURL = strings.TrimSuffix(*state.CustomURL, "/")
		}
		installURL = fmt.Sprintf("%s/apps/%s/installations/new", baseWebURL, conversion.Slug)
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/git?success=true&source_id=%s&install_url=%s", gitSource.ID, installURL), http.StatusFound)
}

func (h *GitHubAppHandlers) InstallCallback(w http.ResponseWriter, r *http.Request) {
	installationID := r.URL.Query().Get("installation_id")
	setupAction := r.URL.Query().Get("setup_action")
	sourceID := r.URL.Query().Get("source_id")

	if setupAction != "install" || installationID == "" || sourceID == "" {
		http.Redirect(w, r, "/dashboard/git?error=invalid_install", http.StatusFound)
		return
	}

	source, err := h.service.GetGitSource(r.Context(), sourceID)
	if err != nil {
		slog.Error("Failed to get git source", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=source_not_found", http.StatusFound)
		return
	}

	source.GitHubInstallationID = &installationID
	source.UpdatedAt = time.Now()

	if err := h.service.UpdateGitSourceDirect(r.Context(), source); err != nil {
		slog.Error("Failed to update git source", "error", err)
		http.Redirect(w, r, "/dashboard/git?error=update_failed", http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/git?success=true&source_id=%s&installed=true", sourceID), http.StatusFound)
}
