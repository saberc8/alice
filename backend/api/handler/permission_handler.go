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

package handler

import (
	"net/http"
	"strconv"

	"alice/api/model"
	"alice/domain/rbac/service"
	"alice/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService service.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req service.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	permission, err := h.permissionService.CreatePermission(c.Request.Context(), &req)
	if err != nil {
		logger.Errorf("创建权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(permission))
}

// GetPermission 获取权限
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "权限ID不能为空"))
		return
	}

	permission, err := h.permissionService.GetPermission(c.Request.Context(), id)
	if err != nil {
		logger.Errorf("获取权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(permission))
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := &service.ListPermissionsRequest{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.permissionService.ListPermissions(c.Request.Context(), req)
	if err != nil {
		logger.Errorf("获取权限列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "权限ID不能为空"))
		return
	}

	var req service.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	req.ID = id

	if err := h.permissionService.UpdatePermission(c.Request.Context(), &req); err != nil {
		logger.Errorf("更新权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "权限ID不能为空"))
		return
	}

	if err := h.permissionService.DeletePermission(c.Request.Context(), id); err != nil {
		logger.Errorf("删除权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// AssignPermissionsToRole 为角色分配权限
func (h *PermissionHandler) AssignPermissionsToRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.permissionService.AssignPermissionsToRole(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		logger.Errorf("为角色分配权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RemovePermissionsFromRole 移除角色权限
func (h *PermissionHandler) RemovePermissionsFromRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	var req struct {
		PermissionIDs []string `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.permissionService.RemovePermissionsFromRole(c.Request.Context(), roleID, req.PermissionIDs); err != nil {
		logger.Errorf("移除角色权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetRolePermissions 获取角色权限
func (h *PermissionHandler) GetRolePermissions(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	permissions, err := h.permissionService.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		logger.Errorf("获取角色权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(permissions))
}

// GetUserPermissions 获取用户权限
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	permissions, err := h.permissionService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		logger.Errorf("获取用户权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(permissions))
}

// CheckUserPermission 检查用户权限
func (h *PermissionHandler) CheckUserPermission(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	resource := c.Query("resource")
	action := c.Query("action")

	if resource == "" || action == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "资源和操作不能为空"))
		return
	}

	hasPermission, err := h.permissionService.CheckUserPermission(c.Request.Context(), userID, resource, action)
	if err != nil {
		logger.Errorf("检查用户权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	result := map[string]bool{
		"has_permission": hasPermission,
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}
