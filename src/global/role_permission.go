package global

type Permission struct {
	ID   uint
	Name string
}

type Role struct {
	ID   uint
	Name string
}

// 权限集合
var (
	PermissionManageUsers  = Permission{ID: 1, Name: "manage_users"}
	PermissionManageAdmins = Permission{ID: 2, Name: "manage_admins"}
)

var Permissions = []Permission{PermissionManageUsers, PermissionManageAdmins}

// 角色集合
var (
	RoleUser       = Role{ID: 1, Name: "user"}
	RoleAdmin      = Role{ID: 2, Name: "admin"}
	RoleSuperAdmin = Role{ID: 3, Name: "super_admin"}
)

var Roles = []Role{RoleUser, RoleAdmin, RoleSuperAdmin}

// RolePermissionMap 角色与权限的映射
var RolePermissionMap = []struct {
	Role        Role
	Permissions []Permission
}{
	{Role: RoleUser, Permissions: []Permission{}},
	{Role: RoleAdmin, Permissions: []Permission{PermissionManageUsers}},
	{Role: RoleSuperAdmin, Permissions: []Permission{PermissionManageUsers, PermissionManageAdmins}},
}
