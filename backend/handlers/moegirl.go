package handlers

import (
	"seiyuu-chat/models"
	"seiyuu-chat/services"
	"seiyuu-chat/utils"

	"github.com/gin-gonic/gin"
)

// MoegirlHandler 萌娘百科集成API处理器
type MoegirlHandler struct {
	moegirlService *services.MoegirlService
	aiService      *services.AIService
}

// NewMoegirlHandler 创建萌娘百科处理器实例
func NewMoegirlHandler(moegirlService *services.MoegirlService, aiService *services.AIService) *MoegirlHandler {
	return &MoegirlHandler{
		moegirlService: moegirlService,
		aiService:      aiService,
	}
}

// GetRawData 获取萌娘百科原始数据（管理员）
// GET /api/admin/moegirl/:name
func (h *MoegirlHandler) GetRawData(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if name == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 从萌娘百科获取原始数据
	rawData, err := h.moegirlService.GetRawData(ctx, name)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 返回响应
	response := &models.MoegirlRawDataResponse{
		Success: true,
		Data:    rawData,
	}

	utils.SuccessResponse(c, response)
}

// ProcessProfile AI处理资料转Markdown（管理员）
// POST /api/admin/process-profile
func (h *MoegirlHandler) ProcessProfile(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.ProcessProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if req.RawText == "" || req.SeiyuuName == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 使用AI处理资料
	processedProfile, err := h.aiService.ProcessSeiyuuProfile(ctx, req.RawText, req.SeiyuuName)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 返回响应
	response := &models.ProcessProfileResponse{
		Success: true,
		Data:    processedProfile,
	}

	utils.SuccessResponse(c, response)
}

// SearchSeiyuu 搜索声优页面（管理员）
// GET /api/admin/moegirl/search
func (h *MoegirlHandler) SearchSeiyuu(c *gin.Context) {
	ctx := c.Request.Context()
	keyword := c.Query("keyword")

	if keyword == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 搜索萌娘百科
	results, err := h.moegirlService.SearchSeiyuu(ctx, keyword)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"results": results,
	})
}
