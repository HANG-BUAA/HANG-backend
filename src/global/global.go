package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	RDB         *gorm.DB
	RedisClient *redis.Client
)

const (
	OptimisticLockMaxRetries = 3
)
