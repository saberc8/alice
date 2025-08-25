package service

import (
	"context"
	"errors"
	"strings"
	"time"

	chatentity "alice/domain/chat/entity"
	chatrepo "alice/domain/chat/repository"
)

type GroupService interface {
	Create(ownerID uint, name string, memberIDs []uint, avatar string) (*chatentity.Group, error)
	Search(name string, limit int) ([]*chatentity.Group, error)
	Join(groupID, userID uint) error
	ListMessages(groupID uint, page, pageSize int) ([]*chatentity.GroupMessage, int64, error)
	SendMessage(groupID, senderID uint, msgType, content string) (*chatentity.GroupMessage, error)
	IsMember(groupID, userID uint) (bool, error)
	Get(groupID uint) (*chatentity.Group, error)
	UpdateGroup(operatorID, groupID uint, name, avatar string) (*chatentity.Group, error)
	ListUserGroups(userID uint, page, pageSize int) ([]*chatentity.Group, int64, error)
	UpdateLastRead(ctx context.Context, groupID, userID, msgID uint) error
	CountUnread(ctx context.Context, groupID, userID uint) (int64, error)
	ListMemberIDs(groupID uint) ([]uint, error)
	AddMembers(operatorID, groupID uint, userIDs []uint) error
	RemoveMember(operatorID, groupID, targetUserID uint) error
	ListMembers(groupID uint) ([]uint, error)
}

type groupServiceImpl struct{ repo chatrepo.GroupRepository }

func NewGroupService(r chatrepo.GroupRepository) GroupService { return &groupServiceImpl{repo: r} }

func (s *groupServiceImpl) Create(ownerID uint, name string, memberIDs []uint, avatar string) (*chatentity.Group, error) {
	name = strings.TrimSpace(name)
	if ownerID == 0 || name == "" {
		return nil, errors.New("invalid params")
	}
	// ensure owner included
	included := false
	for _, id := range memberIDs {
		if id == ownerID {
			included = true
			break
		}
	}
	if !included {
		memberIDs = append(memberIDs, ownerID)
	}
	if len(memberIDs) < 3 { // owner + at least 2 others as requirement
		return nil, errors.New("at least 3 members including owner")
	}
	g := &chatentity.Group{Name: name, OwnerID: ownerID, Avatar: avatar}
	if err := s.repo.Create(g, memberIDs); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *groupServiceImpl) Search(name string, limit int) ([]*chatentity.Group, error) {
	return s.repo.SearchByName(name, limit)
}
func (s *groupServiceImpl) Join(groupID, userID uint) error {
	ok, err := s.repo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return s.repo.AddMembers(groupID, []uint{userID})
}

func (s *groupServiceImpl) ListMessages(groupID uint, page, pageSize int) ([]*chatentity.GroupMessage, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	// use repository impl via type assertion (quick hack)
	type msgRepo interface {
		ListMessages(groupID uint, offset, limit int) ([]*chatentity.GroupMessage, int64, error)
	}
	if mr, ok := s.repo.(msgRepo); ok {
		return mr.ListMessages(groupID, offset, pageSize)
	}
	return nil, 0, errors.New("messages not supported")
}

func (s *groupServiceImpl) SendMessage(groupID, senderID uint, msgType, content string) (*chatentity.GroupMessage, error) {
	if groupID == 0 || senderID == 0 || content == "" {
		return nil, errors.New("invalid params")
	}
	ok, err := s.repo.IsMember(groupID, senderID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not a member")
	}
	m := &chatentity.GroupMessage{GroupID: groupID, SenderID: senderID, Type: firstNonEmpty(msgType, "text"), Content: content, CreatedAt: time.Now()}
	type msgRepo interface {
		SaveMessage(m *chatentity.GroupMessage) error
	}
	if mr, ok := s.repo.(msgRepo); ok {
		if err := mr.SaveMessage(m); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("save not supported")
	}
	return m, nil
}

func (s *groupServiceImpl) IsMember(groupID, userID uint) (bool, error) {
	return s.repo.IsMember(groupID, userID)
}
func (s *groupServiceImpl) Get(groupID uint) (*chatentity.Group, error) { return s.repo.Get(groupID) }
func (s *groupServiceImpl) UpdateGroup(operatorID, groupID uint, name, avatar string) (*chatentity.Group, error) {
	g, err := s.repo.Get(groupID)
	if err != nil {
		return nil, err
	}
	if g.OwnerID != operatorID {
		return nil, errors.New("no permission")
	}
	changed := false
	if name = strings.TrimSpace(name); name != "" && name != g.Name {
		g.Name = name
		changed = true
	}
	if avatar != "" && avatar != g.Avatar {
		g.Avatar = avatar
		changed = true
	}
	if !changed {
		return g, nil
	}
	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}
func (s *groupServiceImpl) ListUserGroups(userID uint, page, pageSize int) ([]*chatentity.Group, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListUserGroups(userID, offset, pageSize)
}

func (s *groupServiceImpl) UpdateLastRead(ctx context.Context, groupID, userID, msgID uint) error {
	return s.repo.UpdateLastRead(groupID, userID, msgID)
}

func (s *groupServiceImpl) CountUnread(ctx context.Context, groupID, userID uint) (int64, error) {
	return s.repo.CountUnread(groupID, userID)
}

func (s *groupServiceImpl) ListMemberIDs(groupID uint) ([]uint, error) {
	return s.repo.ListMemberIDs(groupID)
}

func (s *groupServiceImpl) AddMembers(operatorID, groupID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}
	g, err := s.repo.Get(groupID)
	if err != nil {
		return err
	}
	if g.OwnerID != operatorID {
		return errors.New("no permission")
	}
	return s.repo.AddMembers(groupID, userIDs)
}

func (s *groupServiceImpl) RemoveMember(operatorID, groupID, targetUserID uint) error {
	g, err := s.repo.Get(groupID)
	if err != nil {
		return err
	}
	if g.OwnerID != operatorID {
		return errors.New("no permission")
	}
	if g.OwnerID == targetUserID {
		return errors.New("cannot remove owner")
	}
	return s.repo.RemoveMember(groupID, targetUserID)
}

func (s *groupServiceImpl) ListMembers(groupID uint) ([]uint, error) {
	return s.repo.ListMemberIDs(groupID)
}
