package applications

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DeploymentSource struct {
	Type     DeploymentSourceType `json:"type"`
	GitRepo  *GitRepoSource       `json:"git_repo,omitempty"`
	Registry *RegistrySource      `json:"registry,omitempty"`
	Upload   *UploadSource        `json:"upload,omitempty"`
}

type DeploymentSourceType string

const (
	DeploymentSourceTypeGit      DeploymentSourceType = "git"
	DeploymentSourceTypeRegistry DeploymentSourceType = "registry"
	DeploymentSourceTypeUpload   DeploymentSourceType = "upload"
)

type GitRepoSource struct {
	URL    string `json:"url"`
	Branch string `json:"branch"`
	Path   string `json:"path,omitempty"`
	Token  string `json:"token,omitempty"` // For private repos
}

type RegistrySource struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
}

type UploadSource struct {
	Filename string `json:"filename"`
	FilePath string `json:"file_path"`
}

type BuildpackConfig struct {
	Type   BuildpackType `json:"type"`
	Config any           `json:"config,omitempty"`
}

type Application struct {
	id               ApplicationID
	name             ApplicationName
	description      string
	projectID        uuid.UUID
	environmentID    uuid.UUID
	deploymentSource DeploymentSource
	domain           string
	buildpack        BuildpackConfig
	envVars          map[string]string
	autoDeploy       bool
	status           ApplicationStatus
	createdAt        time.Time
	updatedAt        time.Time
}

type ApplicationID struct {
	value string
}

func NewApplicationID() ApplicationID {
	return ApplicationID{value: uuid.New().String()}
}

func ApplicationIDFromString(s string) (ApplicationID, error) {
	if s == "" {
		return ApplicationID{}, fmt.Errorf("application ID cannot be empty")
	}
	return ApplicationID{value: s}, nil
}

func (id ApplicationID) String() string {
	return id.value
}

type ApplicationName struct {
	value string
}

func NewApplicationName(name string) (ApplicationName, error) {
	if name == "" {
		return ApplicationName{}, fmt.Errorf("application name cannot be empty")
	}
	if len(name) > 100 {
		return ApplicationName{}, fmt.Errorf("application name cannot exceed 100 characters")
	}
	return ApplicationName{value: name}, nil
}

func (n ApplicationName) String() string {
	return n.value
}

type BuildpackType string

const (
	BuildpackTypeNixpacks      BuildpackType = "nixpacks"
	BuildpackTypeStatic        BuildpackType = "static"
	BuildpackTypeDockerfile    BuildpackType = "dockerfile"
	BuildpackTypeDockerCompose BuildpackType = "docker-compose"
	BuildpackTypeBuildpacks    BuildpackType = "buildpacks"
)

type ApplicationStatus string

const (
	ApplicationStatusCreated   ApplicationStatus = "created"
	ApplicationStatusBuilding  ApplicationStatus = "building"
	ApplicationStatusDeploying ApplicationStatus = "deploying"
	ApplicationStatusRunning   ApplicationStatus = "running"
	ApplicationStatusStopped   ApplicationStatus = "stopped"
	ApplicationStatusFailed    ApplicationStatus = "failed"
)

func NewApplication(
	name ApplicationName,
	description string,
	projectID, environmentID uuid.UUID,
	deploymentSource DeploymentSource,
	buildpack BuildpackConfig,
) *Application {
	now := time.Now()
	return &Application{
		id:               NewApplicationID(),
		name:             name,
		description:      description,
		projectID:        projectID,
		environmentID:    environmentID,
		deploymentSource: deploymentSource,
		buildpack:        buildpack,
		envVars:          make(map[string]string),
		autoDeploy:       true,
		status:           ApplicationStatusCreated,
		createdAt:        now,
		updatedAt:        now,
	}
}

func (a *Application) ID() ApplicationID {
	return a.id
}

func (a *Application) Name() ApplicationName {
	return a.name
}

func (a *Application) Description() string {
	return a.description
}

func (a *Application) ProjectID() uuid.UUID {
	return a.projectID
}

func (a *Application) EnvironmentID() uuid.UUID {
	return a.environmentID
}

func (a *Application) DeploymentSource() DeploymentSource {
	return a.deploymentSource
}

func (a *Application) RepoURL() string {
	if a.deploymentSource.Type == DeploymentSourceTypeGit && a.deploymentSource.GitRepo != nil {
		return a.deploymentSource.GitRepo.URL
	}
	return ""
}

func (a *Application) RepoBranch() string {
	if a.deploymentSource.Type == DeploymentSourceTypeGit && a.deploymentSource.GitRepo != nil {
		return a.deploymentSource.GitRepo.Branch
	}
	return ""
}

func (a *Application) RepoPath() string {
	if a.deploymentSource.Type == DeploymentSourceTypeGit && a.deploymentSource.GitRepo != nil {
		return a.deploymentSource.GitRepo.Path
	}
	return ""
}

func (a *Application) Domain() string {
	return a.domain
}

func (a *Application) BuildpackType() BuildpackType {
	return a.buildpack.Type
}

func (a *Application) Config() string {
	if configJSON, err := json.Marshal(a.buildpack.Config); err == nil {
		return string(configJSON)
	}
	return "{}"
}

func (a *Application) Buildpack() BuildpackConfig {
	return a.buildpack
}

func (a *Application) EnvVars() map[string]string {
	result := make(map[string]string)
	for k, v := range a.envVars {
		result[k] = v
	}
	return result
}

