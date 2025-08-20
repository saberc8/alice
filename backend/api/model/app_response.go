package model

// AppUserInfo 移动端用户信息
type AppUserInfo struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	Bio      string `json:"bio"`
}

// AppAuthResponse 注册后下发 token + 基本资料
type AppAuthResponse struct {
	User  AppUserInfo `json:"user"`
	Token string      `json:"token"`
}

// FriendListResponse 好友列表
type FriendListResponse struct {
	IDs      []uint `json:"ids"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// FriendDetail 详细好友资料（用于列表）
type FriendDetail struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	Bio      string `json:"bio"`
}

// FriendDetailListResponse 返回详细好友资料
type FriendDetailListResponse struct {
	Items    []FriendDetail `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}
