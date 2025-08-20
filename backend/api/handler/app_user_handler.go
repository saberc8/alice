package handler

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apimodel "alice/api/model"
	"alice/application"
	friendsvc "alice/domain/appfriend/service"
	appsvc "alice/domain/appuser/service"
	"alice/infra/config"
	"alice/pkg/logger"
)

type AppUserHandler struct {
	svc       appsvc.AppUserService
	friendSvc friendsvc.FriendService
}

func NewAppUserHandler(svc appsvc.AppUserService, friendSvc friendsvc.FriendService) *AppUserHandler {
	return &AppUserHandler{svc: svc, friendSvc: friendSvc}
}

// AppRegister 移动端注册
// @Summary App 注册
// @Description 使用邮箱注册，无需邮箱验证
// @Tags App
// @Accept json
// @Produce json
// @Param request body model.AppRegisterRequest true "注册请求"
// @Success 200 {object} model.APIResponse{data=model.AppAuthResponse}
// @Failure 400 {object} model.APIResponse
// @Router /app/register [post]
func (h *AppUserHandler) AppRegister(c *gin.Context) {
	var req apimodel.AppRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	u, err := h.svc.Register(req.Email, req.Password, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "failed to issue token"))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.AppAuthResponse{User: apimodel.AppUserInfo{ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar, Gender: u.Gender, Bio: u.Bio}, Token: token}))
}

// AppLogin 移动端登录
// @Summary App 登录
// @Description 使用邮箱+密码登录，返回 JWT
// @Tags App
// @Accept json
// @Produce json
// @Param request body model.AppLoginRequest true "登录请求"
// @Success 200 {object} model.APIResponse{data=model.LoginResponse}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/login [post]
func (h *AppUserHandler) AppLogin(c *gin.Context) {
	var req apimodel.AppLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgInvalidCredentials))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.LoginResponse{Token: token}))
}

// AppProfile 获取移动端用户资料
// @Summary App 当前用户资料
// @Tags App
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.APIResponse{data=model.AppUserInfo}
// @Failure 401 {object} model.APIResponse
// @Router /app/profile [get]
func (h *AppUserHandler) AppProfile(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	u, err := h.svc.GetByID(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, apimodel.ErrorResponse(apimodel.CodeNotFound, apimodel.MsgUserNotFound))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.AppUserInfo{ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar, Gender: u.Gender, Bio: u.Bio}))
}

// AppUpdateProfile 更新移动端用户资料
// @Summary App 更新资料
// @Tags App
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AppUpdateProfileRequest true "更新资料"
// @Success 200 {object} model.APIResponse{data=model.AppUserInfo}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/profile [put]
func (h *AppUserHandler) AppUpdateProfile(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	var req apimodel.AppUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	u, err := h.svc.UpdateProfile(uid, req.Nickname, req.Avatar, req.Gender, req.Bio)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.AppUserInfo{ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar, Gender: u.Gender, Bio: u.Bio}))
}

// RequestFriend 发送好友请求（通过对方邮箱）
// @Summary App 发送好友请求
// @Tags App
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.AddFriendRequest true "添加好友"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/friends/request [post]
func (h *AppUserHandler) RequestFriend(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	var req apimodel.AddFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	if err := h.friendSvc.RequestFriend(uid, req.FriendEmail); err != nil {
		logger.Errorf("request friend failed: %v", err)
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponseWithMessage("request sent", nil))
}

// ListFriends 好友列表（返回详细资料）
// @Summary App 好友列表（含详细资料）
// @Tags App
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页条数"
// @Success 200 {object} model.APIResponse{data=model.FriendDetailListResponse}
// @Failure 401 {object} model.APIResponse
// @Router /app/friends [get]
func (h *AppUserHandler) ListFriends(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	users, total, err := h.friendSvc.ListFriendDetails(uid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, apimodel.MsgInternalError))
		return
	}
	items := make([]apimodel.FriendDetail, 0, len(users))
	for _, u := range users {
		items = append(items, apimodel.FriendDetail{ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar, Gender: u.Gender, Bio: u.Bio})
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.FriendDetailListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}))
}

