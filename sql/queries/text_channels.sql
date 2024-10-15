-- name: CreateTextChannel :one
INSERT INTO text_channels (
        id,
        owner_id,
        server_id,
        language_id,
        channel_name,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: DeleteTextChannel :exec
DELETE FROM text_channels
WHERE id = $1;

-- name: GetServerTextChannels :many
SELECT * FROM text_channels
WHERE server_id = $1;