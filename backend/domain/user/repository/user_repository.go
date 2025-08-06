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
