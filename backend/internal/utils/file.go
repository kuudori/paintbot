package utils

import (
	"PaintBackend/internal/config"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(uploadDir string, file *multipart.FileHeader) (string, string, error) {
	cfg := config.GetConfig()
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", fmt.Errorf("create dir error: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	imageURL := fmt.Sprintf("%s/%s/%s", cfg.BackendDomain, uploadDir, newFilename)
	fullFilePath := filepath.Join(uploadDir, newFilename)

	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("open file error: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullFilePath)
	if err != nil {
		return "", "", fmt.Errorf("create file error: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("save file error: %w", err)
	}

	return imageURL, fullFilePath, nil
}
