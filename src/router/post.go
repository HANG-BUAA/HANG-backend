package router

import (
	"HANG-backend/src/api"
	"HANG-backend/src/middleware"
	"HANG-backend/src/permission"
	"github.com/gin-gonic/gin"
)

func InitPostRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdmin *gin.RouterGroup) {
		postApi := api.NewPostApi()
		{

		}
		rgAuthPost := rgAuth.Group("/posts")
		{
			rgAuthPost.POST("", middleware.Permission(permission.PostPost), postApi.Create)
			rgAuthPost.POST("/:post_id/like", middleware.PostExistence(middleware.URI), postApi.Like)
			rgAuthPost.POST("/:post_id/collect", middleware.PostExistence(middleware.URI), postApi.Collect)
			rgAuthPost.GET("", middleware.CheckPaginationParams(), postApi.List)
			rgAuthPost.GET("/collections", middleware.CheckPaginationParams(), postApi.CollectionList)
		}
	})
}
