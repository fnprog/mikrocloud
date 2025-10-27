package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Analytics AnalyticsConfig `mapstructure:"analytics"`
	Docker    DockerConfig    `mapstructure:"docker"`
	SSL       SSLConfig       `mapstructure:"ssl"`
	Auth      AuthConfig      `mapstructure:"auth"`
	SMTP      SMTPConfig      `mapstructure:"smtp"`
	Proxy     ProxyConfig     `mapstructure:"proxy"`
	Metrics   MetricsConfig   `mapstructure:"metrics"`
	Tunnel    TunnelConfig    `mapstructure:"tunnel"`
}

type ServerConfig struct {
	Host           string   `mapstructure:"host"`
	Port           int      `mapstructure:"port"`
	DataDir        string   `mapstructure:"data_dir"`
	LogLevel       string   `mapstructure:"log_level"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	PublicIP       string   `mapstructure:"public_ip"`
	PublicURL      string   `mapstructure:"public_url"`
}

type DatabaseConfig struct {
	Type      string                  `mapstructure:"type"`      // "sqlite" or "postgres"
	URL       string                  `mapstructure:"url"`       // File path for SQLite, connection string for PostgreSQL
	Container DatabaseContainerConfig `mapstructure:"container"` // Container configuration for auto-start
}

type DatabaseContainerConfig struct {
	AutoStart bool   `mapstructure:"auto_start"` // Auto-start container if type=postgres
	Image     string `mapstructure:"image"`
	Port      int    `mapstructure:"port"`
	Database  string `mapstructure:"database"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
}

type AnalyticsConfig struct {
	Type      string                   `mapstructure:"type"`      // "duckdb" or "clickhouse"
	URL       string                   `mapstructure:"url"`       // File path for DuckDB, connection string for ClickHouse
	Container AnalyticsContainerConfig `mapstructure:"container"` // Container configuration for auto-start
}

