package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels"
)

type SQLiteTunnelRepository struct {
	db *sql.DB
}

func NewSQLiteTunnelRepository(db *sql.DB) *SQLiteTunnelRepository {
	return &SQLiteTunnelRepository{db: db}
}

func (r *SQLiteTunnelRepository) Create(ctx context.Context, tunnel *tunnels.CloudflareTunnel) error {
	config, _ := json.Marshal(tunnel.Config())

	var projectIDStr *string
	if tunnel.ProjectID() != nil {
		idStr := tunnel.ProjectID().String()
		projectIDStr = &idStr
	}

	var lastHealthCheck *time.Time
	if tunnel.LastHealthCheck() != nil {
		lastHealthCheck = tunnel.LastHealthCheck()
	}

	query := `
		INSERT INTO cloudflare_tunnels (
			id, name, project_id, organization_id, tunnel_token, container_id,
			status, last_health_check, health_status,
			error_message, config, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		tunnel.ID().String(),
		tunnel.Name().String(),
		projectIDStr,
		tunnel.OrganizationID().String(),
		tunnel.TunnelToken().String(),
		tunnel.ContainerID(),
		string(tunnel.Status()),
		lastHealthCheck,
		string(tunnel.HealthStatus()),
		tunnel.ErrorMessage(),
		string(config),
		tunnel.CreatedBy().String(),
		tunnel.CreatedAt(),
		tunnel.UpdatedAt(),
	)

	return err
}

func (r *SQLiteTunnelRepository) GetByID(ctx context.Context, id tunnels.TunnelID) (*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id.String())
	return r.scanTunnel(row)
}

func (r *SQLiteTunnelRepository) GetByName(ctx context.Context, organizationID uuid.UUID, name tunnels.TunnelName) (*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE organization_id = ? AND name = ?
	`

	row := r.db.QueryRowContext(ctx, query, organizationID.String(), name.String())
	return r.scanTunnel(row)
}

func (r *SQLiteTunnelRepository) GetByContainerID(ctx context.Context, containerID string) (*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE container_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, containerID)
	return r.scanTunnel(row)
}

func (r *SQLiteTunnelRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE organization_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTunnels(rows)
}

func (r *SQLiteTunnelRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE project_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTunnels(rows)
}

func (r *SQLiteTunnelRepository) ListByStatus(ctx context.Context, status tunnels.TunnelStatus) ([]*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels WHERE status = ? ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTunnels(rows)
}

func (r *SQLiteTunnelRepository) ListAll(ctx context.Context) ([]*tunnels.CloudflareTunnel, error) {
	query := `
		SELECT id, name, project_id, organization_id, tunnel_token, container_id,
			   status, last_health_check, health_status,
			   error_message, config, created_by, created_at, updated_at
		FROM cloudflare_tunnels ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTunnels(rows)
}

func (r *SQLiteTunnelRepository) Update(ctx context.Context, tunnel *tunnels.CloudflareTunnel) error {
	config, _ := json.Marshal(tunnel.Config())

	var projectIDStr *string
	if tunnel.ProjectID() != nil {
		idStr := tunnel.ProjectID().String()
		projectIDStr = &idStr
	}

	query := `
		UPDATE cloudflare_tunnels SET
			name = ?, project_id = ?, organization_id = ?, tunnel_token = ?,
			container_id = ?, status = ?,
			last_health_check = ?, health_status = ?, error_message = ?,
			config = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		tunnel.Name().String(),
		projectIDStr,
		tunnel.OrganizationID().String(),
		tunnel.TunnelToken().String(),
		tunnel.ContainerID(),
		string(tunnel.Status()),
		tunnel.LastHealthCheck(),
		string(tunnel.HealthStatus()),
		tunnel.ErrorMessage(),
		string(config),
		tunnel.UpdatedAt(),
		tunnel.ID().String(),
	)

	return err
}

