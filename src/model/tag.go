package model

import (
	"errors"
	"gorm.io/gorm"
)

const (
	COURSE_TAG = iota + 1
)

type Tag struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(255);not null;unique;index"`
	Type int    `gorm:"type:int;not null"`
}

func (m *Tag) BeforeSave(tx *gorm.DB) error {
	if m.Type < 1 || m.Type > 1 {
		return errors.New("invalid tag type")
	}
	return nil
}
