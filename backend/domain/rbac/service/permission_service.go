package service

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"alice/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PermissionService 权限服务接口
type PermissionService interface {
	// CreatePermission 创建权限
	CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*entity.Permission, error)

	// GetPermission 获取权限
	GetPermission(ctx context.Context, id string) (*entity.Permission, error)

	// ListPermissions 获取权限列表
	ListPermissions(ctx context.Context, req *ListPermissionsRequest) (*ListPermissionsResponse, error)

	// UpdatePermission 更新权限
	UpdatePermission(ctx context.Context, req *UpdatePermissionRequest) error

	// DeletePermission 删除权限
	DeletePermission(ctx context.Context, id string) error

	// AssignPermissionsToRole 为角色分配权限
	AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error

	// RemovePermissionsFromRole 移除角色权限
	RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error

	// GetRolePermissions 获取角色权限
	GetRolePermissions(ctx context.Context, roleID string) ([]*entity.Permission, error)

	// GetUserPermissions 获取用户权限
	GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error)

	// CheckUserPermission 检查用户权限
	CheckUserPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string                  `json:"name" validate:"required,max=100"`
	Code        string                  `json:"code" validate:"required,max=100"`
	Resource    string                  `json:"resource" validate:"required,max=100"`
	Action      string                  `json:"action" validate:"required,max=50"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      entity.PermissionStatus `json:"status,omitempty"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	ID          string                  `json:"id" validate:"required"`
	Name        string                  `json:"name" validate:"required,max=100"`
	Code        string                  `json:"code" validate:"required,max=100"`
	Resource    string                  `json:"resource" validate:"required,max=100"`
	Action      string                  `json:"action" validate:"required,max=50"`
	Description *string                 `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      entity.PermissionStatus `json:"status,omitempty"`
}

// ListPermissionsRequest 权限列表请求
type ListPermissionsRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// ListPermissionsResponse 权限列表响应
type ListPermissionsResponse struct {
	Items    []*entity.Permission `json:"items"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

// permissionService 权限服务实现
type permissionService struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

// CreatePermission 创建权限
func (s *permissionService) CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*entity.Permission, error) {
	// 检查代码是否已存在
	existing, _ := s.permissionRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("权限代码 %s 已存在", req.Code)
	}

	// 创建权限实体
	permission := &entity.Permission{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Code:        req.Code,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		Status:      req.Status,
	}

	if permission.Status == "" {
		permission.Status = entity.PermissionStatusActive
	}

	// 保存到数据库
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		logger.Errorf("创建权限失败: %v", err)
		return nil, fmt.Errorf("创建权限失败: %w", err)
	}

	return permission, nil
}

// GetPermission 获取权限
func (s *permissionService) GetPermission(ctx context.Context, id string) (*entity.Permission, error) {
	permission, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取权限失败: %v", err)
		return nil, fmt.Errorf("获取权限失败: %w", err)
	}

	if permission == nil {
		return nil, fmt.Errorf("权限不存在")
	}

	return permission, nil
}

// ListPermissions 获取权限列表
func (s *permissionService) ListPermissions(ctx context.Context, req *ListPermissionsRequest) (*ListPermissionsResponse, error) {
	offset := (req.Page - 1) * req.PageSize

	permissions, total, err := s.permissionRepo.List(ctx, offset, req.PageSize)
	if err != nil {
		logger.Errorf("获取权限列表失败: %v", err)
		return nil, fmt.Errorf("获取权限列表失败: %w", err)
	}

	return &ListPermissionsResponse{
		Items:    permissions,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdatePermission 更新权限
func (s *permissionService) UpdatePermission(ctx context.Context, req *UpdatePermissionRequest) error {
	// 检查权限是否存在
	existing, err := s.permissionRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.Errorf("获取权限失败: %v", err)
		return fmt.Errorf("获取权限失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("权限不存在")
	}

	// 检查代码是否被其他权限使用
	if existing.Code != req.Code {
		codeExists, _ := s.permissionRepo.GetByCode(ctx, req.Code)
		if codeExists != nil && codeExists.ID != req.ID {
			return fmt.Errorf("权限代码 %s 已被其他权限使用", req.Code)
		}
	}

	// 更新权限信息
	existing.Name = req.Name
	existing.Code = req.Code
	existing.Resource = req.Resource
	existing.Action = req.Action
	existing.Description = req.Description
	existing.Status = req.Status

	if err := s.permissionRepo.Update(ctx, existing); err != nil {
		logger.Errorf("更新权限失败: %v", err)
		return fmt.Errorf("更新权限失败: %w", err)
	}

	return nil
}

// DeletePermission 删除权限
func (s *permissionService) DeletePermission(ctx context.Context, id string) error {
	// 检查权限是否存在
	existing, err := s.permissionRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取权限失败: %v", err)
		return fmt.Errorf("获取权限失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("权限不存在")
	}

	if err := s.permissionRepo.Delete(ctx, id); err != nil {
		logger.Errorf("删除权限失败: %v", err)
		return fmt.Errorf("删除权限失败: %w", err)
	}

	return nil
}

// AssignPermissionsToRole 为角色分配权限
func (s *permissionService) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if err := s.permissionRepo.AssignToRole(ctx, roleID, permissionIDs); err != nil {
		logger.Errorf("为角色分配权限失败: %v", err)
		return fmt.Errorf("为角色分配权限失败: %w", err)
	}

	return nil
}

// RemovePermissionsFromRole 移除角色权限
func (s *permissionService) RemovePermissionsFromRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if err := s.permissionRepo.RemoveFromRole(ctx, roleID, permissionIDs); err != nil {
		logger.Errorf("移除角色权限失败: %v", err)
		return fmt.Errorf("移除角色权限失败: %w", err)
	}

	return nil
}

// GetRolePermissions 获取角色权限
func (s *permissionService) GetRolePermissions(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	permissions, err := s.permissionRepo.GetByRoleID(ctx, roleID)
	if err != nil {
		logger.Errorf("获取角色权限失败: %v", err)
		return nil, fmt.Errorf("获取角色权限失败: %w", err)
	}

	return permissions, nil
}

// GetUserPermissions 获取用户权限
func (s *permissionService) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	permissions, err := s.permissionRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户权限失败: %v", err)
		return nil, fmt.Errorf("获取用户权限失败: %w", err)
	}

	return permissions, nil
}

// CheckUserPermission 检查用户权限
func (s *permissionService) CheckUserPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	hasPermission, err := s.permissionRepo.CheckUserPermission(ctx, userID, resource, action)
	if err != nil {
		logger.Errorf("检查用户权限失败: %v", err)
		return false, fmt.Errorf("检查用户权限失败: %w", err)
	}

	return hasPermission, nil
}
