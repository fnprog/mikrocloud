package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/projects"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dm"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

// Repository defines the interface for project persistence
type Repository interface {
	// Save persists a project
	Save(ctx context.Context, project *projects.Project) error

	// FindByID retrieves a project by its ID
	FindByID(ctx context.Context, id projects.ProjectID) (*projects.Project, error)

	// FindByName retrieves a project by its name
	FindByName(ctx context.Context, name projects.ProjectName) (*projects.Project, error)

	// FindAll retrieves all projects
	FindAll(ctx context.Context) ([]*projects.Project, error)

	// Delete removes a project
	Delete(ctx context.Context, id projects.ProjectID) error

	// Exists checks if a project exists by name
	Exists(ctx context.Context, name projects.ProjectName) (bool, error)
}

// SQLiteProjectRepository implements the project.Repository interface
type SQLiteProjectRepository struct {
	db *sql.DB
}

// NewSQLiteProjectRepository creates a new SQLite-based project repository
func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

// Save persists a project to the database
func (r *SQLiteProjectRepository) Save(ctx context.Context, proj *projects.Project) error {
	// Use Bob query builder for INSERT with ON CONFLICT (upsert)
	query := sqlite.Insert(
		im.Into("projects"),
		im.Values(sqlite.Arg(proj.ID().String()), sqlite.Arg(proj.Name().String()), sqlite.Arg(proj.Description()), sqlite.Arg(proj.CreatedAt().Format(time.RFC3339)), sqlite.Arg(proj.UpdatedAt().Format(time.RFC3339))),
		im.OnConflict("name").DoUpdate(
			im.SetCol("description").ToArg(proj.Description()),
			im.SetCol("updated_at").ToArg(proj.UpdatedAt().Format(time.RFC3339)),
		),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	return nil
}

// FindByID retrieves a project by its ID
func (r *SQLiteProjectRepository) FindByID(ctx context.Context, id projects.ProjectID) (*projects.Project, error) {
	// Use Bob query builder for SELECT
	query := sqlite.Select(
		sm.Columns("id", "name", "description", "created_at", "updated_at"),
		sm.From("projects"),
		sm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row projectRow
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(
		&row.ID, &row.Name, &row.Description, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find project by ID: %w", err)
	}

	return r.mapRowToProject(row)
}

// FindByName retrieves a project by its name
func (r *SQLiteProjectRepository) FindByName(ctx context.Context, name projects.ProjectName) (*projects.Project, error) {
	// Use Bob query builder for SELECT
	query := sqlite.Select(
		sm.Columns("id", "name", "description", "created_at", "updated_at"),
		sm.From("projects"),
		sm.Where(sqlite.Quote("name").EQ(sqlite.Arg(name.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row projectRow
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(
		&row.ID, &row.Name, &row.Description, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", name.String())
		}
		return nil, fmt.Errorf("failed to find project by name: %w", err)
	}

	return r.mapRowToProject(row)
}

// FindAll retrieves all projects
func (r *SQLiteProjectRepository) FindAll(ctx context.Context) ([]*projects.Project, error) {
	// Use Bob query builder for SELECT with ORDER BY
	query := sqlite.Select(
		sm.Columns("id", "name", "description", "created_at", "updated_at"),
		sm.From("projects"),
		sm.OrderBy("created_at").Desc(),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query all projects: %w", err)
	}
	defer rows.Close()

	var projects []*projects.Project
	for rows.Next() {
		var row projectRow
		err := rows.Scan(&row.ID, &row.Name, &row.Description, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		domainProject, err := r.mapRowToProject(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map project: %w", err)
		}

		projects = append(projects, domainProject)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over project rows: %w", err)
	}

	return projects, nil
}

// Delete removes a project
func (r *SQLiteProjectRepository) Delete(ctx context.Context, id projects.ProjectID) error {
	// Use Bob query builder for DELETE
	query := sqlite.Delete(
		dm.From("projects"),
		dm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", id.String())
	}

	return nil
}

// Exists checks if a project exists by name
func (r *SQLiteProjectRepository) Exists(ctx context.Context, name projects.ProjectName) (bool, error) {
	// Use Bob query builder for SELECT COUNT
	query := sqlite.Select(
		sm.Columns("COUNT(*)"),
		sm.From("projects"),
		sm.Where(sqlite.Quote("name").EQ(sqlite.Arg(name.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var count int
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if project exists: %w", err)
	}

	return count > 0, nil
}

// projectRow represents the database row structure
type projectRow struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
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
		projectID, projectName, row.Description, createdAt, updatedAt), nil
}
