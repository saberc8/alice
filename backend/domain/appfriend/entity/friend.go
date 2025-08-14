package entity

import "time"

// FriendRelation 好友关系（双向确认后存两条，或单向 + 状态）
// 这里采用单表两条记录的方式：A 添加 B 确认后，保存 A->B 与 B->A 两条 active 记录
// 也可扩展一个“请求表”，本实现先做直接添加+自动双向建立（简单 MVP）
type FriendRelation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_friend,priority:1;uniqueIndex:ux_user_friend_pair,priority:1"`
	FriendID  uint      `json:"friend_id" gorm:"not null;index:idx_user_friend,priority:2;uniqueIndex:ux_user_friend_pair,priority:2"`
	CreatedAt time.Time `json:"created_at"`
}

func (FriendRelation) TableName() string { return "app_friend_relations" }
