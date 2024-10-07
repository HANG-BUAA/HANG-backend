package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"HANG-backend/src/service"
	"HANG-backend/src/service/dto"
	"HANG-backend/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var iUserLoginRequestDTO dto.UserLoginRequestDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserLoginRequestDTO}).GetError(); err != nil {
		return
	}

	iUserLoginResponseDTO, err := m.Service.Login(&iUserLoginRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Status: http.StatusUnauthorized,
			Code:   global.ERR_CODE_LOGIN_FAILED,
			Msg:    err.Error(),
		})
		return
	}

	m.OK(ResponseJson{
		Data: *iUserLoginResponseDTO,
	})
}

func (m UserApi) SendEmail(c *gin.Context) {
	var iUserSendEmailRequestDTO dto.UserSendEmailRequestDTO

	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &iUserSendEmailRequestDTO}).GetError(); err != nil {
		return
	}

	err := m.Service.SendEmail(&iUserSendEmailRequestDTO)
	if err != nil {
		m.ServerFail(ResponseJson{
			Code: global.ERR_CODE_SEND_EMAIL,
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
	var userRegisterRequestDTO dto.UserRegisterRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &userRegisterRequestDTO}).GetError(); err != nil {
		return
	}

	iUserRegisterResponseDTO, err := m.Service.Register(&userRegisterRequestDTO)

	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_ADD_USER,
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

	user := c.MustGet("user").(*model.User)
	id := user.ID
	file, err := c.FormFile("avatar")
	if err != nil {
		m.Fail(ResponseJson{
			Code: global.ERR_CODE_UPLOAD_MISSING_FIELD,
			Msg:  "custom_error when fetching avatar",
		})
		return
	}
	path, err := utils.UploadFile(file, "user_avatars")
	if err != nil {
		m.ServerFail(ResponseJson{
			Code: global.ERR_CODE_UPLOAD_SERVER_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	iUserUpdateAvatarResponseDTO, err := m.Service.UpdateAvatar(&dto.UserUpdateAvatarRequestDTO{
		ID:  id,
		Url: path,
	})
	if err != nil {
		m.ServerFail(ResponseJson{
			Code: global.ERR_CODE_UPLOAD_SERVER_FAILED,
			Msg:  err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *iUserUpdateAvatarResponseDTO,
	})
}

func (m UserApi) AdminList(c *gin.Context) {
	var adminUserListRequestDTO dto.AdminUserListRequestDTO
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &adminUserListRequestDTO}).GetError(); err != nil {
		return
	}
	adminUserListResponseDTO, err := m.Service.AdminList(&adminUserListRequestDTO)
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: *adminUserListResponseDTO,
	})
}
