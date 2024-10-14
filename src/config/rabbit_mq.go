package config

import (
	"HANG-backend/src/global"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func InitRabbitMq() (*amqp.Channel, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		viper.GetString("rabbit_mq.username"),
		viper.GetString("rabbit_mq.password"),
		viper.GetString("rabbit_mq.host"),
		viper.GetString("rabbit_mq.port"))
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	global.RabbitMqConnection = conn
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 声明 search 交换机
	err = ch.ExchangeDeclare("search",
		"direct",
		true,
		false,
		false,
		false, // todo 生产环境下防止数据丢失，要等待服务器确认
		nil,
	)
	if err != nil {
		return nil, err
	}

	// 声明 log 交换机
	err = ch.ExchangeDeclare(
		"log",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
