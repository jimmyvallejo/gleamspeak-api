-- name: CreateServer :one
INSERT INTO servers (
        id,
        owner_id,
        server_name,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;