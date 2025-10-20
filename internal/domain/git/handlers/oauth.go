package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/api/middleware"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/domain/git"
	"github.com/mikrocloud/mikrocloud/internal/domain/git/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type OAuthState struct {
	UserID                  string
	OrgID                   string
	Provider                git.GitProvider
	Name                    string
	WebhookURL              *string
	AllowPreviewDeployments bool
	CustomURL               *string
	CreatedAt               time.Time
}

type OAuthStateStore struct {
	mu     sync.RWMutex
	states map[string]*OAuthState
}

func NewOAuthStateStore() *OAuthStateStore {
	store := &OAuthStateStore{
		states: make(map[string]*OAuthState),
	}
	go store.cleanup()
	return store
}

func (s *OAuthStateStore) Set(stateID string, state *OAuthState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[stateID] = state
}

func (s *OAuthStateStore) Get(stateID string) (*OAuthState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, exists := s.states[stateID]
	return state, exists
}

func (s *OAuthStateStore) Delete(stateID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, stateID)
}

func (s *OAuthStateStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, state := range s.states {
			if now.Sub(state.CreatedAt) > 10*time.Minute {
				delete(s.states, id)
			}
		}
		s.mu.Unlock()
	}
}

type OAuthHandlers struct {
	service    *service.GitService
	cfg        *config.Config
	stateStore *OAuthStateStore
}

func NewOAuthHandlers(service *service.GitService, cfg *config.Config) *OAuthHandlers {
	return &OAuthHandlers{
		service:    service,
		cfg:        cfg,
		stateStore: NewOAuthStateStore(),
	}
}

func (h *OAuthHandlers) GetStateStore() *OAuthStateStore {
	return h.stateStore
}

