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

func (m UserApi) UpdateUser(c *gin.Context) {
	type UpdateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	user := c.MustGet("user").(*model.User)
	// 解析请求体中的更新数据
	var req UpdateUserRequest
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &req}).GetError(); err != nil {
		return
	}

	// 更新字段（只修改传入的字段）
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		user.Password = req.Password // 注意：通常这里密码应该加密存储
	}

	// 保存更新后的用户数据
	if err := global.RDB.Save(user).Error; err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
	}
	m.OK(ResponseJson{
		Data: gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"student_id": user.StudentID,
			"avatar":     user.Avatar,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
			"deleted_at": user.DeletedAt,
		},
	})
}

func (m UserApi) ListNotification(c *gin.Context) {
	// 获取请求的查询参数
	var query struct {
		NotifierID uint   `form:"notifier_id"`
		Type       string `form:"type"`
		Page       int    `form:"page" binding:"required"`
		PageSize   int    `form:"page_size" binding:"required"`
	}
	if err := m.BuildRequest(BuildRequestOption{Ctx: c, DTO: &query}).GetError(); err != nil {
		return
	}
	user := c.MustGet("user").(*model.User)
	query.NotifierID = user.ID
	// 定义查询条件
	notifierID := query.NotifierID
	notificationType := query.Type
	page := query.Page
	pageSize := query.PageSize

	// 构建查询
	var notifications []model.Notification
	var totalCount int64

	// 使用 global.RDB 作为数据库实例进行查询
	db := global.RDB

	// 1. 基本查询条件：根据 NotifierID 筛选
	queryBuilder := db.Model(&model.Notification{}).Where("notifier_id = ?", notifierID)

	// 2. 如果传入了 Type，添加类型筛选
	if notificationType != "" {
		queryBuilder = queryBuilder.Where("type = ?", notificationType)
	}

	// 3. 获取总记录数
	err := queryBuilder.Count(&totalCount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notifications"})
		return
	}

	// 4. 执行分页查询
	err = queryBuilder.Offset((page - 1) * pageSize).Limit(pageSize).Find(&notifications).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	// 5. 返回响应
	m.OK(ResponseJson{
		Data: gin.H{
			"total_count":   totalCount,
			"page":          page,
			"page_size":     pageSize,
			"notifications": notifications,
		},
	})
}
