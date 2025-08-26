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
	GetByID(ctx context.Context, id uint) (*entity.Role, error)

	// GetByCode 根据代码获取角色
	GetByCode(ctx context.Context, code string) (*entity.Role, error)

	// List 获取角色列表
	List(ctx context.Context, offset, limit int) ([]*entity.Role, int64, error)

	// Search 按条件筛选角色列表（名称/代码/状态 支持模糊匹配名称与代码）
	Search(ctx context.Context, offset, limit int, name, code string, status *entity.RoleStatus) ([]*entity.Role, int64, error)

	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error

	// Delete 删除角色
	Delete(ctx context.Context, id uint) error

	// GetByUserID 根据用户ID获取角色列表
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Role, error)

	// AssignToUser 为用户分配角色
	AssignToUser(ctx context.Context, userID uint, roleIDs []uint) error

	// RemoveFromUser 移除用户角色
	RemoveFromUser(ctx context.Context, userID uint, roleIDs []uint) error
}
