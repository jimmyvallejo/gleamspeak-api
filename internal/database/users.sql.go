// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUserStandard = `-- name: CreateUserStandard :one
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
RETURNING id, email, password, is_active, handle, first_name, last_name, bio, avatar_url, is_verified, created_at, updated_at
`

type CreateUserStandardParams struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	Password  sql.NullString `json:"password"`
	Handle    string         `json:"handle"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (q *Queries) CreateUserStandard(ctx context.Context, arg CreateUserStandardParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUserStandard,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.Handle,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Handle,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.AvatarUrl,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, is_active, handle, first_name, last_name, bio, avatar_url, is_verified, created_at, updated_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Handle,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.AvatarUrl,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, password, is_active, handle, first_name, last_name, bio, avatar_url, is_verified, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Handle,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.AvatarUrl,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserAvatarByID = `-- name: UpdateUserAvatarByID :one
UPDATE users
SET avatar_url = $1
WHERE id = $2
RETURNING id, email, password, is_active, handle, first_name, last_name, bio, avatar_url, is_verified, created_at, updated_at
`

type UpdateUserAvatarByIDParams struct {
	AvatarUrl sql.NullString `json:"avatar_url"`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateUserAvatarByID(ctx context.Context, arg UpdateUserAvatarByIDParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserAvatarByID, arg.AvatarUrl, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Handle,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.AvatarUrl,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserByID = `-- name: UpdateUserByID :one
UPDATE users
SET email = $1,
    handle = $2,
    first_name = $3,
    last_name = $4,
    bio = $5,
    updated_at = $6
WHERE id = $7
RETURNING id, email, password, is_active, handle, first_name, last_name, bio, avatar_url, is_verified, created_at, updated_at
`

type UpdateUserByIDParams struct {
	Email     string         `json:"email"`
	Handle    string         `json:"handle"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Bio       sql.NullString `json:"bio"`
	UpdatedAt time.Time      `json:"updated_at"`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateUserByID(ctx context.Context, arg UpdateUserByIDParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserByID,
		arg.Email,
		arg.Handle,
		arg.FirstName,
		arg.LastName,
		arg.Bio,
		arg.UpdatedAt,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.IsActive,
		&i.Handle,
		&i.FirstName,
		&i.LastName,
		&i.Bio,
		&i.AvatarUrl,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
