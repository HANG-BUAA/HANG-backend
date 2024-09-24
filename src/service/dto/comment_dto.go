package dto

import (
	"gorm.io/gorm"
	"time"
)

type CommentOverviewDTO struct {
	ID                 uint           `json:"id"`
	PostID             uint           `json:"post_id"`
	ReplyToID          uint           `json:"reply_to_id"`           // 回复的评论的id，为0时说明是评论楼主
	ReplyToName        string         `json:"reply_to_name"`         // 回复的评论的主人的名字，要特判为0的情况（楼主），同样也要筛
	IsReplyToAnonymous bool           `json:"is_reply_to_anonymous"` // 回复的帖子/评论是否是匿名的
	UserID             uint           `json:"user_id,omitempty"`
	Username           string         `json:"user_name,omitempty"`   // 该字段要过滤
	UserAvatar         string         `json:"user_avatar,omitempty"` // 该字段要过滤
	Content            string         `json:"content"`
	IsAnonymous        bool           `json:"is_anonymous"`
	LikeNum            int            `json:"like_num"`
	HasLiked           bool           `json:"has_liked"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at,omitempty"`
}

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
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type CommentLikeRequestDTO struct {
	CommentID uint `uri:"comment_id" binding:"required" required_err:"comment_id is Required"`
	UserID    uint
}

type CommentListRequestDTO struct {
	UserID   uint
	Page     int  `json:"page" form:"page" binding:"required,min=1" required_err:"page is Required" min_err:"page could not be lower than 1"`
	PageSize int  `json:"page_size" form:"page_size" binding:"required,min=10,max=100" required_err:"page_size is Required" min_err:"page_size could not be lower than 10" max_err:"page_size could not be greater than 100"`
	PostID   uint `json:"post_id" form:"post_id" binding:"required" required_err:"post_id is Required"`
}

type CommentListResponseDTO struct {
	Comments   []CommentOverviewDTO `json:"comments"`
	Pagination PaginationInfo       `json:"pagination"`
}
