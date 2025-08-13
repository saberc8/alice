package entity

import (
	"time"
)

// RoleStatus 角色状态
type RoleStatus string

const (
	RoleStatusActive   RoleStatus = "active"
	RoleStatusInactive RoleStatus = "inactive"
)

// Role 角色实体
type Role struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Code        string     `json:"code" gorm:"uniqueIndex;not null;size:100"`
	Description *string    `json:"description" gorm:"size:500"`
	Status      RoleStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// IsActive 检查角色是否激活
func (r *Role) IsActive() bool {
	return r.Status == RoleStatusActive
}
