package dto

import (
	"gorm.io/gorm"
	"time"
)

type PostCreateRequestDTO struct {
	UserID      uint
	Title       string `json:"title" form:"title" binding:"required" required_err:"title is Required"`
	Content     string `json:"content" form:"content" binding:"required" required_err:"content is Required"`
	IsAnonymous *bool  `json:"is_anonymous" form:"is_anonymous" binding:"required" required_err:"is_anonymous is Required"`
}

type PostCreateResponseDTO struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	IsAnonymous bool           `json:"is_anonymous"`
	LikeNum     int            `json:"like_num"`
	CollectNum  int            `json:"collect_num"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type PostLikeRequestDTO struct {
	PostID uint `uri:"post_id" binding:"required" required_err:"post_id is Required"`
	UserID uint
}

type PostCollectRequestDTO struct {
	PostID uint `uri:"post_id" binding:"required" required_err:"post_id is Required"`
	UserID uint
}
