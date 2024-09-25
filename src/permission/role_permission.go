package permission

import (
	"HANG-backend/src/global"
	"HANG-backend/src/model"
)

type Role uint
type RoleWeight struct {
	Role        Role
	Description string
	Weight      int
}

// 角色的集合
const (
	User Role = iota + 1
	Admin
	SuperAdmin
)

var Roles = []RoleWeight{
	{User, "普通用户", 1000},
	{Admin, "管理员", 7000},
	{SuperAdmin, "超级管理员", 10000},
}

type Permission uint
type PermissionWeight struct {
	Permission  Permission
	Description string
	Weight      int
}

// 权限的集合，永远不要打乱当前的顺序，有新的权限的话加到最后面
const (
	SetRole     Permission = iota + 1 // 设置其他人的权限
	PostPost                          // 发帖子
	PostComment                       // 发评论
	GetUserList                       // 获取用户列表
)

var Permissions = []PermissionWeight{
	{SetRole, "设置角色", 10000},
	{PostPost, "发帖", 1000},
	{PostComment, "发评论", 1000},
	{GetUserList, "获取用户列表", 7000},
}

var userPermissions = []Permission{
	PostPost, PostComment,
}

var adminPermissions = []Permission{
	PostPost,
	PostComment,
	GetUserList,
}

var superAdminPermissions = []Permission{
	SetRole,
	PostPost,
	PostComment,
	GetUserList,
}

var rolePermissionMap = map[Role][]Permission{
	User:       userPermissions,
	Admin:      adminPermissions,
	SuperAdmin: superAdminPermissions,
}

func InitPermissions() {
	// 先清空原来的表
	global.RDB.Exec("TrUNCATE TABLE permission")
	for _, p := range Permissions {
		// 创建相应表中的权限
		per := model.Permission{
			ID:          uint(p.Permission),
			Description: p.Description,
			Weight:      p.Weight,
		}
		global.RDB.Create(&per)
	}
}

func InitUserPermission(userID uint, role Role) error {
	for _, permission := range rolePermissionMap[role] {
		userPermission := model.UserPermission{
			UserID:       userID,
			PermissionID: uint(permission),
		}
		if err := global.RDB.FirstOrCreate(&userPermission).Error; err != nil {
			return err
		}
	}
	return nil
}
