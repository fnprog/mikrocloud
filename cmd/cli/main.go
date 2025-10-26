package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mikrocloud/mikrocloud/internal/config"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "mikrocloud-cli",
		Short: "Mikrocloud CLI - Manage your PaaS infrastructure",
		Long:  `Mikrocloud CLI provides command-line access to manage projects, applications, databases, and infrastructure.`,
	}
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("CLI failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./mikrocloud.toml)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")

	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(appCmd)
	rootCmd.AddCommand(dbCmd)

	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectDeleteCmd)

	appCmd.AddCommand(appDeployCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appDeleteCmd)

	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbDeleteCmd)
}

func startMikrocloud() error {
	slog.Info("Starting Mikrocloud server...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	runtime := cfg.Docker.Runtime
	port := cfg.Server.Port
	dataDir := cfg.Server.DataDir
	socketPath := cfg.Docker.SocketPath

	if isContainerRunning(runtime, containerName) {
		fmt.Println("✅ Mikrocloud is already running")
		return nil
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	if err := ensureNetwork(runtime); err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}

	if cfg.Queue.Enabled && cfg.Queue.AutoStart {
		if err := startQueue(runtime, cfg); err != nil {
			return fmt.Errorf("failed to start queue: %w", err)
		}
	}

	if cfg.Proxy.Enabled && cfg.Proxy.AutoStart {
		if err := startProxy(runtime, cfg); err != nil {
			return fmt.Errorf("failed to start proxy: %w", err)
		}
	}

	if !isContainerExists(runtime, containerName) {
		slog.Info("Pulling Mikrocloud container image...", "image", containerImage)
		fmt.Printf("⬇️  Pulling %s...\n", containerImage)

		pullCmd := exec.Command(runtime, "pull", containerImage)
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr

		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull container image: %w", err)
		}
	}

	slog.Info("Starting main container...", "name", containerName)
	fmt.Printf("🚀 Starting Mikrocloud container...\n")

	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", containerName,
		"--network", networkName,
		"--restart", "unless-stopped",
		"-p", fmt.Sprintf("%d:%d", port, port),
		"-v", fmt.Sprintf("%s:/app/data", dataDir),
		"-v", fmt.Sprintf("%s:/var/run/docker.sock", socketPath),
		containerImage,
	)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %w\n%s", err, string(output))
	}

	if cfg.Metrics.Enabled && cfg.Metrics.AutoStart {
		if err := startMetrics(runtime, cfg); err != nil {
			slog.Warn("Failed to start metrics", "error", err)
		}
	}

	if cfg.Tunnel.Enabled && cfg.Tunnel.AutoStart && cfg.Tunnel.Token != "" {
		if err := startTunnel(runtime, cfg); err != nil {
			slog.Warn("Failed to start tunnel", "error", err)
		}
	}

	fmt.Println("✅ Mikrocloud started successfully!")
	fmt.Printf("🌐 Dashboard: http://localhost:%d\n", port)

	if cfg.Proxy.Enabled && cfg.Proxy.AutoStart {
		fmt.Printf("🔀 Proxy Dashboard: http://localhost:%d\n", cfg.Proxy.DashboardPort)
	}

	if cfg.Metrics.Enabled && cfg.Metrics.AutoStart {
		fmt.Printf("📊 Metrics: http://localhost:%d\n", cfg.Metrics.Port)
	}

	return nil
}

