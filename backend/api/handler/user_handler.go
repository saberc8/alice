package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"alice/api/model"
	rbacsvc "alice/domain/rbac/service"
	userentity "alice/domain/user/entity"
	"alice/domain/user/service"
	"alice/pkg/logger"
)

type UserHandler struct {
	userService service.UserService
	roleService rbacsvc.RoleService
}

func NewUserHandler(userService service.UserService, roleService rbacsvc.RoleService) *UserHandler {
	return &UserHandler{
		userService: userService,
		roleService: roleService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册请求"
// @Success 200 {object} model.APIResponse{data=model.RegisterResponse}
// @Failure 400 {object} model.APIResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Register bind request failed: %v", err)
		response := model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := h.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		logger.Errorf("Register failed: %v", err)
		response := model.ErrorResponse(model.CodeBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		logger.Errorf("Generate token failed: %v", err)
		response := model.ErrorResponse(model.CodeInternalError, "Generate token failed")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := model.SuccessResponseWithMessage(model.MsgRegisterSuccess, model.RegisterResponse{
		User: model.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Token: token,
	})
	c.JSON(http.StatusOK, response)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取 JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录请求"
// @Success 200 {object} model.APIResponse{data=model.LoginResponse}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Login bind request failed: %v", err)
		response := model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		logger.Errorf("Login failed: %v", err)
		response := model.ErrorResponse(model.CodeUnauthorized, model.MsgInvalidCredentials)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	response := model.SuccessResponseWithMessage(model.MsgLoginSuccess, model.LoginResponse{
		Token: token,
	})
	c.JSON(http.StatusOK, response)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户资料
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.APIResponse{data=model.UserInfo}
// @Failure 401 {object} model.APIResponse
// @Router /auth/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response := model.ErrorResponse(model.CodeUnauthorized, model.MsgUnauthorized)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response := model.ErrorResponse(model.CodeUnauthorized, "Invalid user ID")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	user, err := h.userService.GetUserByID(uid)
	if err != nil {
		logger.Errorf("Get user profile failed: %v", err)
		response := model.ErrorResponse(model.CodeNotFound, model.MsgUserNotFound)
		c.JSON(http.StatusNotFound, response)
		return
	}

	// 查询用户角色（若可用）
	var roles []model.RoleBrief
	if h.roleService != nil {
		// 注意：RoleService 使用的是 string 类型的用户ID，这里统一为字符串
		userIDStr := strconv.FormatUint(uint64(uid), 10)
		if rlist, rerr := h.roleService.GetUserRoles(c.Request.Context(), userIDStr); rerr == nil {
			roles = make([]model.RoleBrief, 0, len(rlist))
			for _, r := range rlist {
				roles = append(roles, model.RoleBrief{ID: r.ID, Name: r.Name, Code: r.Code})
			}
		} else {
			logger.Errorf("Get user roles failed: %v", rerr)
		}
	}

	response := model.SuccessResponse(model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	})
	c.JSON(http.StatusOK, response)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户资料
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileRequest true "更新资料请求"
// @Success 200 {object} model.APIResponse{data=model.UserInfo}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response := model.ErrorResponse(model.CodeUnauthorized, model.MsgUnauthorized)
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response := model.ErrorResponse(model.CodeUnauthorized, "Invalid user ID")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("UpdateProfile bind request failed: %v", err)
		response := model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := h.userService.UpdateProfile(uid, req.Email)
	if err != nil {
		logger.Errorf("Update profile failed: %v", err)
		response := model.ErrorResponse(model.CodeBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := model.SuccessResponseWithMessage("Profile updated successfully", model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
	c.JSON(http.StatusOK, response)
}

// ========== Admin: 用户管理 CRUD ==========

// CreateUser 管理员创建用户
// @Summary 管理员创建用户
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AdminCreateUserRequest true "创建用户"
// @Success 200 {object} model.APIResponse{data=model.UserInfo}
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	status := userentity.UserStatusActive
	if req.Status != "" {
		status = userentity.UserStatus(req.Status)
	}
	user, err := h.userService.CreateUserByAdmin(req.Username, req.Password, req.Email, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(model.UserInfo{ID: user.ID, Username: user.Username, Email: user.Email}))
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.APIResponse{data=model.UserInfo}
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("user_id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	user, err := h.userService.GetUser(uint(id64))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(model.CodeNotFound, model.MsgUserNotFound))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(model.UserInfo{ID: user.ID, Username: user.Username, Email: user.Email}))
}

// ListUsers 用户列表
// @Summary 用户列表
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页条数"
// @Success 200 {object} model.APIResponse{data=model.UserListResponse}
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var q model.ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	if q.Page == 0 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = 10
	}
	users, total, err := h.userService.ListUsers(q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, model.MsgInternalError))
		return
	}
	items := make([]model.UserItem, 0, len(users))
	for _, u := range users {
		items = append(items, model.UserItem{ID: u.ID, Username: u.Username, Email: u.Email, Status: string(u.Status)})
	}
	c.JSON(http.StatusOK, model.SuccessResponse(model.UserListResponse{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}))
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param request body model.AdminUpdateUserRequest true "更新用户"
// @Success 200 {object} model.APIResponse{data=model.UserInfo}
// @Router /users/{user_id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("user_id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	var req model.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	var status userentity.UserStatus
	if req.Status != "" {
		status = userentity.UserStatus(req.Status)
	}
	user, err := h.userService.UpdateUserByAdmin(uint(id64), req.Email, status, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(model.UserInfo{ID: user.ID, Username: user.Username, Email: user.Email}))
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.APIResponse
// @Router /users/{user_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("user_id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(model.CodeBadRequest, model.MsgInvalidRequest))
		return
	}
	if err := h.userService.DeleteUser(uint(id64)); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(model.CodeInternalError, model.MsgInternalError))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponseWithMessage("deleted", nil))
}
