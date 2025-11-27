package redisdb

import (
	"context"

	"github.com/ShopOnGO/media-service/configs"
	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB(conf *configs.Config) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	// Пинг при старте, чтобы убедиться что база жива
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic("Не удалось подключиться к Redis: " + err.Error())
	}

	return &RedisDB{Client: client}
}
