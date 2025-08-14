package chat

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	apimodel "alice/api/model"
	chatservice "alice/domain/chat/service"
	"alice/pkg/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub 管理用户连接与转发
type Hub struct {
	mu    sync.RWMutex
	conns map[uint]*websocket.Conn
	chat  chatservice.ChatService
}

func NewHub(s chatservice.ChatService) *Hub {
	return &Hub{conns: make(map[uint]*websocket.Conn), chat: s}
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
			Content string `json:"content"`
		}
		if err := conn.ReadJSON(&payload); err != nil {
			logger.Infof("ws read closed: %v", err)
			return
		}
		msg, err := h.chat.Send(uid, payload.To, payload.Content, payload.Type)
		if err != nil {
			_ = conn.WriteJSON(gin.H{"error": err.Error()})
			continue
		}
		// 给自己回显
		_ = conn.WriteJSON(msg)
		// 推给对方在线
		h.mu.RLock()
		peer := h.conns[payload.To]
		h.mu.RUnlock()
		if peer != nil {
			_ = peer.WriteJSON(msg)
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
	c.JSON(http.StatusOK, apimodel.SuccessResponse(gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
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
	c.JSON(http.StatusOK, apimodel.SuccessResponse(nil))
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
