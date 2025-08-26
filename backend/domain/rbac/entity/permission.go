package entity

import (
	"time"
)

// PermissionStatus 权限状态
type PermissionStatus string

const (
	PermissionStatusActive   PermissionStatus = "active"
	PermissionStatusInactive PermissionStatus = "inactive"
)

// Permission 权限实体
type Permission struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"not null;size:100"`
	Code string `json:"code" gorm:"uniqueIndex;not null;size:100"`
	// MenuID 绑定的菜单ID（按钮权限归属菜单），可为空
	MenuID      *uint            `json:"menu_id,omitempty" gorm:"index"`
	Resource    string           `json:"resource" gorm:"not null;size:100"`
	Action      string           `json:"action" gorm:"not null;size:50"`
	Description *string          `json:"description" gorm:"size:500"`
	Status      PermissionStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// IsActive 检查权限是否激活
func (p *Permission) IsActive() bool {
	return p.Status == PermissionStatusActive
}
