package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type redisHandler struct {
	client *redis.Client
}

func NewRedisHandler() *redisHandler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &redisHandler{client: redisClient}
}

func (r *redisHandler) setUserInfo(ctx context.Context, userID, userName string) error {
	err := r.client.Set(ctx, userID, userName, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting user info in redis: %v", err)
	}

	return nil
}

func (r *redisHandler) getUserInfo(ctx context.Context, userID string) string {
	name := r.client.Get(ctx, userID).Val()
	return name
}
