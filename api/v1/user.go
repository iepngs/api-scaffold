package v1

// ChangePasswordRequest 修改密码请求结构体
// @Description 用户修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"oldpass123"`       // 旧密码
	NewPassword string `json:"new_password" validate:"required,min=8" example:"newpass123"` // 新密码
}

// ProfileResponse 个人资料响应结构体
// @Description 用户个人资料
type ProfileResponse struct {
	ID       int64  `json:"id" example:"uuid123"`                               // 用户 ID
	Username string `json:"username" example:"user123"`                         // 用户名
	Avatar   string `json:"avatar" example:"https://r2.example.com/avatar.jpg"` // 头像 URL
}

// RealNameAuthRequest 实名认证请求结构体
// @Description 用户实名认证请求
type RealNameAuthRequest struct {
	IDCardNumber string `json:"id_card_number" validate:"required" example:"123456789012345678"`                 // 身份证号
	BankCard     string `json:"bank_card" validate:"required" example:"1234567890123456"`                        // 银行卡号
	IDCardFront  string `json:"id_card_front" validate:"required" example:"https://r2.example.com/id_front.jpg"` // 身份证正面
	IDCardBack   string `json:"id_card_back" validate:"required" example:"https://r2.example.com/id_back.jpg"`   // 身份证反面
}
