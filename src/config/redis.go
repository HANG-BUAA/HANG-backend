package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() (*redis.Client, error) {
	rdClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("db.redis.host") + ":" + viper.GetString("db.redis.port"),
		Password: "",
		DB:       0,
	})

	_, err := rdClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdClient, nil
}