func (a *Application) AutoDeploy() bool {
	return a.autoDeploy
}

func (a *Application) Status() ApplicationStatus {
	return a.status
}

func (a *Application) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Application) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Application) UpdateDescription(description string) {
	a.description = description
	a.updatedAt = time.Now()
}

func (a *Application) SetDeploymentSource(source DeploymentSource) {
	a.deploymentSource = source
	a.updatedAt = time.Now()
}

func (a *Application) SetRepoURL(repoURL string) {
	if a.deploymentSource.Type == DeploymentSourceTypeGit {
		if a.deploymentSource.GitRepo == nil {
			a.deploymentSource.GitRepo = &GitRepoSource{}
		}
		a.deploymentSource.GitRepo.URL = repoURL
		a.updatedAt = time.Now()
	}
}

func (a *Application) SetRepoBranch(branch string) {
	if a.deploymentSource.Type == DeploymentSourceTypeGit {
		if a.deploymentSource.GitRepo == nil {
			a.deploymentSource.GitRepo = &GitRepoSource{}
		}
		a.deploymentSource.GitRepo.Branch = branch
		a.updatedAt = time.Now()
	}
}

func (a *Application) SetRepoPath(path string) {
	if a.deploymentSource.Type == DeploymentSourceTypeGit {
		if a.deploymentSource.GitRepo == nil {
			a.deploymentSource.GitRepo = &GitRepoSource{}
		}
		a.deploymentSource.GitRepo.Path = path
		a.updatedAt = time.Now()
	}
}

func (a *Application) SetDomain(domain string) {
	a.domain = domain
	a.updatedAt = time.Now()
}

func (a *Application) SetBuildpackType(buildpackType BuildpackType) {
	a.buildpack.Type = buildpackType
	a.updatedAt = time.Now()
}

func (a *Application) UpdateConfig(config string) {
	var configData interface{}
	if err := json.Unmarshal([]byte(config), &configData); err == nil {
		a.buildpack.Config = configData
	}
	a.updatedAt = time.Now()
}

func (a *Application) SetBuildpack(buildpack BuildpackConfig) {
	a.buildpack = buildpack
	a.updatedAt = time.Now()
}

func (a *Application) SetEnvVar(key, value string) {
	if a.envVars == nil {
		a.envVars = make(map[string]string)
	}
	a.envVars[key] = value
	a.updatedAt = time.Now()
}

func (a *Application) RemoveEnvVar(key string) {
	if a.envVars != nil {
		delete(a.envVars, key)
		a.updatedAt = time.Now()
	}
}

func (a *Application) SetEnvVars(envVars map[string]string) {
	a.envVars = make(map[string]string)
	for k, v := range envVars {
		a.envVars[k] = v
	}
	a.updatedAt = time.Now()
}

func (a *Application) SetAutoDeploy(autoDeploy bool) {
	a.autoDeploy = autoDeploy
	a.updatedAt = time.Now()
}

func (a *Application) ChangeStatus(status ApplicationStatus) {
	a.status = status
	a.updatedAt = time.Now()
}

func (a *Application) CanDeploy() error {
	switch a.status {
	case ApplicationStatusBuilding:
		return fmt.Errorf("application is currently building")
	case ApplicationStatusDeploying:
		return fmt.Errorf("application is currently deploying")
	default:
		return nil
	}
}

func (a *Application) CanStop() error {
	if a.status != ApplicationStatusRunning {
		return fmt.Errorf("application is not running")
	}
	return nil
}

func ReconstructApplication(
	id ApplicationID,
	name ApplicationName,
	description string,
	projectID, environmentID uuid.UUID,
	deploymentSource DeploymentSource,
	domain string,
	buildpack BuildpackConfig,
	envVars map[string]string,
	autoDeploy bool,
	status ApplicationStatus,
	createdAt, updatedAt time.Time,
) *Application {
	if envVars == nil {
		envVars = make(map[string]string)
	}
	return &Application{
		id:               id,
		name:             name,
		description:      description,
		projectID:        projectID,
		environmentID:    environmentID,
		deploymentSource: deploymentSource,
		domain:           domain,
		buildpack:        buildpack,
		envVars:          envVars,
		autoDeploy:       autoDeploy,
		status:           status,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// Helper functions for creating deployment sources
func NewGitDeploymentSource(url, branch, path, token string) DeploymentSource {
	if branch == "" {
		branch = "main"
	}
	return DeploymentSource{
		Type: DeploymentSourceTypeGit,
		GitRepo: &GitRepoSource{
			URL:    url,
			Branch: branch,
			Path:   path,
			Token:  token,
		},
	}
}

func NewRegistryDeploymentSource(image, tag string) DeploymentSource {
	if tag == "" {
		tag = "latest"
	}
	return DeploymentSource{
		Type: DeploymentSourceTypeRegistry,
		Registry: &RegistrySource{
			Image: image,
			Tag:   tag,
		},
	}
}

func NewUploadDeploymentSource(filename, filePath string) DeploymentSource {
	return DeploymentSource{
		Type: DeploymentSourceTypeUpload,
		Upload: &UploadSource{
			Filename: filename,
			FilePath: filePath,
		},
	}
}

// Helper functions for creating buildpack configs
func NewBuildpackConfig(buildpackType BuildpackType, config interface{}) BuildpackConfig {
	return BuildpackConfig{
		Type:   buildpackType,
		Config: config,
	}
}
