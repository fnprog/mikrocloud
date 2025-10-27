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
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/domain/auth/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type UserOAuthState struct {
	Provider  string
	CreatedAt time.Time
}

type UserOAuthStateStore struct {
	mu     sync.RWMutex
	states map[string]*UserOAuthState
}

func NewUserOAuthStateStore() *UserOAuthStateStore {
	store := &UserOAuthStateStore{
		states: make(map[string]*UserOAuthState),
	}
	go store.cleanup()
	return store
}

func (s *UserOAuthStateStore) Set(stateID string, state *UserOAuthState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[stateID] = state
}

func (s *UserOAuthStateStore) Get(stateID string) (*UserOAuthState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, exists := s.states[stateID]
	return state, exists
}

func (s *UserOAuthStateStore) Delete(stateID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, stateID)
}

func (s *UserOAuthStateStore) cleanup() {
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

type UserOAuthHandlers struct {
	authService *service.AuthService
	cfg         *config.Config
	stateStore  *UserOAuthStateStore
}

func NewUserOAuthHandlers(authService *service.AuthService, cfg *config.Config) *UserOAuthHandlers {
	return &UserOAuthHandlers{
		authService: authService,
		cfg:         cfg,
		stateStore:  NewUserOAuthStateStore(),
	}
}

func (h *UserOAuthHandlers) StartOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	var clientID, redirectURL string
	var scopes []string

	switch provider {
	case "github":
		if h.cfg.Auth.UserOAuth.GitHub.ClientID == "" {
			utils.SendError(w, http.StatusBadRequest, "oauth_not_configured", "GitHub OAuth not configured")
			return
		}
		clientID = h.cfg.Auth.UserOAuth.GitHub.ClientID
		redirectURL = h.cfg.Auth.UserOAuth.GitHub.RedirectURL
		scopes = []string{"read:user", "user:email"}
	case "gitlab":
		if h.cfg.Auth.UserOAuth.GitLab.ClientID == "" {
			utils.SendError(w, http.StatusBadRequest, "oauth_not_configured", "GitLab OAuth not configured")
			return
		}
		clientID = h.cfg.Auth.UserOAuth.GitLab.ClientID
		redirectURL = h.cfg.Auth.UserOAuth.GitLab.RedirectURL
		scopes = []string{"read_user", "read_api"}
	case "google":
		if h.cfg.Auth.UserOAuth.Google.ClientID == "" {
			utils.SendError(w, http.StatusBadRequest, "oauth_not_configured", "Google OAuth not configured")
			return
		}
		clientID = h.cfg.Auth.UserOAuth.Google.ClientID
		redirectURL = h.cfg.Auth.UserOAuth.Google.RedirectURL
		scopes = []string{"openid", "profile", "email"}
	default:
		utils.SendError(w, http.StatusBadRequest, "invalid_provider", "Invalid OAuth provider")
		return
	}

	// Generate state
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		utils.SendError(w, http.StatusInternalServerError, "state_generation_failed", "Failed to generate OAuth state")
		return
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// Store state
	h.stateStore.Set(state, &UserOAuthState{
		Provider:  provider,
		CreatedAt: time.Now(),
	})

	// Build authorization URL
	var authURL string
	switch provider {
	case "github":
		authURL = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
			url.QueryEscape(clientID),
			url.QueryEscape(redirectURL),
			url.QueryEscape(strings.Join(scopes, " ")),
			url.QueryEscape(state))
	case "gitlab":
		authURL = fmt.Sprintf("https://gitlab.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code",
			url.QueryEscape(clientID),
			url.QueryEscape(redirectURL),
			url.QueryEscape(strings.Join(scopes, " ")),
			url.QueryEscape(state))
	case "google":
		authURL = fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code&access_type=offline",
			url.QueryEscape(clientID),
			url.QueryEscape(redirectURL),
			url.QueryEscape(strings.Join(scopes, " ")),
			url.QueryEscape(state))
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *UserOAuthHandlers) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		slog.Error("OAuth callback error", "provider", provider, "error", errorParam)
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	if code == "" || state == "" {
		utils.SendError(w, http.StatusBadRequest, "missing_parameters", "Missing code or state parameter")
		return
	}

	// Verify state
	storedState, exists := h.stateStore.Get(state)
	if !exists {
		utils.SendError(w, http.StatusBadRequest, "invalid_state", "Invalid or expired state")
		return
	}
	h.stateStore.Delete(state)

	if storedState.Provider != provider {
		utils.SendError(w, http.StatusBadRequest, "state_mismatch", "State provider mismatch")
		return
	}

	// Exchange code for token
	var clientID, clientSecret, tokenURL, redirectURL string
	switch provider {
	case "github":
		clientID = h.cfg.Auth.UserOAuth.GitHub.ClientID
		clientSecret = h.cfg.Auth.UserOAuth.GitHub.ClientSecret
		tokenURL = "https://github.com/login/oauth/access_token"
		redirectURL = h.cfg.Auth.UserOAuth.GitHub.RedirectURL
	case "gitlab":
		clientID = h.cfg.Auth.UserOAuth.GitLab.ClientID
		clientSecret = h.cfg.Auth.UserOAuth.GitLab.ClientSecret
		tokenURL = "https://gitlab.com/oauth/token"
		redirectURL = h.cfg.Auth.UserOAuth.GitLab.RedirectURL
	case "google":
		clientID = h.cfg.Auth.UserOAuth.Google.ClientID
		clientSecret = h.cfg.Auth.UserOAuth.Google.ClientSecret
		tokenURL = "https://oauth2.googleapis.com/token"
		redirectURL = h.cfg.Auth.UserOAuth.Google.RedirectURL
	default:
		utils.SendError(w, http.StatusBadRequest, "invalid_provider", "Invalid OAuth provider")
		return
	}

	// Exchange code for access token
	tokenResp, err := h.exchangeCodeForToken(tokenURL, clientID, clientSecret, code, redirectURL)
	if err != nil {
		slog.Error("Failed to exchange code for token", "error", err, "provider", provider)
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get user info
	userInfo, err := h.getUserInfo(provider, tokenResp.AccessToken)
	if err != nil {
		slog.Error("Failed to get user info", "error", err, "provider", provider)
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Create or login user
	ctx := r.Context()
	result, err := h.authService.OAuthLogin(ctx, service.OAuthLoginCommand{
		Provider:   provider,
		ProviderID: userInfo.ID,
		Email:      userInfo.Email,
		Name:       userInfo.Name,
		Username:   userInfo.Username,
		AvatarURL:  userInfo.AvatarURL,
	})

	if err != nil {
		slog.Error("OAuth login failed", "error", err, "provider", provider, "email", userInfo.Email)
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24,
	})

	// Redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type UserInfo struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	Username  *string `json:"login,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

func (h *UserOAuthHandlers) exchangeCodeForToken(tokenURL, clientID, clientSecret, code, redirectURL string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func (h *UserOAuthHandlers) getUserInfo(provider, accessToken string) (*UserInfo, error) {
	var userURL string
	switch provider {
	case "github":
		userURL = "https://api.github.com/user"
	case "gitlab":
		userURL = "https://gitlab.com/api/v4/user"
	case "google":
		userURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	// For GitHub, get email separately if not provided
	if provider == "github" && userInfo.Email == "" {
		emails, err := h.getGitHubEmails(accessToken)
		if err == nil && len(emails) > 0 {
			userInfo.Email = emails[0].Email
		}
	}

	return &userInfo, nil
}

func (h *UserOAuthHandlers) getGitHubEmails(accessToken string) ([]GitHubEmail, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("email request failed with status %d", resp.StatusCode)
	}

	var emails []GitHubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	// Filter for primary email
	var primaryEmails []GitHubEmail
	for _, email := range emails {
		if email.Primary {
			primaryEmails = append(primaryEmails, email)
		}
	}

	return primaryEmails, nil
}

type GitHubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func RegisterUserOAuthRoutes(r chi.Router, authService *service.AuthService, cfg *config.Config) {
	handler := NewUserOAuthHandlers(authService, cfg)

	r.Route("/auth/oauth", func(r chi.Router) {
		r.Get("/{provider}/start", handler.StartOAuth)
		r.Get("/{provider}/callback", handler.Callback)
	})
}
