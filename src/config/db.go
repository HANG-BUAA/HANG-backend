package config

import (
    "HANG-backend/src/model"
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
    db, err := gorm.Open(mysql.Open(""), &gorm.Config{
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
