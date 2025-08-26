package service

import (
	"alice/domain/rbac/entity"
	"alice/domain/rbac/repository"
	"alice/pkg/logger"
	"context"
	"fmt"
)

// RoleService 角色服务接口
type RoleService interface {
	// CreateRole 创建角色
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error)

	// GetRole 获取角色
	GetRole(ctx context.Context, id uint) (*entity.Role, error)

	// ListRoles 获取角色列表
	ListRoles(ctx context.Context, req *ListRolesRequest) (*ListRolesResponse, error)

	// UpdateRole 更新角色
	UpdateRole(ctx context.Context, req *UpdateRoleRequest) error

	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, id uint) error

	// AssignRolesToUser 为用户分配角色
	AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error

	// RemoveRolesFromUser 移除用户角色
	RemoveRolesFromUser(ctx context.Context, userID uint, roleIDs []uint) error

	// GetUserRoles 获取用户角色
	GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error)
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string            `json:"name" validate:"required,max=100"`
	Code        string            `json:"code" validate:"required,max=100"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      entity.RoleStatus `json:"status,omitempty"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID          uint              `json:"id" validate:"required"`
	Name        string            `json:"name" validate:"required,max=100"`
	Code        string            `json:"code" validate:"required,max=100"`
	Description *string           `json:"description,omitempty" validate:"omitempty,max=500"`
	Status      entity.RoleStatus `json:"status,omitempty"`
}

// ListRolesRequest 角色列表请求
type ListRolesRequest struct {
	Page     int                `json:"page" validate:"min=1"`
	PageSize int                `json:"page_size" validate:"min=1,max=100"`
	Name     string             `json:"name,omitempty"`
	Code     string             `json:"code,omitempty"`
	Status   *entity.RoleStatus `json:"status,omitempty"`
}

// ListRolesResponse 角色列表响应
type ListRolesResponse struct {
	Items    []*entity.Role `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// roleService 角色服务实现
type roleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

// CreateRole 创建角色
func (s *roleService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error) {
	// 检查代码是否已存在
	existing, _ := s.roleRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("角色代码 %s 已存在", req.Code)
	}

	// 创建角色实体
	role := &entity.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	if role.Status == "" {
		role.Status = entity.RoleStatusActive
	}

	// 保存到数据库
	if err := s.roleRepo.Create(ctx, role); err != nil {
		logger.Errorf("创建角色失败: %v", err)
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	return role, nil
}

// GetRole 获取角色
func (s *roleService) GetRole(ctx context.Context, id uint) (*entity.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取角色失败: %v", err)
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}

	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	return role, nil
}

// ListRoles 获取角色列表
func (s *roleService) ListRoles(ctx context.Context, req *ListRolesRequest) (*ListRolesResponse, error) {
	offset := (req.Page - 1) * req.PageSize

	var (
		roles []*entity.Role
		total int64
		err   error
	)

	// 如果存在过滤条件则使用 Search
	if req.Name != "" || req.Code != "" || req.Status != nil {
		roles, total, err = s.roleRepo.Search(ctx, offset, req.PageSize, req.Name, req.Code, req.Status)
	} else {
		roles, total, err = s.roleRepo.List(ctx, offset, req.PageSize)
	}
	if err != nil {
		logger.Errorf("获取角色列表失败: %v", err)
		return nil, fmt.Errorf("获取角色列表失败: %w", err)
	}

	return &ListRolesResponse{
		Items:    roles,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateRole 更新角色
func (s *roleService) UpdateRole(ctx context.Context, req *UpdateRoleRequest) error {
	// 检查角色是否存在
	existing, err := s.roleRepo.GetByID(ctx, req.ID)
	if err != nil {
		logger.Errorf("获取角色失败: %v", err)
		return fmt.Errorf("获取角色失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("角色不存在")
	}

	// 检查代码是否被其他角色使用
	if existing.Code != req.Code {
		codeExists, _ := s.roleRepo.GetByCode(ctx, req.Code)
		if codeExists != nil && codeExists.ID != req.ID {
			return fmt.Errorf("角色代码 %s 已被其他角色使用", req.Code)
		}
	}

	// 更新角色信息
	existing.Name = req.Name
	existing.Code = req.Code
	existing.Description = req.Description
	existing.Status = req.Status

	if err := s.roleRepo.Update(ctx, existing); err != nil {
		logger.Errorf("更新角色失败: %v", err)
		return fmt.Errorf("更新角色失败: %w", err)
	}

	return nil
}

// DeleteRole 删除角色
func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	// 检查角色是否存在
	existing, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("获取角色失败: %v", err)
		return fmt.Errorf("获取角色失败: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("角色不存在")
	}

	if err := s.roleRepo.Delete(ctx, id); err != nil {
		logger.Errorf("删除角色失败: %v", err)
		return fmt.Errorf("删除角色失败: %w", err)
	}

	return nil
}

// AssignRolesToUser 为用户分配角色
func (s *roleService) AssignRolesToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := s.roleRepo.AssignToUser(ctx, userID, roleIDs); err != nil {
		logger.Errorf("为用户分配角色失败: %v", err)
		return fmt.Errorf("为用户分配角色失败: %w", err)
	}

	return nil
}

// RemoveRolesFromUser 移除用户角色
func (s *roleService) RemoveRolesFromUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := s.roleRepo.RemoveFromUser(ctx, userID, roleIDs); err != nil {
		logger.Errorf("移除用户角色失败: %v", err)
		return fmt.Errorf("移除用户角色失败: %w", err)
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *roleService) GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error) {
	roles, err := s.roleRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取用户角色失败: %v", err)
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	return roles, nil
}
