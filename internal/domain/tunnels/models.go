package tunnels

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CloudflareTunnel struct {
	id              TunnelID
	name            TunnelName
	projectID       *uuid.UUID
	organizationID  uuid.UUID
	tunnelToken     TunnelToken
	containerID     string
	status          TunnelStatus
	lastHealthCheck *time.Time
	healthStatus    HealthStatus
	errorMessage    string
	config          map[string]string
	createdBy       uuid.UUID
	createdAt       time.Time
	updatedAt       time.Time
}

type TunnelID struct {
	value string
}

func NewTunnelID() TunnelID {
	return TunnelID{value: uuid.Must(uuid.NewV7()).String()}
}

func TunnelIDFromString(s string) (TunnelID, error) {
	if s == "" {
		return TunnelID{}, fmt.Errorf("tunnel ID cannot be empty")
	}
	return TunnelID{value: s}, nil
}

func (id TunnelID) String() string {
	return id.value
}

type TunnelName struct {
	value string
}

func NewTunnelName(name string) (TunnelName, error) {
	if name == "" {
		return TunnelName{}, fmt.Errorf("tunnel name cannot be empty")
	}
	if len(name) > 64 {
		return TunnelName{}, fmt.Errorf("tunnel name cannot exceed 64 characters")
	}
	return TunnelName{value: name}, nil
}

func (n TunnelName) String() string {
	return n.value
}

type TunnelToken struct {
	value string
}

func NewTunnelToken(token string) (TunnelToken, error) {
	if token == "" {
		return TunnelToken{}, fmt.Errorf("tunnel token cannot be empty")
	}
	return TunnelToken{value: token}, nil
}

func (t TunnelToken) String() string {
	return t.value
}

type TunnelStatus string

const (
	TunnelStatusStopped  TunnelStatus = "stopped"
	TunnelStatusStarting TunnelStatus = "starting"
	TunnelStatusRunning  TunnelStatus = "running"
	TunnelStatusError    TunnelStatus = "error"
	TunnelStatusStopping TunnelStatus = "stopping"
)

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

func NewCloudflareTunnel(
	name TunnelName,
	projectID *uuid.UUID,
	organizationID uuid.UUID,
	tunnelToken TunnelToken,
	createdBy uuid.UUID,
) (*CloudflareTunnel, error) {
	now := time.Now()
	return &CloudflareTunnel{
		id:             NewTunnelID(),
		name:           name,
		projectID:      projectID,
		organizationID: organizationID,
		tunnelToken:    tunnelToken,
		status:         TunnelStatusStopped,
		healthStatus:   HealthStatusUnknown,
		config:         make(map[string]string),
		createdBy:      createdBy,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func (t *CloudflareTunnel) ID() TunnelID {
	return t.id
}

func (t *CloudflareTunnel) Name() TunnelName {
	return t.name
}

func (t *CloudflareTunnel) ProjectID() *uuid.UUID {
	return t.projectID
}

func (t *CloudflareTunnel) OrganizationID() uuid.UUID {
	return t.organizationID
}

func (t *CloudflareTunnel) TunnelToken() TunnelToken {
	return t.tunnelToken
}

func (t *CloudflareTunnel) ContainerID() string {
	return t.containerID
}

func (t *CloudflareTunnel) Status() TunnelStatus {
	return t.status
}

func (t *CloudflareTunnel) LastHealthCheck() *time.Time {
	return t.lastHealthCheck
}

func (t *CloudflareTunnel) HealthStatus() HealthStatus {
	return t.healthStatus
}

func (t *CloudflareTunnel) ErrorMessage() string {
	return t.errorMessage
}

func (t *CloudflareTunnel) Config() map[string]string {
	config := make(map[string]string)
	for k, v := range t.config {
		config[k] = v
	}
	return config
}

func (t *CloudflareTunnel) CreatedBy() uuid.UUID {
	return t.createdBy
}

func (t *CloudflareTunnel) CreatedAt() time.Time {
	return t.createdAt
}

func (t *CloudflareTunnel) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *CloudflareTunnel) SetContainerID(containerID string) {
	t.containerID = containerID
	t.updatedAt = time.Now()
}

func (t *CloudflareTunnel) ChangeStatus(status TunnelStatus) {
	t.status = status
	t.updatedAt = time.Now()
}

func (t *CloudflareTunnel) SetError(message string) {
	t.status = TunnelStatusError
	t.errorMessage = message
	t.healthStatus = HealthStatusUnhealthy
	t.updatedAt = time.Now()
}

func (t *CloudflareTunnel) UpdateHealth(status HealthStatus) {
	now := time.Now()
	t.healthStatus = status
	t.lastHealthCheck = &now
	t.updatedAt = now
}

func (t *CloudflareTunnel) SetConfigValue(key, value string) {
	t.config[key] = value
	t.updatedAt = time.Now()
}

func (t *CloudflareTunnel) GetContainerName() string {
	return fmt.Sprintf("cloudflared-%s", t.name.String())
}

func ReconstructCloudflareTunnel(
	id TunnelID,
	name TunnelName,
	projectID *uuid.UUID,
	organizationID uuid.UUID,
	tunnelToken TunnelToken,
	containerID string,
	status TunnelStatus,
	lastHealthCheck *time.Time,
	healthStatus HealthStatus,
	errorMessage string,
	config map[string]string,
	createdBy uuid.UUID,
	createdAt, updatedAt time.Time,
) *CloudflareTunnel {
	return &CloudflareTunnel{
		id:              id,
		name:            name,
		projectID:       projectID,
		organizationID:  organizationID,
		tunnelToken:     tunnelToken,
		containerID:     containerID,
		status:          status,
		lastHealthCheck: lastHealthCheck,
		healthStatus:    healthStatus,
		errorMessage:    errorMessage,
		config:          config,
		createdBy:       createdBy,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}
