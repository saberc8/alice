package entity

import "time"

// MomentComment 动态评论
type MomentComment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	MomentID  uint      `json:"moment_id" gorm:"index"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (MomentComment) TableName() string { return "app_moment_comments" }