func (r *SQLiteTunnelRepository) Delete(ctx context.Context, id tunnels.TunnelID) error {
	query := `DELETE FROM cloudflare_tunnels WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *SQLiteTunnelRepository) Exists(ctx context.Context, id tunnels.TunnelID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM cloudflare_tunnels WHERE id = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	return exists, err
}

func (r *SQLiteTunnelRepository) ExistsByName(ctx context.Context, organizationID uuid.UUID, name tunnels.TunnelName) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM cloudflare_tunnels WHERE organization_id = ? AND name = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, organizationID.String(), name.String()).Scan(&exists)
	return exists, err
}

func (r *SQLiteTunnelRepository) scanTunnel(row *sql.Row) (*tunnels.CloudflareTunnel, error) {
	var (
		id              string
		name            string
		projectID       *string
		organizationID  string
		tunnelToken     string
		containerID     string
		status          string
		lastHealthCheck *time.Time
		healthStatus    string
		errorMessage    string
		configJSON      string
		createdBy       string
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := row.Scan(
		&id, &name, &projectID, &organizationID, &tunnelToken, &containerID,
		&status, &lastHealthCheck, &healthStatus,
		&errorMessage, &configJSON, &createdBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tunnel not found")
		}
		return nil, err
	}

	return r.reconstructTunnel(
		id, name, projectID, organizationID, tunnelToken, containerID,
		status, lastHealthCheck, healthStatus,
		errorMessage, configJSON, createdBy, createdAt, updatedAt,
	)
}

func (r *SQLiteTunnelRepository) scanTunnels(rows *sql.Rows) ([]*tunnels.CloudflareTunnel, error) {
	var result []*tunnels.CloudflareTunnel

	for rows.Next() {
		var (
			id              string
			name            string
			projectID       *string
			organizationID  string
			tunnelToken     string
			containerID     string
			status          string
			lastHealthCheck *time.Time
			healthStatus    string
			errorMessage    string
			configJSON      string
			createdBy       string
			createdAt       time.Time
			updatedAt       time.Time
		)

		err := rows.Scan(
			&id, &name, &projectID, &organizationID, &tunnelToken, &containerID,
			&status, &lastHealthCheck, &healthStatus,
			&errorMessage, &configJSON, &createdBy, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		tunnel, err := r.reconstructTunnel(
			id, name, projectID, organizationID, tunnelToken, containerID,
			status, lastHealthCheck, healthStatus,
			errorMessage, configJSON, createdBy, createdAt, updatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, tunnel)
	}

	return result, rows.Err()
}

func (r *SQLiteTunnelRepository) reconstructTunnel(
	id, name string,
	projectID *string,
	organizationID, tunnelToken, containerID string,
	status string,
	lastHealthCheck *time.Time,
	healthStatus, errorMessage, configJSON, createdBy string,
	createdAt, updatedAt time.Time,
) (*tunnels.CloudflareTunnel, error) {
	tunnelID, err := tunnels.TunnelIDFromString(id)
	if err != nil {
		return nil, err
	}

	tunnelName, err := tunnels.NewTunnelName(name)
	if err != nil {
		return nil, err
	}

	var projID *uuid.UUID
	if projectID != nil {
		parsed, err := uuid.Parse(*projectID)
		if err != nil {
			return nil, err
		}
		projID = &parsed
	}

	orgID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, err
	}

	token, err := tunnels.NewTunnelToken(tunnelToken)
	if err != nil {
		return nil, err
	}

	createdByID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, err
	}

	var config map[string]string
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		config = make(map[string]string)
	}

	return tunnels.ReconstructCloudflareTunnel(
		tunnelID,
		tunnelName,
		projID,
		orgID,
		token,
		containerID,
		tunnels.TunnelStatus(status),
		lastHealthCheck,
		tunnels.HealthStatus(healthStatus),
		errorMessage,
		config,
		createdByID,
		createdAt,
		updatedAt,
	), nil
}
