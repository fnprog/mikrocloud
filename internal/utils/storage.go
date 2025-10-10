package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/webp"
)

const (
	MaxAvatarSize     = 2 * 1024 * 1024
	AvatarStoragePath = "./storage/avatars"
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type FileValidationError struct {
	Message string
}

func (e *FileValidationError) Error() string {
	return e.Message
}

func ValidateImageFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > MaxAvatarSize {
		return &FileValidationError{Message: fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", MaxAvatarSize)}
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return &FileValidationError{Message: "file type not allowed, must be JPEG, PNG, GIF, or WebP"}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return &FileValidationError{Message: "invalid image file"}
	}

	validFormats := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"gif":  true,
		"webp": true,
	}

	if !validFormats[format] {
		return &FileValidationError{Message: "unsupported image format"}
	}

	return nil
}

func SaveAvatar(fileHeader *multipart.FileHeader, userID string) (string, error) {
	if err := ValidateImageFile(fileHeader); err != nil {
		return "", err
	}

	if err := os.MkdirAll(AvatarStoragePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}
	hashString := hex.EncodeToString(hash.Sum(nil))

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		contentType := fileHeader.Header.Get("Content-Type")
		switch contentType {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	filename := fmt.Sprintf("%s-%s%s", userID, hashString[:16], ext)
	filePath := filepath.Join(AvatarStoragePath, filename)

	file.Seek(0, 0)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fmt.Sprintf("/storage/avatars/%s", filename), nil
}

func DeleteAvatar(avatarURL string) error {
	if avatarURL == "" || !strings.HasPrefix(avatarURL, "/storage/avatars/") {
		return nil
	}

	filename := strings.TrimPrefix(avatarURL, "/storage/avatars/")
	filePath := filepath.Join(AvatarStoragePath, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete avatar file: %w", err)
	}

	return nil
}
