package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mikrocloud/mikrocloud/internal/domain/git"
	"github.com/mikrocloud/mikrocloud/internal/domain/git/service"
	"github.com/mikrocloud/mikrocloud/internal/utils"
)

type WebhookHandlers struct {
	service *service.GitService
}

func NewWebhookHandlers(service *service.GitService) *WebhookHandlers {
	return &WebhookHandlers{
		service: service,
	}
}

type WebhookEvent struct {
	SourceID   string
	Provider   git.GitProvider
	EventType  WebhookEventType
	Repository string
	Branch     string
	Commit     string
	CommitMsg  string
	Author     string
	IsPR       bool
	PRNumber   int
	PRAction   string
	PRBranch   string
}

type WebhookEventType string

const (
	EventTypePush   WebhookEventType = "push"
	EventTypePR     WebhookEventType = "pull_request"
	EventTypeBranch WebhookEventType = "branch"
)

func (h *WebhookHandlers) HandleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "source_id")
	if sourceID == "" {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "source_id is required")
		return
	}

	gitSource, err := h.service.GetGitSource(r.Context(), sourceID)
	if err != nil {
		slog.Error("Failed to get git source", "error", err, "source_id", sourceID)
		utils.SendError(w, http.StatusNotFound, "source_not_found", "Git source not found")
		return
	}

	if gitSource.Provider != git.GitProviderGitHub {
		utils.SendError(w, http.StatusBadRequest, "invalid_provider", "Source is not a GitHub provider")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_body", "Failed to read request body")
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		signature = r.Header.Get("X-Hub-Signature")
	}

	if !h.verifyGitHubSignature(body, signature, gitSource.AccessToken) {
		slog.Warn("Invalid GitHub webhook signature", "source_id", sourceID)
		utils.SendError(w, http.StatusUnauthorized, "invalid_signature", "Invalid webhook signature")
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	webhookEvent, err := h.parseGitHubEvent(event, body)
	if err != nil {
		slog.Error("Failed to parse GitHub event", "error", err, "event", event)
		utils.SendError(w, http.StatusBadRequest, "parse_error", err.Error())
		return
	}

	if webhookEvent == nil {
		utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Event ignored"})
		return
	}

	webhookEvent.SourceID = sourceID
	webhookEvent.Provider = git.GitProviderGitHub

	if err := h.processWebhookEvent(r.Context(), webhookEvent, gitSource); err != nil {
		slog.Error("Failed to process webhook event", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "processing_error", "Failed to process webhook event")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Webhook processed successfully"})
}

func (h *WebhookHandlers) HandleGitLabWebhook(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "source_id")
	if sourceID == "" {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "source_id is required")
		return
	}

	gitSource, err := h.service.GetGitSource(r.Context(), sourceID)
	if err != nil {
		slog.Error("Failed to get git source", "error", err, "source_id", sourceID)
		utils.SendError(w, http.StatusNotFound, "source_not_found", "Git source not found")
		return
	}

	if gitSource.Provider != git.GitProviderGitLab {
		utils.SendError(w, http.StatusBadRequest, "invalid_provider", "Source is not a GitLab provider")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_body", "Failed to read request body")
		return
	}
	defer r.Body.Close()

	token := r.Header.Get("X-Gitlab-Token")
	if token != gitSource.AccessToken {
		slog.Warn("Invalid GitLab webhook token", "source_id", sourceID)
		utils.SendError(w, http.StatusUnauthorized, "invalid_token", "Invalid webhook token")
		return
	}

	event := r.Header.Get("X-Gitlab-Event")
	webhookEvent, err := h.parseGitLabEvent(event, body)
	if err != nil {
		slog.Error("Failed to parse GitLab event", "error", err, "event", event)
		utils.SendError(w, http.StatusBadRequest, "parse_error", err.Error())
		return
	}

	if webhookEvent == nil {
		utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Event ignored"})
		return
	}

	webhookEvent.SourceID = sourceID
	webhookEvent.Provider = git.GitProviderGitLab

	if err := h.processWebhookEvent(r.Context(), webhookEvent, gitSource); err != nil {
		slog.Error("Failed to process webhook event", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "processing_error", "Failed to process webhook event")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Webhook processed successfully"})
}

