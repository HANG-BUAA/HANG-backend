package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitChatRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup) {
		chatApi := api.NewChatApi()
		{

		}
		rgAuthChat := rgAuth.Group("/chat")
		{
			rgAuthChat.POST("/messages", chatApi.CreateMessage)
			rgAuthChat.GET("/messages", chatApi.ListMessage)
			rgAuthChat.GET("/poll", chatApi.LongPollingNewMessages)
			rgAuthChat.GET("/friends", chatApi.ListFriends)
		}
	})
}
