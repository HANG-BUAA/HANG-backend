package model

type Course struct {
	ID uint `gorm:"primaryKey;autoIncrement; not null"`
}
