package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels"
	"github.com/mikrocloud/mikrocloud/internal/domain/tunnels/repository"
	"github.com/mikrocloud/mikrocloud/pkg/containers/tunnel"
)

type TunnelService struct {
	repo              repository.TunnelRepository
	cloudflaredMgr    *tunnel.CloudflaredManager
	healthCheckTicker *time.Ticker
	stopHealthCheck   chan struct{}
}

func NewTunnelService(repo repository.TunnelRepository, cloudflaredMgr *tunnel.CloudflaredManager) *TunnelService {
	return &TunnelService{
		repo:           repo,
		cloudflaredMgr: cloudflaredMgr,
	}
}

type CreateTunnelRequest struct {
	Name           string
	ProjectID      *uuid.UUID
	OrganizationID uuid.UUID
	TunnelToken    string
	CreatedBy      uuid.UUID
}

func (s *TunnelService) CreateTunnel(ctx context.Context, req CreateTunnelRequest) (*tunnels.CloudflareTunnel, error) {
	name, err := tunnels.NewTunnelName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("invalid tunnel name: %w", err)
	}

	exists, err := s.repo.ExistsByName(ctx, req.OrganizationID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check tunnel existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("tunnel with name %s already exists", req.Name)
	}

	token, err := tunnels.NewTunnelToken(req.TunnelToken)
	if err != nil {
		return nil, fmt.Errorf("invalid tunnel token: %w", err)
	}

	tun, err := tunnels.NewCloudflareTunnel(
		name,
		req.ProjectID,
		req.OrganizationID,
		token,
		req.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tunnel: %w", err)
	}

	if err := s.repo.Create(ctx, tun); err != nil {
		return nil, fmt.Errorf("failed to persist tunnel: %w", err)
	}

	return tun, nil
}

