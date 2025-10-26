package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type BuildWorkspace struct {
	Path string
}

func CreateBuildWorkspace(buildID string, buildDir string) (*BuildWorkspace, error) {
	if buildDir == "" {
		return nil, fmt.Errorf("build directory cannot be empty")
	}

	if buildID == "" {
		return nil, fmt.Errorf("build ID cannot be empty")
	}

	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create build directory: %w", err)
	}

	workspacePath := filepath.Join(buildDir, buildID)

	if err := os.MkdirAll(workspacePath, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return &BuildWorkspace{Path: workspacePath}, nil
}

func (w *BuildWorkspace) SaveUploadedFile(filename string, content []byte) error {
	if w.Path == "" {
		return fmt.Errorf("workspace path is empty")
	}

	sanitizedFilename := filepath.Base(filename)
	filePath := filepath.Join(w.Path, sanitizedFilename)

	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (w *BuildWorkspace) SaveInlineDockerfile(content string) error {
	if w.Path == "" {
		return fmt.Errorf("workspace path is empty")
	}

	filePath := filepath.Join(w.Path, "Dockerfile")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to save Dockerfile: %w", err)
	}

	return nil
}

func (w *BuildWorkspace) SaveInlineCompose(content string) error {
	if w.Path == "" {
		return fmt.Errorf("workspace path is empty")
	}

	filePath := filepath.Join(w.Path, "docker-compose.yml")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to save docker-compose.yml: %w", err)
	}

	return nil
}

func (w *BuildWorkspace) ExtractZip(zipPath string) error {
	if w.Path == "" {
		return fmt.Errorf("workspace path is empty")
	}

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if err := w.extractZipFile(file); err != nil {
			return fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}
	}

	if err := os.Remove(zipPath); err != nil {
		return fmt.Errorf("failed to remove zip file: %w", err)
	}

	return nil
}

func (w *BuildWorkspace) extractZipFile(file *zip.File) error {
	filePath := filepath.Join(w.Path, file.Name)

	if !strings.HasPrefix(filePath, filepath.Clean(w.Path)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, file.Mode())
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

func (w *BuildWorkspace) Cleanup() error {
	if w.Path == "" {
		return nil
	}

	if err := os.RemoveAll(w.Path); err != nil {
		return fmt.Errorf("failed to cleanup workspace: %w", err)
	}

	return nil
}

func CleanupWorkspace(workspacePath string) error {
	if workspacePath == "" {
		return nil
	}

	if err := os.RemoveAll(workspacePath); err != nil {
		return fmt.Errorf("failed to cleanup workspace: %w", err)
	}

	return nil
}
