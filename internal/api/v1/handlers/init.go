package handlers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
)

type DBInterface interface {
	CreateUserStandard(ctx context.Context, params database.CreateUserStandardParams) (database.User, error)
	GetRoleIDByName(ctx context.Context, name string) (uuid.UUID, error)
	CreateUserRoles(ctx context.Context, params database.CreateUserRolesParams) (database.UserRole, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUserByID(ctx context.Context, arg database.UpdateUserByIDParams) (database.User, error)
	UpdateUserAvatarByID(ctx context.Context, arg database.UpdateUserAvatarByIDParams) (database.User, error)
	
	CreateServer(ctx context.Context, arg database.CreateServerParams) (database.Server, error)
	CreateUserServer(ctx context.Context, arg database.CreateUserServerParams) (database.UserServer, error)
	UpdateServerMemberCount(ctx context.Context, arg database.UpdateServerMemberCountParams) (database.UpdateServerMemberCountRow, error)
	UpdateServerBannerByID(ctx context.Context, arg database.UpdateServerBannerByIDParams) (database.UpdateServerBannerByIDRow, error)
	UpdateServerIconByID(ctx context.Context, arg database.UpdateServerIconByIDParams) (database.UpdateServerIconByIDRow, error)
	UpdateServerByID(ctx context.Context, arg database.UpdateServerByIDParams) (database.Server, error)
	DeleteServer(ctx context.Context, id uuid.UUID) error
	
	GetUserServers(ctx context.Context, userID uuid.UUID) ([]database.GetUserServersRow, error)
	GetUserServer(ctx context.Context, arg database.GetUserServerParams) (database.UserServer, error)
	GetOneServerByID(ctx context.Context, id uuid.UUID) (database.Server, error)
	GetOneServerByCode(ctx context.Context, inviteCode string) (database.Server, error)
	GetRecentServers(ctx context.Context) ([]database.GetRecentServersRow, error)
	DeleteUserServer(ctx context.Context, arg database.DeleteUserServerParams) error

	
	CreateTextChannel(ctx context.Context, arg database.CreateTextChannelParams) (database.TextChannel, error)
	GetServerTextChannels(ctx context.Context, serverID uuid.UUID) ([]database.TextChannel, error)
	GetLanguageIDByName(ctx context.Context, language string) (uuid.UUID, error)
	
	CreateTextMessage(ctx context.Context, arg database.CreateTextMessageParams) (database.TextMessage, error)
	GetChannelTextMessages(ctx context.Context, channelID uuid.UUID) ([]database.GetChannelTextMessagesRow, error)

	CreateVoiceChannel(ctx context.Context, arg database.CreateVoiceChannelParams) (database.VoiceChannel, error)
	GetServerVoiceChannels(ctx context.Context, serverID uuid.UUID) ([]database.GetServerVoiceChannelsRow, error)
	LeaveVoiceChannelByUser(ctx context.Context, userID uuid.UUID) error
}

type Handlers struct {
	DB  DBInterface
	RDB *redis.RedisClient
	JWT string
	S3  *s3.Client
}

func NewHandlers(db DBInterface, rdb *redis.RedisClient, jwt string, s3 *s3.Client) *Handlers {
	return &Handlers{
		DB:  db,
		RDB: rdb,
		JWT: jwt,
		S3:  s3,
	}
}
