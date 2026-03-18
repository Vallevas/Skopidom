// Package storage provides file storage implementations for the inventory system.
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
// Two-layer validation is applied:
//  1. File extension must be in the allowed list.
//  2. Actual content type is detected via magic bytes — a PHP script renamed
//     to .jpg is rejected even though the extension is valid.
func (s *LocalStorage) Save(
	_ context.Context,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("storage: file extension %q not allowed", ext)
	}

	// Read the first 512 bytes — http.DetectContentType inspects magic bytes,
	// not the filename, so it catches renamed malicious files.
	headerBuf := make([]byte, 512)
	bytesRead, err := file.Read(headerBuf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("storage: read file header: %w", err)
	}

	contentType := http.DetectContentType(headerBuf[:bytesRead])
	if !allowedContentTypes[contentType] {
		return "", fmt.Errorf("storage: content type %q not allowed", contentType)
	}

	// Rewind to the beginning before copying the full file to disk.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("storage: rewind file: %w", err)
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

// allowedExtensions is the whitelist of accepted image file extensions.
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// allowedContentTypes is the whitelist of MIME types detected via magic bytes.
// Note: http.DetectContentType does not distinguish .jpeg from .jpg —
// both return "image/jpeg". WebP requires Go 1.20+.
var allowedContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

