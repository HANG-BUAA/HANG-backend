package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdmin *gin.RouterGroup) {
		userApi := api.NewUserApi()
		{
			rgPublic.POST("/login", userApi.Login)
			rgPublic.POST("/send-email", userApi.SendEmail)
			rgPublic.POST("/register", userApi.Register)
		}

		rgAuthUser := rgAuth.Group("/users")
		{
			rgAuthUser.PUT("/update-avatar", userApi.UploadAvatar)
			rgAuthUser.PUT("", userApi.UpdateUser)
			rgAuthUser.GET("/notifications", userApi.ListNotification)
		}
		//rgAdminUser := rgAdmin.Group("/users")
		//{
		//
		//}
	})
}
