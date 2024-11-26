package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

const (
	SHAHE_CAMPUS = iota + 1
	XUEYUANLU_CAMPUS
	HANGZHOU_CAMPUS
)

// Course 课程
type Course struct {
	ID        string   `gorm:"primaryKey;type:varchar(100);not null;unique;index"`
	Name      string   `gorm:"type:varchar(100);not null;unique;index"`
	Credits   *float32 `gorm:"type:decimal(4,2);index"`
	Campus    *int     `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Course) BeforeSave(tx *gorm.DB) error {
	if m.Campus != nil {
		if *m.Campus < 1 || *m.Campus > 3 {
			return errors.New("invalid campus")
		}
	}
	return nil
}

type CourseTag struct {
	CourseID string `gorm:"primaryKey"`
	TagID    uint   `gorm:"primaryKey"`
}

// CourseReview 课程评价
type CourseReview struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;not null"`
	CourseID    string `gorm:"type:varchar(100);index;not null"`
	UserID      uint   `gorm:"index;not null"`
	Content     string `gorm:"type:text;not null"`
	Score       int    `gorm:"type:int;not null;check:score >= 1 and score <= 5"`
	LikeNum     int    `gorm:"default:0;index; not null"`
	LikeVersion int    `gorm:"default:0;not null"` // 喜欢操作的乐观锁版本号
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CourseReviewLike struct {
	CourseReviewID uint `gorm:"primaryKey"`
	UserID         uint `gorm:"primaryKey"`
}

const (
	MaterialSource_BHPAN int = iota + 1 // 北航云盘
	MaterialSource_BLOG                 // 博客
)

type CourseMaterial struct {
	ID          uint   `gorm:"primaryKey;autoIncrement;not null"`
	CourseID    string `gorm:"type:varchar(100);index;not null"`
	UserID      uint   `gorm:"index;not null"`
	Link        string `gorm:"type:varchar(1024);not null"` // 资料链接
	Source      int    `gorm:"type:int;index"`
	Description string `gorm:"type:text;not null"`
	IsApproved  bool   `gorm:"default:false;not null"`
	IsOfficial  bool   `gorm:"default:false;not null"`
	LikeNum     int    `gorm:"default:0;index; not null"`
	LikeVersion int    `gorm:"default:0;not null"` // 喜欢操作的乐观锁版本号
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type CourseMaterialLike struct {
	CourseMaterialID uint `gorm:"primaryKey"`
	UserID           uint `gorm:"primaryKey"`
}

func (m *CourseMaterial) BeforeCreate(tx *gorm.DB) error {
	if m.Source < 1 || m.Source > 2 {
		return errors.New("invalid material source")
	}
	return nil
}
