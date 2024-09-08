package storage

import (
	queries "PaintBackend/constants"
	database "PaintBackend/internal/storage/models"
	"context"
	"database/sql"
)

type PostgresFileRepository struct {
	db *sql.DB
}

func NewPostgresFileRepository(db *sql.DB) database.FileRepository {
	return &PostgresFileRepository{db: db}
}

func (p *PostgresFileRepository) Create(ctx context.Context, file *database.DbFile) error {
	return p.db.QueryRowContext(ctx, queries.InsertFileQuery, file.ChatID, file.FileUrl, file.FilePath).Scan(&file.ID)
}

func (p *PostgresFileRepository) GetFilesURLByChatID(ctx context.Context, chatID int64, offset int) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, queries.FileNamesQuery, chatID, (offset-1)*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fileNames := make([]string, 0, 10) // query limit = 10
	for rows.Next() {
		var fileName string
		if err := rows.Scan(&fileName); err != nil {
			return nil, err
		}
		fileNames = append(fileNames, fileName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fileNames, nil
}

func (p *PostgresFileRepository) GetFileById(ctx context.Context, id int64) (*database.DbFile, error) {
	var file database.DbFile

	err := p.db.QueryRowContext(ctx, queries.FileByIdQuery, id).Scan(&file.ID, &file.ChatID, &file.FileUrl, &file.FilePath)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (p *PostgresFileRepository) Delete(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, queries.DeleteFileQuery, id)
	return err
}
