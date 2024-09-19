package middleware

import (
	"HANG-backend/src/api"
	"HANG-backend/src/global"
	"HANG-backend/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	TOKEN_NAME   = "Authorization"
	TOKEN_PREFIX = "Bearer "
)

func tokenErr(c *gin.Context) {
	api.Fail(c, api.ResponseJson{
		Status: http.StatusUnauthorized,
		Code:   global.ERR_CODE_INVALID_TOKEN,
		Msg:    "Invalid Token",
	})
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader(TOKEN_NAME)

		// token 不存在，直接返回
		if token == "" || !strings.HasPrefix(token, TOKEN_PREFIX) {
			tokenErr(c)
			return
		}

		// token 无法解析，直接返回
		token = token[len(TOKEN_PREFIX):]
		iJwtCustomClaims, err := utils.ParseToken(token)
		nUserId := iJwtCustomClaims.ID
		if err != nil || nUserId == 0 {
			tokenErr(c)
			return
		}

		// todo 判断 token 是否过期与续期

		// 把 id 存到 context 中
		c.Set("id", iJwtCustomClaims.ID)
		c.Next()
	}
}