type AnalyticsContainerConfig struct {
	AutoStart  bool   `mapstructure:"auto_start"` // Auto-start container if type=clickhouse
	Image      string `mapstructure:"image"`
	HTTPPort   int    `mapstructure:"http_port"`
	NativePort int    `mapstructure:"native_port"`
	Database   string `mapstructure:"database"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
}

type ProxyConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	AutoStart     bool   `mapstructure:"auto_start"`
	Image         string `mapstructure:"image"`
	HTTPPort      int    `mapstructure:"http_port"`
	HTTPSPort     int    `mapstructure:"https_port"`
	DashboardPort int    `mapstructure:"dashboard_port"`
}

type MetricsConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	AutoStart bool   `mapstructure:"auto_start"`
	Image     string `mapstructure:"image"`
	Port      int    `mapstructure:"port"`
}

type TunnelConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	AutoStart bool   `mapstructure:"auto_start"`
	Token     string `mapstructure:"token"`
}

type DockerConfig struct {
	Runtime     string `mapstructure:"runtime"` // "docker" or "podman"
	SocketPath  string `mapstructure:"socket_path"`
	Rootless    bool   `mapstructure:"rootless"`
	BuildDir    string `mapstructure:"build_dir"`    // Directory for build workspaces
	NetworkMode string `mapstructure:"network_mode"` // Network mode for proxy: "bridge" or "host"
}

type SSLConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	ACMEEmail string `mapstructure:"acme_email"`
	Staging   bool   `mapstructure:"staging"`
	CertsDir  string `mapstructure:"certs_dir"`
}

type SMTPConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	FromEmail string `mapstructure:"from_email"`
	FromName  string `mapstructure:"from_name"`
}

type AuthConfig struct {
	JWTSecret string          `mapstructure:"jwt_secret"`
	Enabled   bool            `mapstructure:"enabled"`
	UserOAuth UserOAuthConfig `mapstructure:"user_oauth"`
	GitOAuth  GitOAuthConfig  `mapstructure:"git_oauth"`
}

type GitOAuthConfig struct {
	GitHub    GitHubOAuthConfig    `mapstructure:"github"`
	GitLab    GitLabOAuthConfig    `mapstructure:"gitlab"`
	Bitbucket BitbucketOAuthConfig `mapstructure:"bitbucket"`
}

type UserOAuthConfig struct {
	GitHub GitHubOAuthConfig `mapstructure:"github"`
	GitLab GitLabOAuthConfig `mapstructure:"gitlab"`
	Google GoogleOAuthConfig `mapstructure:"google"`
}

type GitHubOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type GitLabOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type BitbucketOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	Tenant       string `mapstructure:"tenant"`
}

func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	setDefaults()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand environment variables in paths
	cfg.Server.DataDir = expandEnvVars(cfg.Server.DataDir)
	cfg.Database.URL = expandEnvVars(cfg.Database.URL)
	cfg.Analytics.URL = expandEnvVars(cfg.Analytics.URL)
	cfg.SSL.CertsDir = expandEnvVars(cfg.SSL.CertsDir)
	cfg.Docker.BuildDir = expandEnvVars(cfg.Docker.BuildDir)

	return &cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 3000)
	viper.SetDefault("server.data_dir", "${HOME}/.local/share/mikrocloud")
	viper.SetDefault("server.log_level", "info")
	viper.SetDefault("server.allowed_origins", []string{"*"})
	viper.SetDefault("server.public_ip", "")
	viper.SetDefault("server.public_url", "")

	// Database defaults - SQLite database path
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.url", "${HOME}/.local/share/mikrocloud/mikrocloud.db")
	viper.SetDefault("database.container.auto_start", false)
	viper.SetDefault("database.container.image", "postgres:16-alpine")
	viper.SetDefault("database.container.port", 5432)
	viper.SetDefault("database.container.database", "mikrocloud")
	viper.SetDefault("database.container.user", "mikrocloud")
	viper.SetDefault("database.container.password", "mikrocloud")

	// Analytics defaults - DuckDB database path
	viper.SetDefault("analytics.type", "duckdb")
	viper.SetDefault("analytics.url", "${HOME}/.local/share/mikrocloud/analytics.duckdb")
	viper.SetDefault("analytics.container.auto_start", false)
	viper.SetDefault("analytics.container.image", "clickhouse/clickhouse-server:latest")
	viper.SetDefault("analytics.container.http_port", 8123)
	viper.SetDefault("analytics.container.native_port", 9000)
	viper.SetDefault("analytics.container.database", "mikrocloud_analytics")
	viper.SetDefault("analytics.container.user", "default")
	viper.SetDefault("analytics.container.password", "")

	// Proxy defaults
	viper.SetDefault("proxy.enabled", true)
	viper.SetDefault("proxy.auto_start", true)
	viper.SetDefault("proxy.image", "traefik:v3.0")
	viper.SetDefault("proxy.http_port", 80)
	viper.SetDefault("proxy.https_port", 443)
	viper.SetDefault("proxy.dashboard_port", 8080)

	// Metrics defaults
	viper.SetDefault("metrics.enabled", false)
	viper.SetDefault("metrics.auto_start", false)
	viper.SetDefault("metrics.image", "prom/prometheus:latest")
	viper.SetDefault("metrics.port", 9090)

	// Tunnel defaults
	viper.SetDefault("tunnel.enabled", false)
	viper.SetDefault("tunnel.auto_start", false)
	viper.SetDefault("tunnel.token", "")

	// Docker defaults
	viper.SetDefault("docker.runtime", "docker")
	viper.SetDefault("docker.socket_path", "/var/run/docker.sock")
	viper.SetDefault("docker.rootless", false)
	viper.SetDefault("docker.build_dir", "${HOME}/.local/share/mikrocloud/builds")
	viper.SetDefault("docker.network_mode", "bridge")

	// SSL defaults
	viper.SetDefault("ssl.enabled", false)
	viper.SetDefault("ssl.staging", true)
	viper.SetDefault("ssl.certs_dir", "${HOME}/.local/share/mikrocloud/certs")

	// SMTP defaults
	viper.SetDefault("smtp.enabled", false)
	viper.SetDefault("smtp.host", "")
	viper.SetDefault("smtp.port", 587)
	viper.SetDefault("smtp.username", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.from_email", "")
	viper.SetDefault("smtp.from_name", "Mikrocloud")

	// Auth defaults
	viper.SetDefault("auth.enabled", false)
	viper.SetDefault("auth.jwt_secret", "")

	// User OAuth defaults
	viper.SetDefault("auth.user_oauth.github.client_id", "")
	viper.SetDefault("auth.user_oauth.github.client_secret", "")
	viper.SetDefault("auth.user_oauth.github.redirect_url", "")
	viper.SetDefault("auth.user_oauth.gitlab.client_id", "")
	viper.SetDefault("auth.user_oauth.gitlab.client_secret", "")
	viper.SetDefault("auth.user_oauth.gitlab.redirect_url", "")
	viper.SetDefault("auth.user_oauth.google.client_id", "")
	viper.SetDefault("auth.user_oauth.google.client_secret", "")
	viper.SetDefault("auth.user_oauth.google.redirect_url", "")
	viper.SetDefault("auth.user_oauth.google.tenant", "")

	// Git OAuth defaults - these should be set via environment variables
	viper.SetDefault("auth.git_oauth.github.client_id", "")
	viper.SetDefault("auth.git_oauth.github.client_secret", "")
	viper.SetDefault("auth.git_oauth.github.redirect_url", "")
	viper.SetDefault("auth.git_oauth.gitlab.client_id", "")
	viper.SetDefault("auth.git_oauth.gitlab.client_secret", "")
	viper.SetDefault("auth.git_oauth.gitlab.redirect_url", "")
	viper.SetDefault("auth.git_oauth.bitbucket.client_id", "")
	viper.SetDefault("auth.git_oauth.bitbucket.client_secret", "")
	viper.SetDefault("auth.git_oauth.bitbucket.redirect_url", "")
}

func expandEnvVars(path string) string {
	if strings.Contains(path, "${") {
		return os.ExpandEnv(path)
	}
	return path
}

func (c *Config) GetPublicURL() string {
	if c.Server.PublicURL != "" {
		return strings.TrimSuffix(c.Server.PublicURL, "/")
	}

	if c.Server.PublicIP != "" {
		publicIP := c.Server.PublicIP
		if strings.HasPrefix(publicIP, "http://") || strings.HasPrefix(publicIP, "https://") {
			return strings.TrimSuffix(publicIP, "/")
		}
		return fmt.Sprintf("https://%s", publicIP)
	}

	return fmt.Sprintf("http://%s:%d", c.Server.Host, c.Server.Port)
}
