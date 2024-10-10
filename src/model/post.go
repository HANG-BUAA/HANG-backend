package model

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;not null"`
	UserID         uint           `gorm:"not null"`
	Title          string         `gorm:"type: varchar(100) not null"`
	Content        string         `gorm:"type: text not null"`
	IsAnonymous    bool           `gorm:"index;not null"`
	LikeNum        int            `gorm:"default:0;index; not null"`
	CollectNum     int            `gorm:"default:0;index; not null"`
	CollectVersion int            `gorm:"default:0;not null"` // 收藏操作乐观锁版本号
	LikeVersion    int            `gorm:"default:0;not null"` // 喜欢操作的乐观锁版本号
	CreatedAt      time.Time      `gorm:"index"`
	UpdatedAt      time.Time      `gorm:"index"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type PostLike struct {
	UserID    uint      `gorm:"primaryKey;index"`
	PostID    uint      `gorm:"primaryKey;index"`
	CreatedAt time.Time `gorm:"index"`
}

type PostCollect struct {
	UserID    uint      `gorm:"primaryKey;index"`
	PostID    uint      `gorm:"primaryKey;index"`
	CreatedAt time.Time `gorm:"index"`
}
