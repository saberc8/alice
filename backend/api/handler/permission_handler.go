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
// @Summary 创建权限
// @Description 创建一个新的权限
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreatePermissionRequest true "创建权限请求"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /permissions [post]
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
// @Summary 获取权限详情
// @Description 根据ID获取权限详情
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "权限ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /permissions/{id} [get]
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
// @Summary 权限列表
// @Description 分页获取权限列表
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /permissions [get]
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
// @Summary 更新权限
// @Description 更新指定权限
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "权限ID"
// @Param request body service.UpdatePermissionRequest true "更新权限请求"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /permissions/{id} [put]
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
// @Summary 删除权限
// @Description 删除指定权限
// @Tags Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "权限ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /permissions/{id} [delete]
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
// @Summary 为角色分配权限
// @Description 为指定角色批量分配权限
// @Tags RolePermissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body model.AssignIDsRequest true "权限ID集合 (permission_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/permissions [post]
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
// @Summary 移除角色权限
// @Description 从角色移除一组权限
// @Tags RolePermissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Param request body model.AssignIDsRequest true "权限ID集合 (permission_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/permissions [delete]
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
// @Summary 获取角色权限
// @Description 获取指定角色的权限列表
// @Tags RolePermissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "角色ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id}/permissions [get]
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
// @Summary 获取用户权限
// @Description 获取指定用户的权限列表
// @Tags UserPermissions
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/permissions [get]
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
// @Summary 检查用户权限
// @Description 检查用户是否拥有某资源操作权限
// @Tags UserPermissions
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "用户ID"
// @Param resource query string true "资源标识"
// @Param action query string true "操作标识"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/permissions/check [get]
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
