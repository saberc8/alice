package service

import (
	"errors"
	"strings"

	friendrepo "alice/domain/appfriend/repository"
	appentity "alice/domain/appuser/entity"
	apprepo "alice/domain/appuser/repository"
)

var (
	ErrFriendUserNotFound = errors.New("user not found")
)

type FriendService interface {
	RequestFriend(userID uint, friendEmail string) error
	AcceptRequest(userID uint, requestID uint) error
	DeclineRequest(userID uint, requestID uint) error
	ListPending(userID uint, page, pageSize int) ([]uint, []uint, int64, error)
	RemoveFriend(userID uint, friendID uint) error
	ListFriendIDs(userID uint, page, pageSize int) ([]uint, int64, error)
	ListFriendDetails(userID uint, page, pageSize int) ([]*appentity.AppUser, int64, error)
}

type friendServiceImpl struct {
	appUserRepo apprepo.AppUserRepository
	repo        friendrepo.FriendRepository
}

func NewFriendService(appUserRepo apprepo.AppUserRepository, repo friendrepo.FriendRepository) FriendService {
	return &friendServiceImpl{appUserRepo: appUserRepo, repo: repo}
}

func (s *friendServiceImpl) RequestFriend(userID uint, friendEmail string) error {
	email := strings.ToLower(strings.TrimSpace(friendEmail))
	f, err := s.appUserRepo.GetByEmail(email)
	if err != nil || f == nil {
		return ErrFriendUserNotFound
	}
	return s.repo.CreateRequest(userID, f.ID)
}

func (s *friendServiceImpl) AcceptRequest(userID uint, requestID uint) error {
	requesterID, addresseeID, err := s.repo.AcceptRequest(requestID)
	if err != nil {
		return err
	}
	if addresseeID != userID {
		return errors.New("permission denied")
	}
	if err := s.repo.AddRelation(requesterID, addresseeID); err != nil {
		return err
	}
	if err := s.repo.AddRelation(addresseeID, requesterID); err != nil {
		return err
	}
	return nil
}

func (s *friendServiceImpl) DeclineRequest(userID uint, requestID uint) error {
	return s.repo.DeclineRequest(requestID)
}

func (s *friendServiceImpl) ListPending(userID uint, page, pageSize int) ([]uint, []uint, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.GetPendingRequests(userID, offset, pageSize)
}

func (s *friendServiceImpl) RemoveFriend(userID uint, friendID uint) error {
	if err := s.repo.RemoveRelation(userID, friendID); err != nil {
		return err
	}
	if err := s.repo.RemoveRelation(friendID, userID); err != nil {
		return err
	}
	return nil
}

func (s *friendServiceImpl) ListFriendIDs(userID uint, page, pageSize int) ([]uint, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.ListFriends(userID, offset, pageSize)
}

// ListFriendDetails 先取 ID 再批量查询资料
func (s *friendServiceImpl) ListFriendDetails(userID uint, page, pageSize int) ([]*appentity.AppUser, int64, error) {
	ids, total, err := s.ListFriendIDs(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	users, err := s.appUserRepo.GetByIDs(ids)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
