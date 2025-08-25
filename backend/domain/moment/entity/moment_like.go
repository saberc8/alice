package entity

import "time"

// MomentLike 动态点赞
type MomentLike struct {
	MomentID  uint      `json:"moment_id" gorm:"primaryKey;autoIncrement:false"`
	UserID    uint      `json:"user_id" gorm:"primaryKey;autoIncrement:false;index:idx_moment_like_user,unique"`
	CreatedAt time.Time `json:"created_at"`
}

func (MomentLike) TableName() string { return "app_moment_likes" }
