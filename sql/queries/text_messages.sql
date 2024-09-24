-- name: CreateTextMessage :one
INSERT INTO text_messages (
        id,
        owner_id,
        channel_id,
        message,
        image,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: GetChannelTextMessages :many
SELECT t.id,
    t.owner_id,
    t.channel_id,
    t.message,
    t.image,
    t.created_at,
    t.updated_at,
    u.handle,
    u.avatar_url
FROM text_messages t
    INNER JOIN users u ON t.owner_id = u.id
WHERE t.channel_id = $1;
