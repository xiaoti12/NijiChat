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
}

// NewMoegirlHandler 创建萌娘百科处理器实例
func NewMoegirlHandler(moegirlService *services.MoegirlService) *MoegirlHandler {
	return &MoegirlHandler{
		moegirlService: moegirlService,
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
