package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	client *redis.Client
}

func NewAuthRepository() *AuthRepository {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &AuthRepository{client: redisClient}
}

func (r *AuthRepository) SetUserInfo(ctx context.Context, userID string, userName string) error {
	err := r.client.Set(ctx, userID, userName, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting user info in redis: %v", err)
	}

	return nil
}

func (r *AuthRepository) GetUserInfo(ctx context.Context, userID string) string {
	name := r.client.Get(ctx, userID).Val()
	return name
}
