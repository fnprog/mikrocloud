package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"

	"github.com/mikrocloud/mikrocloud/internal/infrastructure/repository"
)

type Database struct {
	db                    *sql.DB
	ProjectRepository     *repository.SQLiteProjectRepository
	EnvironmentRepository *repository.SQLiteEnvironmentRepository
	ServiceRepository     *repository.SQLiteServiceRepository
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
	projectRepo := repository.NewSQLiteProjectRepository(db)
	environmentRepo := repository.NewSQLiteEnvironmentRepository(db)
	serviceRepo := repository.NewSQLiteServiceRepository(db)

	return &Database{
		db:                    db,
		ProjectRepository:     projectRepo,
		EnvironmentRepository: environmentRepo,
		ServiceRepository:     serviceRepo,
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
