package dto

import (
	"gorm.io/gorm"
	"time"
)

type UserLoginRequestDTO struct {
	Username  string `json:"username" form:"username"`
	StudentID string `json:"student_id" form:"student_id"`
	Password  string `json:"password" form:"password" binding:"required" required_err:"password is Required"`
}

type UserLoginResponseDTO struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	StudentID string         `json:"student_id"`
	Role      uint           `json:"role"`
	Token     string         `json:"token"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type UserRegisterRequestDTO struct {
	StudentID        string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
	Password         string `json:"password" form:"password" binding:"required" required_err:"password is Required"`
	VerificationCode string `json:"verification_code" form:"verification_code" binding:"required" required_err:"verification_code is Required"`
}

type UserRegisterResponseDTO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	Role      uint   `json:"role"`
}

type UserSendEmailRequestDTO struct {
	StudentID string `json:"student_id" form:"student_id" binding:"required" required_err:"student_id is Required"`
}

type UserUpdateAvatarRequestDTO struct {
	ID  uint
	Url string
}

type UserUpdateAvatarResponseDTO struct {
	Url string `json:"url"`
}
