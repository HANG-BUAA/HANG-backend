package api

import (
	"HANG-backend/src/utils"
	"github.com/gin-gonic/gin"
)

type LogApi struct {
	BaseApi
}

func NewLogApi() LogApi {
	return LogApi{
		BaseApi: BaseApi{},
	}
}

func (m LogApi) ListKeywords(c *gin.Context) {
	m.Ctx = c
	keywords, err := utils.ListAllKeywords()
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: keywords,
	})
}
