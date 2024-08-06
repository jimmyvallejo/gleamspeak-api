-- name: GetRoleIDByName :one
SELECT id
FROM roles
WHERE name = $1;