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

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *entity.Menu) error

	// GetByID 根据ID获取菜单
	GetByID(ctx context.Context, id string) (*entity.Menu, error)

	// GetByCode 根据代码获取菜单
	GetByCode(ctx context.Context, code string) (*entity.Menu, error)

	// List 获取菜单列表
	List(ctx context.Context) ([]*entity.Menu, error)

	// GetTree 获取菜单树
	GetTree(ctx context.Context) ([]*entity.Menu, error)

	// Update 更新菜单
	Update(ctx context.Context, menu *entity.Menu) error

	// Delete 删除菜单
	Delete(ctx context.Context, id string) error

	// GetByUserID 根据用户ID获取菜单列表
	GetByUserID(ctx context.Context, userID string) ([]*entity.Menu, error)

	// GetTreeByUserID 根据用户ID获取菜单树
	GetTreeByUserID(ctx context.Context, userID string) ([]*entity.Menu, error)

	// GetByRoleID 根据角色ID获取菜单列表
	GetByRoleID(ctx context.Context, roleID string) ([]*entity.Menu, error)

	// AssignToRole 为角色分配菜单
	AssignToRole(ctx context.Context, roleID string, menuIDs []string) error

	// RemoveFromRole 移除角色菜单
	RemoveFromRole(ctx context.Context, roleID string, menuIDs []string) error

	// GetChildren 获取子菜单
	GetChildren(ctx context.Context, parentID string) ([]*entity.Menu, error)
}
