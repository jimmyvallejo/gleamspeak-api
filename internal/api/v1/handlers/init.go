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
	GetUserServers(ctx context.Context, userID uuid.UUID) ([]database.GetUserServersRow, error)
	GetUserServer(ctx context.Context, arg database.GetUserServerParams) (database.UserServer, error)
	DeleteUserServer(ctx context.Context, arg database.DeleteUserServerParams) error
	CreateTextChannel(ctx context.Context, arg database.CreateTextChannelParams) (database.TextChannel, error)
	GetServerTextChannels(ctx context.Context, serverID uuid.UUID) ([]database.TextChannel, error)
	GetLanguageIDByName(ctx context.Context, language string) (uuid.UUID, error)
	CreateTextMessage(ctx context.Context, arg database.CreateTextMessageParams) (database.TextMessage, error)
	GetChannelTextMessages(ctx context.Context, channelID uuid.UUID) ([]database.GetChannelTextMessagesRow, error)
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
