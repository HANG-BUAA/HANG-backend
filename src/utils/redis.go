package utils

import (
    "HANG-backend/src/global"
    "context"
    "time"
)

func SetRedis(key string, value any, expiration time.Duration) error {
    return global.RedisClient.Set(context.Background(), key, value, expiration).Err()
}

func GetRedis(key string) (any, error) {
    return global.RedisClient.Get(context.Background(), key).Result()
}

func DelRedis(key ...string) error  {
    return global.RedisClient.Del(context.Background(), key...).Err()
}
