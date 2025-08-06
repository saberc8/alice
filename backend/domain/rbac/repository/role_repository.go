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
	"alice/domain/rbac/entity"
	"context"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error

	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, id string) (*entity.Role, error)

	// GetByCode 根据代码获取角色
	GetByCode(ctx context.Context, code string) (*entity.Role, error)

	// List 获取角色列表
	List(ctx context.Context, offset, limit int) ([]*entity.Role, int64, error)

	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error

	// Delete 删除角色
	Delete(ctx context.Context, id string) error

	// GetByUserID 根据用户ID获取角色列表
	GetByUserID(ctx context.Context, userID string) ([]*entity.Role, error)

	// AssignToUser 为用户分配角色
	AssignToUser(ctx context.Context, userID string, roleIDs []string) error

	// RemoveFromUser 移除用户角色
	RemoveFromUser(ctx context.Context, userID string, roleIDs []string) error
}
