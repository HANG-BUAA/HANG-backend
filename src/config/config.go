package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")     // 本地开发路径
	viper.AddConfigPath("/root/config/") // Docker 容器路径
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal custom_error config file: %s \n", err))
	}
}
