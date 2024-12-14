package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitLogRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup) {
		logApi := api.NewLogApi()
		{

		}
		rgAdminLog := rgAdminGroup.Group("/logs")
		{
			rgAdminLog.GET("/keywords", logApi.ListKeywords)
		}
	})
}
