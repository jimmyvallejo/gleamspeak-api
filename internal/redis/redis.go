package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewClient() (*RedisClient, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	addr := os.Getenv("REDIS_PORT_TEST")

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &RedisClient{rdb: rdb}, nil
}

func (r *RedisClient) GetJSON(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func (r *RedisClient) SetJson(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, json, expiration).Err()
}
