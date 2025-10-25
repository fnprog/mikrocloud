package service

import (
	"context"
	"fmt"

	"github.com/mikrocloud/mikrocloud/internal/domain/applications"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments/logs"
	"github.com/mikrocloud/mikrocloud/internal/domain/deployments/repository"
	"github.com/mikrocloud/mikrocloud/internal/domain/users"
	"github.com/mikrocloud/mikrocloud/pkg/containers"
	"github.com/mikrocloud/mikrocloud/pkg/containers/build"
	"github.com/mikrocloud/mikrocloud/pkg/containers/manager"
	services "github.com/mikrocloud/mikrocloud/pkg/containers/service"
)

type BuildService interface {
	BuildImage(ctx context.Context, request build.BuildRequest) (*build.BuildResult, error)
}

type ApplicationService interface {
	GetApplication(ctx context.Context, id applications.ApplicationID) (*applications.Application, error)
}

type DeploymentService struct {
	repo             repository.DeploymentRepository
	containerService *services.ContainerService
}

func NewDeploymentService(repo repository.DeploymentRepository, containerService *services.ContainerService) *DeploymentService {
	return &DeploymentService{
		repo:             repo,
		containerService: containerService,
	}
}

type CreateDeploymentCommand struct {
	ApplicationID    applications.ApplicationID
	IsProduction     bool
	TriggeredBy      *users.UserID
	TriggerType      deployments.TriggerType
	ImageTag         string
	GitCommitHash    string
	GitCommitMessage string
	GitBranch        string
	GitAuthorName    string
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, cmd CreateDeploymentCommand) (*deployments.Deployment, error) {
	// Get the next deployment number for this application
	deploymentNumber, err := s.getNextDeploymentNumber(ctx, cmd.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next deployment number: %w", err)
	}

	deployment := deployments.NewDeployment(
		cmd.ApplicationID,
		deploymentNumber,
		cmd.IsProduction,
		cmd.TriggeredBy,
		cmd.TriggerType,
		cmd.ImageTag,
	)

	if cmd.GitCommitHash != "" || cmd.GitCommitMessage != "" || cmd.GitBranch != "" || cmd.GitAuthorName != "" {
		deployment.SetGitInfo(cmd.GitCommitHash, cmd.GitCommitMessage, cmd.GitBranch, cmd.GitAuthorName)
	}

	if err := s.repo.Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return deployment, nil
}

