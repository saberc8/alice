package repository

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"context"

	"gorm.io/gorm"
)

// roleRepositoryImpl 角色仓储实现
type roleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepositoryImpl{
		db: db,
	}
}

// Create 创建角色
func (r *roleRepositoryImpl) Create(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID 根据ID获取角色
func (r *roleRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据代码获取角色
func (r *roleRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *roleRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.Role, int64, error) {
	var roles []*entity.Role
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&roles).Error

	return roles, total, err
}

// Update 更新角色
func (r *roleRepositoryImpl) Update(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// Delete 删除角色
func (r *roleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}

		// 删除角色权限关联
		if err := tx.Where("role_id = ?", id).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}

		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&entity.RoleMenu{}).Error; err != nil {
			return err
		}

		// 删除角色
		return tx.Where("id = ?", id).Delete(&entity.Role{}).Error
	})
}

// GetByUserID 根据用户ID获取角色列表
func (r *roleRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entity.Role, error) {
	var roles []*entity.Role

	err := r.db.WithContext(ctx).
		Select("roles.*").
		Table("roles").
		Joins("INNER JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = ?", userID, entity.RoleStatusActive).
		Find(&roles).Error

	return roles, err
}

// AssignToUser 为用户分配角色
func (r *roleRepositoryImpl) AssignToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除用户现有角色
		if err := tx.Where("user_id = ?", userID).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}

		// 创建新的关联关系
		var userRoles []entity.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, entity.UserRole{UserID: userID, RoleID: roleID})
		}

		return tx.Create(&userRoles).Error
	})
}

// RemoveFromUser 移除用户角色
func (r *roleRepositoryImpl) RemoveFromUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&entity.UserRole{}).Error
}
