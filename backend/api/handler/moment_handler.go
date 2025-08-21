package handler

import (
	apimodel "alice/api/model"
	"alice/application"
	momentservice "alice/domain/moment/service"
	"alice/infra/config"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MomentHandler struct{ svc momentservice.MomentService }

func NewMomentHandler(svc momentservice.MomentService) *MomentHandler {
	return &MomentHandler{svc: svc}
}

// PostMoment 发布朋友圈动态
// @Summary App 发布动态
// @Tags App
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.CreateMomentRequest true "动态内容"
// @Success 200 {object} model.APIResponse{data=model.MomentItem}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/moments [post]
func (h *MomentHandler) PostMoment(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	var req apimodel.CreateMomentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	m, err := h.svc.Publish(uid, req.Content, req.Images)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	u, _ := application.AppUserSvc.GetByID(uid)
	// 补全图片 URL（安卓需完整 http(s) 才能显示）
	imgs := m.ParseImages()
	if len(imgs) > 0 {
		fullImgs := make([]string, 0, len(imgs))
		for _, p := range imgs {
			fullImgs = append(fullImgs, fullAvatarURL(p))
		}
		imgs = fullImgs
	}
	item := apimodel.MomentItem{ID: m.ID, UserID: m.UserID, Nickname: u.Nickname, Avatar: fullAvatarURL(u.Avatar), Content: m.Content, Images: imgs, CreatedAt: m.CreatedAt.Unix()}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(item))
}

// ListMoments 全部动态（时间倒序）
// @Summary App 动态列表
// @Tags App
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} model.APIResponse{data=model.MomentListResponse}
// @Failure 401 {object} model.APIResponse
// @Router /app/moments [get]
func (h *MomentHandler) ListMoments(c *gin.Context) {
	_, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	page, pageSize := 1, 20
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		}
	}
	if v := c.Query("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil {
			pageSize = ps
		}
	}
	list, total, err := h.svc.ListAll(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, apimodel.MsgInternalError))
		return
	}
	items := make([]apimodel.MomentItem, 0, len(list))
	for _, m := range list {
		u, _ := application.AppUserSvc.GetByID(m.UserID)
		imgs := m.ParseImages()
		if len(imgs) > 0 {
			fullImgs := make([]string, 0, len(imgs))
			for _, p := range imgs {
				fullImgs = append(fullImgs, fullAvatarURL(p))
			}
			imgs = fullImgs
		}
		items = append(items, apimodel.MomentItem{ID: m.ID, UserID: m.UserID, Nickname: u.Nickname, Avatar: fullAvatarURL(u.Avatar), Content: m.Content, Images: imgs, CreatedAt: m.CreatedAt.Unix()})
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.MomentListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}))
}

// ListUserMoments 查看某个用户的动态
// @Summary App 用户动态
// @Tags App
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页"
// @Success 200 {object} model.APIResponse{data=model.MomentListResponse}
// @Failure 401 {object} model.APIResponse
// @Router /app/users/{user_id}/moments [get]
func (h *MomentHandler) ListUserMoments(c *gin.Context) {
	_, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uidStr := c.Param("user_id")
	uid64, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	page, pageSize := 1, 20
	if v := c.Query("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		}
	}
	if v := c.Query("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil {
			pageSize = ps
		}
	}
	list, total, err := h.svc.ListByUser(uint(uid64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, apimodel.MsgInternalError))
		return
	}
	items := make([]apimodel.MomentItem, 0, len(list))
	for _, m := range list {
		u, _ := application.AppUserSvc.GetByID(m.UserID)
		imgs := m.ParseImages()
		if len(imgs) > 0 {
			fullImgs := make([]string, 0, len(imgs))
			for _, p := range imgs {
				fullImgs = append(fullImgs, fullAvatarURL(p))
			}
			imgs = fullImgs
		}
		items = append(items, apimodel.MomentItem{ID: m.ID, UserID: m.UserID, Nickname: u.Nickname, Avatar: fullAvatarURL(u.Avatar), Content: m.Content, Images: imgs, CreatedAt: m.CreatedAt.Unix()})
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(apimodel.MomentListResponse{Items: items, Total: total, Page: page, PageSize: pageSize}))
}

// UploadImage 上传动态图片（单文件）
// @Summary App 上传动态图片
// @Tags App
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} model.APIResponse{data=object}
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/moments/images [post]
func (h *MomentHandler) UploadImage(c *gin.Context) {
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
	if maxBytes > 0 && int64(len(data)) > maxBytes {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "file too large"))
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "only image allowed"))
		return
	}
	if !appMimeAllowed(cfg.Minio.AllowedMIMEs, contentType) {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "mime not allowed"))
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if len(ext) > 10 {
		ext = ""
	}
	objectName := "moment-" + strconv.FormatUint(uint64(uid), 10) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	bucket := "app-moment-images"
	_, err = application.ObjectStore.PutObject(c.Request.Context(), bucket, objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	relative := "/" + bucket + "/" + objectName
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"path": relative, "url": fullAvatarURL(relative)}))
}

// DeleteMoment 删除自己的动态
// @Summary App 删除动态
// @Tags App
// @Security BearerAuth
// @Param moment_id path int true "动态ID"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.APIResponse
// @Failure 401 {object} model.APIResponse
// @Router /app/moments/{moment_id} [delete]
func (h *MomentHandler) DeleteMoment(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	midStr := c.Param("moment_id")
	mid, err := strconv.ParseUint(midStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, apimodel.MsgInvalidRequest))
		return
	}
	if err := h.svc.Delete(uid, uint(mid)); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponseWithMessage("deleted", nil))
}

// fullAvatarURL 复制自 AppUserHandler（避免循环引用），后续可抽公共
func fullAvatarURL(raw string) string {
	if raw == "" {
		return ""
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return raw
	}
	cfg := config.Load()
	base := cfg.Minio.BaseURL
	if base == "" {
		scheme := "http"
		if cfg.Minio.UseSSL {
			scheme = "https"
		}
		base = scheme + "://" + cfg.Minio.Endpoint
	}
	if strings.HasSuffix(base, "/") {
		base = strings.TrimRight(base, "/")
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return base + raw
}
