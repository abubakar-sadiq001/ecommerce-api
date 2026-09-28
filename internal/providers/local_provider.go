package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{
		basePath: basePath,
	}
}

// UploadFile saves an uploaded file to the local storage path and returns its URL.
func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	// 1. join the base path with the file path
	fullPath := filepath.Join(p.basePath, path)

	// 2. create a new directory with permission for the specific file using its path name
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	// 3. Open the file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// create destination
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// read from source to destination
	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil

}

func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}
