package api

import (
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERR_CODE_ADD_USER   = 10011
	ERR_CODE_LOGIN      = 10012
	ERR_CODE_SEND_EMAIL = 10013
)

type UserApi struct {
	BaseApi
	Service *service.UserService
}

func NewUserApi() UserApi {
	return UserApi{
		BaseApi: NewBaseApi(),
		Service: service.NewUserService(),
	}
}

func (m UserApi) Login(c *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginDTO}).GetError(); err != nil {
		return
	}

	iUser, token, err := m.Service.Login(iUserLoginDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Status: http.StatusUnauthorized,
			Code:   ERR_CODE_LOGIN,
			Msg:    err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: gin.H{
			"token": token,
			"user":  iUser,
		},
	})
}

func (m UserApi) SendEmail(c *gin.Context) {
	var iUserSendEmailDTO dto.UserSendEmailDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserSendEmailDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.SendEmail(iUserSendEmailDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_SEND_EMAIL,
			Msg:  err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: gin.H{
			"status": "Send email successfully",
		},
	})
}

func (m UserApi) Register(c *gin.Context) {
	var iUserRegisterDTO dto.UserRegisterDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserRegisterDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.Register(&iUserRegisterDTO)

	if err != nil {
		m.ServerFail(ResponseJson{
			Code: ERR_CODE_ADD_USER,
			Msg:  err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: iUserRegisterDTO,
	})
}