func (s *DeploymentService) CreateAndExecuteDeployment(ctx context.Context, cmd CreateDeploymentCommand, appService ApplicationService) (*deployments.Deployment, error) {
	// Create the deployment record
	deployment, err := s.CreateDeployment(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Start the build process in the background with inherited context
	go func() {
		bgCtx := context.WithoutCancel(ctx)
		if err := s.executeBuildAndDeploy(bgCtx, deployment.ID(), appService); err != nil {
			s.FailDeployment(bgCtx, deployment.ID(), err.Error())
		}
	}()

	return deployment, nil
}

func (s *DeploymentService) executeBuildAndDeploy(ctx context.Context, deploymentID deployments.DeploymentID, appService ApplicationService) error {
	// Start build phase
	if err := s.StartBuild(ctx, deploymentID); err != nil {
		return fmt.Errorf("failed to start build: %w", err)
	}

	// Get deployment and application details
	deployment, err := s.repo.GetByID(ctx, deploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	app, err := appService.GetApplication(ctx, deployment.ApplicationID())
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	// Convert application config to build request
	buildRequest, err := s.createBuildRequest(deployment, app)
	if err != nil {
		return fmt.Errorf("failed to create build request: %w", err)
	}

	// Add real-time log callback
	buildRequest.LogCallback = func(log string) {
		s.AppendBuildLogs(ctx, deploymentID, log)
	}

	// Execute the build
	buildResult, err := s.containerService.BuildImage(ctx, *buildRequest)
	if err != nil {
		s.AppendBuildLogs(ctx, deploymentID, fmt.Sprintf("Build failed: %v", err))
		return fmt.Errorf("build failed: %w", err)
	}

	// Update deployment with any remaining build logs
	if buildResult.BuildLogs != "" {
		s.AppendBuildLogs(ctx, deploymentID, buildResult.BuildLogs)
	}

	if !buildResult.Success {
		s.AppendBuildLogs(ctx, deploymentID, fmt.Sprintf("Build failed: %s", buildResult.Error))
		return fmt.Errorf("build failed: %s", buildResult.Error)
	}

	// Complete build phase
	if err := s.CompleteBuild(ctx, deploymentID); err != nil {
		return fmt.Errorf("failed to complete build: %w", err)
	}

	// Start deploy phase
	if err := s.StartDeploy(ctx, deploymentID); err != nil {
		return fmt.Errorf("failed to start deploy: %w", err)
	}

	// Set image information
	if buildResult.ImageTag != "" {
		if err := s.SetImageDigest(ctx, deploymentID, buildResult.ImageTag); err != nil {
			return fmt.Errorf("failed to set image digest: %w", err)
		}
	}

	// Deploy the container
	if err := s.deployContainer(ctx, deploymentID, deployment, app, buildResult.ImageTag); err != nil {
		s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Container deployment failed: %v", err))
		return fmt.Errorf("container deployment failed: %w", err)
	}

	// Complete deploy phase
	if err := s.CompleteDeploy(ctx, deploymentID); err != nil {
		return fmt.Errorf("failed to complete deploy: %w", err)
	}

	return nil
}

func (s *DeploymentService) createBuildRequest(deployment *deployments.Deployment, app *applications.Application) (*build.BuildRequest, error) {
	// Extract deployment source information
	deploymentSource := app.DeploymentSource()

	var gitRepo, gitBranch, contextRoot string
	var environment map[string]string

	// Handle different deployment source types
	switch deploymentSource.Type {
	case applications.DeploymentSourceTypeGit:
		if deploymentSource.GitRepo != nil {
			gitRepo = deploymentSource.GitRepo.URL
			gitBranch = deploymentSource.GitRepo.Branch
			contextRoot = deploymentSource.GitRepo.Path
		}
	case applications.DeploymentSourceTypeRegistry:
		// For registry deployments, we don't need to build
		return nil, fmt.Errorf("registry deployments don't require building")
	case applications.DeploymentSourceTypeUpload:
		// For upload deployments, the build context is different
		// TODO: Implement upload-based builds
		return nil, fmt.Errorf("upload deployments not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported deployment source type: %s", deploymentSource.Type)
	}

	// Get environment variables
	environment = app.EnvVars()
	// Convert buildpack config to build request format
	buildpackConfig := app.Buildpack()
	if buildpackConfig == nil {
		return nil, fmt.Errorf("application has no buildpack configuration")
	}

	var buildpackType build.BuildpackType

	switch buildpackConfig.BuildpackType() {
	case applications.BuildpackTypeNixpacks:
		buildpackType = build.Nixpacks
	case applications.BuildpackTypeStatic:
		buildpackType = build.Static
	case applications.BuildpackTypeDockerfile:
		buildpackType = build.DockerfileType
	case applications.BuildpackTypeDockerCompose:
		buildpackType = build.DockerCompose
	default:
		return nil, fmt.Errorf("unsupported buildpack type: %s", buildpackConfig.BuildpackType())
	}

	// Create build request
	buildRequest := &build.BuildRequest{
		ID:            deployment.ID().String(),
		ImageTag:      deployment.ImageTag(),
		GitRepo:       gitRepo,
		GitBranch:     gitBranch,
		ContextRoot:   contextRoot,
		BuildpackType: buildpackType,
		Environment:   environment,
	}

	// Set buildpack-specific configurations based on the config field
	// Note: We now use typed configs from the BuildConfig

	switch buildpackConfig.BuildpackType() {
	case applications.BuildpackTypeNixpacks:
		if nixConfig := buildpackConfig.NixpacksConfig(); nixConfig != nil {
			nixpacksConfig := &build.NixpacksConfig{
				StartCommand: nixConfig.StartCommand,
				BuildCommand: nixConfig.BuildCommand,
				Variables:    nixConfig.Variables,
			}
			buildRequest.NixpacksConfig = nixpacksConfig
		}
	case applications.BuildpackTypeStatic:
		if staticConfig := buildpackConfig.StaticConfig(); staticConfig != nil {
			// TODO: Complete here
			staticBuildConfig := &build.StaticConfig{
				OutputDir:   staticConfig.OutputDir,
				NginxConfig: staticConfig.NginxConfig,
			}
			buildRequest.StaticConfig = staticBuildConfig
		}
	case applications.BuildpackTypeDockerfile:
		if containerConfig := buildpackConfig.DockerfileConfig(); containerConfig != nil {
			containerfileConfig := &build.ContainerfileConfig{
				ContainerfilePath: containerConfig.DockerfilePath,
				BuildArgs:         containerConfig.BuildArgs,
				Target:            containerConfig.Target,
			}
			buildRequest.ContainerfileConfig = containerfileConfig
		}
	case applications.BuildpackTypeDockerCompose:
		if composeConfig := buildpackConfig.ComposeConfig(); composeConfig != nil {
			composeRequestConfig := &build.ComposeConfig{
				ComposeFile: composeConfig.ComposeFile,
				Service:     composeConfig.Service,
			}
			buildRequest.ComposeConfig = composeRequestConfig
		}
	}

	return buildRequest, nil
}

func (s *DeploymentService) getNextDeploymentNumber(ctx context.Context, applicationID applications.ApplicationID) (int, error) {
	latest, err := s.repo.GetLatestByApplication(ctx, applicationID)
	if err != nil {
		// If no deployments exist, start with 1
		return 1, nil
	}
	return latest.DeploymentNumber() + 1, nil
}

func (s *DeploymentService) GetDeployment(ctx context.Context, id deployments.DeploymentID) (*deployments.Deployment, error) {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deployment, nil
}

func (s *DeploymentService) GetDeploymentWithMetadata(ctx context.Context, id deployments.DeploymentID) (*repository.DeploymentWithMetadata, error) {
	deploymentWithMeta, err := s.repo.GetByIDWithMetadata(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deploymentWithMeta, nil
}

func (s *DeploymentService) ListDeploymentsByApplicationWithMetadata(ctx context.Context, applicationID applications.ApplicationID) ([]*repository.DeploymentWithMetadata, error) {
	deploymentsWithMeta, err := s.repo.ListByApplicationWithMetadata(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return deploymentsWithMeta, nil
}

func (s *DeploymentService) StartBuild(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.StartBuild()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CompleteBuild(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.CompleteBuild()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) StartDeploy(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.StartDeploy()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CompleteDeploy(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.CompleteDeploy()

	if err := s.parseFinalLogs(deployment); err != nil {
		return fmt.Errorf("failed to parse logs: %w", err)
	}

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) parseFinalLogs(deployment *deployments.Deployment) error {
	parser := logs.NewLogParser(deployment.StartedAt())

	buildLogs := deployment.BuildLogs()
	if buildLogs != "" {
		parsedBuildLogs, err := parser.ParseAndSerialize(buildLogs)
		if err != nil {
			return fmt.Errorf("failed to parse build logs: %w", err)
		}
		if parsedBuildLogs != "" {
			deployment.SetBuildLogs(parsedBuildLogs)
		}
	}

	deployLogs := deployment.DeployLogs()
	if deployLogs != "" {
		parsedDeployLogs, err := parser.ParseAndSerialize(deployLogs)
		if err != nil {
			return fmt.Errorf("failed to parse deploy logs: %w", err)
		}
		if parsedDeployLogs != "" {
			deployment.SetDeployLogs(parsedDeployLogs)
		}
	}

	return nil
}

func (s *DeploymentService) FailDeployment(ctx context.Context, id deployments.DeploymentID, errorMessage string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Fail(errorMessage)

	if err := s.parseFinalLogs(deployment); err != nil {
		return fmt.Errorf("failed to parse logs: %w", err)
	}

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) CancelDeployment(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Cancel()

	if err := s.parseFinalLogs(deployment); err != nil {
		return fmt.Errorf("failed to parse logs: %w", err)
	}

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) StopDeployment(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.Stop()

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) SetContainerID(ctx context.Context, id deployments.DeploymentID, containerID string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.SetContainerID(containerID)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) SetImageDigest(ctx context.Context, id deployments.DeploymentID, digest string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.SetImageDigest(digest)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) AppendBuildLogs(ctx context.Context, id deployments.DeploymentID, logs string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.AppendBuildLogs(logs)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) AppendDeployLogs(ctx context.Context, id deployments.DeploymentID, logs string) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	deployment.AppendDeployLogs(logs)

	if err := s.repo.Update(ctx, deployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) ListDeployments(ctx context.Context) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) ListDeploymentsByApplication(ctx context.Context, applicationID applications.ApplicationID) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.ListByApplication(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by application: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) GetLatestDeploymentByApplication(ctx context.Context, applicationID applications.ApplicationID) (*deployments.Deployment, error) {
	deployment, err := s.repo.GetLatestByApplication(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest deployment: %w", err)
	}
	return deployment, nil
}

func (s *DeploymentService) ListDeploymentsByStatus(ctx context.Context, status deployments.DeploymentStatus) ([]*deployments.Deployment, error) {
	deployments, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments by status: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentService) deployContainer(ctx context.Context, deploymentID deployments.DeploymentID, deployment *deployments.Deployment, app *applications.Application, imageTag string) error {
	s.AppendDeployLogs(ctx, deploymentID, "Starting container deployment...")

	containerName := containers.SanitizeDockerName(fmt.Sprintf("%s-%d", app.Name().String(), deployment.DeploymentNumber()))

	ports := make(map[string]string)
	if len(app.PortMappings()) > 0 {
		for _, mapping := range app.PortMappings() {
			ports[fmt.Sprintf("%d", mapping.ContainerPort)] = fmt.Sprintf("%d", mapping.HostPort)
		}
		s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Configured %d port mapping(s)", len(app.PortMappings())))
	} else {
		s.AppendDeployLogs(ctx, deploymentID, "No port mappings configured - container will not expose ports to host")
	}

	labels := make(map[string]string)
	domain := app.Domain()
	if domain == "" {
		domain = app.GeneratedDomain()
	}

	if domain != "" {
		labels["traefik.enable"] = "true"
		labels["traefik.http.routers."+containerName+".rule"] = fmt.Sprintf("Host(`%s`)", domain)
		labels["traefik.http.routers."+containerName+".entrypoints"] = "web"

		if len(app.ExposedPorts()) > 0 {
			labels["traefik.http.services."+containerName+".loadbalancer.server.port"] = fmt.Sprintf("%d", app.ExposedPorts()[0])
			s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Configured Traefik routing: %s -> port %d", domain, app.ExposedPorts()[0]))
		} else {
			labels["traefik.http.services."+containerName+".loadbalancer.server.port"] = "8080"
			s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Configured Traefik routing: %s -> port 8080 (default)", domain))
		}
	}

	containerConfig := manager.ContainerConfig{
		Image:         imageTag,
		Name:          containerName,
		Ports:         ports,
		Environment:   app.EnvVars(),
		Networks:      []string{},
		RestartPolicy: "unless-stopped",
		AutoRemove:    false,
		Labels:        labels,
	}

	s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Creating container: %s from image: %s", containerName, imageTag))

	containerID, err := s.containerService.CreateContainer(ctx, containerConfig)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Container created with ID: %s", containerID))

	if err := s.SetContainerID(ctx, deploymentID, containerID); err != nil {
		return fmt.Errorf("failed to update deployment with container ID: %w", err)
	}

	s.AppendDeployLogs(ctx, deploymentID, "Starting container...")

	if err := s.containerService.StartContainer(ctx, containerID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	containerInfo, err := s.containerService.InspectContainer(ctx, containerID)
	if err != nil {
		s.AppendDeployLogs(ctx, deploymentID, "Warning: Failed to inspect container, but it may be running")
	} else {
		s.AppendDeployLogs(ctx, deploymentID, fmt.Sprintf("Container started successfully. State: %s, Status: %s", containerInfo.State, containerInfo.Status))
	}

	s.AppendDeployLogs(ctx, deploymentID, "Container deployment completed successfully")
	return nil
}

func (s *DeploymentService) DeleteDeployment(ctx context.Context, id deployments.DeploymentID) error {
	deployment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	if deployment.ContainerID() != "" {
		if err := s.containerService.StopContainer(ctx, deployment.ContainerID()); err != nil {
		}
		if err := s.containerService.DeleteContainer(ctx, deployment.ContainerID()); err != nil {
		}
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	return nil
}

func (s *DeploymentService) RecreateContainer(ctx context.Context, applicationID applications.ApplicationID, getApp func(context.Context, applications.ApplicationID) (*applications.Application, error)) error {
	latestDeployment, err := s.repo.GetLatestByApplication(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("no deployment found for application: %w", err)
	}

	if latestDeployment.Status() != deployments.DeploymentStatusRunning {
		return fmt.Errorf("latest deployment is not running, cannot recreate container")
	}

	if latestDeployment.ContainerID() == "" {
		return fmt.Errorf("deployment has no container ID")
	}

	app, err := getApp(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	oldContainerID := latestDeployment.ContainerID()
	if err := s.containerService.StopContainer(ctx, oldContainerID); err != nil {
		return fmt.Errorf("failed to stop old container: %w", err)
	}

	if err := s.containerService.DeleteContainer(ctx, oldContainerID); err != nil {
		return fmt.Errorf("failed to delete old container: %w", err)
	}

	imageTag := latestDeployment.ImageTag()
	if imageTag == "" {
		return fmt.Errorf("deployment has no image tag")
	}

	if err := s.deployContainer(ctx, latestDeployment.ID(), latestDeployment, app, imageTag); err != nil {
		return fmt.Errorf("failed to deploy new container: %w", err)
	}

	return nil
}
