package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Logger             *zap.SugaredLogger
	RDB                *gorm.DB
	RedisClient        *redis.Client
	RabbitMqChannel    *amqp.Channel
	RabbitMqConnection *amqp.Connection
)

const (
	OptimisticLockMaxRetries = 3
)
