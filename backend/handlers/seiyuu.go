package handlers

import (
	"seiyuu-chat/models"
	"seiyuu-chat/router"
	"seiyuu-chat/services"
	"seiyuu-chat/utils"
)

// SeiyuuHandler 声优相关API处理器
type SeiyuuHandler struct {
	seiyuuService *services.SeiyuuService
}

// NewSeiyuuHandler 创建声优处理器实例
func NewSeiyuuHandler(seiyuuService *services.SeiyuuService) *SeiyuuHandler {
	return &SeiyuuHandler{
		seiyuuService: seiyuuService,
	}
}

// GetAllSeiyuu 获取所有已发布的声优列表
// GET /api/seiyuu
func (h *SeiyuuHandler) GetAllSeiyuu(c *router.Context) {
	ctx := c.Request.Context()

	seiyuus, err := h.seiyuuService.GetAllSeiyuu(ctx)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessResponse(c, seiyuus)
}

// GetSeiyuuByID 获取声优详细资料
// GET /api/seiyuu/:id
func (h *SeiyuuHandler) GetSeiyuuByID(c *router.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	seiyuu, err := h.seiyuuService.GetSeiyuuByID(ctx, id)
	if err != nil {
		if err == models.ErrSeiyuuNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	// 只返回已发布的声优
	if !seiyuu.IsPublic() {
		utils.NotFoundError(c, models.ErrSeiyuuNotFound)
		return
	}

	utils.SuccessResponse(c, seiyuu)
}

// GetAllSeiyuuAdmin 获取所有声优列表（管理员）
// GET /api/admin/seiyuu
func (h *SeiyuuHandler) GetAllSeiyuuAdmin(c *router.Context) {
	ctx := c.Request.Context()

	seiyuus, err := h.seiyuuService.GetAllSeiyuuAdmin(ctx)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessResponse(c, seiyuus)
}

// CreateSeiyuu 创建声优（管理员）
// POST /api/admin/seiyuu
func (h *SeiyuuHandler) CreateSeiyuu(c *router.Context) {
	ctx := c.Request.Context()

	var req models.CreateSeiyuuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if err := utils.ValidateSeiyuuName(req.Name); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	if err := utils.ValidateMarkdown(req.ProfileMarkdown); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	if err := utils.ValidateTags(req.Tags); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 创建声优
	seiyuu, err := h.seiyuuService.CreateSeiyuu(ctx, &req)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "声优创建成功", seiyuu)
}

// UpdateSeiyuu 更新声优资料（管理员）
// PUT /api/admin/seiyuu/:id
func (h *SeiyuuHandler) UpdateSeiyuu(c *router.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	var req models.UpdateSeiyuuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if req.Name != nil {
		if err := utils.ValidateSeiyuuName(*req.Name); err != nil {
			utils.BadRequestError(c, err)
			return
		}
	}

	if req.ProfileMarkdown != nil {
		if err := utils.ValidateMarkdown(*req.ProfileMarkdown); err != nil {
			utils.BadRequestError(c, err)
			return
		}
	}

	if req.Tags != nil {
		if err := utils.ValidateTags(*req.Tags); err != nil {
			utils.BadRequestError(c, err)
			return
		}
	}

	// 更新声优
	seiyuu, err := h.seiyuuService.UpdateSeiyuu(ctx, id, &req)
	if err != nil {
		if err == models.ErrSeiyuuNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "声优更新成功", seiyuu)
}

// DeleteSeiyuu 删除声优（管理员）
// DELETE /api/admin/seiyuu/:id
func (h *SeiyuuHandler) DeleteSeiyuu(c *router.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	err := h.seiyuuService.DeleteSeiyuu(ctx, id)
	if err != nil {
		if err == models.ErrSeiyuuNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "声优删除成功", nil)
}

// GetMoegirlRawData 获取萌娘百科原始数据（管理员）
// GET /api/admin/seiyuu/moegirl/:name
func (h *SeiyuuHandler) GetMoegirlRawData(c *router.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")

	if name == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 从萌娘百科获取原始数据
	rawData, err := h.seiyuuService.GetMoegirlRawData(ctx, name)
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
