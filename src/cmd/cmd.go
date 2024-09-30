package cmd

import (
	"HANG-backend/src/config"
	"HANG-backend/src/global"
	"HANG-backend/src/permission"
	"HANG-backend/src/router"
	"HANG-backend/src/utils"
	"fmt"
)

func Start() {
	var initErr error
	// 初始化系统配置文件
	config.InitConfig()

	// 初始化日志组件
	global.Logger = config.InitLogger()

	// 初始化数据库连接
	db, err := config.InitDB()
	global.RDB = db
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}
	// 初始化权限
	permission.InitPermissions()

	// 初始化 Redis 连接
	rdClient, err := config.InitRedis()
	global.RedisClient = rdClient
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}

	// 初始化 rabbit_mq 连接
	rabbit, err := config.InitRabbitMq()
	global.RabbitMqChannel = rabbit
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}

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
	global.RabbitMqChannel.Close()
	global.RabbitMqConnection.Close()
	fmt.Println("=======Clean==========")
}
