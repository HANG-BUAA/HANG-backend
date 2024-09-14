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
}

func (m *User) Encrypt() error {
	stHash, err := utils.Encrypt(m.Password)
	if err == nil {
		m.Password = stHash
	}
	return err
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	return m.Encrypt()
}
