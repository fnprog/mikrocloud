package tunnel

import (
	"context"
	"fmt"

	"github.com/mikrocloud/mikrocloud/pkg/containers/manager"
)

type CloudflaredManager struct {
	containerManager manager.ContainerManager
}

func NewCloudflaredManager(containerManager manager.ContainerManager) *CloudflaredManager {
	return &CloudflaredManager{
		containerManager: containerManager,
	}
}

type CloudflaredConfig struct {
	TunnelToken   string
	ContainerName string
	TargetURL     string
	NetworkMode   string
	Networks      []string
}

func (cm *CloudflaredManager) StartTunnel(ctx context.Context, config CloudflaredConfig) (string, error) {
	if err := cm.containerManager.PullImage(ctx, "cloudflare/cloudflared:latest"); err != nil {
		return "", fmt.Errorf("failed to pull cloudflared image: %w", err)
	}

	containerConfig := manager.ContainerConfig{
		Image:         "cloudflare/cloudflared:latest",
		Name:          config.ContainerName,
		RestartPolicy: "unless-stopped",
		NetworkMode:   config.NetworkMode,
		Networks:      config.Networks,
		Command: []string{
			"tunnel",
			"--no-autoupdate",
			"run",
			"--token",
			config.TunnelToken,
		},
		Labels: map[string]string{
			"mikrocloud.type":    "cloudflared",
			"mikrocloud.managed": "true",
		},
	}

	if config.TargetURL != "" {
		containerConfig.Environment = map[string]string{
			"TUNNEL_URL": config.TargetURL,
		}
	}

	containerID, err := cm.containerManager.Create(ctx, containerConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create cloudflared container: %w", err)
	}

	if err := cm.containerManager.Start(ctx, containerID); err != nil {
		cm.containerManager.Delete(ctx, containerID)
		return "", fmt.Errorf("failed to start cloudflared container: %w", err)
	}

	return containerID, nil
}

func (cm *CloudflaredManager) StopTunnel(ctx context.Context, containerID string) error {
	if err := cm.containerManager.Stop(ctx, containerID); err != nil {
		return fmt.Errorf("failed to stop cloudflared container: %w", err)
	}
	return nil
}

func (cm *CloudflaredManager) DeleteTunnel(ctx context.Context, containerID string) error {
	if err := cm.containerManager.Delete(ctx, containerID); err != nil {
		return fmt.Errorf("failed to delete cloudflared container: %w", err)
	}
	return nil
}

func (cm *CloudflaredManager) GetTunnelStatus(ctx context.Context, containerID string) (*manager.ContainerInfo, error) {
	info, err := cm.containerManager.Inspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect cloudflared container: %w", err)
	}
	return info, nil
}

func (cm *CloudflaredManager) RestartTunnel(ctx context.Context, containerID string) error {
	if err := cm.containerManager.Restart(ctx, containerID); err != nil {
		return fmt.Errorf("failed to restart cloudflared container: %w", err)
	}
	return nil
}

func (cm *CloudflaredManager) IsHealthy(ctx context.Context, containerID string) (bool, error) {
	info, err := cm.GetTunnelStatus(ctx, containerID)
	if err != nil {
		return false, err
	}
	return info.State == "running", nil
}
