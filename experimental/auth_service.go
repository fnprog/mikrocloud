package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"github.com/mikrocloud/mikrocloud/internal/domain"
	"github.com/mikrocloud/mikrocloud/pkg/auth"
)

// AuthService handles authentication operations
type AuthService struct {
	repos      domain.RepositoryManager
	jwtManager *auth.JWTManager
}

// NewAuthService creates a new authentication service
func NewAuthService(repos domain.RepositoryManager, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		repos:      repos,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Password        string `json:"password" binding:"required,min=8"`
	PreferredCloud  string `json:"preferred_cloud,omitempty"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	PreferredCloud string    `json:"preferred_cloud"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.repos.User().GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse preferred cloud provider
	var preferredCloud domain.CloudProvider
	switch req.PreferredCloud {
	case "aws":
		preferredCloud = domain.CloudProviderAWS
	case "azure":
		preferredCloud = domain.CloudProviderAzure
	case "gcp":
		preferredCloud = domain.CloudProviderGCP
	default:
		preferredCloud = domain.CloudProviderAWS // Default to AWS
	}

	// Create user
	user := &domain.User{
		Email:          req.Email,
		Name:           req.Name,
		PasswordHash:   string(hashedPassword),
		PreferredCloud: preferredCloud,
		Role:           "user", // Default role
		IsActive:       true,
	}

	if err := s.repos.User().Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtManager.Generate(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:      s.userToResponse(user),
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiration
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.repos.User().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repos.User().Update(ctx, user); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Generate JWT token
	token, err := s.jwtManager.Generate(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:      s.userToResponse(user),
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// RefreshToken generates a new token from an existing one
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (*AuthResponse, error) {
	// Verify and extract claims from existing token
	claims, err := s.jwtManager.Verify(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user to ensure they still exist and are active
	user, err := s.repos.User().GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new token
	newToken, err := s.jwtManager.Generate(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		User:      s.userToResponse(user),
		Token:     newToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// GetProfile returns user profile information
func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.repos.User().GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return s.userToResponse(user), nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*UserResponse, error) {
	user, err := s.repos.User().GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok && name != "" {
		user.Name = name
	}

	if preferredCloud, ok := updates["preferred_cloud"].(string); ok {
		switch preferredCloud {
		case "aws":
			user.PreferredCloud = domain.CloudProviderAWS
		case "azure":
			user.PreferredCloud = domain.CloudProviderAzure
		case "gcp":
			user.PreferredCloud = domain.CloudProviderGCP
		}
	}

	if err := s.repos.User().Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return s.userToResponse(user), nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.repos.User().GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.repos.User().Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GenerateAPIKey generates a new API key for the user
func (s *AuthService) GenerateAPIKey(ctx context.Context, userID uuid.UUID, name string) (string, error) {
	// Generate random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}

	apiKey := "mk_" + hex.EncodeToString(keyBytes)

	// TODO: Store API key in database with metadata
	// For now, just return the key
	return apiKey, nil
}

// ValidateAPIKey validates an API key and returns user information
func (s *AuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*domain.User, error) {
	// TODO: Implement API key validation from database
	// For now, return an error
	return nil, fmt.Errorf("API key validation not yet implemented")
}

// Helper methods

func (s *AuthService) userToResponse(user *domain.User) *UserResponse {
	var preferredCloud string
	switch user.PreferredCloud {
	case domain.CloudProviderAWS:
		preferredCloud = "aws"
	case domain.CloudProviderAzure:
		preferredCloud = "azure"
	case domain.CloudProviderGCP:
		preferredCloud = "gcp"
	}

	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		PreferredCloud: preferredCloud,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
	}
}