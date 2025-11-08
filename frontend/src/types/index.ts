/**
 * 类型定义入口文件
 * 统一导出所有类型定义
 */

// 声优相关类型
export type {
  Seiyuu,
  SeiyuuStatus,
  SeiyuuGroup,
  SeiyuuBrief,
  CreateSeiyuuRequest,
  UpdateSeiyuuRequest,
  SeiyuuRelationship,
  CreateRelationshipRequest,
  UpdateRelationshipRequest,
  GenerateRelationshipRequest,
  GenerateRelationshipResponse
} from './seiyuu'

// 聊天相关类型
export type {
  Message,
  MessageType,
  Room,
  RoomType,
  Conversation,
  DualTheater,
  GroupTheater,
  TheaterStatus,
  TheaterMessage,
  SeiyuuConversationGroup,
  ChatSettings
} from './chat'

// AI模型相关类型
export type {
  AIModelConfig,
  AIModelType,
  AICallOptions,
  AICallOptionsDual,
  AIResponse,
  AIStreamEvent,
  GeminiConfig,
  OpenAIConfig,
  ClaudeConfig,
  AIModelTestResult
} from './ai'

// API接口相关类型
export type {
  ApiResponse,
  PaginationParams,
  PaginatedResponse,
  SeiyuuListResponse,
  SeiyuuDetailResponse,
  GroupListResponse,
  SchedulerRequest,
  SchedulerResponse,
  MoegirlRawDataResponse,
  ProcessProfileRequest,
  ProcessProfileResponse,
  AdminLoginRequest,
  AdminLoginResponse,
  WebDAVConfig,
  WebDAVSyncStatus,
  UserConfig
} from './api'
