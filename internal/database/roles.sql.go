// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: roles.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getRoleIDByName = `-- name: GetRoleIDByName :one
SELECT id
FROM roles
WHERE name = $1
`

func (q *Queries) GetRoleIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getRoleIDByName, name)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
