/*
 * Copyright 2025 alice Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
