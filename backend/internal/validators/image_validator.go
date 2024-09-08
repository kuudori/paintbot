package validators

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

func IsJPG(file *multipart.FileHeader) bool {

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}
