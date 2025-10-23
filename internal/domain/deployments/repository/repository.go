package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/im"
)

type DeploymentWithMetadata struct {
	Deployment *deployments.Deployment
	Username   *string
}

type DeploymentRepository interface {
	Create(ctx context.Context, deployment *deployments.Deployment) error
	GetByID(ctx context.Context, id deployments.DeploymentID) (*deployments.Deployment, error)
	GetByIDWithMetadata(ctx context.Context, id deployments.DeploymentID) (*DeploymentWithMetadata, error)
	Update(ctx context.Context, deployment *deployments.Deployment) error
	Delete(ctx context.Context, id deployments.DeploymentID) error
	List(ctx context.Context) ([]*deployments.Deployment, error)
	ListWithMetadata(ctx context.Context) ([]*DeploymentWithMetadata, error)
	ListByApplication(ctx context.Context, applicationID applications.ApplicationID) ([]*deployments.Deployment, error)
	ListByApplicationWithMetadata(ctx context.Context, applicationID applications.ApplicationID) ([]*DeploymentWithMetadata, error)
	GetLatestByApplication(ctx context.Context, applicationID applications.ApplicationID) (*deployments.Deployment, error)
	ListByStatus(ctx context.Context, status deployments.DeploymentStatus) ([]*deployments.Deployment, error)
}

type sqliteDeploymentRepository struct {
	db *sql.DB
}

func NewSQLiteDeploymentRepository(db *sql.DB) DeploymentRepository {
	return &sqliteDeploymentRepository{db: db}
}

func (r *sqliteDeploymentRepository) Create(ctx context.Context, deployment *deployments.Deployment) error {
	var triggeredBy *string
	if deployment.TriggeredBy() != nil {
		userID := deployment.TriggeredBy().String()
		triggeredBy = &userID

		var exists bool
		err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check user existence: %w", err)
		}
		if !exists {
			return fmt.Errorf("triggered_by user does not exist: %s", userID)
		}
	}

	var appExists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM applications WHERE id = ?)", deployment.ApplicationID().String()).Scan(&appExists)
	if err != nil {
		return fmt.Errorf("failed to check application existence: %w", err)
	}
	if !appExists {
		return fmt.Errorf("application does not exist: %s", deployment.ApplicationID().String())
	}

	query := sqlite.Insert(
		im.Into("deployments"),
		im.Values(
			sqlite.Arg(deployment.ID().String()),
			sqlite.Arg(deployment.ApplicationID().String()),
			sqlite.Arg(deployment.DeploymentNumber()),
			sqlite.Arg(boolToInt(deployment.IsProduction())),
			sqlite.Arg(triggeredBy),
			sqlite.Arg(string(deployment.TriggerType())),
			sqlite.Arg(string(deployment.Status())),
			sqlite.Arg(deployment.ContainerID()),
			sqlite.Arg(deployment.ImageTag()),
			sqlite.Arg(deployment.ImageDigest()),
			sqlite.Arg(deployment.GitCommitHash()),
			sqlite.Arg(deployment.GitCommitMessage()),
			sqlite.Arg(deployment.GitBranch()),
			sqlite.Arg(deployment.GitAuthorName()),
			sqlite.Arg(deployment.BuildLogs()),
			sqlite.Arg(deployment.DeployLogs()),
			sqlite.Arg(deployment.ErrorMessage()),
			sqlite.Arg(deployment.StartedAt().Format(time.RFC3339)),
			sqlite.Arg(formatTimePtr(deployment.BuildStartedAt())),
			sqlite.Arg(formatTimePtr(deployment.BuildCompletedAt())),
			sqlite.Arg(formatTimePtr(deployment.DeployStartedAt())),
			sqlite.Arg(formatTimePtr(deployment.DeployCompletedAt())),
			sqlite.Arg(formatTimePtr(deployment.StoppedAt())),
			sqlite.Arg(deployment.BuildDurationSeconds()),
			sqlite.Arg(deployment.DeployDurationSeconds()),
			sqlite.Arg(deployment.UpdatedAt().Format(time.RFC3339)),
		),
	)

	queryStr, args, err := query.Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

