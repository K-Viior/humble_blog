package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   5,
	})
	ping := RDB.Ping(context.Background())
	if _, err := ping.Result(); err != nil {
		zap.L().Error(err.Error())
	}
	return
}
