package repository

import (
	chatentity "alice/domain/chat/entity"
)

type GroupRepository interface {
	Create(g *chatentity.Group, memberIDs []uint) error
	AddMembers(groupID uint, userIDs []uint) error
	ListUserGroups(userID uint, offset, limit int) ([]*chatentity.Group, int64, error)
	Get(id uint) (*chatentity.Group, error)
	Update(g *chatentity.Group) error
	SearchByName(q string, limit int) ([]*chatentity.Group, error)
	IsMember(groupID, userID uint) (bool, error)
	ListMemberIDs(groupID uint) ([]uint, error)
	GetLastRead(groupID, userID uint) (uint, error)
	UpdateLastRead(groupID, userID, msgID uint) error
	CountUnread(groupID, userID uint) (int64, error)
	RemoveMember(groupID, userID uint) error
}
