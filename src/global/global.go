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
	PageSize int = 20 // 分页大小
)
