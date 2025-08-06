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

package service

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"alice/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MenuService 菜单服务接口
type MenuService interface {
	// CreateMenu 创建菜单
	CreateMenu(ctx context.Context, req *CreateMenuRequest) (*entity.Menu, error)

	// GetMenu 获取菜单
	GetMenu(ctx context.Context, id string) (*entity.Menu, error)

	// ListMenus 获取菜单列表
	ListMenus(ctx context.Context) ([]*entity.Menu, error)

	// GetMenuTree 获取菜单树
	GetMenuTree(ctx context.Context) ([]*entity.Menu, error)

	// UpdateMenu 更新菜单
	UpdateMenu(ctx context.Context, req *UpdateMenuRequest) error

	// DeleteMenu 删除菜单
	DeleteMenu(ctx context.Context, id string) error

	// AssignMenusToRole 为角色分配菜单
	AssignMenusToRole(ctx context.Context, roleID string, menuIDs []string) error

	// RemoveMenusFromRole 移除角色菜单
	RemoveMenusFromRole(ctx context.Context, roleID string, menuIDs []string) error

	// GetRoleMenus 获取角色菜单
	GetRoleMenus(ctx context.Context, roleID string) ([]*entity.Menu, error)

	// GetUserMenus 获取用户菜单
	GetUserMenus(ctx context.Context, userID string) ([]*entity.Menu, error)

	// GetUserMenuTree 获取用户菜单树
	GetUserMenuTree(ctx context.Context, userID string) ([]*entity.Menu, error)
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID    *string           `json:"parent_id,omitempty"`
	Name        string            `json:"name" validate:"required,max=100"`
	Code        string            `json:"code" validate:"required,max=100"`
	Path        *string           `json:"path,omitempty" validate:"omitempty,max=200"`
	Type        entity.MenuType   `json:"type" validate:"required"`
	Order       int               `json:"order"`
	Status      entity.MenuStatus `json:"status,omitempty"`
	Meta        entity.MenuMeta   `json:"meta,omitempty"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ID          string            `json:"id" validate:"required"`
	ParentID    *string           `json:"parent_id,omitempty"`
	Name        string            `json:"name" validate:"required,max=100"`
	Code        string            `json:"code" validate:"required,max=100"`
	Path        *string           `json:"path,omitempty" validate:"omitempty,max=200"`
	Type        entity.MenuType   `json:"type" validate:"required"`
	Order       int               `json:"order"`
	Status      entity.MenuStatus `json:"status,omitempty"`
	Meta        entity.MenuMeta   `json:"meta,omitempty"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
}

// menuService 菜单服务实现
type menuService struct {
	menuRepo repository.MenuRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

// CreateMenu 创建菜单
func (s *menuService) CreateMenu(ctx context.Context, req *CreateMenuRequest) (*entity.Menu, error) {
	// 检查代码是否已存在
	existing, _ := s.menuRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("菜单代码 %s 已存在", req.Code)
	}

	// 检查父菜单是否存在
	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.menuRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("获取父菜单失败: %w", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("父菜单不存在")
		}
	}

	// 创建菜单实体
	menu := &entity.Menu{
		ID:          uuid.New().String(),
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Path:        req.Path,
		Type:        req.Type,
		Order:       req.Order,
		Status:      req.Status,
		Meta:        req.Meta,
		Description: req.Description,
	}

	if menu.Status == "" {
		menu.Status = entity.MenuStatusActive
	}

	// 保存到数据库
	if err := s.menuRepo.Create(ctx, menu); err != nil {
		logger.Errorf("创建菜单失败: %v", err)
		return nil, fmt.Errorf("创建菜单失败: %w", err)
	}

	return menu, nil
}

// GetMenu 获取菜单
func (s *menuService) GetMenu(ctx context.Context, id string) (*entity.Menu, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取菜单失败: %v", err)
		return nil, fmt.Errorf("获取菜单失败: %w", err)
	}

	if menu == nil {
		return nil, fmt.Errorf("菜单不存在")
	}

	return menu, nil
}

// ListMenus 获取菜单列表
func (s *menuService) ListMenus(ctx context.Context) ([]*entity.Menu, error) {
	menus, err := s.menuRepo.List(ctx)
	if err != nil {
		logger.Errorf("获取菜单列表失败: %v", err)
		return nil, fmt.Errorf("获取菜单列表失败: %w", err)
	}

	return menus, nil
}

