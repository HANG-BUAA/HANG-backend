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
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			paginationErr(c)
		}
		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize > maxSize {
			paginationErr(c)
		}
		c.Set("page", page)
		c.Set("page_size", pageSize)
		c.Next()
	}
}
