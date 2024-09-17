package api

import (
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ERR_CODE_ADD_USER    = 10011
	ERR_CODE_LOGIN       = 10012
	ERR_CODE_SEND_EMAIL  = 10013
	ERR_CODE_UPLOAD_FILE = 10014
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
	var iUserLoginRequestDTO dto.UserLoginRequestDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginRequestDTO}).GetError(); err != nil {
		return
	}

	iUserLoginResponseDTO, err := m.Service.Login(&iUserLoginRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Status: http.StatusUnauthorized,
			Code:   ERR_CODE_LOGIN,
			Msg:    err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *iUserLoginResponseDTO,
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
	var iUserSendEmailRequestDTO dto.UserSendEmailRequestDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserSendEmailRequestDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.SendEmail(&iUserSendEmailRequestDTO)
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
	var iUserRegisterRequestDTO dto.UserRegisterRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserRegisterRequestDTO}).GetError(); err != nil {
		return
	}

	iUserRegisterResponseDTO, err := m.Service.Register(&iUserRegisterRequestDTO)

	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_ADD_USER,
			Msg:  err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *iUserRegisterResponseDTO,
	})
}

func (m UserApi) UploadAvatar(c *gin.Context) {
	// 这个接口实际上没有 json 数据要传入，只有一个文件要从表单里传输，在下面的逻辑里实现
	// 此处还要 BuildRequest 的原因是把 c(*gin.Context) 绑定到 UserApi 上
	if err := m.BuildRequest(BuildRequestOption{Ctx: c}).GetError(); err != nil {
		return
	}

	id, _ := c.Get("id") // 这里经过中间件的处理，一定有 id
	file, err := c.FormFile("avatar")
	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_UPLOAD_FILE,
			Msg:  "error when fetching avatar",
		})
		return
	}
	path, err := utils.UploadFile(file, "user_avatars")
	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_UPLOAD_FILE,
			Msg:  err.Error(),
		})
		return
	}
	iUserUpdateAvatarResponseDTO, err := m.Service.UpdateAvatar(&dto.UserUpdateAvatarRequestDTO{
		ID:  id.(uint),
		Url: path,
	})
	if err != nil {
		m.Fail(ResponseJson{
			Code: ERR_CODE_UPLOAD_FILE,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *iUserUpdateAvatarResponseDTO,
	})
}
