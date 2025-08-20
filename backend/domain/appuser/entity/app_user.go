package entity

import "time"

// AppUserStatus 移动端用户状态
type AppUserStatus string

const (
	AppUserStatusActive   AppUserStatus = "active"
	AppUserStatusInactive AppUserStatus = "inactive"
	AppUserStatusBanned   AppUserStatus = "banned"
)

// AppUser 移动端用户表（独立于后台管理用户）
type AppUser struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	Username     string        `json:"username" gorm:"default:''"`
	Email        string        `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string        `json:"-" gorm:"not null"`
	Nickname     string        `json:"nickname" gorm:"default:''"`
	Avatar       string        `json:"avatar" gorm:"default:''"`
	Gender       string        `json:"gender" gorm:"type:varchar(10);default:''"` // male / female / other （可枚举）
	Bio          string        `json:"bio" gorm:"default:''"`
	Status       AppUserStatus `json:"status" gorm:"not null;default:'active'"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func (AppUser) TableName() string { return "app_users" }

func (u *AppUser) IsActive() bool { return u.Status == AppUserStatusActive }
