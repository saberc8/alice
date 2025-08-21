package entity

import (
	"strings"
	"time"
)

// Moment 朋友圈动态（公开，可被所有用户查看；后续可扩展可见范围）
type Moment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Images    string    `json:"images" gorm:"type:text;default:''"` // 逗号分隔的相对路径 /bucket/object
	CreatedAt time.Time `json:"created_at"`
}

func (Moment) TableName() string { return "app_moments" }

// ParseImages 将存储字段解析为 slice
func (m *Moment) ParseImages() []string {
	if m.Images == "" {
		return []string{}
	}
	parts := []string{}
	for _, p := range strings.Split(m.Images, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}
