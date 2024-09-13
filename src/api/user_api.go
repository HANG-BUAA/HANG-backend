package api

import (
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
)

type UserApi struct{}

func NewUserApi() UserApi {
	return UserApi{}
}

func (m UserApi) Login(ctx *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO
	errs := ctx.ShouldBind(&iUserLoginDTO)
	if errs != nil {
		Fail(ctx, ResponseJson{
			Msg: errs.Error(),
		})
		return
	}
	OK(ctx, ResponseJson{
		Data: iUserLoginDTO,
	})
}
