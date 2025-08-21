package repository

import (
	momententity "alice/domain/moment/entity"
	momentrepo "alice/domain/moment/repository"

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
