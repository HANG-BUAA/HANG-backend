package router

import (
	"HANG-backend/src/api"
	"HANG-backend/src/middleware"
	"HANG-backend/src/permission"
	"github.com/gin-gonic/gin"
)

func InitCommentRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup) {
		commentApi := api.NewCommentApi()
		{

		}
		rgAuthComment := rgAuth.Group("/comments")
		{
			rgAuthComment.POST("", middleware.Permission(permission.PostComment), commentApi.Create)
			rgAuthComment.POST("/:comment_id/like", middleware.CommentExistence(middleware.URI), commentApi.Like)
			rgAuthComment.POST("/:comment_id/unlike", middleware.CommentExistence(middleware.URI), commentApi.Unlike)
			rgAuthComment.GET("", middleware.CheckPaginationParams(), commentApi.List)
		}
	})
}
