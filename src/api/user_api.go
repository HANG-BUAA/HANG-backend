package api

import (
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	BaseApi
}

func NewUserApi() UserApi {
	return UserApi{
		BaseApi: NewBaseApi(),
	}
}

func (m UserApi) Login(c *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginDTO}).GetError(); err != nil {
		return
	}

	m.OK(ResponseJson{
		Data: iUserLoginDTO,
	})
}
