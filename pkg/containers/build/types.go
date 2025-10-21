package build

import (
	"fmt"
)

type BuildpackType string

const (
	Nixpacks       BuildpackType = "nixpacks"
	Static         BuildpackType = "static"
	DockerfileType BuildpackType = "dockerfile"
	DockerCompose  BuildpackType = "docker-compose"
)

type BuildRequest struct {
	ID            string
	GitRepo       string
	GitBranch     string
	ContextRoot   string
	BuildpackType BuildpackType
	Environment   map[string]string
	ImageTag      string

	// Buildpack-specific configurations
	NixpacksConfig      *NixpacksConfig      `json:"nixpacks_config,omitempty"`
	StaticConfig        *StaticConfig        `json:"static_config,omitempty"`
	ContainerfileConfig *ContainerfileConfig `json:"dockerfile_config,omitempty"`
	ComposeConfig       *ComposeConfig       `json:"compose_config,omitempty"`

	// Optional callback for streaming logs in real-time
	LogCallback func(log string) `json:"-"`
}

// Result of an image build
type BuildResult struct {
	Success   bool
	ImageTag  string
	BuildLogs string
	Error     string
}

// Abstraction over all buildpacks
type BuildpackConfig interface {
	GetBuildCommands() []string
	Validate() error
}

type NixpacksConfig struct {
	StartCommand string            `json:"start_command,omitempty"`
	BuildCommand string            `json:"build_command,omitempty"`
	Variables    map[string]string `json:"variables,omitempty"`
}

func (n *NixpacksConfig) GetBuildCommands() []string {
	commands := []string{
		"nixpacks build . --name app",
	}

	if n.BuildCommand != "" {
		commands = append(commands, n.BuildCommand)
	}

	return commands
}

func (n *NixpacksConfig) Validate() error {
	return nil // Nixpacks handles most validation
}

type StaticConfig struct {
	OutputDir   string `json:"output_dir,omitempty"`
	NginxConfig string `json:"nginx_config,omitempty"`
	IS_SPA      string `json:"is_spa,omitempty"`
}

func (s *StaticConfig) GetBuildCommands() []string {
	outputDir := s.OutputDir
	if outputDir == "" {
		outputDir = "dist"
	}

	commands := []string{
		"echo 'Building with Nixpacks...'",
		fmt.Sprintf("docker build -t ${IMAGE_TAG} -f- . <<'EOF'\nFROM nginx:alpine\nCOPY %s /usr/share/nginx/html\n%sEXPOSE 80\nCMD [\"nginx\", \"-g\", \"daemon off;\"]\nEOF", outputDir, s.getNginxConfigDockerfile()),
	}

	return commands
}

func (s *StaticConfig) getNginxConfigDockerfile() string {
	if s.NginxConfig != "" {
		return fmt.Sprintf("COPY %s /etc/nginx/conf.d/default.conf\n", s.NginxConfig)
	}
	return ""
}

func (s *StaticConfig) Validate() error {
	if s.OutputDir == "" {
		s.OutputDir = "dist"
	}
	return nil
}

type ContainerfileConfig struct {
	ContainerfilePath string            `json:"dockerfile_path,omitempty"`
	BuildArgs         map[string]string `json:"build_args,omitempty"`
	Target            string            `json:"target,omitempty"`
}

func (c *ContainerfileConfig) GetBuildCommands() []string {
	containerfilePath := c.ContainerfilePath
	if containerfilePath == "" {
		containerfilePath = "Dockerfile"
	}

	buildCmd := fmt.Sprintf("docker build -t ${IMAGE_TAG} -f %s", containerfilePath)

	// Add build args
	for key, value := range c.BuildArgs {
		buildCmd += fmt.Sprintf(" --build-arg %s=%s", key, value)
	}

	// Add target if specified
	if c.Target != "" {
		buildCmd += fmt.Sprintf(" --target %s", c.Target)
	}

	buildCmd += " ."

	return []string{buildCmd}
}

func (d *ContainerfileConfig) Validate() error {
	if d.ContainerfilePath == "" {
		d.ContainerfilePath = "Dockerfile"
	}
	return nil
}

type ComposeConfig struct {
	ComposeFile string `json:"compose_file,omitempty"`
	Service     string `json:"service,omitempty"`
}

func (c *ComposeConfig) GetBuildCommands() []string {
	composeFile := c.ComposeFile
	if composeFile == "" {
		composeFile = "docker-compose.yml"
	}

	buildCmd := fmt.Sprintf("docker compose -f %s build", composeFile)

	if c.Service != "" {
		buildCmd += fmt.Sprintf(" %s", c.Service)
	}

	// Tag the built image
	tagCmd := ""
	if c.Service != "" {
		tagCmd = fmt.Sprintf("docker tag $(docker compose -f %s images -q %s) ${IMAGE_TAG}", composeFile, c.Service)
	}

	commands := []string{buildCmd}
	if tagCmd != "" {
		commands = append(commands, tagCmd)
	}

	return commands
}

func (c *ComposeConfig) Validate() error {
	if c.ComposeFile == "" {
		c.ComposeFile = "docker-compose.yml"
	}
	return nil
}
