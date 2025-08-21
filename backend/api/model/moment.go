package model

// CreateMomentRequest 发送动态
type CreateMomentRequest struct {
	Content string   `json:"content" binding:"required"` // 文本内容
	Images  []string `json:"images"`                     // 已上传后的相对路径 /bucket/object，最多9张
}

// MomentItem 动态条目
type MomentItem struct {
	ID        uint     `json:"id"`
	UserID    uint     `json:"user_id"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Content   string   `json:"content"`
	Images    []string `json:"images"`
	CreatedAt int64    `json:"created_at"`
}

// MomentListResponse 列表响应
type MomentListResponse struct {
	Items    []MomentItem `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}
