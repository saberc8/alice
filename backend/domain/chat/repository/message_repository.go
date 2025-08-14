package repository

import (
	chatentity "alice/domain/chat/entity"
)

type MessageRepository interface {
	Save(msg *chatentity.Message) error
	ListConversation(a, b uint, offset, limit int) ([]*chatentity.Message, int64, error)
	MarkRead(a, b uint, beforeID uint) error
}