func stopMikrocloud() error {
	slog.Info("Stopping Mikrocloud server...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	runtime := cfg.Docker.Runtime

	if !isContainerRunning(runtime, containerName) {
		fmt.Println("⚠️  Mikrocloud is not running")
		return nil
	}

	fmt.Printf("⏹️  Stopping Mikrocloud containers...\n")

	stopContainer(runtime, containerName)
	stopContainer(runtime, tunnelName)
	stopContainer(runtime, metricsName)
	stopContainer(runtime, proxyName)
	stopContainer(runtime, queueName)

	fmt.Println("✅ Mikrocloud stopped")
	return nil
}

func stopContainer(runtime, name string) {
	if isContainerRunning(runtime, name) {
		_ = exec.Command(runtime, "stop", name).Run()
	}
}

func statusMikrocloud() error {
	slog.Info("Checking Mikrocloud status...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	runtime := cfg.Docker.Runtime
	port := cfg.Server.Port

	if !isContainerExists(runtime, containerName) {
		fmt.Println("❌ Mikrocloud container does not exist")
		fmt.Printf("   Run 'mikrocloud-cli start' to create and start the container\n")
		return nil
	}

	if isContainerRunning(runtime, containerName) {
		fmt.Println("✅ Mikrocloud is running")
		fmt.Printf("🌐 Dashboard: http://localhost:%d\n", port)

		inspectCmd := exec.Command(runtime, "inspect", "--format", "{{.State.StartedAt}}", containerName)
		output, err := inspectCmd.Output()
		if err == nil {
			fmt.Printf("   Started: %s\n", strings.TrimSpace(string(output)))
		}

		if isContainerRunning(runtime, queueName) {
			fmt.Println("✅ Queue is running")
		}
		if isContainerRunning(runtime, proxyName) {
			fmt.Printf("✅ Proxy is running (Dashboard: http://localhost:%d)\n", cfg.Proxy.DashboardPort)
		}
		if isContainerRunning(runtime, metricsName) {
			fmt.Printf("✅ Metrics is running (http://localhost:%d)\n", cfg.Metrics.Port)
		}
		if isContainerRunning(runtime, tunnelName) {
			fmt.Println("✅ Tunnel is running")
		}
	} else {
		fmt.Println("⏹️  Mikrocloud container exists but is not running")
		fmt.Printf("   Run 'mikrocloud-cli start' to start it\n")
	}

	return nil
}

func isContainerExists(runtime, name string) bool {
	cmd := exec.Command(runtime, "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", name), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == name
}

func isContainerRunning(runtime, name string) bool {
	cmd := exec.Command(runtime, "ps", "--filter", fmt.Sprintf("name=^%s$", name), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == name
}

func ensureNetwork(runtime string) error {
	cmd := exec.Command(runtime, "network", "inspect", networkName)
	if err := cmd.Run(); err == nil {
		return nil
	}

	fmt.Printf("🌐 Creating network %s...\n", networkName)
	createCmd := exec.Command(runtime, "network", "create", networkName)
	return createCmd.Run()
}

func startQueue(runtime string, cfg *config.Config) error {
	if isContainerRunning(runtime, queueName) {
		fmt.Println("✅ Queue is already running")
		return nil
	}

	if !isContainerExists(runtime, queueName) {
		fmt.Printf("⬇️  Pulling dragonfly image...\n")
		pullCmd := exec.Command(runtime, "pull", "docker.dragonflydb.io/dragonflydb/dragonfly:latest")
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull dragonfly: %w", err)
		}
	}

	fmt.Printf("🚀 Starting queue container...\n")
	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", queueName,
		"--network", networkName,
		"--restart", "unless-stopped",
		"-p", "6379:6379",
		"docker.dragonflydb.io/dragonflydb/dragonfly:latest",
		"--logtostderr",
	)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start queue: %w\n%s", err, string(output))
	}

	return waitForHealthy(runtime, queueName, 30)
}

