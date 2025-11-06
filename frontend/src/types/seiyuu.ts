/**
 * 声优相关类型定义
 */

// 声优基础信息
export interface Seiyuu {
  id: string                    // 声优唯一ID
  name: string                  // 声优姓名
  avatar_url?: string           // 头像图片URL
  profile_markdown: string      // 完整Markdown格式资料
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
  tags: string[]
}

// 声优更新请求
export interface UpdateSeiyuuRequest {
  name?: string
  avatar_url?: string
  profile_markdown?: string
  tags?: string[]
  status?: SeiyuuStatus
}
