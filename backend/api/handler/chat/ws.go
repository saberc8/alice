package chat

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	apimodel "alice/api/model"
	"alice/application"
	appentity "alice/domain/appuser/entity"
	appuserservice "alice/domain/appuser/service"
	chatentity "alice/domain/chat/entity"
	chatservice "alice/domain/chat/service"
	"alice/infra/config"
	"alice/pkg/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub 管理用户连接与转发
type Hub struct {
	mu        sync.RWMutex
	conns     map[uint]*websocket.Conn
	chat      chatservice.ChatService
	appUserSv appuserservice.AppUserService
	// group service via application package (quick access)
}

func NewHub(s chatservice.ChatService, appUserSv appuserservice.AppUserService) *Hub {
	return &Hub{conns: make(map[uint]*websocket.Conn), chat: s, appUserSv: appUserSv}
}

// WS 处理 WebSocket 连接
func (h *Hub) WS(c *gin.Context) {
	uid, err := getAppUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("ws upgrade failed: %v", err)
		return
	}
	// 注册连接
	h.mu.Lock()
	h.conns[uid] = conn
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.conns, uid)
		h.mu.Unlock()
		_ = conn.Close()
	}()

	for {
		var payload struct {
			Type    string `json:"type"`
			To      uint   `json:"to"`
			GroupID uint   `json:"group_id"`
			Content string `json:"content"`
			MsgType string `json:"msg_type"`
		}
		if err := conn.ReadJSON(&payload); err != nil {
			logger.Infof("ws read closed: %v", err)
			return
		}
		// 群聊消息
		if payload.GroupID > 0 {
			gm, err := application.GroupSvc.SendMessage(payload.GroupID, uid, payload.MsgType, payload.Content)
			if err != nil {
				_ = conn.WriteJSON(gin.H{"error": err.Error()})
				continue
			}
			// 富化 sender 信息
			users, _ := h.appUserSv.GetByIDs([]uint{gm.SenderID})
			var sender gin.H
			if len(users) > 0 {
				sender = h.userInfoMap(users)[gm.SenderID]
			}
			resp := gin.H{
				"type":         "group_message",
				"group_id":     gm.GroupID,
				"message_id":   gm.ID,
				"sender_id":    gm.SenderID,
				"content":      gm.Content,
				"message_type": gm.Type,
				"created_at":   gm.CreatedAt,
				"sender":       sender,
			}
			// 只推送给群成员在线连接
			memberIDs, _ := application.GroupSvc.ListMemberIDs(gm.GroupID)
			memberSet := make(map[uint]struct{}, len(memberIDs))
			for _, mid := range memberIDs {
				memberSet[mid] = struct{}{}
			}
			// 自己先回显
			_ = conn.WriteJSON(resp)
			h.mu.RLock()
			for mid, c2 := range h.conns {
				if mid == uid { // 已回显
					continue
				}
				if _, ok := memberSet[mid]; ok {
					_ = c2.WriteJSON(resp)
				}
			}
			h.mu.RUnlock()
			continue
		}
		// 私聊消息
		msg, err := h.chat.Send(uid, payload.To, payload.Content, firstNonEmpty(payload.MsgType, payload.Type))
		if err != nil {
			_ = conn.WriteJSON(gin.H{"error": err.Error()})
			continue
		}
		enriched := h.enrichSingleMessage(msg)
		_ = conn.WriteJSON(enriched)
		h.mu.RLock()
		peer := h.conns[payload.To]
		h.mu.RUnlock()
		if peer != nil {
			_ = peer.WriteJSON(enriched)
		}
	}
}

