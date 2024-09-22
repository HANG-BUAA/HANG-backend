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
	Content     string `gorm:"type: text not null"`
	IsAnonymous bool   `gorm:"index;not null"`
	DisplayName string `gorm:"index;size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *Comment) BeforeCreate(tx *gorm.DB) error {
	// 生成给前端展示的名字
	// 首先判断这条评论是不是楼主发的
	var iPost Post
	if err := tx.First(&iPost, m.PostID).Error; err != nil {
		return err
	}
	op := iPost.UserID == m.UserID
	if m.IsAnonymous {
		if op {
			m.DisplayName = "楼主"
		} else {
			m.DisplayName = "匿名用户"
		}
	} else {
		var iUser User
		err := tx.Where("id = ?", m.UserID).First(&iUser).Error
		if err != nil {
			return err
		}
		m.DisplayName = iUser.UserName
		if op {
			m.DisplayName += "(楼主)"
		}
	}
	return nil
}
