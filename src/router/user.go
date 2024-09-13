package router

import (
    "HANG-backend/src/api"
    "github.com/gin-gonic/gin"
)

func InitUserRoutes() {
    RegisterRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
        userApi := api.NewUserApi()
        {
            rgPublic.POST("/login", userApi.Login)
            rgPublic.POST("/register", userApi.Register)
        }

        rgAuthUser := rgAuth.Group("/user")
        {
            rgAuthUser.GET("", func(context *gin.Context) {

            })
        }
    })
}
