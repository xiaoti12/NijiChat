package models

import (
	"time"
)

// SeiyuuStatus 声优状态枚举
type SeiyuuStatus string

const (
	SeiyuuStatusPending  SeiyuuStatus = "pending"  // 待审核
	SeiyuuStatusActive   SeiyuuStatus = "active"   // 已发布
	SeiyuuStatusInactive SeiyuuStatus = "inactive" // 已下线
)

// Seiyuu 声优数据模型
type Seiyuu struct {
	ID              string       `json:"id"`                // 声优唯一ID (UUID)
	Name            string       `json:"name"`              // 声优姓名
	AvatarURL       string       `json:"avatar_url"`        // 头像图片URL
	ProfileMarkdown string       `json:"profile_markdown"`  // 完整Markdown格式资料
	RawProfileData  string       `json:"raw_profile_data"`  // 原始资料数据
	Tags            []string     `json:"tags"`              // 标签列表
	Status          SeiyuuStatus `json:"status"`            // 状态
	CreatedAt       time.Time    `json:"created_at"`        // 创建时间
	UpdatedAt       time.Time    `json:"updated_at"`        // 更新时间
}

// SeiyuuBrief 声优简要信息（用于列表展示）
type SeiyuuBrief struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	AvatarURL string   `json:"avatar_url"`
	Tags      []string `json:"tags"`
}

// CreateSeiyuuRequest 创建声优请求
type CreateSeiyuuRequest struct {
	Name            string   `json:"name" binding:"required"`
	AvatarURL       string   `json:"avatar_url"`
	ProfileMarkdown string   `json:"profile_markdown" binding:"required"`
	RawProfileData  string   `json:"raw_profile_data"`  // 原始资料数据（可选，未提供则使用ProfileMarkdown）
	Tags            []string `json:"tags"`
}

// UpdateSeiyuuRequest 更新声优请求
type UpdateSeiyuuRequest struct {
	Name            *string       `json:"name"`
	AvatarURL       *string       `json:"avatar_url"`
	ProfileMarkdown *string       `json:"profile_markdown"`
	RawProfileData  *string       `json:"raw_profile_data"`  // 原始资料数据
	Tags            *[]string     `json:"tags"`
	Status          *SeiyuuStatus `json:"status"`
}

// ToBrief 转换为简要信息
func (s *Seiyuu) ToBrief() *SeiyuuBrief {
	return &SeiyuuBrief{
		ID:        s.ID,
		Name:      s.Name,
		AvatarURL: s.AvatarURL,
		Tags:      s.Tags,
	}
}

// Validate 验证声优数据
func (s *Seiyuu) Validate() error {
	if s.Name == "" {
		return ErrInvalidSeiyuuName
	}
	if s.ProfileMarkdown == "" {
		return ErrInvalidSeiyuuProfile
	}
	if s.Status != SeiyuuStatusPending && s.Status != SeiyuuStatusActive && s.Status != SeiyuuStatusInactive {
		return ErrInvalidSeiyuuStatus
	}
	return nil
}

// IsActive 检查声优是否已发布
func (s *Seiyuu) IsActive() bool {
	return s.Status == SeiyuuStatusActive
}

// IsPublic 检查声优是否可公开访问
func (s *Seiyuu) IsPublic() bool {
	return s.Status == SeiyuuStatusActive
}
