package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type DBInterface interface {
	CreateUserStandard(ctx context.Context, params database.CreateUserStandardParams) (database.User, error)
	GetRoleIDByName(ctx context.Context, name string) (uuid.UUID, error)
	CreateUserRoles(ctx context.Context, params database.CreateUserRolesParams) (database.UserRole, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	UpdateUserByID(ctx context.Context, arg database.UpdateUserByIDParams) (database.User, error)
	CreateServer(ctx context.Context, arg database.CreateServerParams) (database.Server, error)
	CreateUserServer(ctx context.Context, arg database.CreateUserServerParams) (database.UserServer, error)
}

type Handlers struct {
	DB  DBInterface
	JWT string
}

func NewHandlers(db DBInterface, jwt string) *Handlers {
	return &Handlers{
		DB:  db,
		JWT: jwt,
	}
}
