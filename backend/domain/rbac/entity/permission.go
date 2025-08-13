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
	ID          string           `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Name        string           `json:"name" gorm:"not null;size:100"`
	Code        string           `json:"code" gorm:"uniqueIndex;not null;size:100"`
	Resource    string           `json:"resource" gorm:"not null;size:100"`
	Action      string           `json:"action" gorm:"not null;size:50"`
	Description *string          `json:"description" gorm:"size:500"`
	Status      PermissionStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// IsActive 检查权限是否激活
func (p *Permission) IsActive() bool {
	return p.Status == PermissionStatusActive
}
