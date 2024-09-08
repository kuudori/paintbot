package constants

const (
	FileNamesQuery = `
        SELECT file_url FROM files 
        WHERE chat_id = $1 
        ORDER BY id desc
        OFFSET $2 LIMIT 10
    `
	FileByIdQuery   = `SELECT * FROM files WHERE id = $1`
	InsertFileQuery = `INSERT INTO files (chat_id, file_url, file_path) values ($1, $2, $3) RETURNING id`
	DeleteFileQuery = `DELETE FROM files WHERE id = $1`
)
