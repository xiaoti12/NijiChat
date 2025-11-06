package models

import (
	"time"
)

// SeiyuuGroup 声优群组数据模型
type SeiyuuGroup struct {
	ID                string    `json:"id"`                  // 群组唯一ID
	Name              string    `json:"name"`                // 群组名称
	Description       string    `json:"description"`         // 群组描述和规则
	MemberIDs         []string  `json:"member_ids"`          // 成员声优ID列表
	IsDiscussionMode  bool      `json:"is_discussion_mode"`  // 是否群讨论模式
	CreatedAt         time.Time `json:"created_at"`          // 创建时间
	UpdatedAt         time.Time `json:"updated_at"`          // 更新时间
}

// CreateGroupRequest 创建群组请求
type CreateGroupRequest struct {
	Name             string   `json:"name" binding:"required"`
	Description      string   `json:"description"`
	MemberIDs        []string `json:"member_ids" binding:"required"`
	IsDiscussionMode bool     `json:"is_discussion_mode"`
}

// UpdateGroupRequest 更新群组请求
type UpdateGroupRequest struct {
	Name             *string   `json:"name"`
	Description      *string   `json:"description"`
	MemberIDs        *[]string `json:"member_ids"`
	IsDiscussionMode *bool     `json:"is_discussion_mode"`
}

// Validate 验证群组数据
func (g *SeiyuuGroup) Validate() error {
	if g.Name == "" {
		return ErrInvalidGroupName
	}
	if len(g.MemberIDs) == 0 {
		return ErrInvalidGroupMembers
	}
	// 群讨论模式至少需要3个成员
	if g.IsDiscussionMode && len(g.MemberIDs) < 3 {
		return ErrInvalidGroupMembersCount
	}
	return nil
}

// HasMember 检查群组是否包含指定声优
func (g *SeiyuuGroup) HasMember(seiyuuID string) bool {
	for _, id := range g.MemberIDs {
		if id == seiyuuID {
			return true
		}
	}
	return false
}

// MemberCount 获取成员数量
func (g *SeiyuuGroup) MemberCount() int {
	return len(g.MemberIDs)
}
