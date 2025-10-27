package server

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"

	"github.com/mikrocloud/mikrocloud/internal/api"
	"github.com/mikrocloud/mikrocloud/internal/api/deps"
	"github.com/mikrocloud/mikrocloud/internal/config"
	"github.com/mikrocloud/mikrocloud/internal/database"
	"github.com/mikrocloud/mikrocloud/internal/domain/proxy"
	"github.com/mikrocloud/mikrocloud/internal/domain/servers"
	proxyContainers "github.com/mikrocloud/mikrocloud/pkg/containers/proxy"
)

type Server struct {
	config     *config.Config
	deps       *deps.Dependencies
	router     *chi.Mux
	staticFS   fs.FS
	traefikSvc *proxyContainers.TraefikService
}

func New(cfg *config.Config, staticFS fs.FS) *Server {
	db, err := database.New(cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	dependencies, err := deps.NewDependencies(cfg, db)
	if err != nil {
		slog.Error("Failed to setup services", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	return &Server{
		config:   cfg,
		deps:     dependencies,
		router:   router,
		staticFS: staticFS,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.setupMiddlewares()
	s.setupAPIRoutes()

	if err := s.setupStaticRoutes(); err != nil {
		return fmt.Errorf("failed to setup Static routes: %w", err)
	}

	if err := s.initializeControlPlaneServer(ctx); err != nil {
		slog.Warn("Failed to initialize control plane server", "error", err)
		// Don't fail the server if this fails, just log the warning
	}

	if err := s.setupDependencies(ctx); err != nil {
		slog.Error("Failed to start Traefik service", "error", err)
		return err
	}

	s.setupBackgroundTasks(ctx)

	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	slog.Info("Starting Mikrocloud server", "address", addr)

	// Start server in a goroutine
	serverChan := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverChan <- err
		}
	}()

	// Wait for interrupt signal or server error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-sigChan:
		slog.Info("Received shutdown signal", "signal", sig)

		// Graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Stop Traefik service first
		if s.traefikSvc != nil {
			slog.Info("Stopping Traefik service")
			if err := s.traefikSvc.Stop(shutdownCtx); err != nil {
				slog.Error("Error stopping Traefik service", "error", err)
			}
		}

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		slog.Info("Server shutdown complete")
		return nil
	}
}

func (s *Server) setupMiddlewares() {
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: s.config.Server.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// TODO: Pipe to a file instead of clogging the term
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				slog.Info("HTTP Request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"duration", time.Since(start),
					"bytes", ww.BytesWritten(),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	})
}

func (s *Server) setupDependencies(ctx context.Context) error {
	// Start Traefik proxy if configured
	if s.config.Proxy.Enabled && s.config.Proxy.AutoStart {
		slog.Info("Starting Traefik proxy container", "image", s.config.Proxy.Image)
		globalConfig := proxy.NewTraefikGlobalConfig()

		if err := s.deps.TraefikService.Start(ctx, globalConfig); err != nil {
			return fmt.Errorf("failed to start Traefik container: %w", err)
		}

		s.traefikSvc = s.deps.TraefikService
		slog.Info("Traefik proxy container started successfully")
	}

	return nil
}

func (s *Server) setupAPIRoutes() {
	s.router.Route("/api", func(r chi.Router) {
		api.SetupRoutes(r, s.deps)
	})
}

func (s *Server) setupStaticRoutes() error {
	frontendFS := s.staticFS

	if frontendFS == nil {
		return fmt.Errorf("no frontend assets available")
	}

	// Catch-all handler for everything non-API
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// 404 For any non-found api routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"error":"API endpoint not found"}`)
			return
		}

		// Try to serve the file from the embedded FS first
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		file, err := frontendFS.Open(path)
		if err == nil {
			defer file.Close()

			// Set appropriate content type
			s.setContentType(w, path)

			// Set caching based on file type
			if strings.Contains(path, "/assets/") || strings.Contains(path, "/_app/") {
				w.Header().Set("Cache-Control", "max-age=31536000, immutable") // 1 year for versioned assets
			} else if strings.HasSuffix(path, ".html") {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // No cache for HTML
			} else {
				w.Header().Set("Cache-Control", "max-age=86400") // 1 day for other static files
			}

			// Use ServeContent for better performance
			if stat, err := file.Stat(); err == nil {
				http.ServeContent(w, r, path, stat.ModTime(), file.(io.ReadSeeker))
			} else {
				io.Copy(w, file)
			}
			return
		}

		// File not found, serve index.html for SPA routing
		s.serveIndexHTML(frontendFS)(w, r)
	})

	return nil
}

