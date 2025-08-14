package repository

import (
	"gorm.io/gorm"

	"alice/domain/user/entity"
	"alice/domain/user/repository"
)

// userRepositoryImpl 用户仓储实现
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create 创建用户
func (r *userRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepositoryImpl) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepositoryImpl) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepositoryImpl) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

// List 分页获取用户列表
func (r *userRepositoryImpl) List(offset, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	if err := r.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}
