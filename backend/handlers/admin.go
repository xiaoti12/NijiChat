package handlers

import (
	"seiyuu-chat/models"
	"seiyuu-chat/router"
	"seiyuu-chat/utils"
)

// AdminHandler 管理员功能API处理器
type AdminHandler struct {
	// 移除数据库依赖，使用硬编码账号认证
}

// NewAdminHandler 创建管理员处理器实例
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// Login 管理员登录
// POST /api/admin/login
func (h *AdminHandler) Login(c *router.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据格式
	if err := utils.ValidateUsername(req.Username); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 使用硬编码账号验证
	if !utils.ValidateAdminCredentials(req.Username, req.Password) {
		utils.UnauthorizedError(c, models.ErrInvalidCredentials)
		return
	}

	// 生成JWT token
	jwtSecret := utils.GetEnv("ADMIN_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-jwt-secret-key"
	}

	// 生成管理员ID和获取用户名
	adminID := utils.GenerateAdminID(req.Username)
	username := utils.GetAdminUsername(req.Username)

	token, expiresAt, err := utils.GenerateJWT(adminID, username, jwtSecret)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 返回登录响应
	response := &models.AdminLoginResponse{
		Token:     token,
		Username:  username,
		ExpiresAt: expiresAt,
	}

	utils.SuccessWithMessage(c, "登录成功", response)
}


// GetProfile 获取当前管理员信息
// GET /api/admin/profile
func (h *AdminHandler) GetProfile(c *router.Context) {
	// 从context中获取管理员ID（由认证中间件设置）
	adminID := c.GetString("admin_id")
	username := c.GetString("username")

	if adminID == "" {
		utils.UnauthorizedError(c, models.ErrUnauthorized)
		return
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"id":       adminID,
		"username": username,
	})
}

// HealthCheck 健康检查
// GET /api/health
func (h *AdminHandler) HealthCheck(c *router.Context) {
	utils.SuccessResponse(c, map[string]interface{}{
		"status":  "ok",
		"service": "seiyuu-chat-backend",
		"version": "1.0.0",
	})
}
