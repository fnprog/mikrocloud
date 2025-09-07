# Container Management System

This document describes the container abstraction layer and build system for mikrocloud.

## Overview

The system provides:

1. **Container Manager Abstraction** - Unified interface for Docker and Podman
2. **Build System** - Isolated container builds with multiple buildpack support
3. **API Layer** - REST endpoints for container and build operations

## Configuration

Configure the container runtime in your configuration file:

```toml
[docker]
runtime = "docker"  # or "podman"
socket_path = "/var/run/docker.sock"  # or "/run/user/1000/podman/podman.sock" for rootless
rootless = false
build_dir = "/tmp/mikrocloud/builds"
```

## Container Operations

### List Containers

```bash
curl -X GET http://localhost:3000/api/v1/containers
```

### Create Container

```bash
curl -X POST http://localhost:3000/api/v1/containers \
  -H "Content-Type: application/json" \
  -d '{
    "image": "nginx:alpine",
    "name": "my-nginx",
    "ports": {"8080": "80"},
    "environment": {"ENV": "production"},
    "restart_policy": "unless-stopped"
  }'
```

### Start/Stop/Restart Container

```bash
# Start
curl -X POST http://localhost:3000/api/v1/containers/{id}/start

# Stop
curl -X POST http://localhost:3000/api/v1/containers/{id}/stop

# Restart
curl -X POST http://localhost:3000/api/v1/containers/{id}/restart
```

### Delete Container

```bash
curl -X DELETE http://localhost:3000/api/v1/containers/{id}
```

## Build Operations

### Nixpacks Build

```bash
curl -X POST http://localhost:3000/api/v1/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "id": "build-001",
    "git_repo": "https://github.com/user/app.git",
    "git_branch": "main",
    "buildpack_type": "nixpacks",
    "image_tag": "myapp:latest",
    "nixpacks_config": {
      "start_command": "npm start",
      "variables": {"NODE_ENV": "production"}
    }
  }'
```

### Static Site Build

```bash
curl -X POST http://localhost:3000/api/v1/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "id": "build-002",
    "git_repo": "https://github.com/user/static-site.git",
    "buildpack_type": "static",
    "image_tag": "mysite:latest",
    "static_config": {
      "build_command": "npm run build",
      "output_dir": "dist"
    }
  }'
```

### Dockerfile Build

```bash
curl -X POST http://localhost:3000/api/v1/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "id": "build-003",
    "git_repo": "https://github.com/user/docker-app.git",
    "buildpack_type": "dockerfile",
    "image_tag": "dockerapp:latest",
    "dockerfile_config": {
      "dockerfile_path": "Dockerfile",
      "build_args": {"VERSION": "1.0.0"}
    }
  }'
```

### Docker Compose Build

```bash
curl -X POST http://localhost:3000/api/v1/images/build \
  -H "Content-Type: application/json" \
  -d '{
    "id": "build-004",
    "git_repo": "https://github.com/user/compose-app.git",
    "buildpack_type": "docker-compose",
    "image_tag": "composeapp:latest",
    "compose_config": {
      "compose_file": "docker-compose.yml",
      "service": "web"
    }
  }'
```

## Runtime Information

Get information about the configured container runtime:

```bash
curl -X GET http://localhost:3000/api/v1/runtime/info
```

Response:

```json
{
  "runtime": "docker",
  "socket_path": "/var/run/docker.sock",
  "rootless": false,
  "build_dir": "/tmp/mikrocloud/builds"
}
```

## Architecture

### Container Manager Interface

- Abstracts Docker and Podman operations
- Provides unified API for container lifecycle management
- Supports both rootful and rootless configurations

### Build System

- Isolated build environments using containers
- Clones Git repositories into temporary workspaces
- Supports multiple buildpack types
- Connects to host container socket for building

### Buildpack Types

1. **Nixpacks** - Automatic language detection and building
2. **Static** - Static site builds with Nginx serving
3. **Dockerfile** - Traditional Docker builds
4. **Docker Compose** - Multi-service applications

### Security

- Builds run in isolated containers
- No direct host access except for container socket
- Temporary build directories are cleaned up after use
- Support for rootless container runtimes

## Integration

The container service integrates with the existing mikrocloud architecture:

```go
// Service initialization
containerService, err := services.NewContainerService(cfg)
if err != nil {
    log.Fatal(err)
}

// Use in application deployment
buildResult, err := containerService.BuildImage(ctx, buildRequest)
if err != nil {
    return err
}

// Deploy the built image
containerConfig := manager.ContainerConfig{
    Image: buildResult.ImageTag,
    Name:  "app-instance",
    Ports: map[string]string{"80": "8080"},
}

containerID, err := containerService.CreateContainer(ctx, containerConfig)
```

This provides a complete container management and build system that can work with both Docker and Podman, supporting the various deployment scenarios needed by mikrocloud.
