// Package service contains the Service aggregate following DDD principles
package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service represents the core service entity (replaces Application)
type Service struct {
	id            ServiceID
	name          ServiceName
	projectID     uuid.UUID
	environmentID uuid.UUID
	gitURL        *GitURL
	buildConfig   *BuildConfig
	environment   map[string]string
	status        ServiceStatus
	createdAt     time.Time
	updatedAt     time.Time
}

// ServiceID is a value object for service identification
type ServiceID struct {
	value uuid.UUID
}

func NewServiceID() ServiceID {
	return ServiceID{value: uuid.New()}
}

func ServiceIDFromString(s string) (ServiceID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ServiceID{}, fmt.Errorf("invalid service ID: %w", err)
	}
	return ServiceID{value: id}, nil
}

func (id ServiceID) String() string {
	return id.value.String()
}

func (id ServiceID) UUID() uuid.UUID {
	return id.value
}

// ServiceName is a value object that enforces naming rules
type ServiceName struct {
	value string
}

func NewServiceName(name string) (ServiceName, error) {
	if name == "" {
		return ServiceName{}, fmt.Errorf("service name cannot be empty")
	}

	if len(name) > 50 {
		return ServiceName{}, fmt.Errorf("service name cannot exceed 50 characters")
	}

	// Additional validation rules for DNS-safe names
	// TODO: Add regex validation for DNS-compatible names

	return ServiceName{value: name}, nil
}

func (n ServiceName) String() string {
	return n.value
}

// GitURL is a value object for Git repository URLs
type GitURL struct {
	value       string
	branch      string
	contextRoot string
}

func NewGitURL(url, branch, contextRoot string) (*GitURL, error) {
	if url == "" {
		return nil, fmt.Errorf("git URL cannot be empty")
	}

	if branch == "" {
		branch = "main" // Default branch
	}

	return &GitURL{
		value:       url,
		branch:      branch,
		contextRoot: contextRoot,
	}, nil
}

func (g *GitURL) URL() string {
	return g.value
}

func (g *GitURL) Branch() string {
	return g.branch
}

func (g *GitURL) ContextRoot() string {
	return g.contextRoot
}

// BuildConfig represents build configuration for the service
type BuildConfig struct {
	buildpackType BuildpackType
	nixpacks      *NixpacksConfig
	static        *StaticConfig
	dockerfile    *DockerfileConfig
	compose       *ComposeConfig
}

type BuildpackType string

const (
	BuildpackNixpacks      BuildpackType = "nixpacks"
	BuildpackStatic        BuildpackType = "static"
	BuildpackDockerfile    BuildpackType = "dockerfile"
	BuildpackDockerCompose BuildpackType = "docker-compose"
)

