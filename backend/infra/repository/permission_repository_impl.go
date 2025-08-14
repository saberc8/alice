package repository

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// permissionRepositoryImpl 权限仓储实现
type permissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepositoryImpl{
		db: db,
	}
}

// Create 创建权限
func (r *permissionRepositoryImpl) Create(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *permissionRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据代码获取权限
func (r *permissionRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Permission, error) {
	var permission entity.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表
func (r *permissionRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.Permission, int64, error) {
	var permissions []*entity.Permission
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&permissions).Error

	return permissions, total, err
}

// Update 更新权限
func (r *permissionRepositoryImpl) Update(ctx context.Context, permission *entity.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// Delete 删除权限
func (r *permissionRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("permission_id = ?", id).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 删除权限
		return tx.Where("id = ?", id).Delete(&entity.Permission{}).Error
	})
}

// GetByRoleID 根据角色ID获取权限列表
func (r *permissionRepositoryImpl) GetByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	var permissions []*entity.Permission

	err := r.db.WithContext(ctx).
		Select("permissions.*").
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.status = ?", roleID, entity.PermissionStatusActive).
		Find(&permissions).Error

	return permissions, err
}

// GetByUserID 根据用户ID获取权限列表
func (r *permissionRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*entity.Permission, error) {
	var permissions []*entity.Permission

	err := r.db.WithContext(ctx).
		Select("DISTINCT permissions.*").
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.status = ? AND roles.status = ?",
			userID, entity.PermissionStatusActive, entity.RoleStatusActive).
		Find(&permissions).Error

	return permissions, err
}

// AssignToRole 为角色分配权限
func (r *permissionRepositoryImpl) AssignToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除角色现有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 创建新的关联关系
		var rolePermissions []entity.RolePermission
		for _, permissionID := range permissionIDs {
			rolePermissions = append(rolePermissions, entity.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			})
		}

		return tx.Create(&rolePermissions).Error
	})
}

// RemoveFromRole 移除角色权限
func (r *permissionRepositoryImpl) RemoveFromRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&entity.RolePermission{}).Error
}

// CheckUserPermission 检查用户是否有指定权限
func (r *permissionRepositoryImpl) CheckUserPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ? AND permissions.status = ? AND roles.status = ?",
			userID, resource, action, entity.PermissionStatusActive, entity.RoleStatusActive).
		Count(&count).Error

	return count > 0, err
}

// CheckUserPermissionByCode 根据权限码检查用户权限
func (r *permissionRepositoryImpl) CheckUserPermissionByCode(ctx context.Context, userID, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.code = ? AND permissions.status = ? AND roles.status = ?",
			userID, code, entity.PermissionStatusActive, entity.RoleStatusActive).
		Count(&count).Error

	return count > 0, err
}

// GetByMenuIDs 根据菜单ID集合获取权限列表
func (r *permissionRepositoryImpl) GetByMenuIDs(ctx context.Context, menuIDs []string) ([]*entity.Permission, error) {
	if len(menuIDs) == 0 {
		return []*entity.Permission{}, nil
	}
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Where("menu_id IN ? AND status = ?", menuIDs, entity.PermissionStatusActive).
		Order("created_at ASC").
		Find(&permissions).Error
	return permissions, err
}

// GetByUserIDAndMenuIDs 根据用户和菜单ID集合获取权限列表
func (r *permissionRepositoryImpl) GetByUserIDAndMenuIDs(ctx context.Context, userID string, menuIDs []string) ([]*entity.Permission, error) {
	if len(menuIDs) == 0 {
		return []*entity.Permission{}, nil
	}
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Select("DISTINCT permissions.*").
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND permissions.menu_id IN ? AND permissions.status = ? AND roles.status = ?",
			userID, menuIDs, entity.PermissionStatusActive, entity.RoleStatusActive).
		Find(&permissions).Error
	return permissions, err
}

// GetByRoleIDAndMenuIDs 根据角色和菜单ID集合获取权限列表
func (r *permissionRepositoryImpl) GetByRoleIDAndMenuIDs(ctx context.Context, roleID string, menuIDs []string) ([]*entity.Permission, error) {
	if len(menuIDs) == 0 {
		return []*entity.Permission{}, nil
	}
	var permissions []*entity.Permission
	err := r.db.WithContext(ctx).
		Select("permissions.*").
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.menu_id IN ? AND permissions.status = ?",
			roleID, menuIDs, entity.PermissionStatusActive).
		Find(&permissions).Error
	return permissions, err
}
