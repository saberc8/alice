package chat

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	chatentity "alice/domain/chat/entity"
	chatrepo "alice/domain/chat/repository"
)

type groupRepositoryImpl struct{ db *gorm.DB }

func NewGroupRepository(db *gorm.DB) chatrepo.GroupRepository { return &groupRepositoryImpl{db: db} }

func (r *groupRepositoryImpl) Create(g *chatentity.Group, memberIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(g).Error; err != nil {
			return err
		}
		if len(memberIDs) == 0 {
			return nil
		}
		ms := make([]*chatentity.GroupMember, 0, len(memberIDs))
		now := time.Now()
		for _, id := range memberIDs {
			role := "member"
			if id == g.OwnerID {
				role = "owner"
			}
			ms = append(ms, &chatentity.GroupMember{GroupID: g.ID, UserID: id, Role: role, JoinedAt: now})
		}
		return tx.Create(&ms).Error
	})
}

func (r *groupRepositoryImpl) AddMembers(groupID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}
	ms := make([]*chatentity.GroupMember, 0, len(userIDs))
	now := time.Now()
	for _, id := range userIDs {
		ms = append(ms, &chatentity.GroupMember{GroupID: groupID, UserID: id, Role: "member", JoinedAt: now})
	}
	return r.db.Clauses().Create(&ms).Error
}

func (r *groupRepositoryImpl) ListUserGroups(userID uint, offset, limit int) ([]*chatentity.Group, int64, error) {
	var total int64
	q := r.db.Model(&chatentity.Group{}).Joins("JOIN app_chat_group_members gm ON gm.group_id = app_chat_groups.id").Where("gm.user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []*chatentity.Group
	if err := q.Order("app_chat_groups.id DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *groupRepositoryImpl) Get(id uint) (*chatentity.Group, error) {
	var g chatentity.Group
	if err := r.db.First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *groupRepositoryImpl) Update(g *chatentity.Group) error {
	return r.db.Model(&chatentity.Group{}).Where("id=?", g.ID).Updates(map[string]any{
		"name":   g.Name,
		"avatar": g.Avatar,
	}).Error
}

func (r *groupRepositoryImpl) SearchByName(qs string, limit int) ([]*chatentity.Group, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	qs = strings.TrimSpace(qs)
	if qs == "" {
		return []*chatentity.Group{}, nil
	}
	var rows []*chatentity.Group
	if err := r.db.Where("name ILIKE ?", "%"+qs+"%").Order("id DESC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *groupRepositoryImpl) IsMember(groupID, userID uint) (bool, error) {
	var cnt int64
	if err := r.db.Model(&chatentity.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *groupRepositoryImpl) ListMemberIDs(groupID uint) ([]uint, error) {
	var members []chatentity.GroupMember
	if err := r.db.Select("user_id").Where("group_id=?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.UserID)
	}
	return ids, nil
}

func (r *groupRepositoryImpl) GetLastRead(groupID, userID uint) (uint, error) {
	var cursor chatentity.GroupReadCursor
	if err := r.db.Where("group_id=? AND user_id=?", groupID, userID).First(&cursor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return cursor.LastReadMsgID, nil
}

func (r *groupRepositoryImpl) UpdateLastRead(groupID, userID, msgID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var cursor chatentity.GroupReadCursor
		if err := tx.Where("group_id=? AND user_id=?", groupID, userID).First(&cursor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				cursor.GroupID = groupID
				cursor.UserID = userID
				cursor.LastReadMsgID = msgID
				return tx.Create(&cursor).Error
			}
			return err
		}
		if cursor.LastReadMsgID < msgID { // only move forward
			return tx.Model(&cursor).Update("last_read_msg_id", msgID).Error
		}
		return nil
	})
}

func (r *groupRepositoryImpl) CountUnread(groupID, userID uint) (int64, error) {
	lastRead, err := r.GetLastRead(groupID, userID)
	if err != nil {
		return 0, err
	}
	var cnt int64
	if err := r.db.Model(&chatentity.GroupMessage{}).Where("group_id=? AND id > ?", groupID, lastRead).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *groupRepositoryImpl) RemoveMember(groupID, userID uint) error {
	return r.db.Where("group_id=? AND user_id=?", groupID, userID).Delete(&chatentity.GroupMember{}).Error
}

// Helpers for group messages (inline here for brevity)
func (r *groupRepositoryImpl) SaveMessage(m *chatentity.GroupMessage) error {
	return r.db.Create(m).Error
}
func (r *groupRepositoryImpl) ListMessages(groupID uint, offset, limit int) ([]*chatentity.GroupMessage, int64, error) {
	var total int64
	q := r.db.Model(&chatentity.GroupMessage{}).Where("group_id = ?", groupID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []*chatentity.GroupMessage
	if err := q.Order("id DESC").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// Expose additional interface via type assertion in service (quick approach). In production you'd split repos.
var ErrNotOwner = errors.New("not owner")
