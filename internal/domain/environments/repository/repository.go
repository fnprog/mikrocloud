package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/environments"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dm"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
)

// Repository interface for environment persistence
type Repository interface {
	Save(environment *environments.Environment) error
	FindByID(id environments.EnvironmentID) (*environments.Environment, error)
	FindByName(name environments.EnvironmentName, projectID string) (*environments.Environment, error)
	FindByProjectID(projectID string) ([]*environments.Environment, error)
	Update(environment *environments.Environment) error
	Delete(id environments.EnvironmentID) error
	List() ([]*environments.Environment, error)
}

// SQLiteEnvironmentRepository implements the environment.Repository interface
type SQLiteEnvironmentRepository struct {
	db *sql.DB
}

// NewSQLiteEnvironmentRepository creates a new SQLite-based environment repository
func NewSQLiteEnvironmentRepository(db *sql.DB) *SQLiteEnvironmentRepository {
	return &SQLiteEnvironmentRepository{db: db}
}

// Save persists an environment to the database
func (r *SQLiteEnvironmentRepository) Save(env *environments.Environment) error {
	ctx := context.Background()

	// Marshal variables to JSON
	variablesJSON, err := json.Marshal(env.Variables())
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	// Use Bob query builder for INSERT with ON CONFLICT
	query := sqlite.Insert(
		im.Into("environments", "id", "name", "project_id", "description", "variables", "created_at", "updated_at"),
		im.Values(sqlite.Arg(
			env.ID().String(),
			env.Name().String(),
			env.ProjectID().String(),
			env.Description(),
			string(variablesJSON),
			env.CreatedAt().Format(time.RFC3339),
			env.UpdatedAt().Format(time.RFC3339),
		)),
		im.OnConflict("id").DoUpdate(
			im.SetCol("description").ToArg(env.Description()),
			im.SetCol("variables").ToArg(string(variablesJSON)),
			im.SetCol("updated_at").ToArg(env.UpdatedAt().Format(time.RFC3339)),
		),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to save environment: %w", err)
	}

	return nil
}

// FindByID retrieves an environment by its ID
func (r *SQLiteEnvironmentRepository) FindByID(id environments.EnvironmentID) (*environments.Environment, error) {
	ctx := context.Background()

	// Use Bob query builder for SELECT
	query := sqlite.Select(
		sm.Columns("id", "name", "project_id", "description", "variables", "created_at", "updated_at"),
		sm.From("environments"),
		sm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row environmentRow
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(
		&row.ID, &row.Name, &row.ProjectID, &row.Description, &row.Variables, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("environment not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find environment by ID: %w", err)
	}

	return r.mapRowToEnvironment(row)
}

// FindByName retrieves an environment by its name within a project
func (r *SQLiteEnvironmentRepository) FindByName(name environments.EnvironmentName, projectID string) (*environments.Environment, error) {
	ctx := context.Background()

	// Use Bob query builder for SELECT
	query := sqlite.Select(
		sm.Columns("id", "name", "project_id", "description", "variables", "created_at", "updated_at"),
		sm.From("environments"),
		sm.Where(sqlite.Quote("name").EQ(sqlite.Arg(name.String())).And(sqlite.Quote("project_id").EQ(sqlite.Arg(projectID)))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row environmentRow
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(
		&row.ID, &row.Name, &row.ProjectID, &row.Description, &row.Variables, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("environment not found: %s in project %s", name.String(), projectID)
		}
		return nil, fmt.Errorf("failed to find environment by name: %w", err)
	}

	return r.mapRowToEnvironment(row)
}

// FindByProjectID retrieves all environments in a project
func (r *SQLiteEnvironmentRepository) FindByProjectID(projectID string) ([]*environments.Environment, error) {
	ctx := context.Background()

	// Use Bob query builder for SELECT with WHERE and ORDER BY
	query := sqlite.Select(
		sm.Columns("id", "name", "project_id", "description", "variables", "created_at", "updated_at"),
		sm.From("environments"),
		sm.Where(sqlite.Quote("project_id").EQ(sqlite.Arg(projectID))),
		sm.OrderBy("created_at").Asc(),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query environments by project ID: %w", err)
	}
	defer rows.Close()

	var environments []*environments.Environment
	for rows.Next() {
		var row environmentRow
		err := rows.Scan(&row.ID, &row.Name, &row.ProjectID, &row.Description, &row.Variables, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}

		domainEnvironment, err := r.mapRowToEnvironment(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map environment: %w", err)
		}

		environments = append(environments, domainEnvironment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over environment rows: %w", err)
	}

	return environments, nil
}

// Update updates an existing environment
func (r *SQLiteEnvironmentRepository) Update(env *environments.Environment) error {
	ctx := context.Background()

	// Marshal variables to JSON
	variablesJSON, err := json.Marshal(env.Variables())
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	// Use Bob query builder for UPDATE
	query := sqlite.Update(
		um.Table("environments"),
		um.SetCol("description").ToArg(env.Description()),
		um.SetCol("variables").ToArg(string(variablesJSON)),
		um.SetCol("updated_at").ToArg(env.UpdatedAt().Format(time.RFC3339)),
		um.Where(sqlite.Quote("id").EQ(sqlite.Arg(env.ID().String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to update environment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("environment not found: %s", env.ID().String())
	}

	return nil
}

// Delete removes an environment
func (r *SQLiteEnvironmentRepository) Delete(id environments.EnvironmentID) error {
	ctx := context.Background()

	// Use Bob query builder for DELETE
	query := sqlite.Delete(
		dm.From("environments"),
		dm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("environment not found: %s", id.String())
	}

	return nil
}

// List retrieves all environments
func (r *SQLiteEnvironmentRepository) List() ([]*environments.Environment, error) {
	ctx := context.Background()

	// Use Bob query builder for SELECT with ORDER BY
	query := sqlite.Select(
		sm.Columns("id", "name", "project_id", "description", "variables", "created_at", "updated_at"),
		sm.From("environments"),
		sm.OrderBy("created_at").Desc(),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query all environments: %w", err)
	}
	defer rows.Close()

	var environments []*environments.Environment
	for rows.Next() {
		var row environmentRow
		err := rows.Scan(&row.ID, &row.Name, &row.ProjectID, &row.Description, &row.Variables, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan environment row: %w", err)
		}

		domainEnvironment, err := r.mapRowToEnvironment(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map environment: %w", err)
		}

		environments = append(environments, domainEnvironment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over environment rows: %w", err)
	}

	return environments, nil
}

// environmentRow represents the database row structure
type environmentRow struct {
	ID          string
	Name        string
	ProjectID   string
	Description string
	Variables   string
	CreatedAt   string
	UpdatedAt   string
}

// mapRowToEnvironment converts a database row to a domain Environment
func (r *SQLiteEnvironmentRepository) mapRowToEnvironment(row environmentRow) (*environments.Environment, error) {
	// Parse environment ID
	envID, err := environments.EnvironmentIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid environment ID: %w", err)
	}

	// Parse environment name
	envName, err := environments.NewEnvironmentName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid environment name: %w", err)
	}

	// Parse project ID
	projectID := uuid.MustParse(row.ProjectID)

	// Parse variables
	var variables map[string]string
	if row.Variables != "" {
		err = json.Unmarshal([]byte(row.Variables), &variables)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
		}
	} else {
		variables = make(map[string]string)
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

	// Reconstruct environment from persistence
	env := environments.ReconstructEnvironment(
		envID, envName, projectID, row.Description, variables, createdAt, updatedAt)

	return env, nil
}
