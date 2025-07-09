package config

import (
	"github.com/knadh/koanf"
	"github.com/redis/go-redis/v9"
)

func NewRedis(k *koanf.Koanf) *redis.Client {
	host := k.String("redis.host")
	database := k.Int("redis.database")
	password := k.String("redis.password")

	var client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       database,
	})

	return client
}
