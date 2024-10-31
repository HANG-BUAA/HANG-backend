package router

import (
	"HANG-backend/src/api"
	"HANG-backend/src/middleware"
	"HANG-backend/src/permission"
	"github.com/gin-gonic/gin"
)

func InitCourseRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup) {
		courseApi := api.NewCourseApi()
		{

		}
		rgAuthCourse := rgAuth.Group("/courses")
		{
			rgAuthCourse.POST("/reviews", middleware.Permission(permission.ReviewCourse), courseApi.CreateReview)
			rgAuthCourse.POST("/reviews/:review_id/like", middleware.CourseReviewExistence(middleware.URI), courseApi.LikeReview)
			rgAuthCourse.GET("", middleware.CheckPaginationParams(), courseApi.ListCourse)
			rgAuthCourse.GET("/reviews", middleware.CheckPaginationParams(), courseApi.ListReview)
			rgAuthCourse.GET("/:course_id", middleware.CourseExistence(middleware.URI), courseApi.Retrieve)
		}
		rgAdminCourse := rgAdminGroup.Group("/courses")
		{
			rgAdminCourse.POST("", middleware.Permission(permission.CreateCourse), courseApi.CreateCourse)
		}
	})
}
