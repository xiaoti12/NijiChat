package middleware

import (
	"os"
	"seiyuu-chat/models"
	"seiyuu-chat/router"
	"seiyuu-chat/utils"
)

// AdminAuthMiddleware JWT认证中间件 - 仅管理员可访问
func AdminAuthMiddleware() router.HandlerFunc {
	return func(c *router.Context) {
		// 开发环境跳过认证
		if utils.IsDevelopmentMode() {
			// 设置默认管理员信息用于开发环境
			c.Set("admin_id", "dev-admin")
			c.Set("username", "developer")
			c.Next()
			return
		}

		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedError(c, models.ErrUnauthorized)
			c.Abort()
			return
		}

		// 提取token
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.UnauthorizedError(c, models.ErrInvalidToken)
			c.Abort()
			return
		}

		// 从环境变量获取JWT密钥
		jwtSecret := os.Getenv("ADMIN_JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "dev-jwt-secret-key" // 开发环境默认值
		}

		// 验证token
		claims, err := utils.ValidateJWT(token, jwtSecret)
		if err != nil {
			utils.UnauthorizedError(c, models.ErrTokenExpired)
			c.Abort()
			return
		}

		// 将管理员信息存入context
		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetAdminID 从context中获取管理员ID
func GetAdminID(c *router.Context) string {
	if adminID, exists := c.Get("admin_id"); exists {
		if id, ok := adminID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUsername 从context中获取用户名
func GetUsername(c *router.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}
