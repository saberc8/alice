package repository

import (
	"gorm.io/gorm"

	friendentity "alice/domain/appfriend/entity"
	friendrepo "alice/domain/appfriend/repository"
)

type friendRepositoryImpl struct{ db *gorm.DB }

func NewFriendRepository(db *gorm.DB) friendrepo.FriendRepository {
	return &friendRepositoryImpl{db: db}
}

func (r *friendRepositoryImpl) AddRelation(userID, friendID uint) error {
	if userID == friendID {
		return nil
	}
	// ensure unique pair (userID, friendID)
	rel := &friendentity.FriendRelation{UserID: userID, FriendID: friendID}
	// ignore duplicate
	return r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).FirstOrCreate(rel).Error
}

func (r *friendRepositoryImpl) RemoveRelation(userID, friendID uint) error {
	return r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&friendentity.FriendRelation{}).Error
}

func (r *friendRepositoryImpl) ListFriends(userID uint, offset, limit int) ([]uint, int64, error) {
	var ids []uint
	var total int64
	q := r.db.Model(&friendentity.FriendRelation{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []friendentity.FriendRelation
	if err := q.Offset(offset).Limit(limit).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	ids = make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.FriendID)
	}
	return ids, total, nil
}

// Friend Requests
func (r *friendRepositoryImpl) CreateRequest(requesterID, addresseeID uint) error {
	if requesterID == addresseeID {
		return nil
	}
	req := &friendentity.FriendRequest{RequesterID: requesterID, AddresseeID: addresseeID, Status: friendentity.FriendRequestPending}
	// one active pending per requester->addressee
	return r.db.Where("requester_id = ? AND addressee_id = ? AND status = ?", requesterID, addresseeID, friendentity.FriendRequestPending).FirstOrCreate(req).Error
}

func (r *friendRepositoryImpl) AcceptRequest(requestID uint) (uint, uint, error) {
	var req friendentity.FriendRequest
	if err := r.db.First(&req, requestID).Error; err != nil {
		return 0, 0, err
	}
	req.Status = friendentity.FriendRequestAccepted
	if err := r.db.Save(&req).Error; err != nil {
		return 0, 0, err
	}
	return req.RequesterID, req.AddresseeID, nil
}

func (r *friendRepositoryImpl) DeclineRequest(requestID uint) error {
	return r.db.Model(&friendentity.FriendRequest{}).Where("id = ?", requestID).Update("status", friendentity.FriendRequestDeclined).Error
}

// AreFriends 检查是否互为好友（需要同时存在 A->B 与 B->A 记录）
func (r *friendRepositoryImpl) AreFriends(a, b uint) (bool, error) {
	if a == 0 || b == 0 || a == b {
		return false, nil
	}
	var cnt int64
	if err := r.db.Model(&friendentity.FriendRelation{}).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", a, b, b, a).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt >= 2, nil
}

func (r *friendRepositoryImpl) GetPendingRequests(addresseeID uint, offset, limit int) ([]uint, []uint, int64, error) {
	var list []friendentity.FriendRequest
	var total int64
	q := r.db.Model(&friendentity.FriendRequest{}).Where("addressee_id = ? AND status = ?", addresseeID, friendentity.FriendRequestPending)
	if err := q.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("id DESC").Find(&list).Error; err != nil {
		return nil, nil, 0, err
	}
	reqIDs := make([]uint, 0, len(list))
	requesterIDs := make([]uint, 0, len(list))
	for _, r := range list {
		reqIDs = append(reqIDs, r.ID)
		requesterIDs = append(requesterIDs, r.RequesterID)
	}
	return reqIDs, requesterIDs, total, nil
}