func (h *OAuthHandlers) StartOAuth(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	orgID := middleware.GetOrgID(r)
	if userID == "" {
		utils.SendError(w, http.StatusUnauthorized, "unauthorized", "User not authenticated")
		return
	}

	provider := git.GitProvider(r.URL.Query().Get("provider"))
	name := r.URL.Query().Get("name")
	webhookURL := r.URL.Query().Get("webhook_url")
	allowPreview := r.URL.Query().Get("allow_preview") == "true"
	customURL := r.URL.Query().Get("custom_url")

	if provider == "" || name == "" {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "provider and name are required")
		return
	}

	stateID, err := generateRandomString(32)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "state_generation_failed", "Failed to generate OAuth state")
		return
	}

	state := &OAuthState{
		UserID:                  userID,
		OrgID:                   orgID,
		Provider:                provider,
		Name:                    name,
		AllowPreviewDeployments: allowPreview,
		CreatedAt:               time.Now(),
	}

	if webhookURL != "" {
		state.WebhookURL = &webhookURL
	}

	if customURL != "" {
		state.CustomURL = &customURL
	}

	h.stateStore.Set(stateID, state)

	authURL, err := h.buildAuthURL(provider, stateID, customURL)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "auth_url_failed", err.Error())
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *OAuthHandlers) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	stateID := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		http.Redirect(w, r, fmt.Sprintf("/dashboard/git-sources?error=%s", url.QueryEscape(errorParam)), http.StatusFound)
		return
	}

	if code == "" || stateID == "" {
		http.Redirect(w, r, "/dashboard/git-sources?error=missing_parameters", http.StatusFound)
		return
	}

	state, exists := h.stateStore.Get(stateID)
	if !exists {
		http.Redirect(w, r, "/dashboard/git-sources?error=invalid_state", http.StatusFound)
		return
	}

	defer h.stateStore.Delete(stateID)

	token, refreshToken, expiresAt, err := h.exchangeCodeForToken(state.Provider, code, state.CustomURL)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err)
		http.Redirect(w, r, "/dashboard/git-sources?error=token_exchange_failed", http.StatusFound)
		return
	}

	createReq := git.CreateGitSourceRequest{
		Provider:                state.Provider,
		Name:                    state.Name,
		AccessToken:             token,
		RefreshToken:            refreshToken,
		CustomURL:               state.CustomURL,
		WebhookURL:              state.WebhookURL,
		AllowPreviewDeployments: state.AllowPreviewDeployments,
	}

	if expiresAt != nil {
		createReq.TokenExpiresAt = expiresAt
	}

	gitSource, err := h.service.CreateGitSource(r.Context(), state.OrgID, state.UserID, createReq)
	if err != nil {
		slog.Error("Failed to create git source", "error", err)
		http.Redirect(w, r, "/dashboard/git-sources?error=creation_failed", http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/git-sources?success=true&source_id=%s", gitSource.ID), http.StatusFound)
}

func (h *OAuthHandlers) buildAuthURL(provider git.GitProvider, state, customURL string) (string, error) {
	switch provider {
	case git.GitProviderGitHub:
		return h.buildGitHubAuthURL(state, customURL)
	case git.GitProviderGitLab:
		return h.buildGitLabAuthURL(state, customURL)
	case git.GitProviderBitbucket:
		return h.buildBitbucketAuthURL(state)
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (h *OAuthHandlers) buildGitHubAuthURL(state, customURL string) (string, error) {
	baseURL := "https://github.com"
	if customURL != "" {
		baseURL = strings.TrimSuffix(customURL, "/")
	}

	clientID := h.cfg.Auth.GitOAuth.GitHub.ClientID
	if clientID == "" {
		return "", fmt.Errorf("GitHub OAuth client ID not configured")
	}

	redirectURL := h.cfg.Auth.GitOAuth.GitHub.RedirectURL
	scopes := "repo,read:user,user:email,admin:repo_hook"

	params := url.Values{
		"client_id":    {clientID},
		"redirect_uri": {redirectURL},
		"scope":        {scopes},
		"state":        {state},
	}

	return fmt.Sprintf("%s/login/oauth/authorize?%s", baseURL, params.Encode()), nil
}

func (h *OAuthHandlers) buildGitLabAuthURL(state, customURL string) (string, error) {
	baseURL := "https://gitlab.com"
	if customURL != "" {
		baseURL = strings.TrimSuffix(customURL, "/")
	}

	clientID := h.cfg.Auth.GitOAuth.GitLab.ClientID
	if clientID == "" {
		return "", fmt.Errorf("GitLab OAuth client ID not configured")
	}

	redirectURL := h.cfg.Auth.GitOAuth.GitLab.RedirectURL
	scopes := "api,read_user,write_repository"

	params := url.Values{
		"client_id":     {clientID},
		"redirect_uri":  {redirectURL},
		"response_type": {"code"},
		"scope":         {scopes},
		"state":         {state},
	}

	return fmt.Sprintf("%s/oauth/authorize?%s", baseURL, params.Encode()), nil
}

func (h *OAuthHandlers) buildBitbucketAuthURL(state string) (string, error) {
	clientID := h.cfg.Auth.GitOAuth.Bitbucket.ClientID
	if clientID == "" {
		return "", fmt.Errorf("Bitbucket OAuth client ID not configured")
	}

	redirectURL := h.cfg.Auth.GitOAuth.Bitbucket.RedirectURL
	scopes := "repository,repository:admin,webhook"

	params := url.Values{
		"client_id":     {clientID},
		"response_type": {"code"},
		"redirect_uri":  {redirectURL},
		"scope":         {scopes},
		"state":         {state},
	}

	return fmt.Sprintf("https://bitbucket.org/site/oauth2/authorize?%s", params.Encode()), nil
}

func (h *OAuthHandlers) exchangeCodeForToken(provider git.GitProvider, code string, customURL *string) (string, string, *time.Time, error) {
	switch provider {
	case git.GitProviderGitHub:
		return h.exchangeGitHubToken(code, customURL)
	case git.GitProviderGitLab:
		return h.exchangeGitLabToken(code, customURL)
	case git.GitProviderBitbucket:
		return h.exchangeBitbucketToken(code)
	default:
		return "", "", nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (h *OAuthHandlers) exchangeGitHubToken(code string, customURL *string) (string, string, *time.Time, error) {
	baseURL := "https://github.com"
	if customURL != nil && *customURL != "" {
		baseURL = strings.TrimSuffix(*customURL, "/")
	}

	tokenURL := fmt.Sprintf("%s/login/oauth/access_token", baseURL)

	data := url.Values{
		"client_id":     {h.cfg.Auth.GitOAuth.GitHub.ClientID},
		"client_secret": {h.cfg.Auth.GitOAuth.GitHub.ClientSecret},
		"code":          {code},
		"redirect_uri":  {h.cfg.Auth.GitOAuth.GitHub.RedirectURL},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return "", "", nil, fmt.Errorf("no access token in response")
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &exp
	}

	return result.AccessToken, result.RefreshToken, expiresAt, nil
}

func (h *OAuthHandlers) exchangeGitLabToken(code string, customURL *string) (string, string, *time.Time, error) {
	baseURL := "https://gitlab.com"
	if customURL != nil && *customURL != "" {
		baseURL = strings.TrimSuffix(*customURL, "/")
	}

	tokenURL := fmt.Sprintf("%s/oauth/token", baseURL)

	data := url.Values{
		"client_id":     {h.cfg.Auth.GitOAuth.GitLab.ClientID},
		"client_secret": {h.cfg.Auth.GitOAuth.GitLab.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {h.cfg.Auth.GitOAuth.GitLab.RedirectURL},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		CreatedAt    int64  `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return "", "", nil, fmt.Errorf("no access token in response")
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &exp
	}

	return result.AccessToken, result.RefreshToken, expiresAt, nil
}

func (h *OAuthHandlers) exchangeBitbucketToken(code string) (string, string, *time.Time, error) {
	tokenURL := "https://bitbucket.org/site/oauth2/access_token"

	data := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(h.cfg.Auth.GitOAuth.Bitbucket.ClientID, h.cfg.Auth.GitOAuth.Bitbucket.ClientSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.AccessToken == "" {
		return "", "", nil, fmt.Errorf("no access token in response")
	}

	var expiresAt *time.Time
	if result.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
		expiresAt = &exp
	}

	return result.AccessToken, result.RefreshToken, expiresAt, nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func RegisterOAuthRoutes(r chi.Router, handlers *OAuthHandlers) {
	r.Get("/auth/git/oauth/start", handlers.StartOAuth)
}
