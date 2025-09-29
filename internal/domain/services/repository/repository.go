package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/services"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dm"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"github.com/stephenafamo/bob/dialect/sqlite/um"
)

// Repository interface for service persistence
type Repository interface {
	Save(service *services.Service) error
	FindByID(id services.ServiceID) (*services.Service, error)
	FindByName(name services.ServiceName, environmentID string) (*services.Service, error)
	FindByProjectID(projectID string) ([]*services.Service, error)
	FindByEnvironmentID(environmentID string) ([]*services.Service, error)
	FindByProjectAndEnvironment(projectID, environmentID string) ([]*services.Service, error)
	Update(service *services.Service) error
	Delete(id services.ServiceID) error
	List() ([]*services.Service, error)
}

// SQLiteServiceRepository implements the service.Repository interface
type SQLiteServiceRepository struct {
	db *sql.DB
}

// NewSQLiteServiceRepository creates a new SQLite-based service repository
func NewSQLiteServiceRepository(db *sql.DB) Repository {
	return &SQLiteServiceRepository{db: db}
}

// Save persists a service to the database
func (r *SQLiteServiceRepository) Save(svc *services.Service) error {
	ctx := context.Background()

	// Marshal GitURL to JSON
	var gitURLJSON string
	if svc.GitURL() != nil {
		gitURLData := map[string]string{
			"url":          svc.GitURL().URL(),
			"branch":       svc.GitURL().Branch(),
			"context_root": svc.GitURL().ContextRoot(),
		}
		data, err := json.Marshal(gitURLData)
		if err != nil {
			return fmt.Errorf("failed to marshal git URL: %w", err)
		}
		gitURLJSON = string(data)
	}

	// Marshal BuildConfig to JSON
	var buildConfigJSON string
	if svc.BuildConfig() != nil {
		buildConfigData := map[string]interface{}{
			"buildpack_type": string(svc.BuildConfig().BuildpackType()),
		}
		if svc.BuildConfig().NixpacksConfig() != nil {
			buildConfigData["nixpacks"] = svc.BuildConfig().NixpacksConfig()
		}
		if svc.BuildConfig().StaticConfig() != nil {
			buildConfigData["static"] = svc.BuildConfig().StaticConfig()
		}
		if svc.BuildConfig().DockerfileConfig() != nil {
			buildConfigData["dockerfile"] = svc.BuildConfig().DockerfileConfig()
		}
		if svc.BuildConfig().ComposeConfig() != nil {
			buildConfigData["compose"] = svc.BuildConfig().ComposeConfig()
		}
		data, err := json.Marshal(buildConfigData)
		if err != nil {
			return fmt.Errorf("failed to marshal build config: %w", err)
		}
		buildConfigJSON = string(data)
	}

	// Marshal environment variables to JSON
	envJSON, err := json.Marshal(svc.Environment())
	if err != nil {
		return fmt.Errorf("failed to marshal environment: %w", err)
	}

	// Use Bob query builder for INSERT with ON CONFLICT (upsert)
	query := sqlite.Insert(
		im.Into("services"),
		im.Values(sqlite.Arg(svc.ID().String()), sqlite.Arg(svc.Name().String()), sqlite.Arg(svc.EnvironmentID().String()), sqlite.Arg(svc.ProjectID().String()), sqlite.Arg(gitURLJSON), sqlite.Arg(buildConfigJSON), sqlite.Arg(string(envJSON)), sqlite.Arg(svc.Status().String()), sqlite.Arg(svc.CreatedAt().Format(time.RFC3339)), sqlite.Arg(svc.UpdatedAt().Format(time.RFC3339))),
		im.OnConflict("id").DoUpdate(
			im.SetCol("git_url").ToArg(gitURLJSON),
			im.SetCol("build_config").ToArg(buildConfigJSON),
			im.SetCol("environment").ToArg(string(envJSON)),
			im.SetCol("status").ToArg(svc.Status().String()),
			im.SetCol("updated_at").ToArg(svc.UpdatedAt().Format(time.RFC3339)),
		),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to save service: %w", err)
	}

	return nil
}

