package middleware

import (
	"HANG-backend/src/api"
	"HANG-backend/src/global"
	"github.com/gin-gonic/gin"
	"strconv"
)

const maxSize = 20

func paginationErr(c *gin.Context) {
	api.Fail(c, api.ResponseJson{
		Code: global.ERR_CODE_PAGINATION_PARAM,
		Msg:  "page related parameters error!",
	})
}

func CheckPaginationParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor := c.Query("cursor")
		c.Set("cursor", cursor) // 由于 cursor 的形式不唯一，所以直接以字符串形式存储，到具体的 service 再转换

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize > maxSize {
			paginationErr(c)
		}
		c.Set("page_size", pageSize)
		c.Next()
	}
}
