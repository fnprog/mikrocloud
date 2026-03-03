package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// STRUCTURED OUTPUT - Main Result
// ============================================================================

type DeploymentResult struct {
	Strategy     string        `json:"strategy"`
	Discovery    DiscoveryInfo `json:"discovery"`
	ExposedPorts []PortInfo    `json:"exposedPorts"`
	Build        BuildInfo     `json:"build"`
	Runtime      RuntimeInfo   `json:"runtime"`
	Highlights   []Highlight   `json:"highlights"`
}

type DiscoveryInfo struct {
	FoundFiles []string `json:"foundFiles"`
	Runtime    *string  `json:"runtime"`
	Framework  *string  `json:"framework"`
	IsMonorepo bool     `json:"isMonorepo"`
	Confidence float64  `json:"confidence"`
}

type PortInfo struct {
	Port      uint16 `json:"port"`
	Source    string `json:"source"`
	Protocol  string `json:"protocol"`
	IsPrimary bool   `json:"isPrimary"`
}

type BuildInfo struct {
	Command           string   `json:"command"`
	Context           string   `json:"context"`
	Dockerfile        *string  `json:"dockerfile,omitempty"`
	BuildArgs         []string `json:"buildArgs"`
	EstimatedDuration string   `json:"estimatedDuration"`
}

type RuntimeInfo struct {
	StartCommand   string   `json:"startCommand"`
	EnvVars        []string `json:"envVars"`
	HealthCheck    *string  `json:"healthCheck,omitempty"`
	MemoryEstimate *string  `json:"memoryEstimate,omitempty"`
}

type Highlight struct {
	Level    string  `json:"level"`
	Category string  `json:"category"`
	Message  string  `json:"message"`
	Action   *string `json:"action,omitempty"`
}

// ============================================================================
// DOCKERFILE STRUCTURES
// ============================================================================

type DockerAnalysis struct {
	BaseImage    *string
	ExposedPorts []uint16
	EnvVars      []string
	Cmd          *string
	BuildArgs    map[string]string
	MultiStage   bool
	StageCount   int
}

// ============================================================================
// COMPOSE STRUCTURES
// ============================================================================

type ComposeFile struct {
	Services map[string]ComposeService `yaml:"services"`
	Networks map[string]interface{}    `yaml:"networks"`
	Volumes  map[string]interface{}    `yaml:"volumes"`
}

type ComposeService struct {
	Image       string                 `yaml:"image"`
	Build       interface{}            `yaml:"build"`
	Ports       []interface{}          `yaml:"ports"`
	DependsOn   []string               `yaml:"depends_on"`
	Environment interface{}            `yaml:"environment"`
	HealthCheck map[string]interface{} `yaml:"healthcheck"`
}

// ============================================================================
// ANALYZER - Main Logic
// ============================================================================

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(path string) (*DeploymentResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path error: %w", err)
	}

	// Case 1: Single Dockerfile
	if !info.IsDir() && strings.HasPrefix(filepath.Base(path), "Dockerfile") {
		return a.analyzeDockerfileOnly(path)
	}

	// Case 2: Single compose file
	if !info.IsDir() && isComposeFile(path) {
		return a.analyzeComposeOnly(path)
	}

	// Case 3: Directory analysis
	if info.IsDir() {
		return a.analyzeDirectory(path)
	}

	return nil, fmt.Errorf("path must be Dockerfile, docker-compose.yml, or directory")
}

// ============================================================================
// CASE 1: Dockerfile Only
// ============================================================================

