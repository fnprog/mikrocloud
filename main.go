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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	// Set up static filesystem
	staticFS, err := fs.Sub(DistFS, "web/dist")
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
