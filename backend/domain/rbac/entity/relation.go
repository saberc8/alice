package entity

import (
	"time"
)

// UserRole 用户角色关联表
type UserRole struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id" gorm:"not null;index"`
	RoleID   uint      `json:"role_id" gorm:"not null;index"`
	CreateAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoleID       uint      `json:"role_id" gorm:"not null;index"`
	PermissionID uint      `json:"permission_id" gorm:"not null;index"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RoleID    uint      `json:"role_id" gorm:"not null;index"`
	MenuID    uint      `json:"menu_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}
