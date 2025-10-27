package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikrocloud/mikrocloud/internal/domain/auth"
	"github.com/mikrocloud/mikrocloud/internal/domain/auth/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	usersRepo "github.com/mikrocloud/mikrocloud/internal/domain/users/repository"
)

// Service errors
var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrWeakPassword          = errors.New("password does not meet requirements")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrSessionExpired        = errors.New("session has expired")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrIncorrectPassword     = errors.New("incorrect password")
)

// AuthService handles authentication business logic
type AuthService struct {
	sessionRepo repository.SessionRepository
	authRepo    repository.AuthRepository
	usersRepo   usersRepo.Repository

	// Configuration
	jwtSecret            string
	sessionDuration      time.Duration
	refreshTokenDuration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(
	sessionRepo repository.SessionRepository,
	authRepo repository.AuthRepository,
	usersRepo usersRepo.Repository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		sessionRepo:          sessionRepo,
		authRepo:             authRepo,
		usersRepo:            usersRepo,
		jwtSecret:            jwtSecret,
		sessionDuration:      24 * time.Hour,      // 24 hours for regular sessions
		refreshTokenDuration: 30 * 24 * time.Hour, // 30 days for refresh tokens
	}
}

// Command types for service operations
type LoginCommand struct {
	Email    string
	Password string
}

type RegisterCommand struct {
	Name     string
	Email    string
	Password string
}

// Result types for service operations
type LoginResult struct {
	User         *users.User
	Token        string
	RefreshToken string
}

type RegisterResult struct {
	User  *users.User
	Token string
}

type OAuthLoginCommand struct {
	Provider   string
	ProviderID string
	Email      string
	Name       string
	Username   *string
	AvatarURL  *string
}

type OAuthLoginResult struct {
	User         *users.User
	Token        string
	RefreshToken string
}

type RefreshTokenResult struct {
	Token        string
	RefreshToken string
}

