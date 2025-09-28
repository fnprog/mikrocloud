package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
)

// Repository defines the interface for project persistence
type Repository interface {
	Save(ctx context.Context, project *projects.Project) error
	FindByID(ctx context.Context, id projects.ProjectID) (*projects.Project, error)
	FindByName(ctx context.Context, name projects.ProjectName) (*projects.Project, error)
	FindAll(ctx context.Context) ([]*projects.Project, error)
	Delete(ctx context.Context, id projects.ProjectID) error
	Exists(ctx context.Context, name projects.ProjectName) (bool, error)
}

// SQLiteProjectRepository implements Repository using SQLite
type SQLiteProjectRepository struct {
	db *sql.DB
}

// NewSQLiteProjectRepository creates a new SQLite project repository
func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

// Save persists a project to the database using raw SQL
func (r *SQLiteProjectRepository) Save(ctx context.Context, project *projects.Project) error {
	query := `
		INSERT OR REPLACE INTO projects (id, name, description, user_id, organization_id, created_by, settings, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		project.ID().String(),
		project.Name().String(),
		project.Description(),
		project.UserID().String(),
		project.OrganizationID().String(),
		project.CreatedBy().String(),
		project.Settings(),
		project.CreatedAt().Format(time.RFC3339),
		project.UpdatedAt().Format(time.RFC3339),
	)

	if err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	return nil
}

// FindByID retrieves a project by its ID using raw SQL
func (r *SQLiteProjectRepository) FindByID(ctx context.Context, id projects.ProjectID) (*projects.Project, error) {
	query := `
		SELECT id, name, description, user_id, organization_id, created_by, settings, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	var row projectRow
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&row.ID,
		&row.Name,
		&row.Description,
		&row.UserID,
		&row.OrganizationID,
		&row.CreatedBy,
		&row.Settings,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find project by ID: %w", err)
	}

	return r.mapRowToProject(row)
}

// FindByName retrieves a project by its name using raw SQL
func (r *SQLiteProjectRepository) FindByName(ctx context.Context, name projects.ProjectName) (*projects.Project, error) {
	query := `
		SELECT id, name, description, user_id, organization_id, created_by, settings, created_at, updated_at
		FROM projects
		WHERE name = ?
	`

	var row projectRow
	err := r.db.QueryRowContext(ctx, query, name.String()).Scan(
		&row.ID,
		&row.Name,
		&row.Description,
		&row.UserID,
		&row.OrganizationID,
		&row.CreatedBy,
		&row.Settings,
		&row.CreatedAt,
		&row.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", name.String())
		}
		return nil, fmt.Errorf("failed to find project by name: %w", err)
	}

	return r.mapRowToProject(row)
}

// FindAll retrieves all projects using raw SQL
func (r *SQLiteProjectRepository) FindAll(ctx context.Context) ([]*projects.Project, error) {
	query := `
		SELECT id, name, description, user_id, organization_id, created_by, settings, created_at, updated_at
		FROM projects
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all projects: %w", err)
	}
	defer rows.Close()

	var projects []*projects.Project
	for rows.Next() {
		var row projectRow
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Description,
			&row.UserID,
			&row.OrganizationID,
			&row.CreatedBy,
			&row.Settings,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		project, err := r.mapRowToProject(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map project: %w", err)
		}

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project rows: %w", err)
	}

	return projects, nil
}

// Delete removes a project from the database using raw SQL
func (r *SQLiteProjectRepository) Delete(ctx context.Context, id projects.ProjectID) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", id.String())
	}

	return nil
}

// Exists checks if a project with the given name exists using raw SQL
func (r *SQLiteProjectRepository) Exists(ctx context.Context, name projects.ProjectName) (bool, error) {
	query := `SELECT COUNT(*) FROM projects WHERE name = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, name.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check project existence: %w", err)
	}

	return count > 0, nil
}

// projectRow represents the database row structure matching the schema
type projectRow struct {
	ID             string
	Name           string
	Description    string
	UserID         string
	OrganizationID string
	CreatedBy      string
	Settings       string
	CreatedAt      string
	UpdatedAt      string
}

// mapRowToProject converts a database row to a domain Project
func (r *SQLiteProjectRepository) mapRowToProject(row projectRow) (*projects.Project, error) {
	// Parse project ID
	projectID := projects.ProjectIDFromUUID(uuid.MustParse(row.ID))

	// Parse project name
	projectName, err := projects.NewProjectName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid project name: %w", err)
	}

	// Parse user IDs
	userID := users.UserIDFromUUID(uuid.MustParse(row.UserID))
	organizationID := users.OrganizationIDFromUUID(uuid.MustParse(row.OrganizationID))
	createdBy := users.UserIDFromUUID(uuid.MustParse(row.CreatedBy))

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, row.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at timestamp: %w", err)
	}

	// Reconstruct project from persistence
	return projects.ReconstructProject(
		projectID, projectName, row.Description, userID, organizationID, createdBy, row.Settings, createdAt, updatedAt), nil
}
