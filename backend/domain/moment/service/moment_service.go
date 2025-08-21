package service

import (
	momententity "alice/domain/moment/entity"
	momentrepo "alice/domain/moment/repository"
	"errors"
	"strings"
)

type MomentService interface {
	Publish(userID uint, content string, images []string) (*momententity.Moment, error)
	ListAll(page, pageSize int) ([]*momententity.Moment, int64, error)
	ListByUser(userID uint, page, pageSize int) ([]*momententity.Moment, int64, error)
	Delete(userID uint, id uint) error
}

type momentServiceImpl struct{ repo momentrepo.MomentRepository }

func NewMomentService(repo momentrepo.MomentRepository) MomentService {
	return &momentServiceImpl{repo: repo}
}

func (s *momentServiceImpl) Publish(userID uint, content string, images []string) (*momententity.Moment, error) {
	if userID == 0 || strings.TrimSpace(content) == "" {
		return nil, errors.New("invalid params")
	}
	if len(images) > 9 {
		images = images[:9]
	}
	// 过滤空串
	filtered := make([]string, 0, len(images))
	for _, img := range images {
		img = strings.TrimSpace(img)
		if img != "" {
			filtered = append(filtered, img)
		}
	}
	m := &momententity.Moment{UserID: userID, Content: content, Images: strings.Join(filtered, ",")}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func normPage(page, pageSize int) (int, int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return page, pageSize, offset
}

func (s *momentServiceImpl) ListAll(page, pageSize int) ([]*momententity.Moment, int64, error) {
	page, pageSize, offset := normPage(page, pageSize)
	return s.repo.ListAll(offset, pageSize)
}

func (s *momentServiceImpl) ListByUser(userID uint, page, pageSize int) ([]*momententity.Moment, int64, error) {
	page, pageSize, offset := normPage(page, pageSize)
	return s.repo.ListByUser(userID, offset, pageSize)
}

func (s *momentServiceImpl) Delete(userID uint, id uint) error {
	if userID == 0 || id == 0 {
		return errors.New("invalid params")
	}
	// 简单直接删除（条件包含 user_id 保证只能删自己）
	return s.repo.Delete(id, userID)
}