// Login authenticates a user and creates a new session
func (s *AuthService) Login(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// Validate input
	if cmd.Email == "" || cmd.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// Create email value object
	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	// Get user by email
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(cmd.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if user.Status() != users.UserStatusActive {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateJWTToken(ctx, user.ID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Generate session token for database tracking
	sessionToken, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session
	session := auth.NewSession(user.ID(), sessionToken, s.sessionDuration)
	if err := s.sessionRepo.SaveSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Generate refresh token
	refreshTokenStr, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshToken := auth.NewRefreshToken(user.ID(), session.ID(), refreshTokenStr, s.refreshTokenDuration)
	if err := s.sessionRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Update last login
	if err := s.authRepo.UpdateLastLogin(ctx, user.ID()); err != nil {
		// Log error but don't fail the login
		slog.Error("failed to update last login", "error", err, "user_id", user.ID())
	}

	return &LoginResult{
		User:         user,
		Token:        token,
		RefreshToken: refreshTokenStr,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, cmd RegisterCommand) (*RegisterResult, error) {
	// Validate input
	if err := s.validateRegisterCommand(cmd); err != nil {
		return nil, err
	}

	// Create email value object
	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	// Check if user already exists
	exists, err := s.authRepo.UserExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Check if this is the first user (setup scenario)
	hasUsers, err := s.authRepo.HasAnyUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if users exist: %w", err)
	}
	isFirstUser := !hasUsers

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := users.NewUserWithName(email, string(passwordHash), cmd.Name)

	// If this is the first user, make them active (setup scenario)
	if isFirstUser {
		user.ChangeStatus(users.UserStatusActive)
	}

	// Save user
	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// If this is the first user, set up default organization and admin role
	if isFirstUser {
		if err := s.setupFirstUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to set up first user: %w", err)
		}
	}

	// Generate JWT token
	token, err := s.generateJWTToken(ctx, user.ID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Generate session token for database tracking
	sessionToken, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session
	session := auth.NewSession(user.ID(), sessionToken, s.sessionDuration)
	if err := s.sessionRepo.SaveSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return &RegisterResult{
		User:  user,
		Token: token,
	}, nil
}

// OAuthLogin authenticates or creates a user via OAuth
func (s *AuthService) OAuthLogin(ctx context.Context, cmd OAuthLoginCommand) (*OAuthLoginResult, error) {
	// Validate input
	if cmd.Email == "" || cmd.Name == "" || cmd.Provider == "" || cmd.ProviderID == "" {
		return nil, errors.New("invalid OAuth login data")
	}

	// Create email value object
	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	// Check if user exists by email
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil && err != ErrUserNotFound {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	// If user doesn't exist, create them
	if err == ErrUserNotFound {
		// Check if this is the first user (setup scenario)
		hasUsers, err := s.authRepo.HasAnyUsers(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to check if users exist: %w", err)
		}
		isFirstUser := !hasUsers

		// Create new user with OAuth data
		user = users.NewUserWithName(email, "", cmd.Name) // Empty password for OAuth users

		// Set OAuth provider info
		user.SetOAuthProvider(cmd.Provider, cmd.ProviderID)

		if cmd.Username != nil {
			username, err := users.NewUsername(*cmd.Username)
			if err == nil {
				user.SetUsername(username)
			}
		}

		if cmd.AvatarURL != nil {
			user.SetAvatarURL(cmd.AvatarURL)
		}

		// If this is the first user, make them active and set up admin
		if isFirstUser {
			user.ChangeStatus(users.UserStatusActive)
		}

		// Save user
		if err := s.authRepo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create OAuth user: %w", err)
		}

		// If this is the first user, set up default organization and admin role
		if isFirstUser {
			if err := s.setupFirstUser(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to set up first OAuth user: %w", err)
			}
		}
	} else {
		// Update existing user's OAuth info if needed
		if user.OAuthProvider() != cmd.Provider || user.OAuthProviderID() != cmd.ProviderID {
			user.SetOAuthProvider(cmd.Provider, cmd.ProviderID)
			if err := s.usersRepo.Save(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to update user OAuth info: %w", err)
			}
		}

		// Check if user is active
		if user.Status() != users.UserStatusActive {
			return nil, ErrInvalidCredentials
		}
	}

	// Generate JWT token
	token, err := s.generateJWTToken(ctx, user.ID().String())
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Generate session token for database tracking
	sessionToken, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session
	session := auth.NewSession(user.ID(), sessionToken, s.sessionDuration)
	if err := s.sessionRepo.SaveSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Generate refresh token
	refreshTokenStr, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshToken := auth.NewRefreshToken(user.ID(), session.ID(), refreshTokenStr, s.refreshTokenDuration)
	if err := s.sessionRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Update last login
	if err := s.authRepo.UpdateLastLogin(ctx, user.ID()); err != nil {
		// Log error but don't fail the login
		slog.Error("failed to update last login", "error", err, "user_id", user.ID())
	}

	return &OAuthLoginResult{
		User:         user,
		Token:        token,
		RefreshToken: refreshTokenStr,
	}, nil
}

// Logout invalidates a user session
func (s *AuthService) Logout(ctx context.Context, jwtToken string) error {
	if jwtToken == "" {
		return ErrInvalidToken
	}

	// Parse and verify JWT token
	token, err := jwt.Parse([]byte(jwtToken), jwt.WithKey(jwa.HS256, []byte(s.jwtSecret)))
	if err != nil {
		return ErrInvalidToken
	}

	// Extract user_id from claims
	userIDClaim, ok := token.Get("user_id")
	if !ok {
		return ErrInvalidToken
	}

	userIDStr, ok := userIDClaim.(string)
	if !ok || userIDStr == "" {
		return ErrInvalidToken
	}

	// Parse user ID
	userID, err := users.UserIDFromString(userIDStr)
	if err != nil {
		return ErrInvalidToken
	}

	// Revoke all sessions for this user
	if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	return nil
}

// RefreshToken creates a new session token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*RefreshTokenResult, error) {
	if refreshTokenStr == "" {
		return nil, ErrInvalidToken
	}

	// Get refresh token
	refreshToken, err := s.sessionRepo.GetRefreshTokenByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Validate refresh token
	if !refreshToken.IsValid() {
		return nil, ErrInvalidToken
	}

	// Mark refresh token as used
	if err := s.sessionRepo.MarkRefreshTokenAsUsed(ctx, refreshToken.ID()); err != nil {
		return nil, fmt.Errorf("failed to mark refresh token as used: %w", err)
	}

	// Generate new session token
	newToken, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create new session
	newSession := auth.NewSession(refreshToken.UserID(), newToken, s.sessionDuration)
	if err := s.sessionRepo.SaveSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Generate new refresh token
	newRefreshTokenStr, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newRefreshToken := auth.NewRefreshToken(refreshToken.UserID(), newSession.ID(), newRefreshTokenStr, s.refreshTokenDuration)
	if err := s.sessionRepo.SaveRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &RefreshTokenResult{
		Token:        newToken,
		RefreshToken: newRefreshTokenStr,
	}, nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(ctx context.Context, userIDStr string) (*users.User, error) {
	userID, err := users.UserIDFromString(userIDStr)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// ValidateSession validates a session token and returns the associated user
func (s *AuthService) ValidateSession(ctx context.Context, token string) (*users.User, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	// Get session by token
	session, err := s.sessionRepo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Validate session
	if !session.IsValid() {
		return nil, ErrSessionExpired
	}

	// Get user
	user, err := s.authRepo.GetUserByID(ctx, session.UserID())
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// CleanupExpiredSessions removes expired sessions and refresh tokens
func (s *AuthService) CleanupExpiredSessions(ctx context.Context) error {
	if err := s.sessionRepo.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	if err := s.sessionRepo.DeleteExpiredRefreshTokens(ctx); err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}

// SetupStatus represents the setup status of the application
type SetupStatus struct {
	HasUsers bool `json:"has_users"`
	IsSetup  bool `json:"is_setup"`
}

// GetSetupStatus checks if the application has been set up (has at least one user)
func (s *AuthService) GetSetupStatus(ctx context.Context) (*SetupStatus, error) {
	// Check if there are any users in the system
	hasUsers, err := s.authRepo.HasAnyUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if users exist: %w", err)
	}

	return &SetupStatus{
		HasUsers: hasUsers,
		IsSetup:  hasUsers,
	}, nil
}

// HasAnyUsers checks if there are any users in the system
func (s *AuthService) HasAnyUsers(ctx context.Context) (bool, error) {
	return s.authRepo.HasAnyUsers(ctx)
}

// Helper methods

func (s *AuthService) validateRegisterCommand(cmd RegisterCommand) error {
	if cmd.Name == "" {
		return errors.New("name is required")
	}
	if cmd.Email == "" {
		return ErrInvalidEmail
	}
	if len(cmd.Password) < 8 {
		return ErrWeakPassword
	}
	// Additional password strength validation could be added here
	return nil
}

func (s *AuthService) generateSecureToken() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64 URL-safe string
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// generateJWTToken generates a JWT token for the given user
func (s *AuthService) generateJWTToken(ctx context.Context, userID string) (string, error) {
	now := time.Now()

	// Parse user ID
	userIDObj, err := users.UserIDFromString(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	// Get user's organizations to include first one as org_id in JWT
	orgs, err := s.usersRepo.FindOrganizationsByUser(ctx, userIDObj)
	if err != nil {
		slog.Error("Failed to get user organizations", "error", err, "user_id", userID)
		return "", fmt.Errorf("failed to get user organizations: %w", err)
	}

	// Use the first organization ID, or empty string if no organizations
	var orgID string
	if len(orgs) > 0 {
		orgID = orgs[0].ID().String()
		slog.Info("Found organizations for user", "user_id", userID, "org_id", orgID, "org_count", len(orgs))
	} else {
		slog.Warn("No organizations found for user", "user_id", userID)
	}

	token, err := jwt.NewBuilder().
		Issuer("mikrocloud").
		IssuedAt(now).
		Expiration(now.Add(s.sessionDuration)).
		Claim("user_id", userID).
		Claim("org_id", orgID).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build JWT: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, []byte(s.jwtSecret)))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return string(signed), nil
}

func (s *AuthService) setupFirstUser(ctx context.Context, user *users.User) error {
	// Find the default organization
	defaultOrg, err := s.usersRepo.FindOrganizationBySlug(ctx, "default")
	if err != nil {
		return fmt.Errorf("failed to find default organization: %w", err)
	}

	// Add user to the default organization as owner
	if err := s.usersRepo.AddOrganizationMember(ctx, defaultOrg.ID(), user.ID(), "owner", nil); err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	// Find admin role and assign it
	adminRoleID, err := s.usersRepo.FindRoleByName(ctx, "admin")
	if err != nil {
		return fmt.Errorf("failed to find admin role: %w", err)
	}

	if err := s.usersRepo.AddUserRole(ctx, user.ID(), adminRoleID, nil); err != nil {
		return fmt.Errorf("failed to assign admin role: %w", err)
	}

	return nil
}

type UpdateProfileCommand struct {
	Name     *string
	Username *string
	Avatar   *string
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID users.UserID, cmd UpdateProfileCommand) (*users.User, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if cmd.Name != nil {
		user.SetName(*cmd.Name)
	}

	if cmd.Username != nil {
		if *cmd.Username != "" {
			username, err := users.NewUsername(*cmd.Username)
			if err != nil {
				return nil, fmt.Errorf("invalid username: %w", err)
			}

			existingUser, err := s.usersRepo.FindByUsername(ctx, *cmd.Username)
			if err == nil && existingUser.ID().String() != userID.String() {
				return nil, ErrUsernameAlreadyExists
			}

			user.SetUsername(username)
		} else {
			user.SetUsername(nil)
		}
	}

	if cmd.Avatar != nil {
		if *cmd.Avatar != "" {
			user.SetAvatarURL(cmd.Avatar)
		} else {
			user.SetAvatarURL(nil)
		}
	}

	if err := s.usersRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return user, nil
}

type UpdateEmailCommand struct {
	Email    string
	Password string
}

type UpdatePasswordCommand struct {
	CurrentPassword string
	NewPassword     string
}

func (s *AuthService) UpdateEmail(ctx context.Context, userID users.UserID, cmd UpdateEmailCommand) (*users.User, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(cmd.Password)); err != nil {
		return nil, ErrIncorrectPassword
	}

	newEmail, err := users.NewEmail(cmd.Email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	exists, err := s.authRepo.UserExistsByEmail(ctx, newEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	var oauthProvider *string
	if p := user.OAuthProvider(); p != "" {
		oauthProvider = &p
	}
	var oauthProviderID *string
	if pid := user.OAuthProviderID(); pid != "" {
		oauthProviderID = &pid
	}

	reconstructed := users.ReconstructUser(
		user.ID(),
		newEmail,
		user.PasswordHash(),
		user.Name(),
		user.Username(),
		user.AvatarURL(),
		oauthProvider,
		oauthProviderID,
		user.Status(),
		user.EmailVerifiedAt(),
		user.LastLoginAt(),
		user.Timezone(),
		user.CreatedAt(),
		time.Now(),
	)

	if err := s.usersRepo.Save(ctx, reconstructed); err != nil {
		return nil, fmt.Errorf("failed to update user email: %w", err)
	}

	return reconstructed, nil
}

type RequestPasswordResetCommand struct {
	Email string
}

type ResetPasswordCommand struct {
	Token    string
	Password string
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID users.UserID, cmd UpdatePasswordCommand) error {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(cmd.CurrentPassword)); err != nil {
		return ErrIncorrectPassword
	}

	if len(cmd.NewPassword) < 8 {
		return ErrWeakPassword
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.authRepo.UpdateUserPassword(ctx, userID, string(newPasswordHash)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) DeleteAccount(ctx context.Context, userID users.UserID) error {
	if err := s.usersRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user account: %w", err)
	}

	return nil
}

// RequestPasswordReset sends a password reset email to the user
func (s *AuthService) RequestPasswordReset(ctx context.Context, cmd RequestPasswordResetCommand) error {
	// Validate input
	if cmd.Email == "" {
		return errors.New("email is required")
	}

	// Create email value object
	email, err := users.NewEmail(cmd.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	// Check if user exists
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == ErrUserNotFound {
			// Don't reveal if user exists or not for security
			return nil
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	// Check if user is active
	if user.Status() != users.UserStatusActive {
		// Don't reveal user status for security
		return nil
	}

	// Generate reset token
	resetToken, err := s.generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token (we'll need to add this to the database)
	// For now, we'll just log it - in a real implementation, you'd store it with expiration
	slog.Info("Password reset requested", "email", cmd.Email, "token", resetToken)

	// TODO: Send email with reset link
	// This will be implemented when we add email sending functionality

	return nil
}

// ResetPassword resets the user's password using a valid reset token
func (s *AuthService) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	// Validate input
	if cmd.Token == "" || cmd.Password == "" {
		return errors.New("token and password are required")
	}

	if len(cmd.Password) < 8 {
		return ErrWeakPassword
	}

	// TODO: Validate reset token from database
	// For now, we'll accept any token for testing
	// In a real implementation, you'd check the token exists and hasn't expired

	// TODO: Get user ID from token and update password
	// For now, we'll just log the action
	slog.Info("Password reset completed", "token", cmd.Token)

	// Hash new password (placeholder - will be used when we implement token validation)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	_ = passwordHash // Placeholder to avoid unused variable error

	return nil
}
