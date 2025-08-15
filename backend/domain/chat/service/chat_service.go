package service

import (
	"errors"

	friendrepo "alice/domain/appfriend/repository"
	chatentity "alice/domain/chat/entity"
	chatrepo "alice/domain/chat/repository"
)

var (
	ErrNotFriends = errors.New("not friends")
)

type ChatService interface {
	Send(senderID, receiverID uint, content string, msgType string) (*chatentity.Message, error)
	History(a, b uint, page, pageSize int) ([]*chatentity.Message, int64, error)
	MarkRead(a, b uint, beforeID uint) error
	RecentConversations(self uint, page, pageSize int) ([]*chatentity.Conversation, int64, error)
}

type chatServiceImpl struct {
	repo       chatrepo.MessageRepository
	friendRepo friendrepo.FriendRepository
}

func NewChatService(repo chatrepo.MessageRepository, friendRepo friendrepo.FriendRepository) ChatService {
	return &chatServiceImpl{repo: repo, friendRepo: friendRepo}
}

func (s *chatServiceImpl) Send(senderID, receiverID uint, content string, msgType string) (*chatentity.Message, error) {
	if senderID == 0 || receiverID == 0 || senderID == receiverID || content == "" {
		return nil, errors.New("invalid params")
	}
	ok, err := s.friendRepo.AreFriends(senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotFriends
	}
	m := &chatentity.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Type:       firstNonEmpty(msgType, "text"),
		Content:    content,
	}
	if err := s.repo.Save(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *chatServiceImpl) History(a, b uint, page, pageSize int) ([]*chatentity.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListConversation(a, b, offset, pageSize)
}

func (s *chatServiceImpl) MarkRead(a, b uint, beforeID uint) error {
	return s.repo.MarkRead(a, b, beforeID)
}

func (s *chatServiceImpl) RecentConversations(self uint, page, pageSize int) ([]*chatentity.Conversation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListRecentConversations(self, offset, pageSize)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
