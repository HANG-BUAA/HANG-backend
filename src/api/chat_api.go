package api

import (
	"time"
)

// -------------- model ----------------

type ChatMessage struct {
	ID         uint      `gorm:"primaryKey;autoIncrement; not null"`
	SenderID   uint      `gorm:"index; not null"`
	ReceiverID uint      `gorm:"index; not null"`
	Content    string    `gorm:"type:text; not null"`
	CreatedAt  time.Time `gorm:"index"`
}

type CharApi struct {
	BaseApi
}

func NewCharApi() CharApi {
	return CharApi{
		BaseApi: NewBaseApi(),
	}
}

//func (m CharApi) CreateMessage(c *gin.Context) {
//	sender := c.MustGet("user").(*model.User)
//	senderID := sender.ID
//}
