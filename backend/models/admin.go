package models

import (
	"time"
)

// Admin 管理员数据模型
type Admin struct {
	ID           string    `json:"id"`            // 管理员ID (UUID)
	Username     string    `json:"username"`      // 用户名
	PasswordHash string    `json:"-"`             // 密码哈希 (不返回给前端)
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Token     string `json:"token"`      // JWT token
	Username  string `json:"username"`   // 用户名
	ExpiresAt int64  `json:"expires_at"` // 过期时间戳
}

// JWTClaims JWT声明
type JWTClaims struct {
	AdminID  string `json:"admin_id"`
	Username string `json:"username"`
	IssuedAt int64  `json:"iat"`
	ExpiresAt int64 `json:"exp"`
}

// Validate 验证管理员数据
func (a *Admin) Validate() error {
	if a.Username == "" {
		return ErrInvalidAdminUsername
	}
	if a.PasswordHash == "" {
		return ErrInvalidAdminPassword
	}
	return nil
}

// ToPublic 返回公开的管理员信息（不包含密码哈希）
func (a *Admin) ToPublic() *Admin {
	return &Admin{
		ID:        a.ID,
		Username:  a.Username,
		CreatedAt: a.CreatedAt,
	}
}
