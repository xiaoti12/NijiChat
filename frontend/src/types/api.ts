/**
 * API接口相关类型定义
 */

import { Seiyuu, SeiyuuGroup } from './seiyuu'

// 通用API响应
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

// 分页请求参数
export interface PaginationParams {
  page: number
  page_size: number
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 声优列表响应
export type SeiyuuListResponse = ApiResponse<Seiyuu[]>

// 声优详情响应
export type SeiyuuDetailResponse = ApiResponse<Seiyuu>

// 群组列表响应
export type GroupListResponse = ApiResponse<SeiyuuGroup[]>

// 智能调度器请求
export interface SchedulerRequest {
  message: string               // 用户消息
  group_id?: string             // 群组ID（可选）
  available_seiyuu?: string[]   // 可选声优列表（可选）
  context?: string              // 上下文对话历史
}

// 智能调度器响应
export interface SchedulerResponse {
  success: boolean
  data?: {
    selected_seiyuu: Seiyuu
    reason: string              // 选择原因
    confidence: number          // 置信度 (0-1)
  }
  error?: string
}

// 萌娘百科原始数据响应
export interface MoegirlRawDataResponse {
  success: boolean
  data?: {
    raw_text: string            // 原始wiki文本
    page_title: string
    page_url: string
  }
  error?: string
}

// AI处理资料请求
export interface ProcessProfileRequest {
  raw_text: string              // 原始文本
  seiyuu_name: string           // 声优名称
}

// AI处理资料响应
export interface ProcessProfileResponse {
  success: boolean
  data?: {
    profile_markdown: string    // 处理后的Markdown格式资料
    suggested_tags: string[]    // 建议的标签
  }
  error?: string
}

// 管理员登录请求
export interface AdminLoginRequest {
  username: string
  password: string
}

// 管理员登录响应
export interface AdminLoginResponse {
  success: boolean
  data?: {
    token: string
    username: string
    expires_at: number
  }
  error?: string
}

// WebDAV配置
export interface WebDAVConfig {
  server_url: string
  username: string
  password: string              // 加密存储
  enabled: boolean
  auto_sync: boolean
  sync_interval: number         // 同步间隔（分钟）
}

// WebDAV同步状态
export interface WebDAVSyncStatus {
  last_sync: number             // 最后同步时间
  status: 'idle' | 'syncing' | 'error'
  error?: string
}

// 用户配置
export interface UserConfig {
  theme: 'light' | 'dark'
  language: 'zh-CN' | 'en-US'
  chat_bubble_style: string
  selected_chat_model?: string  // 选中的聊天模型ID
  selected_light_model?: string // 选中的轻量级模型ID
  webdav?: WebDAVConfig
}