func (h *WebhookHandlers) HandleBitbucketWebhook(w http.ResponseWriter, r *http.Request) {
	sourceID := chi.URLParam(r, "source_id")
	if sourceID == "" {
		utils.SendError(w, http.StatusBadRequest, "invalid_request", "source_id is required")
		return
	}

	gitSource, err := h.service.GetGitSource(r.Context(), sourceID)
	if err != nil {
		slog.Error("Failed to get git source", "error", err, "source_id", sourceID)
		utils.SendError(w, http.StatusNotFound, "source_not_found", "Git source not found")
		return
	}

	if gitSource.Provider != git.GitProviderBitbucket {
		utils.SendError(w, http.StatusBadRequest, "invalid_provider", "Source is not a Bitbucket provider")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "invalid_body", "Failed to read request body")
		return
	}
	defer r.Body.Close()

	event := r.Header.Get("X-Event-Key")
	webhookEvent, err := h.parseBitbucketEvent(event, body)
	if err != nil {
		slog.Error("Failed to parse Bitbucket event", "error", err, "event", event)
		utils.SendError(w, http.StatusBadRequest, "parse_error", err.Error())
		return
	}

	if webhookEvent == nil {
		utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Event ignored"})
		return
	}

	webhookEvent.SourceID = sourceID
	webhookEvent.Provider = git.GitProviderBitbucket

	if err := h.processWebhookEvent(r.Context(), webhookEvent, gitSource); err != nil {
		slog.Error("Failed to process webhook event", "error", err)
		utils.SendError(w, http.StatusInternalServerError, "processing_error", "Failed to process webhook event")
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]string{"message": "Webhook processed successfully"})
}

func (h *WebhookHandlers) verifyGitHubSignature(payload []byte, signature, secret string) bool {
	if signature == "" {
		return false
	}

	var mac hash.Hash
	var prefix string

	if strings.HasPrefix(signature, "sha256=") {
		mac = hmac.New(sha256.New, []byte(secret))
		prefix = "sha256="
	} else if strings.HasPrefix(signature, "sha1=") {
		mac = hmac.New(sha1.New, []byte(secret))
		prefix = "sha1="
	} else {
		return false
	}

	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	expectedSignature := prefix + expectedMAC

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (h *WebhookHandlers) parseGitHubEvent(event string, body []byte) (*WebhookEvent, error) {
	switch event {
	case "push":
		return h.parseGitHubPushEvent(body)
	case "pull_request":
		return h.parseGitHubPREvent(body)
	default:
		return nil, nil
	}
}

func (h *WebhookHandlers) parseGitHubPushEvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		Ref        string `json:"ref"`
		Deleted    bool   `json:"deleted"`
		HeadCommit struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"head_commit"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse push event: %w", err)
	}

	if payload.Deleted {
		return nil, nil
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	return &WebhookEvent{
		EventType:  EventTypePush,
		Repository: payload.Repository.FullName,
		Branch:     branch,
		Commit:     payload.HeadCommit.ID,
		CommitMsg:  payload.HeadCommit.Message,
		Author:     payload.HeadCommit.Author.Name,
		IsPR:       false,
	}, nil
}

func (h *WebhookHandlers) parseGitHubPREvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		Action      string `json:"action"`
		Number      int    `json:"number"`
		PullRequest struct {
			Head struct {
				Ref string `json:"ref"`
				SHA string `json:"sha"`
			} `json:"head"`
			Base struct {
				Ref string `json:"ref"`
			} `json:"base"`
			Title string `json:"title"`
			User  struct {
				Login string `json:"login"`
			} `json:"user"`
		} `json:"pull_request"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse PR event: %w", err)
	}

	if payload.Action != "opened" && payload.Action != "synchronize" && payload.Action != "reopened" {
		return nil, nil
	}

	return &WebhookEvent{
		EventType:  EventTypePR,
		Repository: payload.Repository.FullName,
		Branch:     payload.PullRequest.Base.Ref,
		Commit:     payload.PullRequest.Head.SHA,
		CommitMsg:  payload.PullRequest.Title,
		Author:     payload.PullRequest.User.Login,
		IsPR:       true,
		PRNumber:   payload.Number,
		PRAction:   payload.Action,
		PRBranch:   payload.PullRequest.Head.Ref,
	}, nil
}

func (h *WebhookHandlers) parseGitLabEvent(event string, body []byte) (*WebhookEvent, error) {
	switch event {
	case "Push Hook":
		return h.parseGitLabPushEvent(body)
	case "Merge Request Hook":
		return h.parseGitLabMREvent(body)
	default:
		return nil, nil
	}
}

func (h *WebhookHandlers) parseGitLabPushEvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		Ref     string `json:"ref"`
		After   string `json:"after"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commits"`
		Project struct {
			PathWithNamespace string `json:"path_with_namespace"`
		} `json:"project"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse push event: %w", err)
	}

	if payload.After == "0000000000000000000000000000000000000000" {
		return nil, nil
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	var commitMsg, author string
	if len(payload.Commits) > 0 {
		commitMsg = payload.Commits[len(payload.Commits)-1].Message
		author = payload.Commits[len(payload.Commits)-1].Author.Name
	}

	return &WebhookEvent{
		EventType:  EventTypePush,
		Repository: payload.Project.PathWithNamespace,
		Branch:     branch,
		Commit:     payload.After,
		CommitMsg:  commitMsg,
		Author:     author,
		IsPR:       false,
	}, nil
}

func (h *WebhookHandlers) parseGitLabMREvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		ObjectAttributes struct {
			Action       string `json:"action"`
			IID          int    `json:"iid"`
			Title        string `json:"title"`
			SourceBranch string `json:"source_branch"`
			TargetBranch string `json:"target_branch"`
			LastCommit   struct {
				ID string `json:"id"`
			} `json:"last_commit"`
		} `json:"object_attributes"`
		User struct {
			Name string `json:"name"`
		} `json:"user"`
		Project struct {
			PathWithNamespace string `json:"path_with_namespace"`
		} `json:"project"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse MR event: %w", err)
	}

	action := payload.ObjectAttributes.Action
	if action != "open" && action != "update" && action != "reopen" {
		return nil, nil
	}

	return &WebhookEvent{
		EventType:  EventTypePR,
		Repository: payload.Project.PathWithNamespace,
		Branch:     payload.ObjectAttributes.TargetBranch,
		Commit:     payload.ObjectAttributes.LastCommit.ID,
		CommitMsg:  payload.ObjectAttributes.Title,
		Author:     payload.User.Name,
		IsPR:       true,
		PRNumber:   payload.ObjectAttributes.IID,
		PRAction:   action,
		PRBranch:   payload.ObjectAttributes.SourceBranch,
	}, nil
}