// FindByID retrieves a service by its ID
func (r *SQLiteServiceRepository) FindByID(id services.ServiceID) (*services.Service, error) {
	ctx := context.Background()

	// Use Bob query builder for SELECT
	query := sqlite.Select(
		sm.Columns("id", "name", "environment_id", "project_id", "git_url", "build_config", "environment", "status", "created_at", "updated_at"),
		sm.From("services"),
		sm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var row serviceRow
	err = r.db.QueryRowContext(ctx, queryStr, args...).Scan(
		&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to find service by ID: %w", err)
	}

	return r.mapRowToService(row)
}

// FindByName retrieves a service by its name within an environment
func (r *SQLiteServiceRepository) FindByName(name services.ServiceName, environmentID string) (*services.Service, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, environment_id, project_id, git_url, build_config, environment, status, created_at, updated_at
		FROM services 
		WHERE name = ? AND environment_id = ?`

	var row serviceRow
	err := r.db.QueryRowContext(ctx, query, name.String(), environmentID).Scan(
		&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found: %s in environment %s", name.String(), environmentID)
		}
		return nil, fmt.Errorf("failed to find service by name: %w", err)
	}

	return r.mapRowToService(row)
}

// FindByEnvironmentID retrieves all services in an environment
func (r *SQLiteServiceRepository) FindByEnvironmentID(environmentID string) ([]*services.Service, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, environment_id, project_id, git_url, build_config, environment, status, created_at, updated_at
		FROM services 
		WHERE environment_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services by environment ID: %w", err)
	}
	defer rows.Close()

	var services []*services.Service
	for rows.Next() {
		var row serviceRow
		err := rows.Scan(&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}

		domainService, err := r.mapRowToService(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map service: %w", err)
		}

		services = append(services, domainService)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over service rows: %w", err)
	}

	return services, nil
}

// FindByProjectAndEnvironment retrieves all services in a specific project and environment
func (r *SQLiteServiceRepository) FindByProjectAndEnvironment(projectID, environmentID string) ([]*services.Service, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, environment_id, project_id, git_url, build_config, environment, status, created_at, updated_at
		FROM services 
		WHERE project_id = ? AND environment_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, projectID, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services by project and environment: %w", err)
	}
	defer rows.Close()

	var services []*services.Service
	for rows.Next() {
		var row serviceRow
		err := rows.Scan(&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}

		domainService, err := r.mapRowToService(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map service: %w", err)
		}

		services = append(services, domainService)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over service rows: %w", err)
	}

	return services, nil
}

// FindByProjectID retrieves all services in a project
func (r *SQLiteServiceRepository) FindByProjectID(projectID string) ([]*services.Service, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, environment_id, project_id, git_url, build_config, environment, status, created_at, updated_at
		FROM services 
		WHERE project_id = ?
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services by project ID: %w", err)
	}
	defer rows.Close()

	var services []*services.Service
	for rows.Next() {
		var row serviceRow
		err := rows.Scan(&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}

		domainService, err := r.mapRowToService(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map service: %w", err)
		}

		services = append(services, domainService)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over service rows: %w", err)
	}

	return services, nil
}

// Update updates an existing service
func (r *SQLiteServiceRepository) Update(svc *services.Service) error {
	ctx := context.Background()

	// Marshal GitURL to JSON
	var gitURLJSON string
	if svc.GitURL() != nil {
		gitURLData := map[string]string{
			"url":          svc.GitURL().URL(),
			"branch":       svc.GitURL().Branch(),
			"context_root": svc.GitURL().ContextRoot(),
		}
		data, err := json.Marshal(gitURLData)
		if err != nil {
			return fmt.Errorf("failed to marshal git URL: %w", err)
		}
		gitURLJSON = string(data)
	}

	// Marshal BuildConfig to JSON
	var buildConfigJSON string
	if svc.BuildConfig() != nil {
		buildConfigData := map[string]interface{}{
			"buildpack_type": string(svc.BuildConfig().BuildpackType()),
		}
		if svc.BuildConfig().NixpacksConfig() != nil {
			buildConfigData["nixpacks"] = svc.BuildConfig().NixpacksConfig()
		}
		if svc.BuildConfig().StaticConfig() != nil {
			buildConfigData["static"] = svc.BuildConfig().StaticConfig()
		}
		if svc.BuildConfig().DockerfileConfig() != nil {
			buildConfigData["dockerfile"] = svc.BuildConfig().DockerfileConfig()
		}
		if svc.BuildConfig().ComposeConfig() != nil {
			buildConfigData["compose"] = svc.BuildConfig().ComposeConfig()
		}
		data, err := json.Marshal(buildConfigData)
		if err != nil {
			return fmt.Errorf("failed to marshal build config: %w", err)
		}
		buildConfigJSON = string(data)
	}

	// Marshal environment variables to JSON
	envJSON, err := json.Marshal(svc.Environment())
	if err != nil {
		return fmt.Errorf("failed to marshal environment: %w", err)
	}

	// Use Bob query builder for UPDATE
	query := sqlite.Update(
		um.Table("services"),
		um.SetCol("git_url").ToArg(gitURLJSON),
		um.SetCol("build_config").ToArg(buildConfigJSON),
		um.SetCol("environment").ToArg(string(envJSON)),
		um.SetCol("status").ToArg(svc.Status().String()),
		um.SetCol("updated_at").ToArg(svc.UpdatedAt().Format(time.RFC3339)),
		um.Where(sqlite.Quote("id").EQ(sqlite.Arg(svc.ID().String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("service not found: %s", svc.ID().String())
	}

	return nil
}

// Delete removes a service
func (r *SQLiteServiceRepository) Delete(id services.ServiceID) error {
	ctx := context.Background()

	// Use Bob query builder for DELETE
	query := sqlite.Delete(
		dm.From("services"),
		dm.Where(sqlite.Quote("id").EQ(sqlite.Arg(id.String()))),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("service not found: %s", id.String())
	}

	return nil
}

// List retrieves all services
func (r *SQLiteServiceRepository) List() ([]*services.Service, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, environment_id, project_id, git_url, build_config, environment, status, created_at, updated_at
		FROM services 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all services: %w", err)
	}
	defer rows.Close()

	var services []*services.Service
	for rows.Next() {
		var row serviceRow
		err := rows.Scan(&row.ID, &row.Name, &row.EnvironmentID, &row.ProjectID, &row.GitURL, &row.BuildConfig, &row.Environment, &row.Status, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}

		domainService, err := r.mapRowToService(row)
		if err != nil {
			return nil, fmt.Errorf("failed to map service: %w", err)
		}

		services = append(services, domainService)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over service rows: %w", err)
	}

	return services, nil
}

// serviceRow represents the database row structure
type serviceRow struct {
	ID            string
	Name          string
	EnvironmentID string
	ProjectID     string
	GitURL        string
	BuildConfig   string
	Environment   string
	Status        string
	CreatedAt     string
	UpdatedAt     string
}

// mapRowToService converts a database row to a domain Service
func (r *SQLiteServiceRepository) mapRowToService(row serviceRow) (*services.Service, error) {
	// Parse service ID
	svcID, err := services.ServiceIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid service ID: %w", err)
	}

	// Parse service name
	svcName, err := services.NewServiceName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid service name: %w", err)
	}

	// Parse environment ID
	envID, err := uuid.Parse(row.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid environment ID: %w", err)
	}

	// Parse project ID
	projectID, err := uuid.Parse(row.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	// Parse GitURL
	var gitURL *services.GitURL
	if row.GitURL != "" {
		var gitURLData map[string]string
		err = json.Unmarshal([]byte(row.GitURL), &gitURLData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal git URL: %w", err)
		}
		gitURL, err = services.NewGitURL(
			gitURLData["url"],
			gitURLData["branch"],
			gitURLData["context_root"],
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create git URL: %w", err)
		}
	}

	// Parse BuildConfig
	var buildConfig *services.BuildConfig
	if row.BuildConfig != "" {
		var buildConfigData map[string]interface{}
		err = json.Unmarshal([]byte(row.BuildConfig), &buildConfigData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal build config: %w", err)
		}

		buildpackType := services.BuildpackType(buildConfigData["buildpack_type"].(string))
		buildConfig = services.NewBuildConfig(buildpackType)

		// Set specific config based on buildpack type
		switch buildpackType {
		case services.BuildpackNixpacks:
			if nixpacks, exists := buildConfigData["nixpacks"]; exists {
				configBytes, _ := json.Marshal(nixpacks)
				var nixpacksConfig services.NixpacksConfig
				json.Unmarshal(configBytes, &nixpacksConfig)
				buildConfig.SetNixpacksConfig(&nixpacksConfig)
			}
		case services.BuildpackStatic:
			if static, exists := buildConfigData["static"]; exists {
				configBytes, _ := json.Marshal(static)
				var staticConfig services.StaticConfig
				json.Unmarshal(configBytes, &staticConfig)
				buildConfig.SetStaticConfig(&staticConfig)
			}
		case services.BuildpackDockerfile:
			if dockerfile, exists := buildConfigData["dockerfile"]; exists {
				configBytes, _ := json.Marshal(dockerfile)
				var dockerfileConfig services.DockerfileConfig
				json.Unmarshal(configBytes, &dockerfileConfig)
				buildConfig.SetDockerfileConfig(&dockerfileConfig)
			}
		case services.BuildpackDockerCompose:
			if compose, exists := buildConfigData["compose"]; exists {
				configBytes, _ := json.Marshal(compose)
				var composeConfig services.ComposeConfig
				json.Unmarshal(configBytes, &composeConfig)
				buildConfig.SetComposeConfig(&composeConfig)
			}
		}
	}

	// Parse environment variables
	var environment map[string]string
	if row.Environment != "" {
		err = json.Unmarshal([]byte(row.Environment), &environment)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal environment: %w", err)
		}
	} else {
		environment = make(map[string]string)
	}

	// Parse status
	var status services.ServiceStatus
	switch row.Status {
	case "created":
		status = services.ServiceStatusCreated
	case "building":
		status = services.ServiceStatusBuilding
	case "deploying":
		status = services.ServiceStatusDeploying
	case "running":
		status = services.ServiceStatusRunning
	case "stopped":
		status = services.ServiceStatusStopped
	case "failed":
		status = services.ServiceStatusFailed
	default:
		status = services.ServiceStatusCreated
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

	// Reconstruct service from persistence
	svc := services.ReconstructService(
		svcID, svcName, projectID, envID, gitURL, buildConfig, environment, status, createdAt, updatedAt)

	return svc, nil
}
