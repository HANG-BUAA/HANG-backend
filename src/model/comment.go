package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement; not null"`
	PostID             uint           `gorm:"index; not null"`        //所属的帖子的id
	UserID             uint           `gorm:"index; not null"`        // 发表评论的人的 id
	UserName           string         `gorm:"type:varchar(20);index"` // 发评论的人的名字。如果是匿名评论，该名即为前端展示的名字；否则该字段无效
	UserAvatar         string         `gorm:"type:varchar(255)"`      // 发评论的人的头像，设定同上
	ReplyCommentID     uint           `gorm:"index; not null"`        // 回复的评论的id，为0时表示一级评论（直接回复帖子）；否则为二级评论
	ReplyRootCommentID uint           `gorm:"index; not null"`        // 回复的根一级评论的id，表示该二级评论最终是在哪个一级评论的下面。一级评论的该字段为0 	// 回复的人的id
	ReplyUserName      string         `gorm:"type:varchar(20);index"` // 回复的人的名字，当回复的人所发的是匿名贴/评论时该字段有效
	Content            string         `gorm:"type:text; not null"`    // 评论内容
	IsAnonymous        bool           `gorm:"index; not null"`        // 这条评论是否是匿名的
	IsReplyAnonymous   bool           `gorm:"index; not null"`        // 回复的 帖子/评论 是否是匿名的
	CreatedAt          time.Time      `gorm:"index"`
	UpdatedAt          time.Time      `gorm:"index"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type CommentLike struct {
	UserID    uint `gorm:"primaryKey;index"`
	CommentID uint `gorm:"primaryKey;index"`
}
