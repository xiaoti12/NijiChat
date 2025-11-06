/**
 * AI模型调用服务
 * 统一的AI调用接口，支持多种AI模型
 */

import { useAIModelStore } from '@/stores/aiModelStore'
import { useConfigStore } from '@/stores/configStore'
import type {
  AICallOptions,
  AIResponse,
  AIModelConfig,
  GeminiConfig,
  OpenAIConfig,
  Message
} from '@/types'

/**
 * AI服务类
 */
export class AIService {
  /**
   * 生成AI回复
   */
  async generateReply(options: AICallOptions): Promise<string> {
    const aiModelStore = useAIModelStore()
    const configStore = useConfigStore()

    // 获取AI模型配置
    const modelId = options.model_id || configStore.config.selected_chat_model
    if (!modelId) {
      throw new Error('未配置AI模型，请先在设置中添加AI模型配置')
    }

    const modelConfig = aiModelStore.getModel(modelId)
    if (!modelConfig) {
      throw new Error('未找到指定的AI模型配置')
    }

    // 构建提示词
    const systemPrompt = this.buildSystemPrompt(options.seiyuu_profile)
    const messages = this.buildMessages(systemPrompt, options.message, options.conversation_history)

    // 根据模型类型调用对应API
    switch (modelConfig.type) {
      case 'gemini':
        return await this.callGeminiAPI(modelConfig as GeminiConfig, messages, options)
      case 'openai':
        return await this.callOpenAIAPI(modelConfig as OpenAIConfig, messages, options)
      case 'claude':
        return await this.callClaudeAPI(modelConfig, messages, options)
      default:
        throw new Error(`不支持的AI模型类型: ${modelConfig.type}`)
    }
  }

  /**
   * 构建系统提示词
   */
  private buildSystemPrompt(seiyuuProfile: string): string {
    return `你是一位声优，请根据以下资料扮演这个角色，用第一人称回复用户的消息。

# 角色资料
${seiyuuProfile}

# 扮演要求
1. 完全以这位声优的身份和语气回复
2. 回复要符合声优的性格特点和说话风格
3. 可以谈论声优本人的作品和经历
4. 保持友好、自然的对话氛围
5. 回复要简洁，一般在50-200字之间
6. 不要透露你是AI，要完全沉浸在角色中`
  }

  /**
   * 构建消息列表
   */
  private buildMessages(systemPrompt: string, userMessage: string, history: Message[]): any[] {
    const messages: any[] = []

    // 添加系统提示
    messages.push({
      role: 'system',
      content: systemPrompt
    })

    // 添加历史对话（最近N条）
    if (history && history.length > 0) {
      const recentHistory = history.slice(-10) // 最多10条历史
      for (const msg of recentHistory) {
        messages.push({
          role: msg.sender_id === 'user-1' ? 'user' : 'assistant',
          content: msg.content
        })
      }
    }

    // 添加当前用户消息
    messages.push({
      role: 'user',
      content: userMessage
    })

    return messages
  }

  /**
   * 调用Gemini API
   */
  private async callGeminiAPI(
    config: GeminiConfig,
    messages: any[],
    options: AICallOptions
  ): Promise<string> {
    const url = `${config.api_endpoint}/v1beta/models/${config.model_name}:generateContent?key=${config.api_key}`

    // 转换消息格式为Gemini格式
    const contents = messages
      .filter(m => m.role !== 'system')
      .map(m => ({
        role: m.role === 'user' ? 'user' : 'model',
        parts: [{ text: m.content }]
      }))

    // 系统提示词添加到第一条消息
    const systemMessage = messages.find(m => m.role === 'system')
    if (systemMessage && contents.length > 0) {
      contents[0].parts[0].text = `${systemMessage.content}\n\n${contents[0].parts[0].text}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        contents,
        generationConfig: {
          temperature: options.temperature ?? config.temperature,
          maxOutputTokens: options.max_tokens ?? config.max_tokens
        }
      })
    })

    if (!response.ok) {
      throw new Error(`Gemini API调用失败: ${response.statusText}`)
    }

    const data = await response.json()
    return data.candidates[0]?.content?.parts[0]?.text || ''
  }

  /**
   * 调用OpenAI API
   */
  private async callOpenAIAPI(
    config: OpenAIConfig,
    messages: any[],
    options: AICallOptions
  ): Promise<string> {
    const url = `${config.api_endpoint}/v1/chat/completions`

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${config.api_key}`
      },
      body: JSON.stringify({
        model: config.model_name,
        messages,
        temperature: options.temperature ?? config.temperature,
        max_tokens: options.max_tokens ?? config.max_tokens
      })
    })

    if (!response.ok) {
      throw new Error(`OpenAI API调用失败: ${response.statusText}`)
    }

    const data = await response.json()
    return data.choices[0]?.message?.content || ''
  }

  /**
   * 调用Claude API
   */
  private async callClaudeAPI(
    config: AIModelConfig,
    messages: any[],
    options: AICallOptions
  ): Promise<string> {
    const url = `${config.api_endpoint}/v1/messages`

    // 提取系统提示词
    const systemMessage = messages.find(m => m.role === 'system')
    const userMessages = messages.filter(m => m.role !== 'system')

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'x-api-key': config.api_key,
        'anthropic-version': '2023-06-01'
      },
      body: JSON.stringify({
        model: config.model_name,
        max_tokens: options.max_tokens ?? config.max_tokens,
        temperature: options.temperature ?? config.temperature,
        system: systemMessage?.content,
        messages: userMessages
      })
    })

    if (!response.ok) {
      throw new Error(`Claude API调用失败: ${response.statusText}`)
    }

    const data = await response.json()
    return data.content[0]?.text || ''
  }

  /**
   * 测试AI模型配置
   */
  async testModel(modelConfig: AIModelConfig): Promise<{ success: boolean; message: string }> {
    try {
      const startTime = Date.now()

      // 发送测试消息
      const testMessage = '你好，这是一条测试消息，请简短回复。'
      const systemPrompt = '你是一个友好的助手。'
      const messages = this.buildMessages(systemPrompt, testMessage, [])

      let response: string
      switch (modelConfig.type) {
        case 'gemini':
          response = await this.callGeminiAPI(modelConfig as GeminiConfig, messages, {
            message: testMessage,
            seiyuu_profile: '',
            conversation_history: []
          })
          break
        case 'openai':
          response = await this.callOpenAIAPI(modelConfig as OpenAIConfig, messages, {
            message: testMessage,
            seiyuu_profile: '',
            conversation_history: []
          })
          break
        case 'claude':
          response = await this.callClaudeAPI(modelConfig, messages, {
            message: testMessage,
            seiyuu_profile: '',
            conversation_history: []
          })
          break
        default:
          throw new Error('不支持的模型类型')
      }

      const latency = Date.now() - startTime

      if (response && response.length > 0) {
        return {
          success: true,
          message: `测试成功！响应时间: ${latency}ms`
        }
      } else {
        return {
          success: false,
          message: '测试失败：未收到有效响应'
        }
      }
    } catch (error: any) {
      return {
        success: false,
        message: `测试失败：${error.message}`
      }
    }
  }
}

// 导出单例
export const aiService = new AIService()

// 导出composable
export function useAIService() {
  return aiService
}
