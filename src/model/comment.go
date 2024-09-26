package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement; not null"`
	PostID             uint           `gorm:"index; not null"`
	UserID             uint           `gorm:"index; not null"` // 发表评论的人的 id
	ReplyCommentID     uint           `gorm:"index; not null"` // 回复的评论的id，为0时表示一级评论（直接回复帖子）；否则为二级评论
	ReplyRootCommentID uint           `gorm:"index; not null"` // 回复的根一级评论的id，表示该二级评论最终是在哪个一级评论的下面。一级评论的该字段为0
	ReplyUserID        uint           `gorm:"index; not null"` // 回复的人的id
	Content            string         `gorm:"type:text; not null"`
	IsAnonymous        bool           `gorm:"index; not null"` // 这条评论是否是匿名的
	IsReplyAnonymous   bool           `gorm:"index; not null"` // 回复的 帖子/评论 是否是匿名的
	CreatedAt          time.Time      `gorm:"index"`
	UpdatedAt          time.Time      `gorm:"index"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
