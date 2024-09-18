package model

type Role struct {
    ID          uint         `gorm:"primarykey"`
    Name        string       `gorm:"type:varchar(20);not null;unique"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}
