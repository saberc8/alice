package handler

import (
	"net/http"
	"strconv"

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
// @Summary 创建菜单
// @Description 创建一个新的菜单
// @Tags Menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateMenuRequest true "创建菜单请求"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus [post]
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
// @Summary 获取菜单详情
// @Description 根据ID获取菜单详情
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID格式错误"))
		return
	}
	menu, err := h.menuService.GetMenu(c.Request.Context(), uint(idVal))
	if err != nil {
		logger.Errorf("获取菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menu))
}

// ListMenus 获取菜单列表
// @Summary 菜单列表
// @Description 获取所有菜单列表
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus [get]
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
// @Summary 菜单树
// @Description 获取完整菜单树
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus/tree [get]
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
// @Summary 更新菜单
// @Description 更新指定菜单
// @Tags Menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Param request body service.UpdateMenuRequest true "更新菜单请求"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID格式错误"))
		return
	}

	var req service.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	req.ID = uint(idVal)

	if err := h.menuService.UpdateMenu(c.Request.Context(), &req); err != nil {
		logger.Errorf("更新菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除指定菜单
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "菜单ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "菜单ID格式错误"))
		return
	}
	if err := h.menuService.DeleteMenu(c.Request.Context(), uint(idVal)); err != nil {
		logger.Errorf("删除菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// AssignMenusToRole 为角色分配菜单
// @Summary 为角色分配菜单
// @Description 为指定角色批量分配菜单
// @Tags RoleMenus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body model.AssignIDsRequest true "菜单ID集合 (menu_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/menus [post]
func (h *MenuHandler) AssignMenusToRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	if roleIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	roleVal, errConv := strconv.ParseUint(roleIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	var req struct {
		MenuIDs []uint `json:"menu_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.menuService.AssignMenusToRole(c.Request.Context(), uint(roleVal), req.MenuIDs); err != nil {
		logger.Errorf("为角色分配菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RemoveMenusFromRole 移除角色菜单
// @Summary 移除角色菜单
// @Description 从角色移除一组菜单
// @Tags RoleMenus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body model.AssignIDsRequest true "菜单ID集合 (menu_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/menus [delete]
func (h *MenuHandler) RemoveMenusFromRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	if roleIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	roleVal, errConv := strconv.ParseUint(roleIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	var req struct {
		MenuIDs []uint `json:"menu_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.menuService.RemoveMenusFromRole(c.Request.Context(), uint(roleVal), req.MenuIDs); err != nil {
		logger.Errorf("移除角色菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetRoleMenus 获取角色菜单
// @Summary 获取角色菜单
// @Description 获取指定角色的菜单列表
// @Tags RoleMenus
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/menus [get]
func (h *MenuHandler) GetRoleMenus(c *gin.Context) {
	roleIDStr := c.Param("id")
	if roleIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	roleVal, errConv := strconv.ParseUint(roleIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	menus, err := h.menuService.GetRoleMenus(c.Request.Context(), uint(roleVal))
	if err != nil {
		logger.Errorf("获取角色菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menus))
}

// GetUserMenus 获取用户菜单
// @Summary 获取用户菜单
// @Description 获取指定用户的菜单列表
// @Tags UserMenus
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/menus [get]
func (h *MenuHandler) GetUserMenus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}
	userVal, errConv := strconv.ParseUint(userIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID格式错误"))
		return
	}
	menus, err := h.menuService.GetUserMenus(c.Request.Context(), uint(userVal))
	if err != nil {
		logger.Errorf("获取用户菜单失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(menus))
}

// GetUserMenuTree 获取用户菜单树
// @Summary 用户菜单树
// @Description 获取指定用户的菜单树
// @Tags UserMenus
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/menus/tree [get]
func (h *MenuHandler) GetUserMenuTree(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID不能为空"))
		return
	}
	userVal, errConv := strconv.ParseUint(userIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "用户ID格式错误"))
		return
	}
	tree, err := h.menuService.GetUserMenuTree(c.Request.Context(), uint(userVal))
	if err != nil {
		logger.Errorf("获取用户菜单树失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(tree))
}

// GetRoleMenuTree 获取角色菜单树（按角色注入 perms）
// @Summary 角色菜单树
// @Description 获取指定角色的菜单树，并在 meta.perms 中下发该角色的按钮权限
// @Tags RoleMenus
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/menus/tree [get]
func (h *MenuHandler) GetRoleMenuTree(c *gin.Context) {
	roleIDStr := c.Param("id")
	if roleIDStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	roleVal, errConv := strconv.ParseUint(roleIDStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	tree, err := h.menuService.GetRoleMenuTree(c.Request.Context(), uint(roleVal))
	if err != nil {
		logger.Errorf("获取角色菜单树失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(tree))
}
