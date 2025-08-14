package repository

import (
	appentity "alice/domain/appuser/entity"
)

type AppUserRepository interface {
	Create(user *appentity.AppUser) error
	GetByID(id uint) (*appentity.AppUser, error)
	GetByEmail(email string) (*appentity.AppUser, error)
	GetByIDs(ids []uint) ([]*appentity.AppUser, error)
	Update(user *appentity.AppUser) error
	Delete(id uint) error
	List(offset, limit int) ([]*appentity.AppUser, int64, error)
}
