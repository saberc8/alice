package model

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// AssignIDsRequest 通用批量ID分配请求
type AssignIDsRequest struct {
	RoleIDs       []string `json:"role_ids" example:"1,2"`
	PermissionIDs []string `json:"permission_ids" example:"1,2"`
	MenuIDs       []string `json:"menu_ids" example:"1,2"`
}
