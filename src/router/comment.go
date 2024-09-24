package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
)

func InitCommentRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		commentApi := api.NewCommentApi()
		{

		}
		rgAuthComment := rgAuth.Group("/comments")
		{
			rgAuthComment.POST("", commentApi.Create)
			rgAuthComment.POST("/:comment_id/like", commentApi.Like)
			rgAuthComment.GET("", commentApi.List)
		}
	})
}
