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
SET email = $1,
    handle = $2,
    first_name = $3,
    last_name = $4,
    bio = $5,
    updated_at = $6
WHERE id = $7
RETURNING *;
-- name: UpdateUserAvatarByID :one
UPDATE users
SET avatar_url = $1
WHERE id = $2
RETURNING *;