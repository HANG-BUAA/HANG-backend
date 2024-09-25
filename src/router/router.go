package router

import (
	"HANG-backend/src/api"
	_ "HANG-backend/src/docs"
	"HANG-backend/src/global"
	"HANG-backend/src/middleware"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type IFnRegisterRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup, rgAdminGroup *gin.RouterGroup)

var gfnRoutes []IFnRegisterRoute // 注册各个路由组的函数列表

// RegisterRoute 注册一个路由组（添加到gfnRoutes中）
func RegisterRoute(fn IFnRegisterRoute) {
	if fn == nil {
		return
	}
	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {
	// 利用管道和协程优雅地退出
	// 目的是最后程序终止时执行收尾的 Clean 函数
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	r := gin.Default()
	r.Use(middleware.Cors())    // 挂载跨域中间件
	pingApi := api.NewPingApi() // 测试连通接口
	r.GET("/ping", pingApi.Ping)

	// 公共路由组、用户鉴权路由组，以及管理端路由组
	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")
	rgAdmin := r.Group("/api/v1/admin")
	rgAuth.Use(middleware.Auth()) // 登录鉴权
	rgAdmin.Use(middleware.Auth())

	// 注册基础平台路由（添加到gfnRoutes中）
	initBasePlatformRoutes()

	// 注册自定义验证器
	registerCustomValidator()

	// 循环遍历gfnRoutes，执行其中的函数
	for _, fnRegisterRoute := range gfnRoutes {
		fnRegisterRoute(rgPublic, rgAuth, rgAdmin)
	}

	// 集成 swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务
	port := viper.GetString("server.port")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		global.Logger.Infof("Start Listen: %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error("Start Server Error: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error("Server Shutdown Error: %s", err.Error())
		return
	}
	global.Logger.Info("Server Shutdown Success")
}

func initBasePlatformRoutes() {
	InitUserRoutes()
	InitPostRoutes()
	InitCommentRoutes()
}

// 自定义校验器
func registerCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("first_is_a", func(fl validator.FieldLevel) bool {
			if value, ok := fl.Field().Interface().(string); ok {
				if value != "" && strings.Index(value, "a") == 0 {
					return true
				}
			}
			return false
		})
	}
}
