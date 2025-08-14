package repository

import (
	"gorm.io/gorm"

	appentity "alice/domain/appuser/entity"
	apprepo "alice/domain/appuser/repository"
)

type appUserRepositoryImpl struct{ db *gorm.DB }

func NewAppUserRepository(db *gorm.DB) apprepo.AppUserRepository {
	return &appUserRepositoryImpl{db: db}
}

func (r *appUserRepositoryImpl) Create(user *appentity.AppUser) error { return r.db.Create(user).Error }

func (r *appUserRepositoryImpl) GetByID(id uint) (*appentity.AppUser, error) {
	var u appentity.AppUser
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *appUserRepositoryImpl) GetByEmail(email string) (*appentity.AppUser, error) {
	var u appentity.AppUser
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *appUserRepositoryImpl) GetByIDs(ids []uint) ([]*appentity.AppUser, error) {
	if len(ids) == 0 {
		return []*appentity.AppUser{}, nil
	}
	var list []*appentity.AppUser
	if err := r.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *appUserRepositoryImpl) Update(user *appentity.AppUser) error { return r.db.Save(user).Error }

func (r *appUserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&appentity.AppUser{}, id).Error
}

func (r *appUserRepositoryImpl) List(offset, limit int) ([]*appentity.AppUser, int64, error) {
	var list []*appentity.AppUser
	var total int64
	if err := r.db.Model(&appentity.AppUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
