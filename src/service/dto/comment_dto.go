package dto

import (
	"gorm.io/gorm"
	"time"
)

type CommentAuthorDTO struct {
	UserID     uint   `json:"user_id,omitempty"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
}

type CommentOverviewDTO struct {
	ID                 uint             `json:"id"`
	PostID             uint             `json:"post_id"`
	Author             CommentAuthorDTO `json:"author"`
	ReplyCommentID     uint             `json:"reply_comment_id"`
	ReplyRootCommentID uint             `json:"reply_root_comment_id"`
	ReplyUserName      string           `json:"reply_user_name"`
	Content            string           `json:"content"`
	IsAnonymous        bool             `json:"is_anonymous"`
	IsReplyAnonymous   bool             `json:"is_reply_anonymous"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          gorm.DeletedAt   `json:"deleted_at"`
}

type CommentCreateRequestDTO struct {
	UserID         uint
	PostID         uint   `json:"post_id" form:"post_id" binding:"required" required_err:"post_id is Required"`
	ReplyCommentID uint   `json:"reply_comment_id" form:"reply_comment_id" binding:"required" required_err:"reply_comment_id is Required"`
	Content        string `json:"content" form:"content" binding:"required" required_err:"post_id is Required"`
	IsAnonymous    *bool  `json:"is_anonymous" form:"is_anonymous" binding:"required" required_err:"is_anonymous is Required"`
}

type CommentCreateResponseDTO CommentOverviewDTO