// History 拉取历史记录（REST）
func (h *Hub) History(c *gin.Context) {
	uid, err := getAppUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
		return
	}
	peerID, _ := parseUintParam(c, "peer_id")
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)
	items, total, err := h.chat.History(uid, peerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	enriched := h.enrichMessages(items)
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{
		"items":     enriched,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// Conversations 最近会话列表
func (h *Hub) Conversations(c *gin.Context) {
	uid, err := getAppUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
		return
	}
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)
	items, _, err := h.chat.RecentConversations(uid, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	// 批量获取 peer 用户信息
	peerIDs := make([]uint, 0, len(items))
	for _, it := range items {
		peerIDs = append(peerIDs, it.PeerID)
	}
	users, _ := h.appUserSv.GetByIDs(peerIDs)
	userMap := make(map[uint]gin.H, len(users))
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
		userMap[u.ID] = gin.H{"id": u.ID, "nickname": u.Nickname, "avatar": full(u.Avatar)}
	}
	type convoItem struct {
		Data gin.H
		Ts   time.Time
	}
	all := make([]convoItem, 0, len(items)+8)
	for _, it := range items {
		var ts time.Time
		lm := gin.H{}
		if it.LastMessage != nil {
			lm = gin.H{"id": it.LastMessage.ID, "sender_id": it.LastMessage.SenderID, "receiver_id": it.LastMessage.ReceiverID, "type": it.LastMessage.Type, "content": it.LastMessage.Content, "created_at": it.LastMessage.CreatedAt}
			ts = it.LastMessage.CreatedAt
		}
		all = append(all, convoItem{Ts: ts, Data: gin.H{
			"type":         "private",
			"peer_id":      it.PeerID,
			"last_message": lm,
			"unread_count": it.UnreadCount,
			"peer":         userMap[it.PeerID],
		}})
	}
	// 群会话
	groups, _, _ := application.GroupSvc.ListUserGroups(uid, 1, 200)
	for _, g := range groups {
		var lastMsg *chatentity.GroupMessage
		if mr, ok := application.GroupSvc.(interface {
			ListMessages(uint, int, int) ([]*chatentity.GroupMessage, int64, error)
		}); ok {
			msgs, _, _ := mr.ListMessages(g.ID, 1, 1)
			if len(msgs) > 0 {
				lastMsg = msgs[0]
			}
		}
		lm := gin.H{}
		var ts time.Time
		if lastMsg != nil {
			lm = gin.H{"id": lastMsg.ID, "group_id": lastMsg.GroupID, "sender_id": lastMsg.SenderID, "type": lastMsg.Type, "content": lastMsg.Content, "created_at": lastMsg.CreatedAt}
			ts = lastMsg.CreatedAt
		}
		unread, _ := application.GroupSvc.CountUnread(c.Request.Context(), g.ID, uid)
		all = append(all, convoItem{Ts: ts, Data: gin.H{
			"type":         "group",
			"group":        gin.H{"id": g.ID, "name": g.Name, "avatar": g.Avatar},
			"last_message": lm,
			"unread_count": unread,
		}})
	}
	// 排序 按 ts desc
	sort.SliceStable(all, func(i, j int) bool { return all[i].Ts.After(all[j].Ts) })
	respItems := make([]gin.H, 0, len(all))
	for _, ci := range all {
		respItems = append(respItems, ci.Data)
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"items": respItems, "total": int64(len(all)), "page": page, "page_size": pageSize}))
}

// enrichSingleMessage 富化单条消息（带 sender / receiver 基础信息）
func (h *Hub) enrichSingleMessage(m *chatentity.Message) gin.H {
	if m == nil {
		return gin.H{}
	}
	users, _ := h.appUserSv.GetByIDs([]uint{m.SenderID, m.ReceiverID})
	userMap := h.userInfoMap(users)
	return gin.H{
		"id":          m.ID,
		"sender_id":   m.SenderID,
		"receiver_id": m.ReceiverID,
		"type":        m.Type,
		"content":     m.Content,
		"is_read":     m.IsRead,
		"read_at":     m.ReadAt,
		"created_at":  m.CreatedAt,
		"sender":      userMap[m.SenderID],
		"receiver":    userMap[m.ReceiverID],
	}
}