// GetMenuTree 获取菜单树
func (s *menuService) GetMenuTree(ctx context.Context) ([]*entity.Menu, error) {
	tree, err := s.menuRepo.GetTree(ctx)
	if err != nil {
		logger.Errorf("获取菜单树失败: %v", err)
		return nil, fmt.Errorf("获取菜单树失败: %w", err)
	}

	return tree, nil
}

// UpdateMenu 更新菜单
func (s *menuService) UpdateMenu(ctx context.Context, req *UpdateMenuRequest) error {
	// 检查菜单是否存在
	existing, err := s.menuRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.Errorf("获取菜单失败: %v", err)
		return fmt.Errorf("获取菜单失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("菜单不存在")
	}

	// 检查代码是否被其他菜单使用
	if existing.Code != req.Code {
		codeExists, _ := s.menuRepo.GetByCode(ctx, req.Code)
		if codeExists != nil && codeExists.ID != req.ID {
			return fmt.Errorf("菜单代码 %s 已被其他菜单使用", req.Code)
		}
	}

	// 检查父菜单是否存在
	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.menuRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return fmt.Errorf("获取父菜单失败: %w", err)
		}
		if parent == nil {
			return fmt.Errorf("父菜单不存在")
		}

		// 检查是否为循环引用
		if *req.ParentID == req.ID {
			return fmt.Errorf("不能将自己设置为父菜单")
		}
	}

	// 更新菜单信息
	existing.ParentID = req.ParentID
	existing.Name = req.Name
	existing.Code = req.Code
	existing.Path = req.Path
	existing.Type = req.Type
	existing.Order = req.Order
	existing.Status = req.Status
	existing.Meta = req.Meta
	existing.Description = req.Description

	if err := s.menuRepo.Update(ctx, existing); err != nil {
		logger.Errorf("更新菜单失败: %v", err)
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	return nil
}

// DeleteMenu 删除菜单
func (s *menuService) DeleteMenu(ctx context.Context, id string) error {
	// 检查菜单是否存在
	existing, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取菜单失败: %v", err)
		return fmt.Errorf("获取菜单失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("菜单不存在")
	}

	// 检查是否有子菜单
	children, err := s.menuRepo.GetChildren(ctx, id)
	if err != nil {
		logger.Errorf("获取子菜单失败: %v", err)
		return fmt.Errorf("获取子菜单失败: %w", err)
	}

	if len(children) > 0 {
		return fmt.Errorf("存在子菜单，无法删除")
	}

	if err := s.menuRepo.Delete(ctx, id); err != nil {
		logger.Errorf("删除菜单失败: %v", err)
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	return nil
}

// AssignMenusToRole 为角色分配菜单
func (s *menuService) AssignMenusToRole(ctx context.Context, roleID string, menuIDs []string) error {
	if err := s.menuRepo.AssignToRole(ctx, roleID, menuIDs); err != nil {
		logger.Errorf("为角色分配菜单失败: %v", err)
		return fmt.Errorf("为角色分配菜单失败: %w", err)
	}

	return nil
}

// RemoveMenusFromRole 移除角色菜单
func (s *menuService) RemoveMenusFromRole(ctx context.Context, roleID string, menuIDs []string) error {
	if err := s.menuRepo.RemoveFromRole(ctx, roleID, menuIDs); err != nil {
		logger.Errorf("移除角色菜单失败: %v", err)
		return fmt.Errorf("移除角色菜单失败: %w", err)
	}

	return nil
}

// GetRoleMenus 获取角色菜单
func (s *menuService) GetRoleMenus(ctx context.Context, roleID string) ([]*entity.Menu, error) {
	menus, err := s.menuRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		logger.Errorf("获取角色菜单失败: %v", err)
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}

	return menus, nil
}

// GetUserMenus 获取用户菜单
func (s *menuService) GetUserMenus(ctx context.Context, userID string) ([]*entity.Menu, error) {
	menus, err := s.menuRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户菜单失败: %v", err)
		return nil, fmt.Errorf("获取用户菜单失败: %w", err)
	}

	return menus, nil
}

// GetUserMenuTree 获取用户菜单树
func (s *menuService) GetUserMenuTree(ctx context.Context, userID string) ([]*entity.Menu, error) {
	tree, err := s.menuRepo.GetTreeByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户菜单树失败: %v", err)
		return nil, fmt.Errorf("获取用户菜单树失败: %w", err)
	}

	return tree, nil
}
