-- name: CreateServer :one
INSERT INTO servers (
        id,
        owner_id,
        server_name,
        invite_code,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetOneServerByID :one
SELECT *
FROM servers
WHERE id = $1;
-- name: GetOneServerByCode :one
SELECT *
FROM servers
WHERE invite_code = $1;
-- name: GetRecentServers :many
SELECT s.id,
    s.server_name,
    s.description,
    s.icon_url,
    s.banner_url,
    s.member_count,
    s.created_at,
    s.updated_at,
    u.handle,
    u.avatar_url
FROM servers s
    INNER JOIN users u ON s.owner_id = u.id
WHERE s.is_public = TRUE
ORDER BY s.created_at DESC
LIMIT 10;
-- name: UpdateServerMemberCount :one
UPDATE servers
SET member_count = $2
WHERE id = $1
RETURNING id,
    server_name,
    member_count;
-- name: UpdateServerIconByID :one
UPDATE servers
SET icon_url = $2
WHERE id = $1
RETURNING id,
    server_name,
    icon_url;
-- name: UpdateServerBannerByID :one
UPDATE servers
SET banner_url = $2
WHERE id = $1
RETURNING id,
    server_name,
    banner_url;