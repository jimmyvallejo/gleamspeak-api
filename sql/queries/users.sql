-- name: CreateUserStandard :one
INSERT INTO users (
        id,
        email,
        password,
        handle,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;

-- name: GetUserByID :one 
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one 
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserByID :one
UPDATE users
SET email = $1, handle = $2, updated_at = $3
WHERE id = $4
RETURNING *;