package dto

import (
    "gorm.io/gorm"
    "time"
)

type CommentCreateRequestDTO struct {
    UserID      uint
    PostID      uint   `json:"post_id" form:"post_id" binding:"required" required_err:"post_id is Required"`
    ReplyTo     *uint  `json:"reply_to" form:"reply_to" binding:"required" required_err:"reply_to is Required"`
    Content     string `json:"content" form:"content" binding:"required" required_err:"content is Required"`
    IsAnonymous *bool  `json:"is_anonymous" form:"is_anonymous" binding:"required" required_err:"is_anonymous is Required"`
}

type CommentCreateResponseDTO struct {
    ID          uint           `json:"id"`
    PostID      uint           `json:"post_id"`
    ReplyTo     uint           `json:"reply_to"`
    Content     string         `json:"content"`
    IsAnonymous bool           `json:"is_anonymous"`
    DisplayName string         `json:"display_name"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type CommentListRequestDTO struct {
    PostID uint `json:"post_id" form:"post_id" binding:"required" required_err:"post_id is Required"`
}
