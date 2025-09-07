# Container-in-Container Build System

This document explains how MikroCloud's build system works using a container-in-container architecture, similar to platforms like Coolify.

## Architecture Overview

### The Problem

Traditional build systems often clone source code directly to the host filesystem, which can lead to:

- Security concerns (source code touching host)
- Build environment inconsistencies
- Cleanup issues
- Isolation problems

### The Solution: Container-in-Container

MikroCloud uses a container-in-container approach where:

1. **Build Helper Containers**: Each build runs in its own isolated helper container
2. **Docker Socket Mounting**: Helper containers access the host Docker daemon via socket mounting
3. **Complete Isolation**: Source code never touches the host filesystem
4. **Consistent Environment**: Every build uses the same, reproducible environment

## How It Works

```
Host Machine
├── Docker Daemon (dockerd)
│   └── Docker Socket (/var/run/docker.sock)
│
├── Build Helper Container
│   ├── /workspace/ (isolated workspace)
│   │   ├── source/ (cloned repository)
│   │   ├── Dockerfile (generated or existing)
│   │   └── build artifacts
│   ├── git (for cloning)
│   ├── docker CLI (connects to host daemon)
│   └── build tools (node, nixpacks, etc.)
│
└── Application Containers (created by helper)
    ├── my-app:latest
    ├── my-static-site:latest
    └── my-python-app:latest
```

## Build Process Flow

### 1. Build Request

```go
buildRequest := build.BuildRequest{
    ID:            "unique-build-id",
    GitRepo:       "https://github.com/user/repo.git",
    GitBranch:     "main",
    ContextRoot:   "frontend/", // optional subdirectory
    BuildpackType: build.Static,
    ImageTag:      "my-app:latest",
    Environment: map[string]string{
        "NODE_ENV": "production",
    },
}
```

### 2. Helper Container Creation

The build service creates a helper container with:

- **Base Image**: Depends on buildpack (alpine/git, nixpacks, node, etc.)
- **Socket Mount**: `/var/run/docker.sock:/var/run/docker.sock`
- **Auto-Remove**: Container automatically deleted after build
- **Isolated Workspace**: `/workspace` for all build operations

### 3. Build Script Generation

A shell script is generated that will run inside the helper container:

```bash
set -e
echo 'Starting build process...'

# Clone the repository
echo 'Cloning repository...'
git clone -b main https://github.com/user/repo.git /workspace/source
cd /workspace/source/frontend

# Execute build commands
echo 'Building static site...'
apk add --no-cache docker-cli
printf "FROM node:18-alpine as builder..." > Dockerfile
docker build -t my-app:latest .
```

### 4. Build Execution

1. Helper container starts and executes the build script
2. Repository is cloned inside the container
3. Build tools are installed as needed
4. Docker commands run inside the helper but connect to host daemon
5. Final image is built and tagged on the host
6. Helper container auto-removes itself

## Buildpack Types

### 1. Nixpacks

```go
buildRequest := build.BuildRequest{
    BuildpackType: build.Nixpacks,
    NixpacksConfig: &build.NixpacksConfig{
        StartCommand: "python app.py",
    },
}
```

- Uses `nixpacks/nixpacks:latest` helper image
- Automatically detects language and runtime
- Generates optimized Dockerfile

### 2. Static Sites

```go
buildRequest := build.BuildRequest{
    BuildpackType: build.Static,
    StaticConfig: &build.StaticConfig{
        BuildCommand: "npm run build",
        OutputDir:    "dist",
    },
}
```

- Uses `alpine/git:latest` + Node.js tools
- Generates multi-stage Dockerfile
- Final image serves static files with nginx

### 3. Dockerfile

```go
buildRequest := build.BuildRequest{
    BuildpackType: build.DockerfileType,
    DockerfileConfig: &build.DockerfileConfig{
        DockerfilePath: "Dockerfile",
        BuildArgs: map[string]string{
            "GO_VERSION": "1.21",
        },
        Target: "production",
    },
}
```

- Uses existing Dockerfile in repository
- Supports build args and multi-stage targets
- Uses `alpine/git:latest` + Docker CLI

### 4. Docker Compose

```go
buildRequest := build.BuildRequest{
    BuildpackType: build.DockerCompose,
    ComposeConfig: &build.ComposeConfig{
        ComposeFile: "docker-compose.yml",
        Service:     "web",
    },
}
```

- Uses `alpine/git:latest` + docker-compose
- Can build specific services or all services
- Supports complex multi-service builds

## Security Benefits

### Complete Isolation

- Source code never touches host filesystem
- Build environment is completely sandboxed
- No risk of host contamination

### Controlled Access

- Helper containers only have access to Docker socket
- No privileged access to host system
- Build artifacts are properly contained

### Automatic Cleanup

- Helper containers auto-remove after build
- No build artifacts left on host
- Consistent clean state after each build

## Performance Benefits

### Parallel Builds

- Multiple builds can run simultaneously
- Each build is completely isolated
- No resource conflicts between builds

### Consistent Environment

- Same build environment every time
- No "works on my machine" issues
- Reproducible builds across different hosts

### Efficient Resource Usage

- Containers only exist during build process
- Resources automatically freed after build
- No long-running build environments

## Docker vs Podman Support

The system supports both Docker and Podman:

### Docker

```go
containerManager, _ := manager.NewContainerManager(manager.Docker)
buildService := build.NewBuildService(containerManager, "/var/run/docker.sock")
```

### Podman

```go
containerManager, _ := manager.NewContainerManager(manager.Podman)
buildService := build.NewBuildService(containerManager, "/run/user/1000/podman/podman.sock")
```

Both use the same interface and provide identical functionality, allowing users to choose their preferred container runtime.

## Usage Example

```go
// Initialize the build system
containerManager, _ := manager.NewContainerManager(manager.Docker)
buildService := build.NewBuildService(containerManager, "/var/run/docker.sock")

// Create a build request
request := build.BuildRequest{
    ID:            "my-build-001",
    GitRepo:       "https://github.com/example/my-app.git",
    GitBranch:     "main",
    BuildpackType: build.Static,
    ImageTag:      "my-app:latest",
    StaticConfig: &build.StaticConfig{
        BuildCommand: "npm run build",
        OutputDir:    "dist",
    },
}

// Execute the build
result, err := buildService.BuildImage(context.Background(), request)
if err != nil {
    log.Fatal(err)
}

if result.Success {
    fmt.Printf("Build successful! Image: %s\n", result.ImageTag)
    // Deploy the built image
} else {
    fmt.Printf("Build failed: %s\n", result.Error)
}
```

This architecture provides a secure, scalable, and maintainable way to build container images while keeping the host system clean and isolated.
