package repository

import (
	"alice/domain/rbac/entity"
	"context"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// Create 创建权限
	Create(ctx context.Context, permission *entity.Permission) error

	// GetByID 根据ID获取权限
	GetByID(ctx context.Context, id string) (*entity.Permission, error)

	// GetByCode 根据代码获取权限
	GetByCode(ctx context.Context, code string) (*entity.Permission, error)

	// List 获取权限列表
	List(ctx context.Context, offset, limit int) ([]*entity.Permission, int64, error)

	// Update 更新权限
	Update(ctx context.Context, permission *entity.Permission) error

	// Delete 删除权限
	Delete(ctx context.Context, id string) error

	// GetByRoleID 根据角色ID获取权限列表
	GetByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)

	// GetByUserID 根据用户ID获取权限列表
	GetByUserID(ctx context.Context, userID string) ([]*entity.Permission, error)

	// AssignToRole 为角色分配权限
	AssignToRole(ctx context.Context, roleID string, permissionIDs []string) error

	// RemoveFromRole 移除角色权限
	RemoveFromRole(ctx context.Context, roleID string, permissionIDs []string) error

	// CheckUserPermission 检查用户是否有指定权限
	CheckUserPermission(ctx context.Context, userID, resource, action string) (bool, error)
}
