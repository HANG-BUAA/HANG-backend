package model

type Permission struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `gorm:"type:varchar(50);not null;unique"`
	Roles []Role `gorm:"many2many:permission_roles"`
}
