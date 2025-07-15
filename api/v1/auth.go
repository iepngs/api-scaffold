package v1

// LoginRequest 登录请求结构体
// @Description 用户登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"user123"`  // 用户名
	Password string `json:"password" validate:"required" example:"pass1234"` // 密码
}

// LoginResponse 登录响应结构体
// @Description 用户登录响应
type LoginResponse struct {
	Token  string `json:"token" example:"jwt.token.here"`    // JWT 令牌
	Expire int64  `json:"expire" example:"jwt.token.expire"` // JWT 令牌有效期
}

// LogoutResponse 退出响应结构体
// @Description 用户退出响应
type LogoutResponse struct {
	Message string `json:"message" example:"Logged out successfully"` // 退出消息
}

// 审核状态，approved/rejected