// ListPendingRequests 待处理好友请求
// @Summary App 待处理好友请求
// @Tags App
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页条数"
// @Success 200 {object} model.APIResponse{data=object}
// @Failure 401 {object} model.APIResponse
// @Router /app/friends/requests [get]
func (h *AppUserHandler) ListPendingRequests(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("page_size"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	reqIDs, requesterIDs, total, err := h.friendSvc.ListPending(uid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, apimodel.MsgInternalError))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"request_ids": reqIDs, "requester_ids": requesterIDs, "total": total, "page": page, "page_size": pageSize}))
}

// AcceptFriendRequest 接受好友请求
// @Summary App 接受好友请求
// @Tags App
// @Security BearerAuth
// @Param request_id path int true "请求ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/friends/requests/{request_id}/accept [post]
func (h *AppUserHandler) AcceptFriendRequest(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	ridStr := c.Param("request_id")
	rid, err := strconv.ParseUint(ridStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	if err := h.friendSvc.AcceptRequest(uid, uint(rid)); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponseWithMessage("accepted", nil))
}

// DeclineFriendRequest 拒绝好友请求
// @Summary App 拒绝好友请求
// @Tags App
// @Security BearerAuth
// @Param request_id path int true "请求ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/friends/requests/{request_id}/decline [post]
func (h *AppUserHandler) DeclineFriendRequest(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	_ = uid // reserved for ownership check
	ridStr := c.Param("request_id")
	rid, err := strconv.ParseUint(ridStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	if err := h.friendSvc.DeclineRequest(uid, uint(rid)); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponseWithMessage("declined", nil))
}

// AppUploadAvatar 上传并更新头像 (单步：上传文件到对象存储并立即更新用户头像字段)
// @Summary App 上传头像并更新资料
// @Tags App
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "头像文件"
// @Success 200 {object} model.APIResponse{data=model.AppUserInfo}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/profile/avatar [post]
func (h *AppUserHandler) AppUploadAvatar(c *gin.Context) {
	// 鉴权
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)

	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "storage not initialized"))
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "missing file"))
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "read file failed"))
		return
	}

	cfg := config.Load()
	maxBytes := int64(cfg.Minio.MaxFileSizeMB) * 1024 * 1024
	if maxBytes > 0 && int64(len(data)) > maxBytes { // 复用全局文件大小限制
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "file too large"))
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// 简单推断：根据扩展名
		ext := strings.ToLower(filepath.Ext(header.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		default:
			contentType = "application/octet-stream"
		}
	}

	// 仅允许图片类型
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "only image allowed"))
		return
	}
	// 若配置了 AllowedMIMEs 进一步校验
	if !appMimeAllowed(cfg.Minio.AllowedMIMEs, contentType) {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "mime not allowed"))
		return
	}

	// 生成对象名：avatar-{uid}-{timestamp}{ext}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if len(ext) > 10 { // 防止异常过长
		ext = ""
	}
	objectName := "avatar-" + strconv.FormatUint(uint64(uid), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ext
	bucket := "app-avatars" // 固定 bucket，必要时可放配置

	url, err := application.ObjectStore.PutObject(c.Request.Context(), bucket, objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}

	// 更新用户头像
	u, err := h.svc.UpdateProfile(uid, "", url, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "update profile failed"))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.AppUserInfo{ID: u.ID, Email: u.Email, Nickname: u.Nickname, Avatar: u.Avatar, Gender: u.Gender, Bio: u.Bio}))
}

// appMimeAllowed 与存储 handler 类似：支持 * / 前缀 / 精确；为空表示不限制
func appMimeAllowed(allowed []string, ct string) bool {
	if len(allowed) == 0 {
		return true
	}
	ct = strings.ToLower(ct)
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" {
			continue
		}
		if a == "*" {
			return true
		}
		if strings.HasSuffix(a, "/*") {
			prefix := strings.TrimSuffix(a, "/*") + "/"
			if strings.HasPrefix(ct, prefix) {
				return true
			}
		} else if a == ct {
			return true
		}
	}
	return false
}
