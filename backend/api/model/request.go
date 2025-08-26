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
	RoleIDs       []uint `json:"role_ids" example:"1,2"`
	PermissionIDs []uint `json:"permission_ids" example:"1,2"`
	MenuIDs       []uint `json:"menu_ids" example:"1,2"`
}

// AdminCreateUserRequest 管理员新增用户
type AdminCreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive banned"`
}

// AdminUpdateUserRequest 管理员更新用户
type AdminUpdateUserRequest struct {
	Email    string  `json:"email" binding:"omitempty,email"`
	Status   string  `json:"status" binding:"omitempty,oneof=active inactive banned"`
	Password *string `json:"password" binding:"omitempty,min=6"`
}

// ListUsersQuery 用户列表查询
type ListUsersQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}