func (a *Analyzer) analyzeDockerfileOnly(path string) (*DeploymentResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	analysis := parseDockerfile(string(content))
	highlights := []Highlight{}

	// Extract ports
	var exposedPorts []PortInfo
	if len(analysis.ExposedPorts) == 0 {
		action := "Add EXPOSE <port> to specify which port your app listens on"
		highlights = append(highlights, Highlight{
			Level:    "warning",
			Category: "ports",
			Message:  "No EXPOSE instruction found in Dockerfile",
			Action:   &action,
		})

		exposedPorts = []PortInfo{{
			Port:      8080,
			Source:    "default",
			Protocol:  "http",
			IsPrimary: true,
		}}
	} else {
		for i, port := range analysis.ExposedPorts {
			exposedPorts = append(exposedPorts, PortInfo{
				Port:      port,
				Source:    "dockerfile",
				Protocol:  "http",
				IsPrimary: i == 0,
			})
		}
	}

	// Check for multi-stage
	if analysis.MultiStage {
		highlights = append(highlights, Highlight{
			Level:    "success",
			Category: "optimization",
			Message:  "Using multi-stage build - great for image size!",
		})
	}

	// Check base image
	if analysis.BaseImage != nil && strings.Contains(*analysis.BaseImage, "latest") {
		action := "Pin to specific version for reproducible builds"
		highlights = append(highlights, Highlight{
			Level:    "warning",
			Category: "stability",
			Message:  fmt.Sprintf("Using 'latest' tag: %s", *analysis.BaseImage),
			Action:   &action,
		})
	}

	runtime := detectRuntimeFromBase(analysis.BaseImage)
	dockerfilePath := path
	memEstimate := "512MB"
	buildArgs := make([]string, 0, len(analysis.BuildArgs))
	for k := range analysis.BuildArgs {
		buildArgs = append(buildArgs, k)
	}

	duration := "fast"
	if analysis.MultiStage {
		duration = "medium"
	}

	cmd := "docker run app"
	if analysis.Cmd != nil {
		cmd = *analysis.Cmd
	}

	return &DeploymentResult{
		Strategy: "dockerfile",
		Discovery: DiscoveryInfo{
			FoundFiles: []string{"Dockerfile"},
			Runtime:    runtime,
			Framework:  nil,
			IsMonorepo: false,
			Confidence: 1.0,
		},
		ExposedPorts: exposedPorts,
		Build: BuildInfo{
			Command:           "docker build -t app .",
			Context:           ".",
			Dockerfile:        &dockerfilePath,
			BuildArgs:         buildArgs,
			EstimatedDuration: duration,
		},
		Runtime: RuntimeInfo{
			StartCommand:   cmd,
			EnvVars:        analysis.EnvVars,
			HealthCheck:    nil,
			MemoryEstimate: &memEstimate,
		},
		Highlights: highlights,
	}, nil
}

// ============================================================================
// CASE 2: Docker Compose Only
// ============================================================================

func (a *Analyzer) analyzeComposeOnly(path string) (*DeploymentResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %w", err)
	}

	var compose ComposeFile
	if err := yaml.Unmarshal(content, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse compose file: %w", err)
	}

	highlights := []Highlight{}
	var exposedPorts []PortInfo

	// Find entry service
	var entryService *string
	for name, service := range compose.Services {
		if (name == "web" || name == "app") && len(service.Ports) > 0 {
			entryService = &name
			for i, portDef := range service.Ports {
				if port := extractPortFromCompose(portDef); port > 0 {
					exposedPorts = append(exposedPorts, PortInfo{
						Port:      port,
						Source:    "compose",
						Protocol:  "http",
						IsPrimary: i == 0,
					})
				}
			}
			break
		}
	}

	// If no web/app service, find first with ports
	if entryService == nil {
		for name, service := range compose.Services {
			if len(service.Ports) > 0 {
				entryService = &name
				for i, portDef := range service.Ports {
					if port := extractPortFromCompose(portDef); port > 0 {
						exposedPorts = append(exposedPorts, PortInfo{
							Port:      port,
							Source:    "compose",
							Protocol:  "http",
							IsPrimary: i == 0,
						})
					}
				}
				break
			}
		}
	}

	if entryService != nil {
		highlights = append(highlights, Highlight{
			Level:    "info",
			Category: "services",
			Message:  fmt.Sprintf("Entry service detected: '%s'", *entryService),
		})
	} else {
		action := "Ensure at least one service exposes ports for external access"
		highlights = append(highlights, Highlight{
			Level:    "warning",
			Category: "services",
			Message:  "No service with exposed ports found",
			Action:   &action,
		})

		exposedPorts = append(exposedPorts, PortInfo{
			Port:      8080,
			Source:    "default",
			Protocol:  "http",
			IsPrimary: true,
		})
	}

	// Multi-service setup
	if len(compose.Services) > 1 {
		highlights = append(highlights, Highlight{
			Level:    "info",
			Category: "architecture",
			Message:  fmt.Sprintf("Multi-service setup with %d services", len(compose.Services)),
		})
	}

	filename := filepath.Base(path)
	memEstimate := "1GB"

	return &DeploymentResult{
		Strategy: "dockerCompose",
		Discovery: DiscoveryInfo{
			FoundFiles: []string{filename},
			Runtime:    nil,
			Framework:  nil,
			IsMonorepo: false,
			Confidence: 1.0,
		},
		ExposedPorts: exposedPorts,
		Build: BuildInfo{
			Command:           "docker-compose build",
			Context:           ".",
			Dockerfile:        nil,
			BuildArgs:         []string{},
			EstimatedDuration: "medium",
		},
		Runtime: RuntimeInfo{
			StartCommand:   "docker-compose up -d",
			EnvVars:        []string{},
			HealthCheck:    nil,
			MemoryEstimate: &memEstimate,
		},
		Highlights: highlights,
	}, nil
}

