package api

import (
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

// @Summary 用户登录
// @Description 登录接口，返回用户信息与token
// @Tags 公共接口
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} LoginResponse "登录成功结果"
// @Failure 400 {object} object "登录失败"
// @Router /api/v1/public/login [post]
func (m UserApi) Login(c *gin.Context) {
	var iUserLoginDTO dto.UserLoginDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.Login(&iUserLoginDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Status: http.StatusUnauthorized,
			Code:   ERR_CODE_LOGIN,
			Msg:    err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: iUserLoginDTO,
	})
}

// @Summary 发送验证码
// @Description 验证码存在 redis 里，默认配置5min，对学号拼接 @buaa.edu.cn 发送
// @Tags 公共接口
// @Param student_id body string true "学号"
/*
@Success 200 {object} json {
  "data": {
    "status": "Send email successfully"
  }
}
*/
// @Failure 500 {object} object "服务器端验证码发送失败，可能是邮箱不对或服务问题"
// @Router /api/v1/public/send-email [post]
func (m UserApi) SendEmail(c *gin.Context) {
	var iUserSendEmailDTO dto.UserSendEmailDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserSendEmailDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.SendEmail(iUserSendEmailDTO)
	if err != nil {
		m.ServerFail(ResponseJson{
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

// @Summary 用户注册
// @Description 注册时需要先发送验证码，这里传入的 username 其实就是学号，创建账号成功后默认的初始 username 也即学号
// @Tags 公共接口
// @Param username body string true "用户名（学号）"
// @Param password body string true "密码"
// @Param code body string true "发送的验证码（有效期内）"
// @Success 200 {object} RegisterResponse "注册成功结果"
// @Failure 400 {object} object "注册失败"
// @Router /api/v1/public/register [post]
func (m UserApi) Register(c *gin.Context) {
	var iUserRegisterDTO dto.UserRegisterDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserRegisterDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.Register(&iUserRegisterDTO)

	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_ADD_USER,
			Msg:  err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: iUserRegisterDTO,
	})
}

type LoginResponse struct {
	Data struct {
		ID        int        `json:"ID"`
		CreatedAt time.Time  `json:"CreatedAt"`
		UpdatedAt time.Time  `json:"UpdatedAt"`
		DeletedAt *time.Time `json:"DeletedAt"`
		Username  string     `json:"username"`
		StudentID string     `json:"student_id"`
		Token     string     `json:"token"`
	} `json:"data"`
}

type RegisterResponse struct {
	Data struct {
		ID       int    `json:"ID"`
		Username string `json:"username"`
	} `json:"data"`
}
