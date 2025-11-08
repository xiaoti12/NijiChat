package models

import (
	"time"
)

// SeiyuuRelationship 声优关系数据模型
type SeiyuuRelationship struct {
	ID                      string    `json:"id"`                        // 关系唯一ID (UUID)
	SeiyuuIdA               string    `json:"seiyuu_id_a"`               // 声优A的ID
	SeiyuuIdB               string    `json:"seiyuu_id_b"`               // 声优B的ID
	RelationshipDescription string    `json:"relationship_description"`  // 统一的关系描述（从上帝视角阐述）
	CreatedAt               time.Time `json:"created_at"`                // 创建时间
	UpdatedAt               time.Time `json:"updated_at"`                // 更新时间
}

// SeiyuuRelationshipWithNames 包含声优姓名的关系信息（用于展示）
type SeiyuuRelationshipWithNames struct {
	SeiyuuRelationship
	SeiyuuNameA string `json:"seiyuu_name_a"` // 声优A的姓名
	SeiyuuNameB string `json:"seiyuu_name_b"` // 声优B的姓名
}

// CreateRelationshipRequest 创建关系请求
type CreateRelationshipRequest struct {
	SeiyuuIdA               string `json:"seiyuu_id_a" binding:"required"`               // 声优A的ID
	SeiyuuIdB               string `json:"seiyuu_id_b" binding:"required"`               // 声优B的ID
	RelationshipDescription string `json:"relationship_description" binding:"required"`  // 关系描述
}

// UpdateRelationshipRequest 更新关系请求
type UpdateRelationshipRequest struct {
	RelationshipDescription *string `json:"relationship_description"`  // 关系描述
}

// GenerateRelationshipRequest AI生成关系请求
type GenerateRelationshipRequest struct {
	SeiyuuIdA string `json:"seiyuu_id_a" binding:"required"` // 声优A的ID
	SeiyuuIdB string `json:"seiyuu_id_b" binding:"required"` // 声优B的ID
}

// GeneratedRelationship AI生成的关系信息
type GeneratedRelationship struct {
	RelationshipDescription string `json:"relationship_description"`  // 关系描述
}

// Validate 验证关系数据
func (r *SeiyuuRelationship) Validate() error {
	if r.SeiyuuIdA == "" {
		return ErrInvalidRelationshipSeiyuuA
	}
	if r.SeiyuuIdB == "" {
		return ErrInvalidRelationshipSeiyuuB
	}
	if r.SeiyuuIdA == r.SeiyuuIdB {
		return ErrRelationshipSameSeiyuu
	}
	if r.RelationshipDescription == "" {
		return ErrInvalidRelationshipDescription
	}
	return nil
}

// SwapSeiyuu 交换两个声优的位置（用于处理无序关系）
func (r *SeiyuuRelationship) SwapSeiyuu() {
	r.SeiyuuIdA, r.SeiyuuIdB = r.SeiyuuIdB, r.SeiyuuIdA
}

// NormalizeSeiyuuOrder 标准化声优顺序（确保SeiyuuIdA < SeiyuuIdB，避免重复关系）
func (r *SeiyuuRelationship) NormalizeSeiyuuOrder() {
	if r.SeiyuuIdA > r.SeiyuuIdB {
		r.SwapSeiyuu()
	}
}