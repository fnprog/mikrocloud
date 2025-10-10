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

	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

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

	if isContainerRunning(runtime) {
		fmt.Println("✅ Mikrocloud is already running")
		return nil
	}

	if !isContainerExists(runtime) {
		slog.Info("Pulling Mikrocloud container image...", "image", containerImage)
		fmt.Printf("⬇️  Pulling %s...\n", containerImage)

		pullCmd := exec.Command(runtime, "pull", containerImage)
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr

		if err := pullCmd.Run(); err != nil {
			return fmt.Errorf("failed to pull container image: %w", err)
		}
	}

	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	slog.Info("Starting container...", "name", containerName)
	fmt.Printf("🚀 Starting Mikrocloud container...\n")

	runCmd := exec.Command(runtime, "run",
		"-d",
		"--name", containerName,
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

	fmt.Println("✅ Mikrocloud started successfully!")
	fmt.Printf("🌐 Dashboard: http://localhost:%d\n", port)
	return nil
}

func stopMikrocloud() error {
	slog.Info("Stopping Mikrocloud server...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	runtime := cfg.Docker.Runtime

	if !isContainerRunning(runtime) {
		fmt.Println("⚠️  Mikrocloud is not running")
		return nil
	}

	fmt.Printf("⏹️  Stopping Mikrocloud container...\n")

	stopCmd := exec.Command(runtime, "stop", containerName)
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	fmt.Println("✅ Mikrocloud stopped")
	return nil
}

func statusMikrocloud() error {
	slog.Info("Checking Mikrocloud status...")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	runtime := cfg.Docker.Runtime
	port := cfg.Server.Port

	if !isContainerExists(runtime) {
		fmt.Println("❌ Mikrocloud container does not exist")
		fmt.Printf("   Run 'mikrocloud-cli start' to create and start the container\n")
		return nil
	}

	if isContainerRunning(runtime) {
		fmt.Println("✅ Mikrocloud is running")
		fmt.Printf("🌐 Dashboard: http://localhost:%d\n", port)

		inspectCmd := exec.Command(runtime, "inspect", "--format", "{{.State.StartedAt}}", containerName)
		output, err := inspectCmd.Output()
		if err == nil {
			fmt.Printf("   Started: %s", strings.TrimSpace(string(output)))
		}
	} else {
		fmt.Println("⏹️  Mikrocloud container exists but is not running")
		fmt.Printf("   Run 'mikrocloud-cli start' to start it\n")
	}

	return nil
}

func isContainerExists(runtime string) bool {
	cmd := exec.Command(runtime, "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", containerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == containerName
}

func isContainerRunning(runtime string) bool {
	cmd := exec.Command(runtime, "ps", "--filter", fmt.Sprintf("name=^%s$", containerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == containerName
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
