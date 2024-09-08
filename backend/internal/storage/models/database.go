package database

import (
	"context"
)

type DbFile struct {
	ID       int64  `db:"id"`
	ChatID   int64  `db:"chat_id"`
	FileUrl  string `db:"file_url"`
	FilePath string `db:"file_path"`
}

type FileRepository interface {
	Create(ctx context.Context, file *DbFile) error
	Delete(ctx context.Context, id int64) error
	GetFilesURLByChatID(ctx context.Context, chatID int64, offset int) ([]string, error)
	GetFileById(ctx context.Context, id int64) (*DbFile, error)
}
