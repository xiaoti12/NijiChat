/**
 * AI模型相关类型定义
 */

// AI模型配置
export interface AIModelConfig {
  id: string                    // 配置ID
  name: string                  // 配置名称
  type: AIModelType             // 模型类型
  api_key: string               // API密钥（敏感数据，仅本地存储）
  api_endpoint: string          // API端点
  model_name: string            // 模型名称
  is_lightweight: boolean       // 是否轻量级模型（用于资料处理）
  max_tokens: number            // 最大token数
  temperature: number           // 温度参数 (0-2)
  top_p?: number                // Top P参数
  frequency_penalty?: number    // 频率惩罚
  presence_penalty?: number     // 存在惩罚
}

// AI模型类型
export type AIModelType = 'gemini' | 'openai'

// AI调用选项
export interface AICallOptions {
  message: string               // 用户消息
  seiyuu_profile: string        // 声优人设资料（Markdown格式）
  conversation_history: Message[] // 对话历史
  model_id?: string             // 指定使用的模型ID
  stream?: boolean              // 是否流式响应
  temperature?: number          // 覆盖默认温度
  max_tokens?: number           // 覆盖默认token数
}

// 双声优对话AI调用选项
export interface AICallOptionsDual {
  responder_profile: string     // 响应者声优资料（当前要发言的声优）
  initiator_profile: string     // 发起者声优资料（对话目标）
  relationship_description: string // 两人关系描述
  conversation_history: Message[] // 对话历史
  current_topic?: string        // 当前对话话题
  model_id?: string             // 指定使用的模型ID
  temperature?: number          // 覆盖默认温度
  max_tokens?: number           // 覆盖默认token数
}

// AI响应
export interface AIResponse {
  content: string               // 生成的内容
  finish_reason: string         // 完成原因
  usage?: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
  }
}

// 流式响应事件
export interface AIStreamEvent {
  type: 'start' | 'chunk' | 'end' | 'error'
  content?: string
  error?: string
}

// Gemini API配置
export interface GeminiConfig extends AIModelConfig {
  type: 'gemini'
  model_name: 'gemini-pro' | 'gemini-pro-vision' | string
}

// OpenAI API配置
export interface OpenAIConfig extends AIModelConfig {
  type: 'openai'
  model_name: 'gpt-3.5-turbo' | 'gpt-4' | 'gpt-4-turbo' | string
}


// AI模型测试结果
export interface AIModelTestResult {
  success: boolean
  message: string
  latency?: number              // 响应延迟（毫秒）
  error?: string
}

import { Message } from './chat'
