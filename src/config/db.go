package config

import (
	"HANG-backend/src/global"
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

	if err = db.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{}); err != nil {
		return nil, err
	}

	// 调试模式下，每次才重建权限表和角色表
	if viper.GetBool("mode.develop") {
		if err := initRolePermission(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func initRolePermission(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 初始化权限
		for _, permission := range global.Permissions {
			// 创建或更新权限
			if err := tx.Where(model.Permission{Name: permission.Name}).FirstOrCreate(&model.Permission{
				ID:   permission.ID,
				Name: permission.Name,
			}).Error; err != nil {
				return err
			}
		}

		// 使用 global.RolePermissionMap 初始化角色和权限
		for _, roleData := range global.RolePermissionMap {
			// 创建或查找角色
			if err := tx.Where(model.Role{Name: roleData.Role.Name}).FirstOrCreate(&model.Role{
				ID:   roleData.Role.ID,
				Name: roleData.Role.Name,
			}).Error; err != nil {
				return err
			}

			// 查找角色
			var role model.Role
			if err := tx.Where(model.Role{Name: roleData.Role.Name}).First(&role).Error; err != nil {
				return err
			}

			// 清空已有的权限关联
			if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
				return err
			}

			// 为角色添加权限
			var permissions []model.Permission
			for _, globalPermission := range roleData.Permissions {
				var dbPermission model.Permission
				if err := tx.Where(model.Permission{Name: globalPermission.Name}).First(&dbPermission).Error; err != nil {
					return err
				}
				permissions = append(permissions, dbPermission)
			}

			// 更新角色的权限
			if err := tx.Model(&role).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}
