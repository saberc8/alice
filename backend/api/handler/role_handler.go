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

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), &req)
	if err != nil {
		logger.Errorf("创建角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(role))
}

// GetRole 获取角色
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), id)
	if err != nil {
		logger.Errorf("获取角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(role))
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := &service.ListRolesRequest{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.roleService.ListRoles(c.Request.Context(), req)
	if err != nil {
		logger.Errorf("获取角色列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	req.ID = id

	if err := h.roleService.UpdateRole(c.Request.Context(), &req); err != nil {
		logger.Errorf("更新角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	if err := h.roleService.DeleteRole(c.Request.Context(), id); err != nil {
		logger.Errorf("删除角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// AssignRolesToUser 为用户分配角色
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.roleService.AssignRolesToUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
		logger.Errorf("为用户分配角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RemoveRolesFromUser 移除用户角色
func (h *RoleHandler) RemoveRolesFromUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	var req struct {
		RoleIDs []string `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.roleService.RemoveRolesFromUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
		logger.Errorf("移除用户角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetUserRoles 获取用户角色
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	roles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		logger.Errorf("获取用户角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(roles))
}
