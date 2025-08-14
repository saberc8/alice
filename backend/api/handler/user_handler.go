package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"alice/api/model"
	rbacsvc "alice/domain/rbac/service"
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
