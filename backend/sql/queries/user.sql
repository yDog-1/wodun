-- name: ListUsers :many
SELECT
	id,
	unique_name,
	display_name,
	email
FROM users
ORDER BY id;

-- name: GetUser :one
SELECT
	id,
	unique_name,
	display_name,
	email
FROM users
WHERE unique_name = ?;

-- name: CreateUser :exec
INSERT INTO users (
	unique_name, display_name, email
) VALUES (
	?, ?, ?
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE unique_name = ?;

-- name: UpdateUser :exec
UPDATE users
SET
	unique_name = COALESCE(sqlc.narg('unique_name'), unique_name),
	display_name = COALESCE(sqlc.narg('display_name'), display_name),
	email = COALESCE(sqlc.narg('email'), email)
WHERE unique_name = sqlc.arg('current_unique_name');
