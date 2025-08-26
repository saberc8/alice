package service

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"alice/pkg/logger"
	"context"
	"fmt"
)

// MenuService 菜单服务接口
type MenuService interface {
	// CreateMenu 创建菜单
	CreateMenu(ctx context.Context, req *CreateMenuRequest) (*entity.Menu, error)

	// GetMenu 获取菜单
	GetMenu(ctx context.Context, id uint) (*entity.Menu, error)

	// ListMenus 获取菜单列表
	ListMenus(ctx context.Context) ([]*entity.Menu, error)

	// GetMenuTree 获取菜单树
	GetMenuTree(ctx context.Context) ([]*entity.Menu, error)

	// UpdateMenu 更新菜单
	UpdateMenu(ctx context.Context, req *UpdateMenuRequest) error

	// DeleteMenu 删除菜单
	DeleteMenu(ctx context.Context, id uint) error

	// AssignMenusToRole 为角色分配菜单
	AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error

	// RemoveMenusFromRole 移除角色菜单
	RemoveMenusFromRole(ctx context.Context, roleID uint, menuIDs []uint) error

	// GetRoleMenus 获取角色菜单
	GetRoleMenus(ctx context.Context, roleID uint) ([]*entity.Menu, error)

	// GetUserMenus 获取用户菜单
	GetUserMenus(ctx context.Context, userID uint) ([]*entity.Menu, error)

	// GetUserMenuTree 获取用户菜单树
	GetUserMenuTree(ctx context.Context, userID uint) ([]*entity.Menu, error)

	// GetRoleMenuTree 获取角色菜单树（按角色注入 perms）
	GetRoleMenuTree(ctx context.Context, roleID uint) ([]*entity.Menu, error)
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID    *uint             `json:"parent_id,omitempty"`
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
	ID          uint              `json:"id" validate:"required"`
	ParentID    *uint             `json:"parent_id,omitempty"`
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
	menuRepo       repository.MenuRepository
	permissionRepo repository.PermissionRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo repository.MenuRepository, permissionRepo repository.PermissionRepository) MenuService {
	return &menuService{
		menuRepo:       menuRepo,
		permissionRepo: permissionRepo,
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
	if req.ParentID != nil && *req.ParentID != 0 {
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
func (s *menuService) GetMenu(ctx context.Context, id uint) (*entity.Menu, error) {
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
	// 注入菜单下的权限集合到 meta.perms（全量，不区分角色/用户）
	var zeroUserID uint = 0
	s.attachPermsToMenus(ctx, tree, zeroUserID)
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
	if req.ParentID != nil && *req.ParentID != 0 {
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
func (s *menuService) DeleteMenu(ctx context.Context, id uint) error {
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

	// 先删除挂载在该菜单下的按钮权限（及其角色关联由仓储层的外键/业务逻辑负责）
	// 这里直接调用权限仓储按菜单ID删除（若无方法，则让数据库层通过外键/触发器处理）。
	// 简化处理：查出权限并逐一删除，复用现有 Delete 逻辑。
	if perms, err := s.permissionRepo.GetByMenuIDs(ctx, []uint{id}); err == nil {
		for _, p := range perms {
			_ = s.permissionRepo.Delete(ctx, p.ID)
		}
	}

	if err := s.menuRepo.Delete(ctx, id); err != nil {
		logger.Errorf("删除菜单失败: %v", err)
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	return nil
}

// AssignMenusToRole 为角色分配菜单
func (s *menuService) AssignMenusToRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	if err := s.menuRepo.AssignToRole(ctx, roleID, menuIDs); err != nil {
		logger.Errorf("为角色分配菜单失败: %v", err)
		return fmt.Errorf("为角色分配菜单失败: %w", err)
	}

	return nil
}

// RemoveMenusFromRole 移除角色菜单
func (s *menuService) RemoveMenusFromRole(ctx context.Context, roleID uint, menuIDs []uint) error {
	if err := s.menuRepo.RemoveFromRole(ctx, roleID, menuIDs); err != nil {
		logger.Errorf("移除角色菜单失败: %v", err)
		return fmt.Errorf("移除角色菜单失败: %w", err)
	}

	return nil
}

// GetRoleMenus 获取角色菜单
func (s *menuService) GetRoleMenus(ctx context.Context, roleID uint) ([]*entity.Menu, error) {
	menus, err := s.menuRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		logger.Errorf("获取角色菜单失败: %v", err)
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}

	return menus, nil
}

// GetUserMenus 获取用户菜单
func (s *menuService) GetUserMenus(ctx context.Context, userID uint) ([]*entity.Menu, error) {
	menus, err := s.menuRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户菜单失败: %v", err)
		return nil, fmt.Errorf("获取用户菜单失败: %w", err)
	}

	return menus, nil
}

// GetUserMenuTree 获取用户菜单树
func (s *menuService) GetUserMenuTree(ctx context.Context, userID uint) ([]*entity.Menu, error) {
	tree, err := s.menuRepo.GetTreeByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户菜单树失败: %v", err)
		return nil, fmt.Errorf("获取用户菜单树失败: %w", err)
	}
	// 注入用户在各菜单下的权限集合到 meta.perms
	s.attachPermsToMenus(ctx, tree, userID)
	return tree, nil
}

// GetRoleMenuTree 获取角色菜单树（按角色注入 perms）
func (s *menuService) GetRoleMenuTree(ctx context.Context, roleID uint) ([]*entity.Menu, error) {
	// 角色菜单列表（不筛用户）
	menus, err := s.menuRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		logger.Errorf("获取角色菜单失败: %v", err)
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}
	// 构建树
	tree := s.buildTreeFromFlat(menus)
	// 注入角色在各菜单下的按钮权限
	s.attachRolePermsToMenus(ctx, tree, roleID)
	return tree, nil
}

// buildTreeFromFlat 根据扁平菜单构建树
func (s *menuService) buildTreeFromFlat(menus []*entity.Menu) []*entity.Menu {
	// 简易复用：直接用仓储的构建方式——按 parentID 递归拼装
	// 这里复制逻辑以避免额外仓储调用
	var tree []*entity.Menu
	idMap := make(map[uint]*entity.Menu)
	for _, m := range menus {
		idMap[m.ID] = m
	}
	for _, m := range menus {
		if m.ParentID == nil {
			tree = append(tree, m)
			continue
		}
		if parent, ok := idMap[*m.ParentID]; ok {
			parent.Children = append(parent.Children, m)
		} else {
			tree = append(tree, m) // 兜底：无父的视为根
		}
	}
	return tree
}

// attachRolePermsToMenus 将某角色的权限码注入到菜单的 Meta.Perms
func (s *menuService) attachRolePermsToMenus(ctx context.Context, menus []*entity.Menu, roleID uint) {
	var menuIDs []uint
	var collect func(ms []*entity.Menu)
	collect = func(ms []*entity.Menu) {
		for _, m := range ms {
			menuIDs = append(menuIDs, m.ID)
			if len(m.Children) > 0 {
				collect(m.Children)
			}
		}
	}
	collect(menus)

	perms, err := s.permissionRepo.GetByRoleIDAndMenuIDs(ctx, roleID, menuIDs)
	if err != nil {
		logger.Errorf("查询角色菜单权限失败: %v", err)
		return
	}
	byMenu := make(map[uint][]string)
	for _, p := range perms {
		if p.MenuID == nil {
			continue
		}
		byMenu[*p.MenuID] = append(byMenu[*p.MenuID], p.Code)
	}
	var fill func(ms []*entity.Menu)
	fill = func(ms []*entity.Menu) {
		for _, m := range ms {
			if val := byMenu[m.ID]; len(val) > 0 {
				m.Meta.Perms = val
			}
			if len(m.Children) > 0 {
				fill(m.Children)
			}
		}
	}
	fill(menus)
}

// attachPermsToMenus 将权限码注入到菜单的 Meta.Perms；userID 为空表示全量
func (s *menuService) attachPermsToMenus(ctx context.Context, menus []*entity.Menu, userID uint) {
	// 收集所有菜单ID
	var menuIDs []uint
	var collect func(ms []*entity.Menu)
	collect = func(ms []*entity.Menu) {
		for _, m := range ms {
			menuIDs = append(menuIDs, m.ID)
			if len(m.Children) > 0 {
				collect(m.Children)
			}
		}
	}
	collect(menus)

	// 查询权限
	var perms []*entity.Permission
	var err error
	if userID == 0 {
		perms, err = s.permissionRepo.GetByMenuIDs(ctx, menuIDs)
	} else {
		perms, err = s.permissionRepo.GetByUserIDAndMenuIDs(ctx, userID, menuIDs)
	}
	if err != nil {
		logger.Errorf("查询菜单权限失败: %v", err)
		return
	}

	// 按菜单ID聚合权限码
	byMenu := make(map[uint][]string)
	for _, p := range perms {
		if p.MenuID == nil {
			continue
		}
		code := p.Code
		byMenu[*p.MenuID] = append(byMenu[*p.MenuID], code)
	}

	// 写回到 Meta.Perms
	var fill func(ms []*entity.Menu)
	fill = func(ms []*entity.Menu) {
		for _, m := range ms {
			if perms := byMenu[m.ID]; len(perms) > 0 {
				// 保持其他 Meta 字段不变
				m.Meta.Perms = perms
			}
			if len(m.Children) > 0 {
				fill(m.Children)
			}
		}
	}
	fill(menus)
}
