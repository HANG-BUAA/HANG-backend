package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitUserRoutes() {
	RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		rgPublic.POST("/login", func(context *gin.Context) {
			context.AbortWithStatusJSON(http.StatusOK, gin.H{
				"msg": "Login Success",
			})
		})

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
