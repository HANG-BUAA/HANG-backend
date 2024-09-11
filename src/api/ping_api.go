package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type PingApi struct{}

func NewPingApi() PingApi {
    return PingApi{}
}

// @Summary 测试连接
// @Description ping通此接口，说明项目运行成功
// @Success 200 {string} string "连接成功"
// @router /ping [post]
func (m PingApi) Ping(context *gin.Context) {
    context.String(http.StatusOK, "pong")
}
