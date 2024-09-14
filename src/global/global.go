package global

import (
	"HANG-backend/src/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger      *zap.SugaredLogger
	DB          *gorm.DB
	RedisClient *config.RedisClient
)
