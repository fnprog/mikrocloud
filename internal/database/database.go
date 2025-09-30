package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"

	applicationsRepo "github.com/mikrocloud/mikrocloud/internal/domain/applications/repository"
	authRepo "github.com/mikrocloud/mikrocloud/internal/domain/auth/repository"
	databasesRepo "github.com/mikrocloud/mikrocloud/internal/domain/databases/repository"
	deploymentsRepo "github.com/mikrocloud/mikrocloud/internal/domain/deployments/repository"
	environmentsRepo "github.com/mikrocloud/mikrocloud/internal/domain/environments/repository"
	projectsRepo "github.com/mikrocloud/mikrocloud/internal/domain/projects/repository"
	servicesRepo "github.com/mikrocloud/mikrocloud/internal/domain/services/repository"
	usersRepo "github.com/mikrocloud/mikrocloud/internal/domain/users/repository"
)

type Database struct {
	db                    *sql.DB
	ProjectRepository     projectsRepo.Repository
	ApplicationRepository applicationsRepo.Repository
	DatabaseRepository    databasesRepo.DatabaseRepository
	EnvironmentRepository environmentsRepo.Repository
	TemplateRepository    servicesRepo.TemplateRepository
	UserRepository        usersRepo.Repository
	SessionRepository     authRepo.SessionRepository
	AuthRepository        authRepo.AuthRepository
	DeploymentRepository  deploymentsRepo.DeploymentRepository
}

func New(databaseURL string) (*Database, error) {
	// Ensure data directory exists
	if err := ensureDataDir(databaseURL); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure SQLite connection
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys and WAL mode for better performance
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	slog.Info("SQLite database connection established", "path", databaseURL)

	// Initialize repositories
	projectRepo := projectsRepo.NewSQLiteProjectRepository(db)
	applicationRepo := applicationsRepo.NewSQLiteApplicationRepository(db)
	databaseRepo := databasesRepo.NewSQLiteDatabaseRepository(db)
	environmentRepo := environmentsRepo.NewSQLiteEnvironmentRepository(db)
	templateRepo := servicesRepo.NewSQLiteTemplateRepository(db)
	userRepo := usersRepo.NewSQLiteUserRepository(db)
	sessionRepo := authRepo.NewSQLiteSessionRepository(db)
	authRepository := authRepo.NewSQLiteAuthRepository(db)
	deploymentRepo := deploymentsRepo.NewSQLiteDeploymentRepository(db)

	return &Database{
		db:                    db,
		ProjectRepository:     projectRepo,
		ApplicationRepository: applicationRepo,
		DatabaseRepository:    databaseRepo,
		EnvironmentRepository: environmentRepo,
		TemplateRepository:    templateRepo,
		UserRepository:        userRepo,
		SessionRepository:     sessionRepo,
		AuthRepository:        authRepository,
		DeploymentRepository:  deploymentRepo,
	}, nil
}

func (db *Database) Close() {
	if err := db.db.Close(); err != nil {
		slog.Error("Error closing database", "error", err)
	} else {
		slog.Info("Database connection closed")
	}
}

func (db *Database) DB() *sql.DB {
	return db.db
}

// Health check method
func (db *Database) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// ensureDataDir creates the directory for the SQLite database if it doesn't exist
func ensureDataDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "." || dir == "/" {
		return nil // Current directory or root, no need to create
	}
	return os.MkdirAll(dir, 0755)
}
