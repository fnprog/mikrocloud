package applications

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Application struct {
	id            ApplicationID
	name          ApplicationName
	description   string
	projectID     uuid.UUID
	environmentID uuid.UUID
	repoURL       string
	repoBranch    string
	repoPath      string
	domain        string
	buildpackType BuildpackType
	config        string
	autoDeploy    bool
	status        ApplicationStatus
	createdAt     time.Time
	updatedAt     time.Time
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
	repoURL string,
	buildpackType BuildpackType,
) *Application {
	now := time.Now()
	return &Application{
		id:            NewApplicationID(),
		name:          name,
		description:   description,
		projectID:     projectID,
		environmentID: environmentID,
		repoURL:       repoURL,
		repoBranch:    "main",
		buildpackType: buildpackType,
		config:        "{}",
		autoDeploy:    true,
		status:        ApplicationStatusCreated,
		createdAt:     now,
		updatedAt:     now,
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

func (a *Application) RepoURL() string {
	return a.repoURL
}

func (a *Application) RepoBranch() string {
	return a.repoBranch
}

func (a *Application) RepoPath() string {
	return a.repoPath
}

func (a *Application) Domain() string {
	return a.domain
}

func (a *Application) BuildpackType() BuildpackType {
	return a.buildpackType
}

func (a *Application) Config() string {
	return a.config
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

func (a *Application) SetRepoURL(repoURL string) {
	a.repoURL = repoURL
	a.updatedAt = time.Now()
}

func (a *Application) SetRepoBranch(branch string) {
	a.repoBranch = branch
	a.updatedAt = time.Now()
}

func (a *Application) SetRepoPath(path string) {
	a.repoPath = path
	a.updatedAt = time.Now()
}

func (a *Application) SetDomain(domain string) {
	a.domain = domain
	a.updatedAt = time.Now()
}

func (a *Application) SetBuildpackType(buildpackType BuildpackType) {
	a.buildpackType = buildpackType
	a.updatedAt = time.Now()
}

func (a *Application) UpdateConfig(config string) {
	a.config = config
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
	repoURL, repoBranch, repoPath, domain string,
	buildpackType BuildpackType,
	config string,
	autoDeploy bool,
	status ApplicationStatus,
	createdAt, updatedAt time.Time,
) *Application {
	return &Application{
		id:            id,
		name:          name,
		description:   description,
		projectID:     projectID,
		environmentID: environmentID,
		repoURL:       repoURL,
		repoBranch:    repoBranch,
		repoPath:      repoPath,
		domain:        domain,
		buildpackType: buildpackType,
		config:        config,
		autoDeploy:    autoDeploy,
		status:        status,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}
