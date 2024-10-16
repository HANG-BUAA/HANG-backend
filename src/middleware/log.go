package middleware

import (
	"HANG-backend/src/log"
	"HANG-backend/src/model"
	"github.com/gin-gonic/gin"
	"time"
)

func AccessLogRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// 执行时间
		duration := time.Since(startTime)

		var user *model.User
		if u, exist := c.Get("user"); exist {
			user = u.(*model.User)
		}

		log.PublishAccessLog(
			duration,
			c.Request.Method,
			c.Request.RequestURI,
			extractHeaders(c.Request.Header),
			c.ClientIP(),
			extractHeaders(c.Request.URL.Query()),
			c.Writer.Status(),
			user,
		)

		//accessLog := log.AccessLog{
		//	BaseLog: log.BaseLog{
		//		Type:      log.ACCESS_TYPE,
		//		Level:     log.INFO_LEVEL,
		//		Source:    1, // 主后端
		//		Timestamp: time.Now().Unix(),
		//	},
		//	Request: struct {
		//		Method      string            `json:"method"`
		//		URL         string            `json:"url"`
		//		Headers     map[string]string `json:"headers"`
		//		ClientIP    string            `json:"client_ip"`
		//		QueryParams map[string]string `json:"query_params"`
		//	}{
		//		Method:      c.Request.Method,
		//		URL:         c.Request.RequestURI,
		//		Headers:     extractHeaders(c.Request.Header),
		//		ClientIP:    c.ClientIP(),
		//		QueryParams: extractQueryParams(c.Request.URL.Query()),
		//	},
		//	Response: struct {
		//		StatusCode int `json:"status_code"`
		//	}{StatusCode: c.Writer.Status()},
		//	ExecutionTime: duration.Milliseconds(),
		//	RequestID:     "123", // todo 增加请求编号
		//}
		//
		//if user, exists := c.Get("user"); exists {
		//	accessLog.User.ID = user.(*model.User).ID
		//	accessLog.User.Role = user.(*model.User).Role
		//}
		//
		//log.PublishLog(log.INFO_LEVEL, accessLog)
	}
}

// 提取 HTTP 头信息为 map[string]string
func extractHeaders(headers map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0] // 只记录第一个值
		}
	}
	return result
}

// 提取查询参数为 map[string]string
func extractQueryParams(query map[string][]string) map[string]string {
	result := make(map[string]string)
	for key, values := range query {
		if len(values) > 0 {
			result[key] = values[0] // 只记录第一个值
		}
	}
	return result
}
