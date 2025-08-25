package chat

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
	"alice/infra/config"
)

type GroupHandler struct{}

func NewGroupHandler() *GroupHandler { return &GroupHandler{} }

// CreateGroup 创建群聊
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	var req struct {
		Name      string `json:"name" binding:"required"`
		MemberIDs []uint `json:"member_ids"`
		Avatar    string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	g, err := application.GroupSvc.Create(uid, req.Name, req.MemberIDs, req.Avatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(g))
}

// SearchGroups 按名称搜索
func (h *GroupHandler) SearchGroups(c *gin.Context) {
	q := c.Query("q")
	groups, _ := application.GroupSvc.Search(q, 20)
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"items": groups}))
}

// JoinGroup 加入群
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	if err := application.GroupSvc.Join(uint(gid64), uid); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"group_id": uint(gid64), "joined_at": time.Now().Unix()}))
}

// GroupMessages 历史
func (h *GroupHandler) GroupMessages(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	okMember, err := application.GroupSvc.IsMember(uint(gid64), uid)
	if err != nil || !okMember {
		c.JSON(http.StatusForbidden, apimodel.ErrorResponse(apimodel.CodeForbidden, "not a member"))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	msgs, total, err := application.GroupSvc.ListMessages(uint(gid64), page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	// 收集 sender ids
	idsSet := make(map[uint]struct{})
	for _, m := range msgs {
		if m != nil {
			idsSet[m.SenderID] = struct{}{}
		}
	}
	ids := make([]uint, 0, len(idsSet))
	for id := range idsSet {
		ids = append(ids, id)
	}
	users, _ := application.AppUserSvc.GetByIDs(ids)
	cfg := config.Load()
	full := func(raw string) string {
		if raw == "" {
			return ""
		}
		lw := strings.ToLower(raw)
		if strings.HasPrefix(lw, "http://") || strings.HasPrefix(lw, "https://") {
			return raw
		}
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
	userMap := make(map[uint]gin.H, len(users))
	for _, u := range users {
		if u != nil {
			userMap[u.ID] = gin.H{"id": u.ID, "nickname": u.Nickname, "avatar": full(u.Avatar)}
		}
	}
	enriched := make([]gin.H, 0, len(msgs))
	for _, m := range msgs {
		if m != nil {
			enriched = append(enriched, gin.H{"id": m.ID, "group_id": m.GroupID, "sender_id": m.SenderID, "type": m.Type, "content": m.Content, "created_at": m.CreatedAt, "sender": userMap[m.SenderID]})
		}
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"items": enriched, "total": total, "page": page, "page_size": pageSize}))
}

// UpdateGroup 基础信息修改（仅群主）
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	var req struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	g, err := application.GroupSvc.UpdateGroup(uid, uint(gid64), req.Name, req.Avatar)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(g))
}

// UploadGroupAvatar 上传群头像（仅群主）
func (h *GroupHandler) UploadGroupAvatar(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	g, err := application.GroupSvc.Get(uint(gid64))
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	if g.OwnerID != uid {
		c.JSON(http.StatusForbidden, apimodel.ErrorResponse(apimodel.CodeForbidden, "no permission"))
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
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "read failed"))
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
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if len(ext) > 10 {
		ext = ""
	}
	if application.ObjectStore == nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, "storage not initialized"))
		return
	}
	objectName := "group-avatar-" + strconv.FormatUint(gid64, 10) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	bucket := "app-group-avatars"
	if _, err = application.ObjectStore.PutObject(c.Request.Context(), bucket, objectName, data, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, err.Error()))
		return
	}
	relative := "/" + bucket + "/" + objectName
	// 保存
	if _, err = application.GroupSvc.UpdateGroup(uid, uint(gid64), "", relative); err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, err.Error()))
		return
	}
	full := func(raw string) string {
		if raw == "" {
			return ""
		}
		lw := strings.ToLower(raw)
		if strings.HasPrefix(lw, "http://") || strings.HasPrefix(lw, "https://") {
			return raw
		}
		cfg2 := config.Load()
		base := cfg2.Minio.BaseURL
		if base == "" {
			scheme := "http"
			if cfg2.Minio.UseSSL {
				scheme = "https"
			}
			base = scheme + "://" + cfg2.Minio.Endpoint
		}
		if strings.HasSuffix(base, "/") {
			base = strings.TrimRight(base, "/")
		}
		if !strings.HasPrefix(raw, "/") {
			raw = "/" + raw
		}
		return base + raw
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"path": relative, "url": full(relative)}))
}

// ListMembers 返回群成员基础信息
func (h *GroupHandler) ListMembers(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	okMember, err := application.GroupSvc.IsMember(uint(gid64), uid)
	if err != nil || !okMember {
		c.JSON(http.StatusForbidden, apimodel.ErrorResponse(apimodel.CodeForbidden, "not a member"))
		return
	}
	ids, _ := application.GroupSvc.ListMemberIDs(uint(gid64))
	// 复用 app user service 获取用户
	users, _ := application.AppUserSvc.GetByIDs(ids)
	out := make([]gin.H, 0, len(users))
	cfg := config.Load()
	full := func(raw string) string {
		if raw == "" {
			return ""
		}
		lw := strings.ToLower(raw)
		if strings.HasPrefix(lw, "http://") || strings.HasPrefix(lw, "https://") {
			return raw
		}
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
	for _, u := range users {
		out = append(out, gin.H{"id": u.ID, "nickname": u.Nickname, "avatar": full(u.Avatar)})
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"members": out}))
}

// AddMembers 添加成员（仅群主）
func (h *GroupHandler) AddMembers(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	if err := application.GroupSvc.AddMembers(uid, uint(gid64), req.UserIDs); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"group_id": gid64, "added": req.UserIDs}))
}

// RemoveMember 移除成员（仅群主）
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	gidStr := c.Param("group_id")
	gid64, _ := strconv.ParseUint(gidStr, 10, 64)
	if gid64 == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid group id"))
		return
	}
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	if err := application.GroupSvc.RemoveMember(uid, uint(gid64), req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"group_id": gid64, "removed": req.UserID}))
}

// MarkReadGroup 群聊已读上报（将成员游标推进到 before_msg_id，如果更大）
func (h *GroupHandler) MarkReadGroup(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, apimodel.MsgUnauthorized))
		return
	}
	uid, _ := idAny.(uint)
	var req struct {
		GroupID     uint `json:"group_id" binding:"required"`
		BeforeMsgID uint `json:"before_msg_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.GroupID == 0 || req.BeforeMsgID == 0 {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	okMember, err := application.GroupSvc.IsMember(req.GroupID, uid)
	if err != nil || !okMember {
		c.JSON(http.StatusForbidden, apimodel.ErrorResponse(apimodel.CodeForbidden, "not a member"))
		return
	}
	if err := application.GroupSvc.UpdateLastRead(c.Request.Context(), req.GroupID, uid, req.BeforeMsgID); err != nil {
		c.JSON(http.StatusInternalServerError, apimodel.ErrorResponse(apimodel.CodeInternalError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"group_id": req.GroupID, "before_msg_id": req.BeforeMsgID, "ts": time.Now().Unix()}))
}
