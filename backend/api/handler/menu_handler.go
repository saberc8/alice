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

	"alice/api/model"
	"alice/domain/rbac/service"
	"alice/pkg/logger"

	"github.com/gin-gonic/gin"
)

// MenuHandler 菜单处理器
type MenuHandler struct {
	menuService service.MenuService
}

// NewMenuHandler 创建菜单处理器
func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// CreateMenu 创建菜单
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req service.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	menu, err := h.menuService.CreateMenu(c.Request.Context(), &req)
	if err != nil {
		logger.Errorf("创建菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(menu))
}

// GetMenu 获取菜单
func (h *MenuHandler) GetMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}

	menu, err := h.menuService.GetMenu(c.Request.Context(), id)
	if err != nil {
		logger.Errorf("获取菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menu))
}

// ListMenus 获取菜单列表
func (h *MenuHandler) ListMenus(c *gin.Context) {
	menus, err := h.menuService.ListMenus(c.Request.Context())
	if err != nil {
		logger.Errorf("获取菜单列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menus))
}

// GetMenuTree 获取菜单树
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	tree, err := h.menuService.GetMenuTree(c.Request.Context())
	if err != nil {
		logger.Errorf("获取菜单树失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(tree))
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}

	var req service.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	req.ID = id

	if err := h.menuService.UpdateMenu(c.Request.Context(), &req); err != nil {
		logger.Errorf("更新菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}

	if err := h.menuService.DeleteMenu(c.Request.Context(), id); err != nil {
		logger.Errorf("删除菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// AssignMenusToRole 为角色分配菜单
func (h *MenuHandler) AssignMenusToRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	var req struct {
		MenuIDs []string `json:"menu_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.menuService.AssignMenusToRole(c.Request.Context(), roleID, req.MenuIDs); err != nil {
		logger.Errorf("为角色分配菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RemoveMenusFromRole 移除角色菜单
func (h *MenuHandler) RemoveMenusFromRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	var req struct {
		MenuIDs []string `json:"menu_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.menuService.RemoveMenusFromRole(c.Request.Context(), roleID, req.MenuIDs); err != nil {
		logger.Errorf("移除角色菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetRoleMenus 获取角色菜单
func (h *MenuHandler) GetRoleMenus(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}

	menus, err := h.menuService.GetRoleMenus(c.Request.Context(), roleID)
	if err != nil {
		logger.Errorf("获取角色菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menus))
}

// GetUserMenus 获取用户菜单
func (h *MenuHandler) GetUserMenus(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	menus, err := h.menuService.GetUserMenus(c.Request.Context(), userID)
	if err != nil {
		logger.Errorf("获取用户菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menus))
}

// GetUserMenuTree 获取用户菜单树
func (h *MenuHandler) GetUserMenuTree(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}

	tree, err := h.menuService.GetUserMenuTree(c.Request.Context(), userID)
	if err != nil {
		logger.Errorf("获取用户菜单树失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(tree))
}
