package api

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// -------------- model ----------------

type ChatMessage struct {
	ID         uint           `gorm:"primaryKey;autoIncrement; not null"`
	SenderID   uint           `gorm:"index; not null"`
	ReceiverID uint           `gorm:"index; not null"`
	Content    string         `gorm:"type:text; not null"`
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
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
		Order("created_at DESC"). // 按时间升序排列
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

func (m ChatApi) LongPollingNewMessages(c *gin.Context) {
	m.Ctx = c
	sender := c.MustGet("user").(*model.User)
	senderID := sender.ID

	// 从GET请求的URL参数中获取接收者ID，设置默认值为0，如果获取失败也为0
	receiverID, err := strconv.ParseUint(c.Query("receiver_id"), 10, 32)
	if err != nil {
		receiverID = 0
	}
	// 设置超时时间，比如这里设置为30秒（可根据实际情况调整）
	timeout := time.Second * 30

	var lastMessageTime time.Time
	// 先查询一下当前已有的最新消息时间，用于后续对比判断是否有新消息
	var messages []ChatMessage
	err = global.RDB.
		Where("sender_id =? AND receiver_id =?", receiverID, senderID).
		Order("created_at DESC").
		Limit(1).
		Find(&messages).Error
	if err == nil && len(messages) > 0 {
		lastMessageTime = messages[0].CreatedAt
	}

	start := time.Now()
	for {
		var newMessages []ChatMessage
		// 查询在最后一条消息时间之后的新消息
		err = global.RDB.
			Where("sender_id =? AND receiver_id =?", receiverID, senderID).
			Where("created_at >?", lastMessageTime).
			Order("created_at ASC").
			Find(&newMessages).Error
		if err != nil {
			m.Fail(ResponseJson{
				Msg: err.Error(),
			})
			return
		}
		if len(newMessages) > 0 {
			// 如果有新消息，返回新消息给客户端
			m.OK(ResponseJson{
				Data: gin.H{
					"messages": newMessages,
				},
			})
			return
		}
		if time.Since(start) > timeout {
			// 如果超时了还没有新消息，返回空消息列表给客户端
			m.OK(ResponseJson{
				Data: gin.H{
					"messages": []ChatMessage{},
				},
			})
			return
		}
		time.Sleep(time.Millisecond * 500) // 短暂休眠后再次检查是否有新消息
	}
}

type UserDetail struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	StudentID string `json:"student_id"`
	Avatar    string `json:"avatar"`
}

func (m ChatApi) ListFriends(c *gin.Context) {
	m.Ctx = c
	sender := c.MustGet("user").(*model.User)
	senderID := sender.ID

	var chattedUsers []uint
	// 使用子查询和DISTINCT关键字找出所有和当前用户聊过天的不同用户
	err := global.RDB.Table("chat_message").
		Select("DISTINCT CASE WHEN sender_id =? THEN receiver_id ELSE sender_id END as user_id", senderID).
		Where("sender_id =? OR receiver_id =?", senderID, senderID).
		Find(&chattedUsers).Error
	if err != nil {
		m.Fail(ResponseJson{
			Msg: err.Error(),
		})
		return
	}

	var res []UserDetail
	for _, userID := range chattedUsers {
		var user model.User
		if err := global.RDB.Where("id = ?", userID).Find(&user).Error; err != nil {
			m.Fail(ResponseJson{
				Msg: err.Error(),
			})
			return
		}
		res = append(res, UserDetail{
			ID:        user.ID,
			Username:  user.Username,
			Avatar:    user.Avatar,
			StudentID: user.StudentID,
		})
	}

	m.OK(ResponseJson{
		Data: gin.H{
			"users": res,
		},
	})
}
