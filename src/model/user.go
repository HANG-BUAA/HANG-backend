package model

import (
	"HANG-backend/src/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	StudentID string `gorm:"type:varchar(20);not null"`
	UserName  string `gorm:"type:varchar(20);not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Avatar    string `gorm:"type:varchar(255)"`
	Role      uint   `gorm:"type:int(8);not null)"`
}

func (m *User) Encrypt() error {
	hash, err := utils.Encrypt(m.Password)
	if err == nil {
		m.Password = hash
	}
	return err
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	return m.Encrypt()
}
