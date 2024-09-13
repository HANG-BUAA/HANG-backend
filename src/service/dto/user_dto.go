package dto

import "HANG-backend/src/model"

// UserLoginDTO 用户注册 DTO
type UserLoginDTO struct {
	Username string `json:"username" binding:"required" message:"用户名错误" required_err:"用户名不可缺省"`
	Password string `json:"password" binding:"required" message:"密码不能为空"`
}

// UserRegisterDTO 添加用户 DTO
type UserRegisterDTO struct {
	ID       uint
	Username string `json:"username" form:"username" binding:"required" message:"用户名不能为空"`
	Password string `json:"password,omitempty" form:"password," binding:"required" message:"密码不能为空"`
}

// ConvertToModel DTO 转换成 Model
func (m *UserRegisterDTO) ConvertToModel(iUser *model.User) {
	iUser.UserName = m.Username
	iUser.Password = m.Password
}
