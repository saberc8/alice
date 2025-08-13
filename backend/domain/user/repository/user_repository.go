package repository

import (
	"alice/domain/user/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(user *entity.User) error

	// GetByID 根据ID获取用户
	GetByID(id uint) (*entity.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*entity.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(email string) (*entity.User, error)

	// Update 更新用户
	Update(user *entity.User) error

	// Delete 删除用户
	Delete(id uint) error
}
