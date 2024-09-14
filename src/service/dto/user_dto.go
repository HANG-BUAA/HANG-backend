package dto

import (
	"HANG-backend/src/model"
	"gorm.io/gorm"
)

// UserLoginDTO 用户注册 DTO
type UserLoginDTO struct {
	gorm.Model
	Username  string `json:"username" form:"username" binding:"required" message:"用户名错误" required_err:"用户名不可缺省"`
	Password  string `json:"password,omitempty" form:"password" binding:"required" message:"密码不能为空"`
	StudentID string `json:"student_id"`
	Token     string `json:"token"`
}

// UserRegisterDTO 添加用户 DTO
type UserRegisterDTO struct {
	ID               uint
	Username         string `json:"username" form:"username" binding:"required" message:"user_name is Required"`
	Password         string `json:"password,omitempty" form:"password" binding:"required" message:"password is Required"`
	VerificationCode string `json:"code,omitempty" form:"code" binding:"required" message:"code is Required"`
}

// ConvertToModel DTO 转换成 Model
func (m *UserRegisterDTO) ConvertToModel(iUser *model.User) {
	iUser.UserName = m.Username
	iUser.Password = m.Password
	iUser.StudentID = m.Username
}

type UserSendEmailDTO struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
}