func (h *Hub) enrichMessages(items []*chatentity.Message) []gin.H {
	// 保持原仓库返回顺序：DESC（最新在前）。前端已有逻辑 items.reversed 来得到升序展示（最新在列表底部）。
	if len(items) == 0 {
		return []gin.H{}
	}
	idSet := make(map[uint]struct{}, len(items)*2)
	for _, m := range items {
		if m == nil {
			continue
		}
		idSet[m.SenderID] = struct{}{}
		idSet[m.ReceiverID] = struct{}{}
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	users, _ := h.appUserSv.GetByIDs(ids)
	userMap := h.userInfoMap(users)
	out := make([]gin.H, 0, len(items))
	for _, m := range items { // 不再反转
		if m == nil {
			continue
		}
		out = append(out, gin.H{
			"id":          m.ID,
			"sender_id":   m.SenderID,
			"receiver_id": m.ReceiverID,
			"type":        m.Type,
			"content":     m.Content,
			"is_read":     m.IsRead,
			"read_at":     m.ReadAt,
			"created_at":  m.CreatedAt,
			"sender":      userMap[m.SenderID],
			"receiver":    userMap[m.ReceiverID],
		})
	}
	return out
}

func (h *Hub) userInfoMap(users []*appentity.AppUser) map[uint]gin.H {
	m := make(map[uint]gin.H, len(users))
	if len(users) == 0 {
		return m
	}
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
		if u == nil {
			continue
		}
		m[u.ID] = gin.H{"id": u.ID, "nickname": u.Nickname, "avatar": full(u.Avatar)}
	}
	return m
}

// local helper (duplicate of chat service) to avoid import cycle
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// MarkRead 标记消息为已读（将对方->我，ID <= before_id 的未读置已读）
func (h *Hub) MarkRead(c *gin.Context) {
	uid, err := getAppUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
		return
	}
	var req struct {
		PeerID   uint `json:"peer_id" binding:"required"`
		BeforeID uint `json:"before_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "invalid request"))
		return
	}
	if err := h.chat.MarkRead(uid, req.PeerID, req.BeforeID); err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"peer_id": req.PeerID, "before_id": req.BeforeID, "ts": time.Now().Unix()}))
}

// UploadImage 聊天图片上传（仅允许图片 mime）
func (h *Hub) UploadImage(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
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
	if !appMimeAllowed(cfg.Minio.AllowedMIMEs, contentType) { // 复用 app 头像/动态校验逻辑
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "mime not allowed"))
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if len(ext) > 10 {
		ext = ""
	}
	objectName := "chat-" + strconv.FormatUint(uint64(uid), 10) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	bucket := "app-chat-images"
	_, err = application.ObjectStore.PutObject(c.Request.Context(), bucket, objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	relative := "/" + bucket + "/" + objectName
	// 使用与 Conversations 中相同的 full() 逻辑构造完整 URL
	fullURL := func(raw string) string {
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
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"path": relative, "url": fullURL(relative)}))
}

// UploadVideo 聊天视频上传（仅允许 video mime）
func (h *Hub) UploadVideo(c *gin.Context) {
	idAny, ok := c.Get("app_user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, apimodel.ErrorResponse(apimodel.CodeUnauthorized, "unauthorized"))
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
	if !strings.HasPrefix(strings.ToLower(contentType), "video/") {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "only video allowed"))
		return
	}
	if !appMimeAllowed(cfg.Minio.AllowedMIMEs, contentType) { // 复用校验逻辑
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, "mime not allowed"))
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if len(ext) > 10 {
		ext = ""
	}
	objectName := "chat-video-" + strconv.FormatUint(uint64(uid), 10) + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	bucket := "app-chat-videos"
	_, err = application.ObjectStore.PutObject(c.Request.Context(), bucket, objectName, data, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, apimodel.ErrorResponse(apimodel.CodeBadRequest, err.Error()))
		return
	}
	relative := "/" + bucket + "/" + objectName
	fullURL := func(raw string) string {
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
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{"path": relative, "url": fullURL(relative)}))
}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	v := c.Param(name)
	if v == "" {
		return 0, nil
	}
	// gin 的 Param 转换
	var id uint
	_, err := fmt.Sscan(v, &id)
	return id, err
}

func parseIntQuery(c *gin.Context, name string, def int) int {
	if s := c.Query(name); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return def
}

// getAppUserID 从 context 读取 app_user_id
func getAppUserID(c *gin.Context) (uint, error) {
	v, ok := c.Get("app_user_id")
	if !ok {
		return 0, fmt.Errorf("no app_user_id in context")
	}
	if id, ok := v.(uint); ok {
		return id, nil
	}
	// 有些中间件可能以 float64 存储
	if f, ok := v.(float64); ok {
		return uint(f), nil
	}
	return 0, fmt.Errorf("invalid app_user_id type")
}

// appMimeAllowed 复制自 app_user_handler，保持一致
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
