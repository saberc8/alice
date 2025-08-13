package entity

import (
	"time"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

// User 用户实体
type User struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	Status       UserStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}
