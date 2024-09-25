package dto

import (
	"gorm.io/gorm"
	"time"
)

type UserOverviewDTO struct {
	ID        uint      `json:"id"`
	StudentID string    `json:"student_id"`
	Username  string    `json:"username"`
	Role      uint      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

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

type AdminUserListRequestDTO struct {
	ID        *uint   `form:"id"`
	StudentID *string `form:"student_id"`
	Username  *string `form:"username"`
	Role      *uint   `form:"role"`
	Page      int     `form:"page" binding:"required" required_err:"page is Required"`
	PageSize  int     `form:"page_size" binding:"required" required_err:"page_size is Required"`
}

type AdminUserListResponseDTO struct {
	Users      []UserOverviewDTO `json:"users"`
	Pagination PaginationInfo    `json:"pagination"`
}
