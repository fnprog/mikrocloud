package main

import (
	"context"
	"embed"
	"io/fs"
	"os"

	"github.com/mikrocloud/mikrocloud/cmd"
	"golang.org/x/exp/slog"
)

// Embed the entire dist directory from the web build
//
//go:embed all:web/dist
var DistFS embed.FS

// GetStaticFS returns the embedded filesystem, stripping the "web/dist" prefix
func GetStaticFS() (fs.FS, error) {
	return fs.Sub(DistFS, "web/dist")
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	// Set up static filesystem
	staticFS, err := GetStaticFS()

	if err != nil {
		slog.Error("Failed to get static filesystem", "error", err)
		os.Exit(1)
	}

	cmd.SetStaticFS(staticFS)

	ctx := context.Background()
	if err := cmd.Execute(ctx); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}
