package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     AppConfig.DBHost + ":6379",
		Password: "",
		DB:       0,
	})
}

