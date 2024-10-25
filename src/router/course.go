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
			rgAuthCourse.POST("/reviews", middleware.Permission(permission.ReviewCourse), courseApi.CreateCourseReview)
		}
		rgAdminCourse := rgAdminGroup.Group("/courses")
		{
			rgAdminCourse.POST("", middleware.Permission(permission.CreateCourse), courseApi.CreateCourse)
		}
	})
}
