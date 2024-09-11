package cmd

import (
	"HANG-backend/src/config"
	"HANG-backend/src/global"
	"HANG-backend/src/router"
	"fmt"
)

func Start() {
	// 初始化系统配置文件
	config.InitConfig()

	// 初始化日志组件
	global.Logger = config.InitLogger()

	// 初始化系统路由
	router.InitRouter()
}

func Clean() {
	fmt.Println("=======Clean==========")
}
