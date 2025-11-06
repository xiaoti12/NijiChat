package handlers

import (
	"seiyuu-chat/models"
	"seiyuu-chat/services"
	"seiyuu-chat/utils"

	"github.com/gin-gonic/gin"
)

// SchedulerHandler 智能调度器API处理器
type SchedulerHandler struct {
	schedulerService *services.SchedulerService
}

// NewSchedulerHandler 创建调度器处理器实例
func NewSchedulerHandler(schedulerService *services.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
	}
}

// SelectSeiyuu 智能选择声优
// POST /api/scheduler/select
func (h *SchedulerHandler) SelectSeiyuu(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.SchedulerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证消息不能为空
	if req.Message == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 调用调度服务
	result, err := h.schedulerService.SelectSeiyuu(ctx, &req)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	// 返回成功响应
	response := &models.SchedulerResponse{
		Success: true,
		Data:    result,
	}

	utils.SuccessResponse(c, response)
}
