package model

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"size:255"`
	Weight      int    `gorm:"default:0"`
}

type UserPermission struct {
	UserID       uint `gorm:"primaryKey;index"`
	PermissionID uint `gorm:"primaryKey;index"`
}
