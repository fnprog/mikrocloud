package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels"
)

type TunnelRepository interface {
	Create(ctx context.Context, tunnel *tunnels.CloudflareTunnel) error
	GetByID(ctx context.Context, id tunnels.TunnelID) (*tunnels.CloudflareTunnel, error)
	GetByName(ctx context.Context, organizationID uuid.UUID, name tunnels.TunnelName) (*tunnels.CloudflareTunnel, error)
	GetByContainerID(ctx context.Context, containerID string) (*tunnels.CloudflareTunnel, error)
	ListByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*tunnels.CloudflareTunnel, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*tunnels.CloudflareTunnel, error)
	ListByStatus(ctx context.Context, status tunnels.TunnelStatus) ([]*tunnels.CloudflareTunnel, error)
	ListAll(ctx context.Context) ([]*tunnels.CloudflareTunnel, error)
	Update(ctx context.Context, tunnel *tunnels.CloudflareTunnel) error
	Delete(ctx context.Context, id tunnels.TunnelID) error
	Exists(ctx context.Context, id tunnels.TunnelID) (bool, error)
	ExistsByName(ctx context.Context, organizationID uuid.UUID, name tunnels.TunnelName) (bool, error)
}
