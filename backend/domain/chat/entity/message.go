package entity

import "time"

// Message 私聊消息
type Message struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	SenderID   uint       `json:"sender_id" gorm:"not null;index:idx_conv,priority:1"`
	ReceiverID uint       `json:"receiver_id" gorm:"not null;index:idx_conv,priority:2"`
	Type       string     `json:"type" gorm:"type:varchar(20);not null;default:'text'"`
	Content    string     `json:"content" gorm:"type:text;not null"`
	IsRead     bool       `json:"is_read" gorm:"not null;default:false;index"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (Message) TableName() string { return "app_chat_messages" }
