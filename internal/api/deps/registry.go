package deps

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/database"
	activitiesService "github.com/mikrocloud/mikrocloud/internal/domain/activities/service"
	applicationsService "github.com/mikrocloud/mikrocloud/internal/domain/applications/service"
	authService "github.com/mikrocloud/mikrocloud/internal/domain/auth/service"
	databaseService "github.com/mikrocloud/mikrocloud/internal/domain/databases/service"
	deploymentService "github.com/mikrocloud/mikrocloud/internal/domain/deployments/service"
	diskService "github.com/mikrocloud/mikrocloud/internal/domain/disks/service"
	environmentService "github.com/mikrocloud/mikrocloud/internal/domain/environments/service"
	gitService "github.com/mikrocloud/mikrocloud/internal/domain/git/service"
	organizationsService "github.com/mikrocloud/mikrocloud/internal/domain/organizations/service"
	projectService "github.com/mikrocloud/mikrocloud/internal/domain/projects/service"
	proxyService "github.com/mikrocloud/mikrocloud/internal/domain/proxy/service"
	serversService "github.com/mikrocloud/mikrocloud/internal/domain/servers/service"
	"github.com/mikrocloud/mikrocloud/internal/domain/services/repository"
	templatesService "github.com/mikrocloud/mikrocloud/internal/domain/services/service"
	settingsService "github.com/mikrocloud/mikrocloud/internal/domain/settings/service"
	tunnelService "github.com/mikrocloud/mikrocloud/internal/domain/tunnels/service"
	buildService "github.com/mikrocloud/mikrocloud/pkg/containers/build"
	containerService "github.com/mikrocloud/mikrocloud/pkg/containers/service"

	databaseContainers "github.com/mikrocloud/mikrocloud/pkg/containers/database"
	proxyContainers "github.com/mikrocloud/mikrocloud/pkg/containers/proxy"
	tunnelContainers "github.com/mikrocloud/mikrocloud/pkg/containers/tunnel"

	"github.com/mikrocloud/mikrocloud/internal/domain/domains"
)

type Dependencies struct {
	DB      *database.Database
	Config  *config.Config
	JwtKeys *jwtauth.JWTAuth

	ContainerService    *containerService.ContainerService
	AuthService         *authService.AuthService
	ActivitiesService   *activitiesService.ActivitiesService
	ProjectService      *projectService.ProjectService
	DatabaseService     *databaseService.DatabaseService
	OrganizationService *organizationsService.OrganizationService
	ServerService       *serversService.ServersService
	ApplicationService  *applicationsService.ApplicationService
	DeploymentService   *deploymentService.DeploymentService
	EnvironmentService  *environmentService.EnvironmentService
	GitService          *gitService.GitService
	ProxyService        *proxyService.ProxyService
	TraefikService      *proxyContainers.TraefikService

	BuildService    *buildService.BuildService
	DiskService     *diskService.DiskService
	TemplateService *templatesService.TemplateService
	SettingsService *settingsService.SettingsService
	TunnelService   *tunnelService.TunnelService

	// Sync services
	DatabaseStatusSyncService *databaseService.StatusSyncService
}

func NewDependencies(cfg *config.Config, db *database.Database) (*Dependencies, error) {
	// Init Auth middleware
	tokenAuthSecret := jwtauth.New("HS256", []byte(cfg.Auth.JWTSecret), nil)

	containerService, err := containerService.NewContainerService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init containerService: %s", err)
	}

	// App service helper
	domainGenerator := domains.NewDomainGenerator(cfg.Server.PublicIP)

	activitiesService := activitiesService.NewActivitiesService(db.ActivitiesRepository)
	authService := authService.NewAuthService(db.SessionRepository, db.AuthRepository, db.UserRepository, cfg.Auth.JWTSecret)
	organizationsSvc := organizationsService.NewOrganizationService(db.OrganizationRepository)
	envService := environmentService.NewEnvironmentService(db.EnvironmentRepository)
	projService := projectService.NewProjectService(db.ProjectRepository, db.EnvironmentRepository, activitiesService)

	serversSvc := serversService.NewServersService(db.ServersRepository)
	gitSvc := gitService.NewGitService(db.GitRepository)
	settingsSvc := settingsService.NewSettingsService(db.SettingsRepository)

	proxySvc := proxyService.New(db.ProxyRepository, db.TraefikConfigRepository)
	traefikConfigDir := filepath.Join(cfg.Server.DataDir, "traefik")
	traefikSvc := proxyContainers.NewTraefikService(*containerService, traefikConfigDir, cfg.Docker.NetworkMode)

	deploymentSvc := deploymentService.NewDeploymentService(db.DeploymentRepository, containerService)
	diskSvc := diskService.NewDiskService(db.DiskRepository, db.DiskBackupRepository)
	dbDeploymentSvc := databaseContainers.NewDatabaseDeploymentService(containerService, diskSvc)

	appSvc := applicationsService.NewApplicationService(db.ApplicationRepository, domainGenerator, deploymentSvc)
	databaseSvc := databaseService.NewDatabaseService(db.DatabaseRepository, dbDeploymentSvc, diskSvc)
	quickDeployService := repository.NewQuickDeployService(db.TemplateRepository, appSvc)
	templateSvc := templatesService.NewTemplateService(db.TemplateRepository, quickDeployService)

	dbStatusSyncSvc := databaseService.NewStatusSyncService(databaseSvc, containerService, 29*time.Second)

	cloudflaredMgr := tunnelContainers.NewCloudflaredManager(containerService.GetManager())
	tunnelSvc := tunnelService.NewTunnelService(db.TunnelRepository, cloudflaredMgr)

	return &Dependencies{
		DB:                  db,
		Config:              cfg,
		ContainerService:    containerService,
		AuthService:         authService,
		ActivitiesService:   activitiesService,
		ProjectService:      projService,
		DatabaseService:     databaseSvc,
		OrganizationService: organizationsSvc,
		ServerService:       serversSvc,
		ApplicationService:  appSvc,
		DeploymentService:   deploymentSvc,
		EnvironmentService:  envService,
		GitService:          gitSvc,
		ProxyService:        proxySvc,
		TraefikService:      traefikSvc,
		DiskService:         diskSvc,
		TemplateService:     templateSvc,
		SettingsService:     settingsSvc,
		TunnelService:       tunnelSvc,

		DatabaseStatusSyncService: dbStatusSyncSvc,
		JwtKeys:                   tokenAuthSecret,
	}, nil
}