// ============================================================================
// CASE 3: Directory Analysis
// ============================================================================

func (a *Analyzer) analyzeDirectory(path string) (*DeploymentResult, error) {
	foundFiles := []string{}
	var dockerfilePath, composePath *string

	// Scan directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		foundFiles = append(foundFiles, name)

		if strings.HasPrefix(name, "Dockerfile") {
			fullPath := filepath.Join(path, name)
			dockerfilePath = &fullPath
		}

		if isComposeFile(filepath.Join(path, name)) {
			fullPath := filepath.Join(path, name)
			composePath = &fullPath
		}
	}

	// Priority 1: Compose
	if composePath != nil {
		return a.analyzeComposeOnly(*composePath)
	}

	// Priority 2: Dockerfile
	if dockerfilePath != nil {
		return a.analyzeDockerfileOnly(*dockerfilePath)
	}

	// Priority 3: Static site
	if isStaticSite(path) {
		return a.analyzeStaticSite(path, foundFiles)
	}

	// Priority 4: Nixpacks
	if runtime := detectRuntime(path); runtime != nil {
		return a.analyzeWithNixpacks(path, runtime, foundFiles)
	}

	return nil, fmt.Errorf("could not determine deployment strategy")
}

// ============================================================================
// Static Site Detection
// ============================================================================

func isStaticSite(path string) bool {
	hasHTML := false
	hasRuntime := false

	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".html") || name == "index.html" {
			hasHTML = true
		}

		if name == "package.json" || name == "Cargo.toml" || name == "go.mod" || name == "requirements.txt" {
			hasRuntime = true
		}
	}

	return hasHTML && !hasRuntime
}

func (a *Analyzer) analyzeStaticSite(path string, foundFiles []string) (*DeploymentResult, error) {
	runtime := "static"
	healthCheck := "/"
	memEstimate := "64MB"

	return &DeploymentResult{
		Strategy: "staticSite",
		Discovery: DiscoveryInfo{
			FoundFiles: foundFiles,
			Runtime:    &runtime,
			Framework:  nil,
			IsMonorepo: false,
			Confidence: 0.95,
		},
		ExposedPorts: []PortInfo{{
			Port:      80,
			Source:    "default",
			Protocol:  "http",
			IsPrimary: true,
		}},
		Build: BuildInfo{
			Command:           "# Copy files to nginx",
			Context:           ".",
			Dockerfile:        nil,
			BuildArgs:         []string{},
			EstimatedDuration: "fast",
		},
		Runtime: RuntimeInfo{
			StartCommand:   "nginx -g 'daemon off;'",
			EnvVars:        []string{},
			HealthCheck:    &healthCheck,
			MemoryEstimate: &memEstimate,
		},
		Highlights: []Highlight{
			{
				Level:    "info",
				Category: "deployment",
				Message:  "Static site detected - will deploy with nginx",
			},
			{
				Level:    "success",
				Category: "performance",
				Message:  "Lightweight deployment - minimal resource usage",
			},
		},
	}, nil
}

// ============================================================================
// Runtime Detection for Nixpacks
// ============================================================================

