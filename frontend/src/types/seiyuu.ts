/**
 * 声优相关类型定义
 */

// 声优基础信息
export interface Seiyuu {
  id: string                    // 声优唯一ID
  name: string                  // 声优姓名
  avatar_url?: string           // 头像图片URL
  profile_markdown: string      // 完整Markdown格式资料
  raw_profile_data?: string     // 原始资料数据（用于AI关系生成）
  tags: string[]                // 标签列表 ["萝莉音", "治愈系"]
  status: SeiyuuStatus          // 状态
  created_at: string            // 创建时间
  updated_at: string            // 更新时间
}

// 声优状态枚举
export type SeiyuuStatus = 'pending' | 'active' | 'inactive'

// 声优群组
export interface SeiyuuGroup {
  id: string                    // 群组唯一ID
  name: string                  // 群组名称
  description: string           // 群组描述和规则
  member_ids: string[]          // 成员声优ID列表
  is_discussion_mode: boolean   // 是否群讨论模式
  created_at: string            // 创建时间
  updated_at: string            // 更新时间
}

// 声优简要信息（用于列表展示）
export interface SeiyuuBrief {
  id: string
  name: string
  avatar_url?: string
  tags: string[]
}

// 声优创建请求
export interface CreateSeiyuuRequest {
  name: string
  avatar_url?: string
  profile_markdown: string
  raw_profile_data?: string
  tags: string[]
}

// 声优更新请求
export interface UpdateSeiyuuRequest {
  name?: string
  avatar_url?: string
  profile_markdown?: string
  raw_profile_data?: string
  tags?: string[]
  status?: SeiyuuStatus
}

// 声优关系
export interface SeiyuuRelationship {
  id: string                    // 关系唯一ID
  seiyuu_id_a: string          // 声优A的ID
  seiyuu_id_b: string          // 声优B的ID
  seiyuu_name_a?: string       // 声优A姓名（查询所有关系时返回）
  seiyuu_name_b?: string       // 声优B姓名（查询所有关系时返回）
  relationship_description: string // 关系描述
  created_at: string           // 创建时间
  updated_at: string           // 更新时间
}

// 创建关系请求
export interface CreateRelationshipRequest {
  seiyuu_id_a: string
  seiyuu_id_b: string
  relationship_description: string
}

// 更新关系请求
export interface UpdateRelationshipRequest {
  relationship_description: string
}

// AI生成关系请求
export interface GenerateRelationshipRequest {
  seiyuu_id_a: string
  seiyuu_id_b: string
}

// AI生成关系响应
export interface GenerateRelationshipResponse {
  success: boolean
  message: string
  data?: {
    relationship_description: string
  }
}
