-- name: CreateUserServer :one
INSERT INTO user_servers (user_id, server_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserServers :many
SELECT s.id AS server_id,
    s.server_name,
    s.owner_id,
    s.description,
    s.icon_url,
    s.banner_url,
    s.is_public,
    s.member_count,
    s.server_level,
    s.max_members,
    s.invite_code,
    s.created_at AS server_created_at,
    s.updated_at AS server_updated_at
FROM user_servers us
    JOIN servers s ON us.server_id = s.id
WHERE us.user_id = $1
ORDER BY s.server_name ASC;

-- name: GetUserServer :one
SELECT * FROM user_servers
WHERE user_id = $1 AND server_id = $2;

-- name: DeleteUserServer :exec
DELETE FROM user_servers
WHERE user_id = $1 AND server_id = $2;