type RuntimeInfo struct {
	Runtime      string
	Framework    *string
	DefaultPort  uint16
	StartCommand string
}

func detectRuntime(path string) *RuntimeInfo {
	// Node.js
	if _, err := os.Stat(filepath.Join(path, "package.json")); err == nil {
		framework := detectNodeFramework(path)
		port := detectNodePort(path)
		if port == 0 {
			port = 3000
		}
		return &RuntimeInfo{
			Runtime:      "node",
			Framework:    framework,
			DefaultPort:  port,
			StartCommand: "npm start",
		}
	}

	// Rust
	if _, err := os.Stat(filepath.Join(path, "Cargo.toml")); err == nil {
		port := detectRustPort(path)
		if port == 0 {
			port = 8080
		}
		return &RuntimeInfo{
			Runtime:      "rust",
			Framework:    nil,
			DefaultPort:  port,
			StartCommand: "cargo run --release",
		}
	}

	// Python
	if _, err := os.Stat(filepath.Join(path, "requirements.txt")); err == nil {
		port := detectPythonPort(path)
		if port == 0 {
			port = 8000
		}
		return &RuntimeInfo{
			Runtime:      "python",
			Framework:    nil,
			DefaultPort:  port,
			StartCommand: "python main.py",
		}
	}

	// Go
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		port := detectGoPort(path)
		if port == 0 {
			port = 8080
		}
		return &RuntimeInfo{
			Runtime:      "go",
			Framework:    nil,
			DefaultPort:  port,
			StartCommand: "go run .",
		}
	}

	return nil
}

func (a *Analyzer) analyzeWithNixpacks(path string, runtime *RuntimeInfo, foundFiles []string) (*DeploymentResult, error) {
	highlights := []Highlight{
		{
			Level:    "info",
			Category: "buildpack",
			Message:  fmt.Sprintf("Using Nixpacks for %s runtime", runtime.Runtime),
		},
		{
			Level:    "info",
			Category: "ports",
			Message:  fmt.Sprintf("Detected port %d from source code", runtime.DefaultPort),
			Action:   strPtr("Verify this matches your application's listening port"),
		},
	}

	memEstimate := "512MB"

	return &DeploymentResult{
		Strategy: "nixpacks",
		Discovery: DiscoveryInfo{
			FoundFiles: foundFiles,
			Runtime:    &runtime.Runtime,
			Framework:  runtime.Framework,
			IsMonorepo: false,
			Confidence: 0.85,
		},
		ExposedPorts: []PortInfo{{
			Port:      runtime.DefaultPort,
			Source:    "sourceCode",
			Protocol:  "http",
			IsPrimary: true,
		}},
		Build: BuildInfo{
			Command:           fmt.Sprintf("nixpacks build %s", path),
			Context:           path,
			Dockerfile:        nil,
			BuildArgs:         []string{},
			EstimatedDuration: "medium",
		},
		Runtime: RuntimeInfo{
			StartCommand:   runtime.StartCommand,
			EnvVars:        []string{},
			HealthCheck:    nil,
			MemoryEstimate: &memEstimate,
		},
		Highlights: highlights,
	}, nil
}

// ============================================================================
// Port Detection from Source Code
// ============================================================================

func detectNodePort(path string) uint16 {
	return scanForPort(path, []string{".js", ".ts", ".mjs"}, []string{
		`\.listen\((\d+)`,
		`port:\s*(\d+)`,
		`PORT\s*=\s*(\d+)`,
	})
}

func detectRustPort(path string) uint16 {
	return scanForPort(path, []string{".rs"}, []string{
		`\.bind\(".*:(\d+)"\)`,
		`port:\s*(\d+)`,
	})
}

func detectPythonPort(path string) uint16 {
	return scanForPort(path, []string{".py"}, []string{
		`port=(\d+)`,
		`\.run\(.*port=(\d+)`,
	})
}

func detectGoPort(path string) uint16 {
	return scanForPort(path, []string{".go"}, []string{
		`ListenAndServe\(":(\d+)"`,
		`port.*=.*(\d+)`,
	})
}

