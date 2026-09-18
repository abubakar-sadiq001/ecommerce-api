package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/abubakar-sadiq001/learning-go-shop/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

// UploadImage validates and uploads an image file to our local storage folder.
func (p *UploadService) UploadImage(productID uint, file *multipart.FileHeader) (string, error) {
	// check for valid extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	// get the filepath
	path := fmt.Sprintf("products/%d/%s", productID, file.Filename)

	// upload the file
	return p.provider.UploadFile(file, path)
}

func isValidImageExt(ext string) bool {
	validExtensions := []string{"jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}

	return false
}