func (s *Server) setContentType(w http.ResponseWriter, path string) {
	switch {
	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	case strings.HasSuffix(path, ".json"):
		w.Header().Set("Content-Type", "application/json")
	case strings.HasSuffix(path, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	case strings.HasSuffix(path, ".ico"):
		w.Header().Set("Content-Type", "image/x-icon")
	case strings.HasSuffix(path, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(path, ".gif"):
		w.Header().Set("Content-Type", "image/gif")
	case strings.HasSuffix(path, ".woff"), strings.HasSuffix(path, ".woff2"):
		w.Header().Set("Content-Type", "font/woff")
	case strings.HasSuffix(path, ".ttf"):
		w.Header().Set("Content-Type", "font/ttf")
	}
}

// serveIndexHTML returns a handler that always serves the index.html file
func (s *Server) serveIndexHTML(frontendFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := frontendFS.Open("index.html")
		if err != nil {
			slog.Error("Failed to open index.html", "error", err)
			http.Error(w, "Frontend not available", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		// Read the file content and serve it
		content, err := io.ReadAll(indexFile)
		if err != nil {
			slog.Error("Failed to read index.html", "error", err)
			http.Error(w, "Frontend not available", http.StatusInternalServerError)
			return
		}

		w.Write(content)
	}
}

func (s *Server) setupBackgroundTasks(ctx context.Context) {
	go s.deps.DatabaseStatusSyncService.Start(ctx)
}

func (s *Server) initializeControlPlaneServer(ctx context.Context) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	// Check if organizations exist, if not create default org and system user
	var count int
	err = s.deps.DB.MainDB().DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM organizations").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check organizations: %w", err)
	}

	if count == 0 {
		// Start transaction
		tx, err := s.deps.DB.MainDB().DB().BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback()

		// Generate UUIDs
		systemUserID := uuid.New().String()
		defaultOrgID := uuid.New().String()
		systemRoleID := "role-system-00000000"
		userRoleID := uuid.New().String()
		orgMemberID := uuid.New().String()

		// Dummy password hash (bcrypt of empty string - not usable for login)
		dummyPasswordHash, err := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate dummy password hash: %w", err)
		}

		// Insert system role
		_, err = tx.ExecContext(ctx, `
			INSERT INTO roles (id, name, description, permissions, created_at)
			VALUES (?, 'system', 'System role for automated tasks', '["deployment:create","deployment:read","activity:log"]', CURRENT_TIMESTAMP)
		`, systemRoleID)
		if err != nil {
			return fmt.Errorf("failed to insert system role: %w", err)
		}

		// Insert system user
		_, err = tx.ExecContext(ctx, `
			INSERT INTO users (id, email, password_hash, username, name, status, email_verified_at, timezone, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 'active', CURRENT_TIMESTAMP, 'UTC', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, systemUserID, "system@mikrocloud.local", string(dummyPasswordHash), "mikrocloud-system", "System User")
		if err != nil {
			return fmt.Errorf("failed to insert system user: %w", err)
		}

		// Insert default org
		_, err = tx.ExecContext(ctx, `
			INSERT INTO organizations (id, name, slug, description, owner_id, billing_email, plan, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, 'free', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, defaultOrgID, "Default Organization", "default", "Default organization for mikrocloud", systemUserID, "system@mikrocloud.local")
		if err != nil {
			return fmt.Errorf("failed to insert default org: %w", err)
		}

		// Assign system role to system user
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_roles (id, user_id, role_id, granted_by, granted_at)
			VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		`, userRoleID, systemUserID, systemRoleID, systemUserID)
		if err != nil {
			return fmt.Errorf("failed to assign system role: %w", err)
		}

		// Add system user to org
		_, err = tx.ExecContext(ctx, `
			INSERT INTO organization_members (id, organization_id, user_id, role, invited_by, invited_at, joined_at, status)
			VALUES (?, ?, ?, 'owner', ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'active')
		`, orgMemberID, defaultOrgID, systemUserID, systemUserID)
		if err != nil {
			return fmt.Errorf("failed to add org member: %w", err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		slog.Info("Default organization and system user initialized")
	}

	existingServer, err := s.deps.ServerService.GetServerByHostname(hostname)

	if err == nil && existingServer != nil {
		slog.Info("Control plane server already initialized", "hostname", hostname, "server_id", existingServer.ID())
		return nil
	}

	var orgID string
	err = s.deps.DB.MainDB().DB().QueryRowContext(ctx, "SELECT id FROM organizations LIMIT 1").Scan(&orgID)
	if err != nil {
		return fmt.Errorf("no organization found in database: %w", err)
	}

	defaultOrgID := uuid.MustParse(orgID)

	cpuCores := runtime.NumCPU()
	memoryMB := getSystemMemoryMB()
	diskGB := getDiskSpaceMB() / 1024

	osInfo := runtime.GOOS
	osVersion := runtime.Version()

	ipv4Address := getLocalIPv4()
	ipv6Address := getLocalIPv6()

	createdServer, err := s.deps.ServerService.CreateServer(
		"Control Plane - "+hostname,
		hostname,
		ipv4Address,
		ipv6Address,
		s.config.Server.Port,
		servers.ServerTypeControlPlane,
		defaultOrgID,
	)
	if err != nil {
		return fmt.Errorf("failed to create control plane server: %w", err)
	}

	createdServer.UpdateDescription("Default control plane server for mikrocloud")
	createdServer.UpdateSpecs(&cpuCores, &memoryMB, &diskGB, &osInfo, &osVersion)
	createdServer.AddTag("control-plane")
	createdServer.AddTag("default")

	if err := s.deps.ServerService.UpdateServer(createdServer); err != nil {
		slog.Warn("Failed to update control plane server specs", "error", err)
	}

	slog.Info("Control plane server initialized",
		"server_id", createdServer.ID(),
		"hostname", hostname,
		"ipv4_address", ipv4Address,
		"ipv6_address", ipv6Address,
		"cpu_cores", cpuCores,
		"memory_mb", memoryMB,
	)

	return nil
}

func getLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}

func getLocalIPv6() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() == nil && ipnet.IP.To16() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func getSystemMemoryMB() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Sys / 1024 / 1024)
}

func getDiskSpaceMB() int {
	return 100000
}
