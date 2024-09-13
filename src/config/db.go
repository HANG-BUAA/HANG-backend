package config

import (
	"HANG-backend/src/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB() (*gorm.DB, error) {
	logMode := logger.Info
	if !viper.GetBool("mode.develop") {
		logMode = logger.Warn
	}
	username := viper.GetString("db.mysql.username")
	password := viper.GetString("db.mysql.password")
	host := viper.GetString("db.mysql.host")
	port := viper.GetString("db.mysql.port")
	dbname := viper.GetString("db.mysql.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("db.mysql.maxIdleConn"))
	sqlDB.SetMaxOpenConns(viper.GetInt("db.mysql.maxOpenConn"))
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
