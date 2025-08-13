package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"alice/api/model"
	"alice/domain/user/service"
	"alice/pkg/logger"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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

	response := model.SuccessResponse(model.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
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