func (h *WebhookHandlers) parseBitbucketEvent(event string, body []byte) (*WebhookEvent, error) {
	switch event {
	case "repo:push":
		return h.parseBitbucketPushEvent(body)
	case "pullrequest:created", "pullrequest:updated":
		return h.parseBitbucketPREvent(body)
	default:
		return nil, nil
	}
}

func (h *WebhookHandlers) parseBitbucketPushEvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		Push struct {
			Changes []struct {
				New struct {
					Type   string `json:"type"`
					Name   string `json:"name"`
					Target struct {
						Hash    string `json:"hash"`
						Message string `json:"message"`
						Author  struct {
							User struct {
								DisplayName string `json:"display_name"`
							} `json:"user"`
						} `json:"author"`
					} `json:"target"`
				} `json:"new"`
			} `json:"changes"`
		} `json:"push"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse push event: %w", err)
	}

	if len(payload.Push.Changes) == 0 {
		return nil, nil
	}

	change := payload.Push.Changes[0]
	if change.New.Type != "branch" {
		return nil, nil
	}

	return &WebhookEvent{
		EventType:  EventTypePush,
		Repository: payload.Repository.FullName,
		Branch:     change.New.Name,
		Commit:     change.New.Target.Hash,
		CommitMsg:  change.New.Target.Message,
		Author:     change.New.Target.Author.User.DisplayName,
		IsPR:       false,
	}, nil
}

func (h *WebhookHandlers) parseBitbucketPREvent(body []byte) (*WebhookEvent, error) {
	var payload struct {
		PullRequest struct {
			ID     int    `json:"id"`
			Title  string `json:"title"`
			Source struct {
				Branch struct {
					Name string `json:"name"`
				} `json:"branch"`
				Commit struct {
					Hash string `json:"hash"`
				} `json:"commit"`
			} `json:"source"`
			Destination struct {
				Branch struct {
					Name string `json:"name"`
				} `json:"branch"`
			} `json:"destination"`
			Author struct {
				DisplayName string `json:"display_name"`
			} `json:"author"`
		} `json:"pullrequest"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse PR event: %w", err)
	}

	return &WebhookEvent{
		EventType:  EventTypePR,
		Repository: payload.Repository.FullName,
		Branch:     payload.PullRequest.Destination.Branch.Name,
		Commit:     payload.PullRequest.Source.Commit.Hash,
		CommitMsg:  payload.PullRequest.Title,
		Author:     payload.PullRequest.Author.DisplayName,
		IsPR:       true,
		PRNumber:   payload.PullRequest.ID,
		PRBranch:   payload.PullRequest.Source.Branch.Name,
	}, nil
}

func (h *WebhookHandlers) processWebhookEvent(ctx context.Context, event *WebhookEvent, gitSource *git.GitSource) error {
	slog.Info("Processing webhook event",
		"source_id", event.SourceID,
		"provider", event.Provider,
		"event_type", event.EventType,
		"repository", event.Repository,
		"branch", event.Branch,
		"commit", event.Commit,
		"is_pr", event.IsPR,
		"pr_number", event.PRNumber,
	)

	if event.IsPR && !gitSource.AllowPreviewDeployments {
		slog.Info("Skipping PR deployment - preview deployments disabled", "source_id", event.SourceID)
		return nil
	}

	slog.Warn("Webhook received but deployment triggering not yet implemented",
		"source_id", event.SourceID,
		"repository", event.Repository,
		"branch", event.Branch,
		"commit", event.Commit[:8],
	)

	return nil
}

func RegisterWebhookRoutes(r chi.Router, handlers *WebhookHandlers) {
	r.Post("/webhooks/git/{source_id}/github", handlers.HandleGitHubWebhook)
	r.Post("/webhooks/git/{source_id}/gitlab", handlers.HandleGitLabWebhook)
	r.Post("/webhooks/git/{source_id}/bitbucket", handlers.HandleBitbucketWebhook)
}