func (r *sqliteDeploymentRepository) GetByID(ctx context.Context, id deployments.DeploymentID) (*deployments.Deployment, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, COALESCE(u.username, u.email) AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.id = ?`

	row := deploymentRow{}
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
		&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
		&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
		&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
		&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
		&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
		&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deployment not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return r.mapRowToDeployment(row)
}

func (r *sqliteDeploymentRepository) GetByIDWithMetadata(ctx context.Context, id deployments.DeploymentID) (*DeploymentWithMetadata, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, COALESCE(u.username, u.email) AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.id = ?`

	row := deploymentRow{}
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
		&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
		&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
		&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
		&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
		&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
		&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deployment not found: %s", id.String())
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	deployment, err := r.mapRowToDeployment(row)
	if err != nil {
		return nil, err
	}

	return &DeploymentWithMetadata{
		Deployment: deployment,
		Username:   row.TriggeredByUsername,
	}, nil
}

func (r *sqliteDeploymentRepository) Update(ctx context.Context, deployment *deployments.Deployment) error {
	var triggeredBy *string
	if deployment.TriggeredBy() != nil {
		userID := deployment.TriggeredBy().String()
		triggeredBy = &userID
	}

	// Use direct SQL for update since Bob's update is complex
	query := `UPDATE deployments SET 
		deployment_number = ?, is_production = ?, triggered_by = ?, trigger_type = ?, 
		status = ?, container_id = ?, image_tag = ?, image_digest = ?,
		git_commit_hash = ?, git_commit_message = ?, git_branch = ?, git_author_name = ?,
		build_logs = ?, deploy_logs = ?, error_message = ?, started_at = ?,
		build_started_at = ?, build_completed_at = ?, deploy_started_at = ?, 
		deploy_completed_at = ?, stopped_at = ?, build_duration_seconds = ?,
		deploy_duration_seconds = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		deployment.DeploymentNumber(),
		boolToInt(deployment.IsProduction()),
		triggeredBy,
		string(deployment.TriggerType()),
		string(deployment.Status()),
		deployment.ContainerID(),
		deployment.ImageTag(),
		deployment.ImageDigest(),
		deployment.GitCommitHash(),
		deployment.GitCommitMessage(),
		deployment.GitBranch(),
		deployment.GitAuthorName(),
		deployment.BuildLogs(),
		deployment.DeployLogs(),
		deployment.ErrorMessage(),
		deployment.StartedAt().Format(time.RFC3339),
		formatTimePtr(deployment.BuildStartedAt()),
		formatTimePtr(deployment.BuildCompletedAt()),
		formatTimePtr(deployment.DeployStartedAt()),
		formatTimePtr(deployment.DeployCompletedAt()),
		formatTimePtr(deployment.StoppedAt()),
		deployment.BuildDurationSeconds(),
		deployment.DeployDurationSeconds(),
		deployment.UpdatedAt().Format(time.RFC3339),
		deployment.ID().String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (r *sqliteDeploymentRepository) Delete(ctx context.Context, id deployments.DeploymentID) error {
	query := `DELETE FROM deployments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}
	return nil
}

