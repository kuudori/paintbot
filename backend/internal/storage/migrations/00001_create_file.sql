-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
id SERIAL PRIMARY KEY,
chat_id BIGINT,
file_url VARCHAR(255),
file_path VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files;
-- +goose StatementEnd