func scanForPort(root string, extensions []string, patterns []string) uint16 {
	var port uint16

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || port > 0 {
			return nil
		}

		// Check extension
		ext := filepath.Ext(path)
		validExt := false
		for _, e := range extensions {
			if ext == e {
				validExt = true
				break
			}
		}
		if !validExt {
			return nil
		}

		// Read file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Try patterns
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			if matches := re.FindStringSubmatch(string(content)); len(matches) > 1 {
				var p int
				fmt.Sscanf(matches[1], "%d", &p)
				if p > 1024 && p < 65535 {
					port = uint16(p)
					return filepath.SkipAll
				}
			}
		}

		return nil
	})

	return port
}

func detectNodeFramework(path string) *string {
	content, err := os.ReadFile(filepath.Join(path, "package.json"))
	if err != nil {
		return nil
	}

	contentStr := string(content)
	if strings.Contains(contentStr, `"next"`) {
		return strPtr("nextjs")
	}
	if strings.Contains(contentStr, `"express"`) {
		return strPtr("express")
	}

	return nil
}

// ============================================================================
// Dockerfile Parser
// ============================================================================

func parseDockerfile(content string) *DockerAnalysis {
	analysis := &DockerAnalysis{
		ExposedPorts: []uint16{},
		EnvVars:      []string{},
		BuildArgs:    make(map[string]string),
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	stageCount := 0

	exposeRe := regexp.MustCompile(`^EXPOSE\s+(\d+)`)
	fromRe := regexp.MustCompile(`^FROM\s+(.+)`)
	envRe := regexp.MustCompile(`^ENV\s+(\w+)`)
	cmdRe := regexp.MustCompile(`^CMD\s+(.+)`)
	argRe := regexp.MustCompile(`^ARG\s+(\w+)(?:=(.*))?`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// FROM
		if matches := fromRe.FindStringSubmatch(line); len(matches) > 1 {
			stageCount++
			if analysis.BaseImage == nil {
				analysis.BaseImage = &matches[1]
			}
		}

		// EXPOSE
		if matches := exposeRe.FindStringSubmatch(line); len(matches) > 1 {
			var port uint16
			fmt.Sscanf(matches[1], "%d", &port)
			analysis.ExposedPorts = append(analysis.ExposedPorts, port)
		}

		// ENV
		if matches := envRe.FindStringSubmatch(line); len(matches) > 1 {
			analysis.EnvVars = append(analysis.EnvVars, matches[1])
		}

		// CMD
		if matches := cmdRe.FindStringSubmatch(line); len(matches) > 1 {
			analysis.Cmd = &matches[1]
		}

		// ARG
		if matches := argRe.FindStringSubmatch(line); len(matches) > 1 {
			key := matches[1]
			val := ""
			if len(matches) > 2 {
				val = matches[2]
			}
			analysis.BuildArgs[key] = val
		}
	}

	analysis.MultiStage = stageCount > 1
	analysis.StageCount = stageCount

	return analysis
}

// ============================================================================
// Helpers
// ============================================================================

func isComposeFile(path string) bool {
	base := filepath.Base(path)
	return base == "docker-compose.yml" || base == "docker-compose.yaml" ||
		base == "compose.yml" || base == "compose.yaml"
}

func extractPortFromCompose(portDef interface{}) uint16 {
	switch v := portDef.(type) {
	case string:
		parts := strings.Split(v, ":")
		var port uint16
		fmt.Sscanf(parts[0], "%d", &port)
		return port
	case int:
		return uint16(v)
	}
	return 0
}

func detectRuntimeFromBase(baseImage *string) *string {
	if baseImage == nil {
		return nil
	}

	img := *baseImage
	if strings.HasPrefix(img, "node") {
		return strPtr("node")
	}
	if strings.HasPrefix(img, "rust") {
		return strPtr("rust")
	}
	if strings.HasPrefix(img, "python") {
		return strPtr("python")
	}
	if strings.HasPrefix(img, "golang") {
		return strPtr("go")
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}

// ============================================================================
// MAIN
// ============================================================================

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	analyzer := NewAnalyzer()
	result, err := analyzer.Analyze(path)
	if err != nil {
		errorResponse := map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		}
		jsonData, _ := json.MarshalIndent(errorResponse, "", "  ")
		fmt.Println(string(jsonData))
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON encoding error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
