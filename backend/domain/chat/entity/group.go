package entity

import "time"

// Group 群聊信息
type Group struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(120);not null"`
	OwnerID   uint      `json:"owner_id" gorm:"not null;index"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(255);default:''"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Group) TableName() string { return "app_chat_groups" }

// GroupMember 群成员
type GroupMember struct {
	GroupID  uint      `json:"group_id" gorm:"primaryKey;autoIncrement:false"`
	UserID   uint      `json:"user_id" gorm:"primaryKey;autoIncrement:false;index"`
	Role     string    `json:"role" gorm:"type:varchar(20);not null;default:'member'"` // owner/member
	JoinedAt time.Time `json:"joined_at"`
}

func (GroupMember) TableName() string { return "app_chat_group_members" }

// GroupMessage 群消息
type GroupMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GroupID   uint      `json:"group_id" gorm:"not null;index"`
	SenderID  uint      `json:"sender_id" gorm:"not null;index"`
	Type      string    `json:"type" gorm:"type:varchar(20);not null;default:'text'"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (GroupMessage) TableName() string { return "app_chat_group_messages" }

// GroupReadCursor 记录成员在群里的最后已读消息 ID
type GroupReadCursor struct {
	GroupID       uint      `json:"group_id" gorm:"primaryKey;autoIncrement:false"`
	UserID        uint      `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	LastReadMsgID uint      `json:"last_read_msg_id" gorm:"not null;default:0"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (GroupReadCursor) TableName() string { return "app_chat_group_read_cursors" }