func (s *TunnelService) StartTunnel(ctx context.Context, tunnelID tunnels.TunnelID) error {
	tun, err := s.repo.GetByID(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tun.Status() == tunnels.TunnelStatusRunning {
		return fmt.Errorf("tunnel is already running")
	}

	tun.ChangeStatus(tunnels.TunnelStatusStarting)
	if err := s.repo.Update(ctx, tun); err != nil {
		return fmt.Errorf("failed to update tunnel status: %w", err)
	}

	containerConfig := tunnel.CloudflaredConfig{
		TunnelToken:   tun.TunnelToken().String(),
		ContainerName: tun.GetContainerName(),
		NetworkMode:   "bridge",
	}

	containerID, err := s.cloudflaredMgr.StartTunnel(ctx, containerConfig)
	if err != nil {
		tun.SetError(fmt.Sprintf("failed to start container: %v", err))
		_ = s.repo.Update(ctx, tun)
		return fmt.Errorf("failed to start cloudflared container: %w", err)
	}

	tun.SetContainerID(containerID)
	tun.ChangeStatus(tunnels.TunnelStatusRunning)
	tun.UpdateHealth(tunnels.HealthStatusHealthy)

	if err := s.repo.Update(ctx, tun); err != nil {
		return fmt.Errorf("failed to update tunnel: %w", err)
	}

	return nil
}

func (s *TunnelService) StopTunnel(ctx context.Context, tunnelID tunnels.TunnelID) error {
	tun, err := s.repo.GetByID(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tun.Status() == tunnels.TunnelStatusStopped {
		return fmt.Errorf("tunnel is already stopped")
	}

	if tun.ContainerID() == "" {
		tun.ChangeStatus(tunnels.TunnelStatusStopped)
		return s.repo.Update(ctx, tun)
	}

	tun.ChangeStatus(tunnels.TunnelStatusStopping)
	if err := s.repo.Update(ctx, tun); err != nil {
		return fmt.Errorf("failed to update tunnel status: %w", err)
	}

	if err := s.cloudflaredMgr.StopTunnel(ctx, tun.ContainerID()); err != nil {
		tun.SetError(fmt.Sprintf("failed to stop container: %v", err))
		_ = s.repo.Update(ctx, tun)
		return fmt.Errorf("failed to stop cloudflared container: %w", err)
	}

	tun.ChangeStatus(tunnels.TunnelStatusStopped)
	tun.UpdateHealth(tunnels.HealthStatusUnknown)

	if err := s.repo.Update(ctx, tun); err != nil {
		return fmt.Errorf("failed to update tunnel: %w", err)
	}

	return nil
}

func (s *TunnelService) DeleteTunnel(ctx context.Context, tunnelID tunnels.TunnelID) error {
	tun, err := s.repo.GetByID(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tun.Status() == tunnels.TunnelStatusRunning {
		if err := s.StopTunnel(ctx, tunnelID); err != nil {
			return fmt.Errorf("failed to stop tunnel before deletion: %w", err)
		}
	}

	if tun.ContainerID() != "" {
		if err := s.cloudflaredMgr.DeleteTunnel(ctx, tun.ContainerID()); err != nil {
			return fmt.Errorf("failed to delete container: %w", err)
		}
	}

	if err := s.repo.Delete(ctx, tunnelID); err != nil {
		return fmt.Errorf("failed to delete tunnel from repository: %w", err)
	}

	return nil
}

func (s *TunnelService) GetTunnel(ctx context.Context, tunnelID tunnels.TunnelID) (*tunnels.CloudflareTunnel, error) {
	return s.repo.GetByID(ctx, tunnelID)
}

func (s *TunnelService) ListTunnelsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*tunnels.CloudflareTunnel, error) {
	return s.repo.ListByOrganization(ctx, organizationID)
}

func (s *TunnelService) ListTunnelsByProject(ctx context.Context, projectID uuid.UUID) ([]*tunnels.CloudflareTunnel, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *TunnelService) CheckTunnelHealth(ctx context.Context, tunnelID tunnels.TunnelID) error {
	tun, err := s.repo.GetByID(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tun.ContainerID() == "" {
		return fmt.Errorf("tunnel has no container")
	}

	healthy, err := s.cloudflaredMgr.IsHealthy(ctx, tun.ContainerID())
	if err != nil {
		tun.UpdateHealth(tunnels.HealthStatusUnhealthy)
		tun.SetError(fmt.Sprintf("health check failed: %v", err))
		_ = s.repo.Update(ctx, tun)
		return err
	}

	if healthy {
		tun.UpdateHealth(tunnels.HealthStatusHealthy)
		if tun.Status() != tunnels.TunnelStatusRunning {
			tun.ChangeStatus(tunnels.TunnelStatusRunning)
		}
	} else {
		tun.UpdateHealth(tunnels.HealthStatusUnhealthy)
		if tun.Status() == tunnels.TunnelStatusRunning {
			tun.ChangeStatus(tunnels.TunnelStatusError)
		}
	}

	return s.repo.Update(ctx, tun)
}

func (s *TunnelService) StartHealthCheckMonitor(ctx context.Context, interval time.Duration) {
	s.healthCheckTicker = time.NewTicker(interval)
	s.stopHealthCheck = make(chan struct{})

	go func() {
		for {
			select {
			case <-s.healthCheckTicker.C:
				s.performHealthChecks(ctx)
			case <-s.stopHealthCheck:
				return
			}
		}
	}()
}

func (s *TunnelService) StopHealthCheckMonitor() {
	if s.healthCheckTicker != nil {
		s.healthCheckTicker.Stop()
	}
	if s.stopHealthCheck != nil {
		close(s.stopHealthCheck)
	}
}

func (s *TunnelService) performHealthChecks(ctx context.Context) {
	tunnels, err := s.repo.ListByStatus(ctx, tunnels.TunnelStatusRunning)
	if err != nil {
		return
	}

	for _, tun := range tunnels {
		_ = s.CheckTunnelHealth(ctx, tun.ID())
	}
}

func (s *TunnelService) RestartTunnel(ctx context.Context, tunnelID tunnels.TunnelID) error {
	tun, err := s.repo.GetByID(ctx, tunnelID)
	if err != nil {
		return fmt.Errorf("tunnel not found: %w", err)
	}

	if tun.ContainerID() == "" {
		return fmt.Errorf("tunnel has no container to restart")
	}

	if err := s.cloudflaredMgr.RestartTunnel(ctx, tun.ContainerID()); err != nil {
		tun.SetError(fmt.Sprintf("failed to restart container: %v", err))
		s.repo.Update(ctx, tun)
		return fmt.Errorf("failed to restart cloudflared container: %w", err)
	}

	tun.ChangeStatus(tunnels.TunnelStatusRunning)
	tun.UpdateHealth(tunnels.HealthStatusHealthy)

	if err := s.repo.Update(ctx, tun); err != nil {
		return fmt.Errorf("failed to update tunnel: %w", err)
	}

	return nil
}
