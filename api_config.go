package main

import (
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
)

type APIConfig struct {
	Port      string
	DB        *database.Queries
	RDB       *redis.RedisClient 
	JwtSecret string
}
