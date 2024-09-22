package model

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          uint           `gorm:"primaryKey;autoIncrement;not null"`
	UserID      uint           `gorm:"not null"`
	Title       string         `gorm:"type: varchar(100) not null"`
	Content     string         `gorm:"type: text not null"`
	IsAnonymous bool           `gorm:"index;not null"`
	CreatedAt   time.Time      `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type PostLike struct {
	User   User
	UserID uint `gorm:"primaryKey;index"`
	Post   Post
	PostID uint `gorm:"primaryKey;index"`
}

type PostCollect struct {
	User   User
	UserID uint `gorm:"primaryKey;index"`
	Post   Post
	PostID uint `gorm:"primaryKey;index"`
}
