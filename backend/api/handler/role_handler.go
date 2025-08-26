package handler

import (
	"net/http"
	"strconv"

	"alice/api/model"
	"alice/domain/rbac/entity"
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
// @Summary 创建角色
// @Description 创建一个新的角色
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateRoleRequest true "创建角色请求"
// @Success 201 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles [post]
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
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	role, err := h.roleService.GetRole(c.Request.Context(), uint(idVal))
	if err != nil {
		logger.Errorf("获取角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(role))
}

// ListRoles 获取角色列表
// @Summary 角色列表
// @Description 分页获取角色列表
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	name := c.Query("name")
	code := c.Query("code")
	statusStr := c.Query("status")
	var statusPtr *entity.RoleStatus
	if statusStr != "" {
		st := entity.RoleStatus(statusStr)
		statusPtr = &st
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	req := &service.ListRolesRequest{
		Page:     page,
		PageSize: pageSize,
		Name:     name,
		Code:     code,
		Status:   statusPtr,
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
// @Summary 更新角色
// @Description 更新指定角色
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Param request body service.UpdateRoleRequest true "更新角色请求"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	req.ID = uint(idVal)

	if err := h.roleService.UpdateRole(c.Request.Context(), &req); err != nil {
		logger.Errorf("更新角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "角色ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID不能为空"))
		return
	}
	idVal, errConv := strconv.ParseUint(idStr, 10, 64)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "角色ID格式错误"))
		return
	}
	if err := h.roleService.DeleteRole(c.Request.Context(), uint(idVal)); err != nil {
		logger.Errorf("删除角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// AssignRolesToUser 为用户分配角色
// @Summary 为用户分配角色
// @Description 为指定用户批量分配角色
// @Tags UserRoles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Param request body model.AssignIDsRequest true "角色ID集合 (role_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/roles [post]
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
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

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.roleService.AssignRolesToUser(c.Request.Context(), uint(userVal), req.RoleIDs); err != nil {
		logger.Errorf("为用户分配角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RemoveRolesFromUser 移除用户角色
// @Summary 移除用户角色
// @Description 从用户移除一组角色
// @Tags UserRoles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Param request body model.AssignIDsRequest true "角色ID集合 (role_ids)"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/roles [delete]
func (h *RoleHandler) RemoveRolesFromUser(c *gin.Context) {
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

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求参数失败: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse(http.StatusBadRequest, "请求参数格式错误"))
		return
	}

	if err := h.roleService.RemoveRolesFromUser(c.Request.Context(), uint(userVal), req.RoleIDs); err != nil {
		logger.Errorf("移除用户角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetUserRoles 获取用户角色
// @Summary 获取用户角色
// @Description 获取指定用户的角色列表
// @Tags UserRoles
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 500 {object} model.APIResponse
// @Router /users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
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
	roles, err := h.roleService.GetUserRoles(c.Request.Context(), uint(userVal))
	if err != nil {
		logger.Errorf("获取用户角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(roles))
}
