package log

import (
	"HANG-backend/src/global"
	"encoding/json"
	"github.com/streadway/amqp"
)

type Level int
type Type int

const (
	DEBUG_LEVEL Level = iota + 1
	INFO_LEVEL
	WARNING_LEVEL
	ERROR_LEVEL
	CRITICAL_LEVEL
)

const (
	ACCESS_TYPE Type = iota + 1
	APPLICATION_TYPE
	ERROR_TYPE
	SECURITY_TYPE
)

var levelToStringMap = map[Level]string{
	DEBUG_LEVEL:    "debug",
	INFO_LEVEL:     "info",
	WARNING_LEVEL:  "warning",
	ERROR_LEVEL:    "error",
	CRITICAL_LEVEL: "critical",
}

type BaseLog struct {
	Type      Type  `json:"type"`
	Level     Level `json:"level"`
	Source    int   `json:"source"`
	Timestamp int64 `json:"timestamp"`
}

type AccessLog struct {
	BaseLog
	Request struct {
		Method      string            `json:"method"`
		URL         string            `json:"url"`
		Headers     map[string]string `json:"headers"`
		ClientIP    string            `json:"client_ip"`
		QueryParams map[string]string `json:"query_params"`
	} `json:"request"`
	Response struct {
		StatusCode int `json:"status_code"`
	} `json:"response"`
	User struct {
		ID   uint `json:"user_id,omitempty"`
		Role uint `json:"role,omitempty"`
	} `json:"user,omitempty"`
	ExecutionTime int64  `json:"execution_time"` // 请求处理时间（毫秒）
	RequestID     string `json:"request_id"`     // 请求唯一标识
}

func PublishLog(level Level, log any) {
	go func() {
		body, err := json.Marshal(log)
		if err != nil {
			// todo
			return
		}

		// 发送
		err = global.RabbitMqChannel.Publish(
			"log",
			levelToStringMap[level],
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			// todo
		}
	}()
}
