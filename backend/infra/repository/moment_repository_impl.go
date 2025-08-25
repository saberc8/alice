package repository

import (
	momententity "alice/domain/moment/entity"
	momentrepo "alice/domain/moment/repository"
	"errors"

	"gorm.io/gorm"
)

type momentRepositoryImpl struct{ db *gorm.DB }

func NewMomentRepository(db *gorm.DB) momentrepo.MomentRepository {
	return &momentRepositoryImpl{db: db}
}

func (r *momentRepositoryImpl) Create(m *momententity.Moment) error { return r.db.Create(m).Error }

func (r *momentRepositoryImpl) ListAll(offset, limit int) ([]*momententity.Moment, int64, error) {
	var list []*momententity.Moment
	var total int64
	if err := r.db.Model(&momententity.Moment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *momentRepositoryImpl) ListByUser(userID uint, offset, limit int) ([]*momententity.Moment, int64, error) {
	var list []*momententity.Moment
	var total int64
	if err := r.db.Model(&momententity.Moment{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("user_id = ?", userID).Order("id DESC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *momentRepositoryImpl) Get(id uint) (*momententity.Moment, error) {
	var m momententity.Moment
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *momentRepositoryImpl) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&momententity.Moment{}).Error
}

// Likes
func (r *momentRepositoryImpl) AddLike(momentID, userID uint) error {
	like := &momententity.MomentLike{MomentID: momentID, UserID: userID}
	return r.db.FirstOrCreate(like, like).Error
}

func (r *momentRepositoryImpl) RemoveLike(momentID, userID uint) error {
	return r.db.Where("moment_id = ? AND user_id = ?", momentID, userID).Delete(&momententity.MomentLike{}).Error
}

func (r *momentRepositoryImpl) HasLiked(momentID, userID uint) (bool, error) {
	var cnt int64
	if err := r.db.Model(&momententity.MomentLike{}).Where("moment_id = ? AND user_id = ?", momentID, userID).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (r *momentRepositoryImpl) CountLikes(momentID uint) (int64, error) {
	var cnt int64
	if err := r.db.Model(&momententity.MomentLike{}).Where("moment_id = ?", momentID).Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

// Comments
func (r *momentRepositoryImpl) AddComment(cmt *momententity.MomentComment) error {
	if cmt.MomentID == 0 || cmt.UserID == 0 || cmt.Content == "" {
		return errors.New("invalid params")
	}
	return r.db.Create(cmt).Error
}

func (r *momentRepositoryImpl) ListComments(momentID uint, offset, limit int) ([]*momententity.MomentComment, int64, error) {
	var list []*momententity.MomentComment
	var total int64
	if err := r.db.Model(&momententity.MomentComment{}).Where("moment_id = ?", momentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("moment_id = ?", momentID).Order("id ASC").Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
