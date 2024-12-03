package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	AVALIABLE_PRODUCT_STATUS = iota + 1
	SOLD_PRODUCT_STATUS
	DELETED_PRODUCT_STATUS
)

type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement;not null"`
	UserID      uint    `gorm:"index;not null"`
	Title       string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text;not null"`
	Campus      int     `gorm:"index;not null"`
	Price       float32 `gorm:"index;type:float;not null"`
	Status      int     `gorm:"index;not null"` // 商品状态
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProductTag struct {
	ProductID uint `gorm:"PrimaryKey"`
	TagID     uint `gorm:"PrimaryKey"`
}
