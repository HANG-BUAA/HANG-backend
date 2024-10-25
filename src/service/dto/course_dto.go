package dto

import (
	"gorm.io/gorm"
	"time"
)

type AdminCourseCreateRequestDTO struct {
	ID      string   `json:"id" form:"id" binding:"required" required_err:"id is Required"`
	Name    string   `json:"name" form:"name" binding:"required" required_err:"name is Required"`
	Credits *float32 `json:"credits" form:"credits"`
	Campus  *int     `json:"campus" form:"campus"`
	Tags    []uint   `json:"tags" form:"tags"`
}

type AdminCourseCreateResponseDTO struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Credits   *float32       `json:"credits"`
	Campus    *int           `json:"campus"`
	Tags      []uint         `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
