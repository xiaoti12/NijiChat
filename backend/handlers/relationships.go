package handlers

import (
	"seiyuu-chat/models"
	"seiyuu-chat/services"
	"seiyuu-chat/utils"

	"github.com/gin-gonic/gin"
)

// RelationshipsHandler 声优关系相关API处理器
type RelationshipsHandler struct {
	relationshipService *services.RelationshipService
}

// NewRelationshipsHandler 创建关系处理器实例
func NewRelationshipsHandler(relationshipService *services.RelationshipService) *RelationshipsHandler {
	return &RelationshipsHandler{
		relationshipService: relationshipService,
	}
}

// GenerateRelationship AI辅助生成关系描述
// POST /api/admin/relationships/generate
func (h *RelationshipsHandler) GenerateRelationship(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.GenerateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if req.SeiyuuIdA == "" || req.SeiyuuIdB == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	if req.SeiyuuIdA == req.SeiyuuIdB {
		utils.BadRequestError(c, models.ErrRelationshipSameSeiyuu)
		return
	}

	// 生成关系
	generatedRel, err := h.relationshipService.GenerateRelationship(ctx, req.SeiyuuIdA, req.SeiyuuIdB)
	if err != nil {
		if err == models.ErrSeiyuuNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "关系生成成功", generatedRel)
}

// CreateRelationship 创建关系
// POST /api/admin/relationships
func (h *RelationshipsHandler) CreateRelationship(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 验证数据
	if req.SeiyuuIdA == req.SeiyuuIdB {
		utils.BadRequestError(c, models.ErrRelationshipSameSeiyuu)
		return
	}

	// 创建关系
	relationship, err := h.relationshipService.CreateRelationship(ctx, &req)
	if err != nil {
		if err == models.ErrSeiyuuNotFound {
			utils.NotFoundError(c, err)
			return
		}
		if err == models.ErrRelationshipAlreadyExists {
			utils.BadRequestError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "关系创建成功", relationship)
}

// GetRelationship 获取特定关系或所有关系
// GET /api/admin/relationships?seiyuu_id_a=xxx&seiyuu_id_b=xxx （获取特定关系）
// GET /api/admin/relationships （获取所有关系）
func (h *RelationshipsHandler) GetRelationship(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查是否提供了查询参数来获取特定关系
	seiyuuIdA := c.Query("seiyuu_id_a")
	seiyuuIdB := c.Query("seiyuu_id_b")

	// 如果提供了两个声优ID，获取特定关系
	if seiyuuIdA != "" && seiyuuIdB != "" {
		relationship, err := h.relationshipService.GetRelationship(ctx, seiyuuIdA, seiyuuIdB)
		if err != nil {
			if err == models.ErrRelationshipNotFound {
				utils.NotFoundError(c, err)
				return
			}
			utils.InternalServerError(c, err)
			return
		}
		utils.SuccessResponse(c, relationship)
		return
	}

	// 如果只提供了部分参数，返回错误
	if seiyuuIdA != "" || seiyuuIdB != "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	// 如果没有提供查询参数，获取所有关系
	relationships, err := h.relationshipService.GetAllRelationships(ctx)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessResponse(c, relationships)
}

// UpdateRelationship 更新关系
// PUT /api/admin/relationships/:id
func (h *RelationshipsHandler) UpdateRelationship(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	var req models.UpdateRelationshipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestError(c, err)
		return
	}

	// 更新关系
	relationship, err := h.relationshipService.UpdateRelationship(ctx, id, &req)
	if err != nil {
		if err == models.ErrRelationshipNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "关系更新成功", relationship)
}

// DeleteRelationship 删除关系
// DELETE /api/admin/relationships/:id
func (h *RelationshipsHandler) DeleteRelationship(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	err := h.relationshipService.DeleteRelationship(ctx, id)
	if err != nil {
		if err == models.ErrRelationshipNotFound {
			utils.NotFoundError(c, err)
			return
		}
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessWithMessage(c, "关系删除成功", nil)
}

// GetSeiyuuRelationships 获取声优的所有关系
// GET /api/admin/seiyuu/:id/relationships
func (h *RelationshipsHandler) GetSeiyuuRelationships(c *gin.Context) {
	ctx := c.Request.Context()
	seiyuuId := c.Param("id")

	if seiyuuId == "" {
		utils.BadRequestError(c, models.ErrBadRequest)
		return
	}

	relationships, err := h.relationshipService.GetSeiyuuRelationships(ctx, seiyuuId)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}

	utils.SuccessResponse(c, relationships)
}

