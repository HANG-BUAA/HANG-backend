package router

import (
	"HANG-backend/src/api"
	"HANG-backend/src/middleware"
	"HANG-backend/src/permission"
	"github.com/gin-gonic/gin"
)

func InitTagRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup) {
		tagApi := api.NewTagApi()
		{

		}
		//rgAuthTag := rgAuth.Group("/tags")
		//{
		//
		//}
		rgAdminTag := rgAdminGroup.Group("/tags")
		{
			rgAdminTag.POST("", middleware.Permission(permission.CreateTag), tagApi.AdminCreateTag)
		}
	})
}
