package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"github.com/gin-gonic/gin"
	"strconv"
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

type ChatApi struct {
	BaseApi
}

func NewChatApi() ChatApi {
	return ChatApi{
		BaseApi: NewBaseApi(),
	}
}

func (m ChatApi) CreateMessage(c *gin.Context) {
	m.Ctx = c
	sender := c.MustGet("user").(*model.User)
	senderID := sender.ID

	var request struct {
		ReceiverID uint   `json:"receiver_id"` // 接收者ID
		Content    string `json:"content"`     // 消息内容
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	message := ChatMessage{
		SenderID:   senderID,
		ReceiverID: request.ReceiverID,
		Content:    request.Content,
		CreatedAt:  time.Now(),
	}

	if err := global.RDB.Create(&message).Error; err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}
	m.OK(ResponseJson{
		Data: message,
	})
}

func (m ChatApi) ListMessage(c *gin.Context) {
	m.Ctx = c
	sender := c.MustGet("user").(*model.User)
	senderID := sender.ID

	// 从GET请求的URL参数中获取接收者ID，设置默认值为0，如果获取失败也为0
	receiverID, err := strconv.ParseUint(c.Query("receiver_id"), 10, 32)
	if err != nil {
		receiverID = 0
	}
	// 从GET请求的URL参数中获取页码，设置默认值为1，如果获取失败也为1
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	// 从GET请求的URL参数中获取每页数量，设置默认值为10，如果获取失败也为10
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		pageSize = 10
	}

	var messages []ChatMessage
	var totalCount int64

	// 查询聊天记录
	err = global.RDB.
		Where("(sender_id =? AND receiver_id =?) OR (sender_id =? AND receiver_id =?)", senderID, receiverID, receiverID, senderID).
		Order("created_at DES"). // 按时间升序排列
		Offset((page - 1) * pageSize).Limit(pageSize). // 分页
		Find(&messages).Error
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	// 获取总记录数
	err = global.RDB.
		Table("chat_message").
		Where("(sender_id =? AND receiver_id =?) OR (sender_id =? AND receiver_id =?)", senderID, receiverID, receiverID, senderID).
		Count(&totalCount).Error
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	// 返回响应
	m.OK(ResponseJson{
		Data: gin.H{
			"messages":  messages,
			"total":     totalCount,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
