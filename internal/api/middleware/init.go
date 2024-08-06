package middleware

import (
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
)

type Middleware struct {
	DB  *database.Queries
	RDB *redis.RedisClient
	JWT string
}

func NewMiddleware(db *database.Queries, rdb *redis.RedisClient, jwt string) *Middleware {
	return &Middleware{
		DB:  db,
		RDB: rdb,
		JWT: jwt,
	}
}
