package model

// AppRegisterRequest 移动端注册
type AppRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
}

// AppLoginRequest 移动端登录
type AppLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AppUpdateProfileRequest 更新资料
type AppUpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=30"`
	// Avatar 允许前端直接提交相对路径（/app-avatars/xxx.jpg）或完整 URL；
	// 后端会在下发时自动补全 base-url。这里移除 url 校验以支持相对路径。
	Avatar string `json:"avatar" binding:"omitempty"`
	Gender string `json:"gender" binding:"omitempty,oneof=male female other"`
	Bio    string `json:"bio" binding:"omitempty,max=160"`
}

// AddFriendRequest 添加好友
type AddFriendRequest struct {
	FriendEmail string `json:"friend_email" binding:"required,email"`
}