func (r *sqliteDeploymentRepository) List(ctx context.Context) ([]*deployments.Deployment, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, u.username AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	ORDER BY d.started_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer rows.Close()

	var result []*deployments.Deployment
	for rows.Next() {
		row := deploymentRow{}
		err := rows.Scan(
			&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
			&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
			&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
			&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
			&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
			&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
			&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		deployment, err := r.mapRowToDeployment(row)
		if err != nil {
			return nil, err
		}
		result = append(result, deployment)
	}

	return result, nil
}

func (r *sqliteDeploymentRepository) ListWithMetadata(ctx context.Context) ([]*DeploymentWithMetadata, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, u.username AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	ORDER BY d.started_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer rows.Close()

	var result []*DeploymentWithMetadata
	for rows.Next() {
		row := deploymentRow{}
		err := rows.Scan(
			&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
			&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
			&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
			&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
			&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
			&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
			&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		deployment, err := r.mapRowToDeployment(row)
		if err != nil {
			return nil, err
		}
		result = append(result, &DeploymentWithMetadata{
			Deployment: deployment,
			Username:   row.TriggeredByUsername,
		})
	}

	return result, nil
}

func (r *sqliteDeploymentRepository) ListByApplication(ctx context.Context, applicationID applications.ApplicationID) ([]*deployments.Deployment, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, u.username AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.application_id = ?
	ORDER BY d.deployment_number DESC`

	rows, err := r.db.QueryContext(ctx, query, applicationID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by application: %w", err)
	}
	defer rows.Close()

	var result []*deployments.Deployment
	for rows.Next() {
		row := deploymentRow{}
		err := rows.Scan(
			&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
			&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
			&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
			&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
			&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
			&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
			&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		deployment, err := r.mapRowToDeployment(row)
		if err != nil {
			return nil, err
		}
		result = append(result, deployment)
	}

	return result, nil
}

func (r *sqliteDeploymentRepository) ListByApplicationWithMetadata(ctx context.Context, applicationID applications.ApplicationID) ([]*DeploymentWithMetadata, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, COALESCE(u.username, u.email) AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.application_id = ?
	ORDER BY d.deployment_number DESC`

	rows, err := r.db.QueryContext(ctx, query, applicationID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by application: %w", err)
	}
	defer rows.Close()

	var result []*DeploymentWithMetadata
	for rows.Next() {
		row := deploymentRow{}
		err := rows.Scan(
			&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
			&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
			&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
			&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
			&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
			&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
			&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		deployment, err := r.mapRowToDeployment(row)
		if err != nil {
			return nil, err
		}
		result = append(result, &DeploymentWithMetadata{
			Deployment: deployment,
			Username:   row.TriggeredByUsername,
		})
	}

	return result, nil
}

func (r *sqliteDeploymentRepository) GetLatestByApplication(ctx context.Context, applicationID applications.ApplicationID) (*deployments.Deployment, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, COALESCE(u.username, u.email) AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.application_id = ?
	ORDER BY d.deployment_number DESC
	LIMIT 1`

	row := deploymentRow{}
	err := r.db.QueryRowContext(ctx, query, applicationID.String()).Scan(
		&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
		&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
		&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
		&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
		&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
		&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
		&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no deployments found for application: %s", applicationID.String())
		}
		return nil, fmt.Errorf("failed to get latest deployment: %w", err)
	}

	return r.mapRowToDeployment(row)
}

func (r *sqliteDeploymentRepository) ListByStatus(ctx context.Context, status deployments.DeploymentStatus) ([]*deployments.Deployment, error) {
	query := `SELECT 
		d.id, d.application_id, d.deployment_number, d.is_production, d.triggered_by,
		d.trigger_type, d.status, d.container_id, d.image_tag, d.image_digest,
		d.git_commit_hash, d.git_commit_message, d.git_branch, d.git_author_name,
		d.build_logs, d.deploy_logs, d.error_message, d.started_at,
		d.build_started_at, d.build_completed_at, d.deploy_started_at,
		d.deploy_completed_at, d.stopped_at, d.build_duration_seconds,
		d.deploy_duration_seconds, d.updated_at, u.username AS triggered_by_username
	FROM deployments d
	LEFT JOIN users u ON d.triggered_by = u.id
	WHERE d.status = ?
	ORDER BY d.started_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by status: %w", err)
	}
	defer rows.Close()

	var result []*deployments.Deployment
	for rows.Next() {
		row := deploymentRow{}
		err := rows.Scan(
			&row.ID, &row.ApplicationID, &row.DeploymentNumber, &row.IsProduction, &row.TriggeredBy,
			&row.TriggerType, &row.Status, &row.ContainerID, &row.ImageTag, &row.ImageDigest,
			&row.GitCommitHash, &row.GitCommitMessage, &row.GitBranch, &row.GitAuthorName,
			&row.BuildLogs, &row.DeployLogs, &row.ErrorMessage, &row.StartedAt,
			&row.BuildStartedAt, &row.BuildCompletedAt, &row.DeployStartedAt,
			&row.DeployCompletedAt, &row.StoppedAt, &row.BuildDurationSeconds,
			&row.DeployDurationSeconds, &row.UpdatedAt, &row.TriggeredByUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deployment row: %w", err)
		}

		deployment, err := r.mapRowToDeployment(row)
		if err != nil {
			return nil, err
		}
		result = append(result, deployment)
	}

	return result, nil
}

type deploymentRow struct {
	ID                    string
	ApplicationID         string
	DeploymentNumber      int
	IsProduction          int
	TriggeredBy           *string
	TriggeredByUsername   *string
	TriggerType           string
	Status                string
	ContainerID           string
	ImageTag              string
	ImageDigest           string
	GitCommitHash         string
	GitCommitMessage      string
	GitBranch             string
	GitAuthorName         string
	BuildLogs             string
	DeployLogs            string
	ErrorMessage          string
	StartedAt             string
	BuildStartedAt        *string
	BuildCompletedAt      *string
	DeployStartedAt       *string
	DeployCompletedAt     *string
	StoppedAt             *string
	BuildDurationSeconds  *int
	DeployDurationSeconds *int
	UpdatedAt             string
}

func (r *sqliteDeploymentRepository) mapRowToDeployment(row deploymentRow) (*deployments.Deployment, error) {
	deploymentID, err := deployments.DeploymentIDFromString(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid deployment ID: %w", err)
	}

	appID, err := applications.ApplicationIDFromString(row.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application ID: %w", err)
	}

	var triggeredBy *users.UserID
	if row.TriggeredBy != nil {
		userID, err := users.UserIDFromString(*row.TriggeredBy)
		if err != nil {
			return nil, fmt.Errorf("invalid triggered by user ID: %w", err)
		}
		triggeredBy = &userID
	}

	startedAt, err := time.Parse(time.RFC3339, row.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid started at time: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, row.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated at time: %w", err)
	}

	var buildStartedAt, buildCompletedAt, deployStartedAt, deployCompletedAt, stoppedAt *time.Time

	if row.BuildStartedAt != nil {
		t, err := time.Parse(time.RFC3339, *row.BuildStartedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid build started at time: %w", err)
		}
		buildStartedAt = &t
	}

	if row.BuildCompletedAt != nil {
		t, err := time.Parse(time.RFC3339, *row.BuildCompletedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid build completed at time: %w", err)
		}
		buildCompletedAt = &t
	}

	if row.DeployStartedAt != nil {
		t, err := time.Parse(time.RFC3339, *row.DeployStartedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid deploy started at time: %w", err)
		}
		deployStartedAt = &t
	}

	if row.DeployCompletedAt != nil {
		t, err := time.Parse(time.RFC3339, *row.DeployCompletedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid deploy completed at time: %w", err)
		}
		deployCompletedAt = &t
	}

	if row.StoppedAt != nil {
		t, err := time.Parse(time.RFC3339, *row.StoppedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid stopped at time: %w", err)
		}
		stoppedAt = &t
	}

	return deployments.ReconstructDeployment(
		deploymentID,
		appID,
		row.DeploymentNumber,
		intToBool(row.IsProduction),
		triggeredBy,
		deployments.TriggerType(row.TriggerType),
		deployments.DeploymentStatus(row.Status),
		row.ContainerID,
		row.ImageTag,
		row.ImageDigest,
		row.GitCommitHash,
		row.GitCommitMessage,
		row.GitBranch,
		row.GitAuthorName,
		row.BuildLogs,
		row.DeployLogs,
		row.ErrorMessage,
		startedAt,
		buildStartedAt,
		buildCompletedAt,
		deployStartedAt,
		deployCompletedAt,
		stoppedAt,
		row.BuildDurationSeconds,
		row.DeployDurationSeconds,
		updatedAt,
	), nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i != 0
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(time.RFC3339)
	return &formatted
}
