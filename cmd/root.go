package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/server"
)

var (
	cfgFile  string
	staticFS fs.FS
	rootCmd  = &cobra.Command{
		Use:   "mikrocloud",
		Short: "Ultra-lightweight Platform as a Service (PaaS)",
		Long:  `Mikrocloud is a next-generation, multi-region Platform as a Service (PaaS) built for ultra-lightweight performance (<50MB memory usage) with enterprise features.`,
	}
)

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

// SetStaticFS sets the static filesystem for the server
func SetStaticFS(fs fs.FS) {
	staticFS = fs
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./mikrocloud.toml)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level (debug, info, warn, error)")

	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	// Add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(migrateCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Mikrocloud server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		srv := server.New(cfg, staticFS)

		return srv.Start(cmd.Context())
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Mikrocloud v%s\n", "0.1.0")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Ensure database directory exists
		dbDir := filepath.Dir(cfg.Database.URL)
		if err := ensureDir(dbDir); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}

		// Open database connection
		db, err := sql.Open("sqlite3", cfg.Database.URL)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		// Set up goose
		goose.SetDialect("sqlite3")

		// Run migrations from embedded migrations directory
		if err := goose.Up(db, "./migrations"); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		slog.Info("Database migrations completed successfully", "database", cfg.Database.URL)
		return nil
	},
}

func ensureDir(dir string) error {
	if dir == "." || dir == "/" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
