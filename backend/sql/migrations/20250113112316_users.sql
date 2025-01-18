-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	unique_name varchar(30) NOT NULL,
	display_name varchar(30) NOT NULL,
	email varchar(255) NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users ADD UNIQUE INDEX (unique_name, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
