package handlers

import (
	"context"
	"database/sql"

	"seiyuu-chat/database"
	"seiyuu-chat/models"
	"seiyuu-chat/router"
	"seiyuu-chat/utils"

	"github.com/google/uuid"
)

// AdminHandler 管理员功能API处理器
type AdminHandler struct {
	db *database.D1Client
}

// NewAdminHandler 创建管理员处理器实例
func NewAdminHandler(db *database.D1Client) *AdminHandler {
	return &AdminHandler{
		db: db,
	}
}

// Login 管理员登录
// POST /api/admin/login
func (h *AdminHandler) Login(c *router.Context) {
	ctx := c.Request.Context()

	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if err := utils.ValidateUsername(req.Username); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 查询管理员
	admin, err := h.getAdminByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.UnauthorizedError(c, models.ErrInvalidCredentials)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, admin.PasswordHash) {
		utils.UnauthorizedError(c, models.ErrInvalidCredentials)
		return
	}

	// 生成JWT token
	jwtSecret := utils.GetEnv("ADMIN_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-jwt-secret-key"
	}

	token, expiresAt, err := utils.GenerateJWT(admin.ID, admin.Username, jwtSecret)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 返回登录响应
	response := &models.AdminLoginResponse{
		Token:     token,
		Username:  admin.Username,
		ExpiresAt: expiresAt,
	}

	utils.SuccessWithMessage(c, "登录成功", response)
}

// CreateAdmin 创建管理员账号（仅用于初始化，实际应该由超级管理员调用）
// POST /api/admin/create
func (h *AdminHandler) CreateAdmin(c *router.Context) {
	ctx := c.Request.Context()

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if err := utils.ValidateUsername(req.Username); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 检查用户名是否已存在
	_, err := h.getAdminByUsername(ctx, req.Username)
	if err == nil {
		utils.BadRequestError(c, models.ErrAdminAlreadyExists)
		return
	}

	// 哈希密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 创建管理员
	id := uuid.New().String()
	query := `
		INSERT INTO admins (id, username, password_hash, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`

	_, err = h.db.Exec(ctx, query, id, req.Username, passwordHash)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "管理员创建成功", map[string]interface{}{
		"id":       id,
		"username": req.Username,
	})
}

// getAdminByUsername 根据用户名获取管理员
func (h *AdminHandler) getAdminByUsername(ctx context.Context, username string) (*models.Admin, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM admins
		WHERE username = ?
	`

	row := h.db.QueryRow(ctx, query, username)

	var admin models.Admin
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &admin, nil
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
