-- name: CreateUserServer :one
INSERT INTO user_servers (user_id, server_id)
VALUES ($1, $2)
RETURNING *;

