package router

import (
	"HANG-backend/src/api"
	"HANG-backend/src/middleware"
	"HANG-backend/src/permission"
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
		}
		rgAdminUser := rgAdmin.Group("/users")
		{
			rgAdminUser.GET("", middleware.Permission(permission.GetUserList), userApi.AdminList)
		}
	})
}
