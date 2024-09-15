package dto

import (
	"HANG-backend/src/utils"
)

//type UserLoginDTO struct {
//	gorm.Model
//	Username  string `json:"username" form:"username" binding:"required" message:"用户名错误" required_err:"用户名不可缺省"`
//	Password  string `json:"password,omitempty" form:"password" binding:"required" message:"密码不能为空"`
//	StudentID string `json:"student_id"`
//	Token     string `json:"token"`
//}
//
//type UserRegisterDTO struct {
//	ID               uint
//	Username         string `json:"username" form:"username" binding:"required" message:"user_name is Required"`
//	Password         string `json:"password,omitempty" form:"password" binding:"required" message:"password is Required"`
//	VerificationCode string `json:"code,omitempty" form:"code" binding:"required" message:"code is Required"`
//}
//
//type UserSendEmailDTO struct {
//	StudentID string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
//}

type UserLoginRequestDTO struct {
	Username  string `json:"username" form:"username"`
	StudentID string `json:"student_id" form:"student_id"`
	Password  string `json:"password" form:"password" binding:"required" required_err:"password is Required"`
}

type UserLoginResponseDTO struct {
	utils.CustomBaseModel
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	Token     string `json:"token"`
}

type UserRegisterRequestDTO struct {
	StudentID        string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
	Password         string `json:"password" form:"password" binding:"required" required_err:"password is Required"`
	VerificationCode string `json:"verification_code" form:"verification_code" binding:"required" required_err:"verificationCode is Required"`
}

type UserRegisterResponseDTO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
}

type UserSendEmailRequestDTO struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
}