func startProxy(runtime string, cfg *config.Config) error {
	if isContainerRunning(runtime, proxyName) {
		fmt.Println("✅ Proxy is already running")
		return nil
	}

	if !isContainerExists(runtime, proxyName) {
		fmt.Printf("⬇️  Pulling %s...\n", cfg.Proxy.Image)
		pullCmd := exec.Command(runtime, "pull", cfg.Proxy.Image)
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull proxy: %w", err)
		}
	}

	fmt.Printf("🚀 Starting proxy container...\n")
	certsDir := fmt.Sprintf("%s/certs", cfg.Server.DataDir)
	_ = os.MkdirAll(certsDir, 0o755)

	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", proxyName,
		"--network", networkName,
		"--restart", "unless-stopped",
		"-p", fmt.Sprintf("%d:80", cfg.Proxy.HTTPPort),
		"-p", fmt.Sprintf("%d:443", cfg.Proxy.HTTPSPort),
		"-p", fmt.Sprintf("%d:8080", cfg.Proxy.DashboardPort),
		"-v", "/var/run/docker.sock:/var/run/docker.sock:ro",
		"-v", fmt.Sprintf("%s:/letsencrypt", certsDir),
		cfg.Proxy.Image,
		"--api.insecure=true",
		"--api.dashboard=true",
		"--providers.docker=true",
		"--providers.docker.network="+networkName,
		"--providers.docker.exposedbydefault=false",
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
	)

	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start proxy: %w\n%s", err, string(output))
	}

	return waitForHealthy(runtime, proxyName, 30)
}

func startMetrics(runtime string, cfg *config.Config) error {
	if isContainerRunning(runtime, metricsName) {
		return nil
	}

	fmt.Printf("⬇️  Pulling %s...\n", cfg.Metrics.Image)
	pullCmd := exec.Command(runtime, "pull", cfg.Metrics.Image)
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("🚀 Starting metrics container...\n")
	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", metricsName,
		"--network", networkName,
		"--restart", "unless-stopped",
		"-p", fmt.Sprintf("%d:9090", cfg.Metrics.Port),
		cfg.Metrics.Image,
	)

	_, err := runCmd.CombinedOutput()
	return err
}

func startTunnel(runtime string, cfg *config.Config) error {
	if isContainerRunning(runtime, tunnelName) {
		return nil
	}

	fmt.Printf("⬇️  Pulling cloudflare/cloudflared...\n")
	pullCmd := exec.Command(runtime, "pull", "cloudflare/cloudflared:latest")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return err
	}

	fmt.Printf("🚀 Starting tunnel...\n")
	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", tunnelName,
		"--network", networkName,
		"--restart", "unless-stopped",
		"cloudflare/cloudflared:latest",
		"tunnel", "--no-autoupdate", "run", "--token", cfg.Tunnel.Token,
	)

	_, err := runCmd.CombinedOutput()
	return err
}

func waitForHealthy(runtime, name string, timeoutSeconds int) error {
	fmt.Printf("⏳ Waiting for %s to be healthy...\n", name)
	for i := 0; i < timeoutSeconds; i++ {
		cmd := exec.Command(runtime, "inspect", "--format", "{{.State.Running}}", name)
		output, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(output)) == "true" {
			fmt.Printf("✅ %s is healthy\n", name)
			return nil
		}
		_ = exec.Command("sleep", "1").Run()
	}
	return fmt.Errorf("timeout waiting for %s to be healthy", name)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("mikrocloud")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.config/mikrocloud")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Using config file", "file", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Mikrocloud CLI v%s\n", "0.1.0")
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Mikrocloud server in the background",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startMikrocloud()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Mikrocloud server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopMikrocloud()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of Mikrocloud services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return statusMikrocloud()
	},
}

const (
	containerName  = "mikrocloud"
	containerImage = "ghcr.io/fnprog/mikrocloud/mikrocloud:latest"
	networkName    = "mikrocloud-network"
	queueName      = "mikrocloud-queue"
	proxyName      = "mikrocloud-proxy"
	metricsName    = "mikrocloud-metrics"
	tunnelName     = "mikrocloud-cloudflared"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Create, list, update, and delete projects`,
}

var projectCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage applications",
	Long:  `Deploy, list, update, and delete applications`,
}

var appDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var appDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an application",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage databases",
	Long:  `Create, list, and delete databases`,
}

var dbCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}

var dbDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("not yet implemented")
	},
}
