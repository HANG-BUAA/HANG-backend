package cmd

import (
	"HANG-backend/src/config"
	"HANG-backend/src/global"
	"HANG-backend/src/router"
	"HANG-backend/src/utils"
	"fmt"
	"time"
)

func Start() {
	var initErr error
	// 初始化系统配置文件
	config.InitConfig()

	// 初始化日志组件
	global.Logger = config.InitLogger()

	// 初始化数据库连接
	db, err := config.InitDB()
	global.DB = db
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}

	// 初始化 Redis 连接
	rdClient, err := config.InitRedis()
	global.RedisClient = rdClient
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}

	_ = global.RedisClient.Set("fxk", "hhh", time.Second*100)
	fmt.Println(global.RedisClient.Get("fxk"))

	if initErr != nil {
		if global.Logger != nil {
			global.Logger.Error(initErr.Error())
		}

		panic(initErr.Error())
	}

	// 初始化系统路由
	router.InitRouter()
}

func Clean() {
	fmt.Println("=======Clean==========")
}
