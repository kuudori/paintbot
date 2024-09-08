package models

import (
	"mime/multipart"
)

type File struct {
	ID       int64  `json:"id" binding:"required"`
	ChatID   int64  `json:"chatId" binding:"required"`
	FileName string `json:"fileName" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

type FileUpload struct {
	Image *multipart.FileHeader `form:"image" binding:"required" validate:"required"`
}
