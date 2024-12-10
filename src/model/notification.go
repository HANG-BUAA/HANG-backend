package model

type Notification struct {
	ID             uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Type           string `gorm:"type:varchar(255); not null" json:"type"`
	OperatorID     uint   `gorm:"not null" json:"operator_id"`
	OperatorName   string `gorm:"type:varchar(255); not null" json:"operator_name"`
	OperatorAvatar string `gorm:"type:varchar(255); not null" json:"operator_avatar"`
	NotifierID     uint   `gorm:"index; not null" json:"notifier_id"`
	EntityID       uint   `gorm:"not null" json:"entity_id"`
}
