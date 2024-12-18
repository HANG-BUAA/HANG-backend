package log

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"encoding/json"
	"github.com/streadway/amqp"
	"time"
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

type ApplicationLogStatus bool

const (
	Success ApplicationLogStatus = true
	Failure ApplicationLogStatus = false
)

type Event string

type ApplicationLog struct {
	BaseLog
	User struct {
		ID        uint   `json:"id"`
		StudentID string `json:"student_id"`
	} `json:"user"`
	Application struct {
		Event    Event                `json:"event"`
		EntityID any                  `json:"entity_id"`
		Status   ApplicationLogStatus `json:"status"`
		Error    struct {
			SysMessage string `json:"sys_message"`
			Info       string `json:"info"`
		} `json:"error"`
	} `json:"application"`
	AdditionalInfo *string `json:"additional_info"`
}

func publishLog(level Level, log any) {
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

func PublishAccessLog(
	duration time.Duration,
	method string,
	url string,
	headers map[string]string,
	clientIP string,
	queryParams map[string]string,
	statusCode int,
	user *model.User,
) {
	accessLog := AccessLog{
		BaseLog: BaseLog{
			Type:      ACCESS_TYPE,
			Level:     INFO_LEVEL,
			Source:    1, // 主后端
			Timestamp: time.Now().Unix(),
		},
		Request: struct {
			Method      string            `json:"method"`
			URL         string            `json:"url"`
			Headers     map[string]string `json:"headers"`
			ClientIP    string            `json:"client_ip"`
			QueryParams map[string]string `json:"query_params"`
		}{Method: method, URL: url, Headers: headers, ClientIP: clientIP, QueryParams: queryParams},
		Response: struct {
			StatusCode int `json:"status_code"`
		}{StatusCode: statusCode},
		ExecutionTime: duration.Milliseconds(),
		RequestID:     "123", // todo 增加请求号
	}

	if user != nil {
		accessLog.User.ID = user.ID
		accessLog.User.Role = user.Role
	}

	publishLog(DEBUG_LEVEL, accessLog)
}

func PublishApplicationLog(
	user *model.User,
	event Event,
	entityID any,
	status ApplicationLogStatus,
	ErrorSysMessage *string,
	ErrorInfo *string,
	additionalInfo *string,
) {
	accessLog := ApplicationLog{
		BaseLog: BaseLog{
			Type:      APPLICATION_TYPE,
			Level:     INFO_LEVEL,
			Source:    1,
			Timestamp: time.Now().Unix(),
		},
		Application: struct {
			Event    Event                `json:"event"`
			EntityID any                  `json:"entity_id"`
			Status   ApplicationLogStatus `json:"status"`
			Error    struct {
				SysMessage string `json:"sys_message"`
				Info       string `json:"info"`
			} `json:"error"`
		}{Event: event, EntityID: entityID, Status: status},
	}
	if user != nil {
		accessLog.User = struct {
			ID        uint   `json:"id"`
			StudentID string `json:"student_id"`
		}{ID: user.ID, StudentID: user.StudentID}
	}
	if ErrorSysMessage != nil {
		accessLog.Application.Error.SysMessage = *ErrorSysMessage
	}
	if ErrorInfo != nil {
		accessLog.Application.Error.Info = *ErrorInfo
	}
	if additionalInfo != nil {
		accessLog.AdditionalInfo = additionalInfo
	}
	publishLog(INFO_LEVEL, accessLog)
}
