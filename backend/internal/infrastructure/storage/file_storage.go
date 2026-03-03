// Package storage provides file storage implementations for the inventory system.
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileStorage defines the contract for storing and deleting uploaded files.
type FileStorage interface {
	// Save stores the uploaded file and returns its public URL path.
	Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)

	// Delete removes a previously stored file by its URL path.
	Delete(ctx context.Context, urlPath string) error
}

// LocalStorage implements FileStorage using the local filesystem.
type LocalStorage struct {
	// baseDir is the absolute directory path where files are stored.
	baseDir string
	// baseURL is the URL prefix served by the HTTP server for static files.
	baseURL string
}

// NewLocalStorage constructs a LocalStorage, creating baseDir if it does not exist.
func NewLocalStorage(baseDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("storage: create base dir: %w", err)
	}
	return &LocalStorage{
		baseDir: baseDir,
		baseURL: strings.TrimRight(baseURL, "/"),
	}, nil
}

// Save writes the uploaded file to disk under a timestamped unique filename.
func (s *LocalStorage) Save(
	_ context.Context,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {
	ext := filepath.Ext(header.Filename)
	if !isAllowedExtension(ext) {
		return "", fmt.Errorf("storage: file type %q not allowed", ext)
	}

	// Build a collision-resistant filename using nanosecond timestamp.
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(s.baseDir, filename)

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("storage: create file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		// Attempt cleanup on write failure.
		_ = os.Remove(destPath)
		return "", fmt.Errorf("storage: write file: %w", err)
	}

	// Return a URL path the frontend can use to fetch the image.
	return fmt.Sprintf("%s/%s", s.baseURL, filename), nil
}

// Delete removes the file referenced by the given URL path from disk.
func (s *LocalStorage) Delete(_ context.Context, urlPath string) error {
	if urlPath == "" {
		return nil
	}

	// Strip the base URL prefix to obtain the filename.
	filename := strings.TrimPrefix(urlPath, s.baseURL+"/")
	if filename == "" || strings.Contains(filename, "/") {
		// Ignore empty paths or paths that traverse directories.
		return nil
	}

	destPath := filepath.Join(s.baseDir, filename)
	if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: delete file: %w", err)
	}
	return nil
}

// allowedExtensions is the set of image extensions the system accepts.
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// isAllowedExtension reports whether the extension is an accepted image format.
func isAllowedExtension(ext string) bool {
	return allowedExtensions[strings.ToLower(ext)]
}
