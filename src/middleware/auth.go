package middleware

import (
	"HANG-backend/src/api"
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/permission"
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

func permissionErr(c *gin.Context) {
	api.Fail(c, api.ResponseJson{
		Status: http.StatusUnauthorized,
		Code:   global.ERR_CODE_PERMISSION_DENIED,
		Msg:    "Permission Denied",
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
		jwtCustomClaims, err := utils.ParseToken(token)
		userId := jwtCustomClaims.ID
		if err != nil || userId == 0 {
			tokenErr(c)
			return
		}

		// todo 判断 token 是否过期与续期

		// 把 id 存到 context 中
		c.Set("id", jwtCustomClaims.ID)
		c.Next()
	}
}

func Permission(permission permission.Permission) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user model.User
		id, _ := c.Get("id")
		if err := global.RDB.First(&user, id).Error; err != nil {
			tokenErr(c)
			return
		}
		// todo 根据 weight先判断一层
		// 判断该用户是否有相应权限
		var userPermission model.UserPermission
		if err := global.RDB.Model(&userPermission).
			Where("user_id = ? AND permission_id = ?", user.ID, uint(permission)).
			First(&userPermission).Error; err != nil {
			permissionErr(c)
			return
		}
		c.Next()
	}
}
