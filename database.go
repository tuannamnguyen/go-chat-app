package main

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type redisHandler struct {
	client *redis.Client
}

func newRedisHandler() *redisHandler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &redisHandler{client: redisClient}
}
