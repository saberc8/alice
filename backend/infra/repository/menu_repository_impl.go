package repository

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// menuRepositoryImpl 菜单仓储实现
type menuRepositoryImpl struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &menuRepositoryImpl{
		db: db,
	}
}

// Create 创建菜单
func (r *menuRepositoryImpl) Create(ctx context.Context, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

// GetByID 根据ID获取菜单
func (r *menuRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

// GetByCode 根据代码获取菜单
func (r *menuRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

// List 获取菜单列表
func (r *menuRepositoryImpl) List(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu

	err := r.db.WithContext(ctx).
		Where("status = ?", entity.MenuStatusActive).
		Order("\"order\" ASC, created_at ASC").
		Find(&menus).Error

	return menus, err
}

// GetTree 获取菜单树
func (r *menuRepositoryImpl) GetTree(ctx context.Context) ([]*entity.Menu, error) {
	// 获取所有菜单
	menus, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return r.buildMenuTree(menus, nil), nil
}

// Update 更新菜单
func (r *menuRepositoryImpl) Update(ctx context.Context, menu *entity.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

// Delete 删除菜单
func (r *menuRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Where("menu_id = ?", id).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}

		// 删除菜单
		return tx.Where("id = ?", id).Delete(&entity.Menu{}).Error
	})
}

// GetByUserID 根据用户ID获取菜单列表
func (r *menuRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*entity.Menu, error) {
	var menus []*entity.Menu

	err := r.db.WithContext(ctx).
		Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Joins("INNER JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND menus.status = ? AND roles.status = ?",
			userID, entity.MenuStatusActive, entity.RoleStatusActive).
		Group("menus.id, menus.parent_id, menus.name, menus.code, menus.path, menus.type, menus.\"order\", menus.status, menus.description, menus.created_at, menus.updated_at").
		Order("menus.\"order\" ASC, menus.created_at ASC").
		Find(&menus).Error

	return menus, err
}

// GetTreeByUserID 根据用户ID获取菜单树
func (r *menuRepositoryImpl) GetTreeByUserID(ctx context.Context, userID string) ([]*entity.Menu, error) {
	// 获取用户菜单
	menus, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return r.buildMenuTree(menus, nil), nil
}

// GetByRoleID 根据角色ID获取菜单列表
func (r *menuRepositoryImpl) GetByRoleID(ctx context.Context, roleID string) ([]*entity.Menu, error) {
	var menus []*entity.Menu

	err := r.db.WithContext(ctx).
		Select("menus.*").
		Table("menus").
		Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = ?", roleID, entity.MenuStatusActive).
		Order("menus.`order` ASC, menus.created_at ASC").
		Find(&menus).Error

	return menus, err
}

// AssignToRole 为角色分配菜单
func (r *menuRepositoryImpl) AssignToRole(ctx context.Context, roleID string, menuIDs []string) error {
	if len(menuIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除角色现有菜单
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}

		// 创建新的关联关系
		var roleMenus []entity.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, entity.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}

		return tx.Create(&roleMenus).Error
	})
}

// RemoveFromRole 移除角色菜单
func (r *menuRepositoryImpl) RemoveFromRole(ctx context.Context, roleID string, menuIDs []string) error {
	if len(menuIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("role_id = ? AND menu_id IN ?", roleID, menuIDs).
		Delete(&entity.RoleMenu{}).Error
}

// GetChildren 获取子菜单
func (r *menuRepositoryImpl) GetChildren(ctx context.Context, parentID string) ([]*entity.Menu, error) {
	var menus []*entity.Menu

	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND status = ?", parentID, entity.MenuStatusActive).
		Order("`order` ASC, created_at ASC").
		Find(&menus).Error

	return menus, err
}

// buildMenuTree 构建菜单树
func (r *menuRepositoryImpl) buildMenuTree(menus []*entity.Menu, parentID *string) []*entity.Menu {
	var tree []*entity.Menu

	for _, menu := range menus {
		// 判断是否为当前父节点的子节点
		isChild := false
		if parentID == nil && menu.ParentID == nil {
			isChild = true
		} else if parentID != nil && menu.ParentID != nil && *parentID == *menu.ParentID {
			isChild = true
		}

		if isChild {
			// 递归查找子节点
			menu.Children = r.buildMenuTree(menus, &menu.ID)
			tree = append(tree, menu)
		}
	}

	return tree
}
