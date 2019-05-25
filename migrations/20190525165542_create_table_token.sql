-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS token (
	id BINARY(16),
	token BINARY(32) NOT NULL,
	created_at DATETIME NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id),
	FOREIGN KEY (id) REFERENCES user(id)
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token;
-- +goose StatementEnd
