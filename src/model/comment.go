package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;not null"`
	PostID      uint   `gorm:"index"`
	ReplyTo     uint   `gorm:"index"` // 回复的评论的id，为0时说明直接回复帖子
	UserID      uint   `gorm:"index"`
	Username    string `gorm:"size:255;index"`
	Content     string `gorm:"type: text not null"`
	IsAnonymous bool   `gorm:"index;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CommentLike struct {
	UserID    uint `gorm:"primaryKey;index"`
	CommentID uint `gorm:"primaryKey;index"`
}
