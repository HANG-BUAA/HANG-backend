package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitPostRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		postApi := api.NewPostApi()
		{

		}
		rgAuthPost := rgAuth.Group("/posts")
		{
			rgAuthPost.POST("", postApi.Create)
			rgAuthPost.POST("/:post_id/like", postApi.Like)
			rgAuthPost.POST("/:post_id/collect", postApi.Collect)
			rgAuthPost.GET("", postApi.List)
		}
	})
}
