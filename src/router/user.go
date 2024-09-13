package router

import (
	"HANG-backend/src/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitUserRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		userApi := api.NewUserApi()
		rgPublic.POST("/login", userApi.Login)

		rgAuthUser := rgAuth.Group("/user")
		rgAuthUser.GET("", func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"data": []map[string]any{
					{"id": 1, "name": "zs"},
					{"id": 2, "name": "ls"},
				},
			})
		})

		rgAuthUser.GET("/:id", func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"id":   1,
				"name": "zs",
			})
		})
	})
}
