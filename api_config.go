package main

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jimmyvallejo/gleamspeak-api/internal/api/v1/handlers"
	"github.com/jimmyvallejo/gleamspeak-api/internal/database"
	"github.com/jimmyvallejo/gleamspeak-api/internal/redis"
)

type APIConfig struct {
	Port     string
	DB       *database.Queries
	RDB      *redis.RedisClient
	Handlers *handlers.Handlers
	JwtSecret string
	S3        *s3.Client
}
