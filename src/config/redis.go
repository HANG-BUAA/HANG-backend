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

//func (c RedisClient) Set(key string, value any, expiration time.Duration) error {
//	return rdClient.Set(context.Background(), key, value, expiration).Err()
//}
//
//func (c RedisClient) Get(key string) (any, error) {
//	return rdClient.Get(context.Background(), key).Result()
//}
//
//func (c RedisClient) Delete(key ...string) error {
//	return rdClient.Del(context.Background(), key...).Err()
//}
