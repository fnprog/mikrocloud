package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mikrocloud/mikrocloud/pkg/containers/build"
	"github.com/mikrocloud/mikrocloud/pkg/containers/manager"
)

func main() {
	ctx := context.Background()

	// Create container manager (Docker or Podman)
	containerManager, err := manager.NewContainerManager(manager.Docker)
	if err != nil {
		log.Fatalf("Failed to create container manager: %v", err)
	}

	// Create build service with Docker socket path
	buildService := build.NewBuildService(containerManager, "/var/run/docker.sock")

	// Example 1: Build a Node.js static site
	staticBuildRequest := build.BuildRequest{
		ID:            "static-build-001",
		GitRepo:       "https://github.com/Rhymond/product-compare-react.git",
		GitBranch:     "main",
		ContextRoot:   "",
		BuildpackType: build.Static,
		ImageTag:      "my-react-app:latest",
		Environment: map[string]string{
			"NODE_ENV": "production",
		},
		StaticConfig: &build.StaticConfig{
			BuildCommand: "npm run build",
			OutputDir:    "dist",
		},
	}

	fmt.Println("Building static site...")
	result, err := buildService.BuildImage(ctx, staticBuildRequest)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ Static build successful! Image: %s\n", result.ImageTag)
	} else {
		fmt.Printf("❌ Static build failed: %s\n", result.Error)
	}

	// Example 2: Build with Dockerfile
	dockerfileBuildRequest := build.BuildRequest{
		ID:            "dockerfile-build-001",
		GitRepo:       "https://github.com/example/my-go-app.git",
		GitBranch:     "main",
		ContextRoot:   "",
		BuildpackType: build.DockerfileType,
		ImageTag:      "my-go-app:latest",
		DockerfileConfig: &build.DockerfileConfig{
			DockerfilePath: "Dockerfile",
			BuildArgs: map[string]string{
				"GO_VERSION": "1.21",
			},
		},
	}

	fmt.Println("Building with Dockerfile...")
	result, err = buildService.BuildImage(ctx, dockerfileBuildRequest)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ Dockerfile build successful! Image: %s\n", result.ImageTag)
	} else {
		fmt.Printf("❌ Dockerfile build failed: %s\n", result.Error)
	}

	// Example 3: Build with Nixpacks
	nixpacksBuildRequest := build.BuildRequest{
		ID:            "nixpacks-build-001",
		GitRepo:       "https://github.com/example/my-python-app.git",
		GitBranch:     "main",
		ContextRoot:   "",
		BuildpackType: build.Nixpacks,
		ImageTag:      "my-python-app:latest",
		NixpacksConfig: &build.NixpacksConfig{
			StartCommand: "python app.py",
		},
	}

	fmt.Println("Building with Nixpacks...")
	result, err = buildService.BuildImage(ctx, nixpacksBuildRequest)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ Nixpacks build successful! Image: %s\n", result.ImageTag)
	} else {
		fmt.Printf("❌ Nixpacks build failed: %s\n", result.Error)
	}

	// Example 4: Build with Docker Compose
	composeBuildRequest := build.BuildRequest{
		ID:            "compose-build-001",
		GitRepo:       "https://github.com/example/my-microservices.git",
		GitBranch:     "main",
		ContextRoot:   "",
		BuildpackType: build.DockerCompose,
		ImageTag:      "my-microservices:latest",
		ComposeConfig: &build.ComposeConfig{
			ComposeFile: "docker-compose.yml",
			Service:     "web", // Build only the 'web' service
		},
	}

	fmt.Println("Building with Docker Compose...")
	result, err = buildService.BuildImage(ctx, composeBuildRequest)
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	if result.Success {
		fmt.Printf("✅ Compose build successful! Image: %s\n", result.ImageTag)
	} else {
		fmt.Printf("❌ Compose build failed: %s\n", result.Error)
	}
}