type NixpacksConfig struct {
	StartCommand string            `json:"start_command,omitempty"`
	BuildCommand string            `json:"build_command,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
}

type StaticConfig struct {
	BuildCommand string `json:"build_command,omitempty"`
	OutputDir    string `json:"output_dir,omitempty"`
	NginxConfig  string `json:"nginx_config,omitempty"`
}

type DockerfileConfig struct {
	DockerfilePath string            `json:"dockerfile_path,omitempty"`
	BuildArgs      map[string]string `json:"build_args,omitempty"`
	Target         string            `json:"target,omitempty"`
}

type ComposeConfig struct {
	ComposeFile string `json:"compose_file,omitempty"`
	Service     string `json:"service,omitempty"`
}

func NewBuildConfig(buildpackType BuildpackType) *BuildConfig {
	return &BuildConfig{
		buildpackType: buildpackType,
	}
}

func (bc *BuildConfig) SetNixpacksConfig(config *NixpacksConfig) {
	bc.buildpackType = BuildpackNixpacks
	bc.nixpacks = config
}

func (bc *BuildConfig) SetStaticConfig(config *StaticConfig) {
	bc.buildpackType = BuildpackStatic
	bc.static = config
}

func (bc *BuildConfig) SetDockerfileConfig(config *DockerfileConfig) {
	bc.buildpackType = BuildpackDockerfile
	bc.dockerfile = config
}

func (bc *BuildConfig) SetComposeConfig(config *ComposeConfig) {
	bc.buildpackType = BuildpackDockerCompose
	bc.compose = config
}

func (bc *BuildConfig) BuildpackType() BuildpackType {
	return bc.buildpackType
}

func (bc *BuildConfig) NixpacksConfig() *NixpacksConfig {
	return bc.nixpacks
}

func (bc *BuildConfig) StaticConfig() *StaticConfig {
	return bc.static
}

func (bc *BuildConfig) DockerfileConfig() *DockerfileConfig {
	return bc.dockerfile
}

func (bc *BuildConfig) ComposeConfig() *ComposeConfig {
	return bc.compose
}

// ServiceStatus represents the current state of a service
type ServiceStatus int

const (
	ServiceStatusCreated ServiceStatus = iota
	ServiceStatusBuilding
	ServiceStatusDeploying
	ServiceStatusRunning
	ServiceStatusStopped
	ServiceStatusFailed
)

func (s ServiceStatus) String() string {
	switch s {
	case ServiceStatusCreated:
		return "created"
	case ServiceStatusBuilding:
		return "building"
	case ServiceStatusDeploying:
		return "deploying"
	case ServiceStatusRunning:
		return "running"
	case ServiceStatusStopped:
		return "stopped"
	case ServiceStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// NewService creates a new service with business rules enforcement
func NewService(name ServiceName, projectID, environmentID uuid.UUID, gitURL *GitURL, buildConfig *BuildConfig) *Service {
	now := time.Now()

	return &Service{
		id:            NewServiceID(),
		name:          name,
		projectID:     projectID,
		environmentID: environmentID,
		gitURL:        gitURL,
		buildConfig:   buildConfig,
		environment:   make(map[string]string),
		status:        ServiceStatusCreated,
		createdAt:     now,
		updatedAt:     now,
	}
}

// Getters
func (s *Service) ID() ServiceID {
	return s.id
}

func (s *Service) Name() ServiceName {
	return s.name
}

func (s *Service) ProjectID() uuid.UUID {
	return s.projectID
}

func (s *Service) EnvironmentID() uuid.UUID {
	return s.environmentID
}

func (s *Service) GitURL() *GitURL {
	return s.gitURL
}

func (s *Service) BuildConfig() *BuildConfig {
	return s.buildConfig
}

func (s *Service) Environment() map[string]string {
	// Return a copy to maintain encapsulation
	env := make(map[string]string)
	for k, v := range s.environment {
		env[k] = v
	}
	return env
}

func (s *Service) Status() ServiceStatus {
	return s.status
}

func (s *Service) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Service) UpdatedAt() time.Time {
	return s.updatedAt
}

// Business methods
func (s *Service) UpdateGitURL(gitURL *GitURL) error {
	if gitURL == nil {
		return fmt.Errorf("git URL cannot be nil")
	}

	s.gitURL = gitURL
	s.updatedAt = time.Now()
	return nil
}

func (s *Service) UpdateBuildConfig(buildConfig *BuildConfig) error {
	if buildConfig == nil {
		return fmt.Errorf("build config cannot be nil")
	}

	s.buildConfig = buildConfig
	s.updatedAt = time.Now()
	return nil
}

func (s *Service) SetEnvironmentVariable(key, value string) error {
	if key == "" {
		return fmt.Errorf("environment variable key cannot be empty")
	}

	s.environment[key] = value
	s.updatedAt = time.Now()
	return nil
}

func (s *Service) RemoveEnvironmentVariable(key string) {
	delete(s.environment, key)
	s.updatedAt = time.Now()
}

func (s *Service) ChangeStatus(status ServiceStatus) {
	s.status = status
	s.updatedAt = time.Now()
}

// CanDeploy checks if the service can be deployed
func (s *Service) CanDeploy() error {
	switch s.status {
	case ServiceStatusBuilding:
		return fmt.Errorf("service is currently building")
	case ServiceStatusDeploying:
		return fmt.Errorf("service is currently deploying")
	default:
		return nil
	}
}

// CanStop checks if the service can be stopped
func (s *Service) CanStop() error {
	if s.status != ServiceStatusRunning {
		return fmt.Errorf("service is not running")
	}
	return nil
}

// ReconstructService recreates a service from persistence data
func ReconstructService(
	id ServiceID,
	name ServiceName,
	projectID uuid.UUID,
	environmentID uuid.UUID,
	gitURL *GitURL,
	buildConfig *BuildConfig,
	environment map[string]string,
	status ServiceStatus,
	createdAt time.Time,
	updatedAt time.Time,
) *Service {
	return &Service{
		id:            id,
		name:          name,
		projectID:     projectID,
		environmentID: environmentID,
		gitURL:        gitURL,
		buildConfig:   buildConfig,
		environment:   environment,
		status:        status,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
