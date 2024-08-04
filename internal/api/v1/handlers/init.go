package handlers

import (
	"context"

	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
)

type DBInterface interface {
	CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error)
}

type Handlers struct {
	DB DBInterface
}

func NewHandlers(db DBInterface) *Handlers {
	return &Handlers{
		DB: db,
	}
}
