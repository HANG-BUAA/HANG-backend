package utils

import (
	"HANG-backend/src/global"
	"encoding/json"
	"github.com/streadway/amqp"
)

type PostMessage struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PublishPostMessage(post PostMessage) error {
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// 发送消息  todo 好像有超时的设置，暂时没研究
	err = global.RabbitMqChannel.Publish(
		"search",
		"post",
		false, // todo 研究一下这个强制消息是什么
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
