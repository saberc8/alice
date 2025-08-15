package chat

import (
	"gorm.io/gorm"

	chatentity "alice/domain/chat/entity"
	chatrepo "alice/domain/chat/repository"
)

type messageRepositoryImpl struct{ db *gorm.DB }

func NewMessageRepository(db *gorm.DB) chatrepo.MessageRepository {
	return &messageRepositoryImpl{db: db}
}

func (r *messageRepositoryImpl) Save(msg *chatentity.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepositoryImpl) ListConversation(a, b uint, offset, limit int) ([]*chatentity.Message, int64, error) {
	var total int64
	q := r.db.Model(&chatentity.Message{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", a, b, b, a,
	)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []*chatentity.Message
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	// 逆序返回按时间正序，前端/调用方可自行处理，这里保持 DESC 以配合分页
	return rows, total, nil
}

func (r *messageRepositoryImpl) MarkRead(a, b uint, beforeID uint) error {
	return r.db.Model(&chatentity.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND id <= ? AND is_read = ?", b, a, beforeID, false).
		Updates(map[string]interface{}{"is_read": true}).Error
}

// ListRecentConversations 聚合最近会话
func (r *messageRepositoryImpl) ListRecentConversations(self uint, offset, limit int) ([]*chatentity.Conversation, int64, error) {
	type row struct {
		PeerID    uint
		LastID    uint
		UnreadCnt int64
	}
	// 1) 找到跟我相关的所有 peer_id（对端）以及该会话的最后一条消息ID
	var rows []row
	sub := r.db.Model(&chatentity.Message{}).
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS peer_id, MAX(id) AS last_id", self).
		Where("sender_id = ? OR receiver_id = ?", self, self).
		Group("peer_id")
	// 分页总数
	var total int64
	if err := r.db.Table("(?) as t", sub).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Table("(?) as t", sub).
		Order("last_id DESC").
		Offset(offset).Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	// 2) 读取最后一条消息详情与未读数
	convs := make([]*chatentity.Conversation, 0, len(rows))
	for _, rrow := range rows {
		var last chatentity.Message
		if err := r.db.First(&last, rrow.LastID).Error; err != nil {
			return nil, 0, err
		}
		// 未读：对方->我 且 is_read=false
		var unread int64
		if err := r.db.Model(&chatentity.Message{}).
			Where("sender_id = ? AND receiver_id = ? AND is_read = ?", rrow.PeerID, self, false).
			Count(&unread).Error; err != nil {
			return nil, 0, err
		}
		convs = append(convs, &chatentity.Conversation{PeerID: rrow.PeerID, LastMessage: &last, UnreadCount: unread})
	}
	return convs, total, nil
}
