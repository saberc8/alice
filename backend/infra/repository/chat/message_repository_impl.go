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
