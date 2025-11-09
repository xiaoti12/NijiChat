/**
 * AI模型调用服务
 * 统一的AI调用接口，支持多种AI模型
 */

import { useAIModelStore } from '@/stores/aiModelStore'
import { useAdminStore } from '@/stores/adminStore'
import { useConfigStore } from '@/stores/configStore'
import type {
  AICallOptions,
  AICallOptionsDual,
  AIModelConfig,
  GeminiConfig,
  OpenAIConfig,
  SeiyuuProfileProcessResult,
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
      default:
        throw new Error(`不支持的AI模型类型: ${modelConfig.type}`)
    }
  }

  /**
   * 生成双声优对话回复
   */
  async generateDualReply(options: AICallOptionsDual): Promise<string> {
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

    // 构建双声优对话提示词
    const systemPrompt = this.buildDualSystemPrompt(
      options.responder_profile,
      options.initiator_profile,
      options.relationship_description,
      options.current_topic
    )

    // 构建消息列表（不包含用户输入，因为是声优之间的对话）
    const messages = this.buildDualMessages(systemPrompt, options.conversation_history)

    // 根据模型类型调用对应API
    switch (modelConfig.type) {
      case 'gemini':
        return await this.callGeminiAPI(modelConfig as GeminiConfig, messages, {
          message: '', // 双声优对话不需要用户消息
          seiyuu_profile: options.responder_profile,
          conversation_history: options.conversation_history,
          temperature: options.temperature,
          max_tokens: options.max_tokens
        })
      case 'openai':
        return await this.callOpenAIAPI(modelConfig as OpenAIConfig, messages, {
          message: '', // 双声优对话不需要用户消息
          seiyuu_profile: options.responder_profile,
          conversation_history: options.conversation_history,
          temperature: options.temperature,
          max_tokens: options.max_tokens
        })
      default:
        throw new Error(`不支持的AI模型类型: ${modelConfig.type}`)
    }
  }

  /**
   * 构建系统提示词
   */
  private buildSystemPrompt(seiyuuProfile: string): string {
    return `# 人物扮演专家

## 定位
你是一个专业的人物扮演助手，能够基于用户提供的人物资料，深入分析并总结人物性格，然后以该角色的身份进行真实、一致的对话和扮演。

## 能力
- **解析资料**：准确理解用户输入的人物资料，包括个人经历、轶事、人际关系等
- **性格总结**：从资料中提取关键信息，推断和描述人物的核心性格特征
- **角色扮演**：以角色的视角、语气和风格进行回应，确保扮演贴合背景和个性

## 知识储备
仅基于用户提供的角色资料进行推断和扮演，不依赖外部知识或添加假设。

## 扮演规则
1. 除非用户明确要求或询问，否则不主动提及或说明角色过往的轶事
2. 保持角色的一致性：语言、态度和行为应符合人物背景
3. 扮演时使用第一人称，自然融入对话情境
4. 不要透露你是AI，要完全沉浸在角色中
5. 避免在句尾添加emoji

## 操作流程
1. 接收用户输入的角色资料
2. 分析资料，识别重要事件、关系和性格线索
3. 总结人物性格特征
4. 以角色身份进行扮演，回应时保持自然流畅

# 角色资料
${seiyuuProfile}

`
  }

  /**
   * 构建双声优对话系统提示词
   */
  private buildDualSystemPrompt(
    responderProfile: string,
    initiatorProfile: string,
    relationshipDescription: string,
    currentTopic?: string
  ): string {
    return `# 双声优剧场对话专家

## 定位
你是一个专业的双人对话模拟助手，能够基于提供的两个角色资料和关系描述，模拟真实的人物互动对话。

## 当前扮演角色
${responderProfile}

## 对话对象
${initiatorProfile}

## 关系背景
${relationshipDescription}

${currentTopic ? `## 当前话题\n${currentTopic}\n` : ''}

## 扮演规则
1. **角色一致性**: 严格按照你的角色资料进行扮演，保持性格、语言风格和行为特征的一致性
2. **关系感知**: 基于提供的关系描述，以合适的态度和方式与对话对象互动
3. **自然对话**: 模拟真实的人际交流，包括情感表达、语气变化和个人观点
4. **上下文连贯**: 根据对话历史自然地推进话题或回应对方
5. **避免说教**: 不要刻意说教或过分解释，保持对话的自然流畅
6. **情感表达**: 根据角色性格和当前情境，适当表达情感和态度
7. **话题推进**: 可以根据角色特点主动引入新话题或深入探讨现有话题

## 注意事项
- 以第一人称进行对话，完全沉浸在角色中
- 直接回复发言，不要添加自己的角色名
- 不要提及你是AI或在进行角色扮演
- 避免在句尾添加emoji
- 回应要简洁自然，避免过长的独白

现在，请以你的角色身份，基于对话历史和当前情境，自然地进行下一句对话。`
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
   * 构建双声优对话消息列表
   */
  private buildDualMessages(systemPrompt: string, history: Message[]): any[] {
    const messages: any[] = []

    // 添加系统提示
    messages.push({
      role: 'system',
      content: systemPrompt
    })

    // 添加历史对话（最近N条）
    if (history && history.length > 0) {
      const recentHistory = history.slice(-15) // 双声优对话保留更多历史上下文
      for (const msg of recentHistory) {
        // 在双声优对话中，所有非系统消息都被视为assistant角色的发言
        // 因为是两个AI角色在对话，不涉及用户输入
        if (msg.sender_id !== 'system') {
          messages.push({
            role: 'assistant',
            content: `${msg.sender_name}: ${msg.content}`
          })
        }
      }
    }

    // 双声优对话中，如果没有历史记录，添加一个引导消息
    if (!history || history.length === 0) {
      messages.push({
        role: 'user',
        content: '请开始你们之间的对话。'
      })
    }

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
    const url = `${config.api_endpoint}/models/${config.model_name}:generateContent?key=${config.api_key}`

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
   * 处理声优资料（管理员功能）
   */
  async processSeiyuuProfile(rawText: string, modelId?: string): Promise<SeiyuuProfileProcessResult> {
    const adminStore = useAdminStore()

    // 如果没有找到模型，尝试重新加载
    if (adminStore.adminAIModels.length === 0) {
      adminStore.reloadAdminAIModels()
      // 等待Vue响应式更新完成
      await new Promise(resolve => setTimeout(resolve, 100))
    }

    // 选择AI模型配置
    let modelConfig: AIModelConfig | undefined
    if (modelId) {
      modelConfig = adminStore.getAdminAIModel(modelId)
      // 如果通过ID没找到，尝试通过名称查找
      if (!modelConfig) {
        modelConfig = adminStore.adminAIModels.find(m => m.name === modelId)
      }
    } else {
      // 使用第一个可用的管理员AI配置
      const availableModels = adminStore.adminAIModels
      if (availableModels.length > 0) {
        modelConfig = availableModels[0]
      } else {
        // 最后一次尝试：直接从localStorage读取
        try {
          const stored = localStorage.getItem('nijichat_admin_ai_models')
          if (stored) {
            const parsed = JSON.parse(stored)
            if (Array.isArray(parsed) && parsed.length > 0) {
              modelConfig = parsed[0]
            }
          }
        } catch (error) {
          console.error('aiService: 从localStorage读取配置失败:', error)
        }
      }
    }

    if (!modelConfig) {
      // 详细的调试信息
      const storedData = localStorage.getItem('nijichat_admin_ai_models')

      let errorMessage = '未找到可用的管理员AI配置。'

      if (!storedData) {
        errorMessage += ' localStorage中没有AI配置数据，请在管理员后台重新添加AI配置。'
      } else {
        try {
          const parsed = JSON.parse(storedData)
          if (!Array.isArray(parsed)) {
            errorMessage += ' 配置数据格式错误，请重新添加AI配置。'
          } else if (parsed.length === 0) {
            errorMessage += ' 配置数据为空，请在管理员后台添加AI配置。'
          } else {
            errorMessage += ` 发现${parsed.length}个配置但无法访问，可能存在状态同步问题。请尝试刷新页面或重新添加配置。`
          }
        } catch (error) {
          errorMessage += ' 配置数据解析失败，请清除数据后重新添加AI配置。'
        }
      }

      throw new Error(errorMessage)
    }

    console.log('aiService: 成功获取到模型配置，开始处理 =', modelConfig.name)

    // 构建声优资料处理提示词
    const systemPrompt = this.buildSeiyuuProfilePrompt()
    const userContent = `角色资料如下：\n${rawText}`

    // 构建消息列表
    const messages = [
      {
        role: 'system',
        content: systemPrompt
      },
      {
        role: 'user',
        content: userContent
      }
    ]

    // 调用AI模型
    let response: string
    switch (modelConfig.type) {
      case 'gemini':
        response = await this.callGeminiAPI(modelConfig as GeminiConfig, messages, {
          message: userContent,
          seiyuu_profile: '',
          conversation_history: [],
          temperature: 0.3, // 资料处理需要较低的温度以确保准确性
          max_tokens: modelConfig.max_tokens
        })
        break
      case 'openai':
        response = await this.callOpenAIAPI(modelConfig as OpenAIConfig, messages, {
          message: userContent,
          seiyuu_profile: '',
          conversation_history: [],
          temperature: 0.3,
          max_tokens: modelConfig.max_tokens
        })
        break
      default:
        throw new Error(`不支持的AI模型类型: ${modelConfig.type}`)
    }

    // 解析AI响应
    try {
      // 清理响应内容：去除markdown代码块包围
      let cleanedResponse = response.trim()

      // 去除 ```json 开头和 ``` 结尾
      if (cleanedResponse.startsWith('```json')) {
        cleanedResponse = cleanedResponse.replace(/^```json\s*/, '')
      }
      if (cleanedResponse.startsWith('```')) {
        cleanedResponse = cleanedResponse.replace(/^```\s*/, '')
      }
      if (cleanedResponse.endsWith('```')) {
        cleanedResponse = cleanedResponse.replace(/\s*```$/, '')
      }

      // 去除可能的其他代码块标记
      cleanedResponse = cleanedResponse
        .replace(/^```[\w]*\s*/, '') // 去除开头的任何代码块标记
        .replace(/\s*```$/, '')      // 去除结尾的代码块标记
        .trim()

      console.log('原始AI响应:', response)
      console.log('清理后的JSON:', cleanedResponse)

      // 尝试解析为JSON
      const result = JSON.parse(cleanedResponse)
      return {
        profile_markdown: result.profile_markdown || response,
        suggested_tags: result.suggested_tags || []
      }
    } catch (error) {
      console.warn('JSON解析失败，使用原始文本:', error)
      // 如果解析失败，返回原始文本
      return {
        profile_markdown: response,
        suggested_tags: []
      }
    }
  }

  /**
   * 构建声优资料处理提示词
   */
  private buildSeiyuuProfilePrompt(): string {
    return `请将以下关于声优的原始资料整理为结构化的Markdown格式。

要求：

总结关键信息：姓名、生日、个人爱好、人际关系、个人轶事等
组织为清晰的Markdown格式
保持客观真实，不添加虚构内容，不允许增加或删除信息，只能整理已有内容
建议3-5个相关标签（例如性格、爱好等）
【重要】不允许增加或删除信息，只能整理已有内容

输出格式要求：
- 直接返回纯JSON格式
- 不要使用 \`\`\`json 或 \`\`\` 包围
- 不要添加任何markdown代码块标记
- 不要添加任何解释文字
- 只输出JSON对象本身

输出示例：
{
  "profile_markdown": "整理后的Markdown文本",
  "suggested_tags": ["标签1", "标签2", "标签3"]
}

现在请处理以下声优资料：`
  }

  /**
   * 生成声优关系描述（管理员功能）
   */
  async generateSeiyuuRelationship(
    seiyuuA: { name: string; raw_profile_data: string },
    seiyuuB: { name: string; raw_profile_data: string },
    modelId?: string
  ): Promise<string> {
    const adminStore = useAdminStore()

    // 如果没有找到模型，尝试重新加载
    if (adminStore.adminAIModels.length === 0) {
      console.log('aiService: generateSeiyuuRelationship 模型列表为空，尝试重新加载...')
      adminStore.reloadAdminAIModels()
    }

    // 选择AI模型配置
    let modelConfig: AIModelConfig | undefined
    if (modelId) {
      modelConfig = adminStore.getAdminAIModel(modelId)
    } else {
      // 使用第一个可用的管理员AI配置
      const availableModels = adminStore.adminAIModels
      if (availableModels.length > 0) {
        modelConfig = availableModels[0]
      } else {
        // 最后一次尝试：直接从localStorage读取
        try {
          const stored = localStorage.getItem('nijichat_admin_ai_models')
          if (stored) {
            const parsed = JSON.parse(stored)
            if (Array.isArray(parsed) && parsed.length > 0) {
              modelConfig = parsed[0]
            }
          }
        } catch (error) {
          console.error('aiService: generateSeiyuuRelationship 从localStorage读取配置失败:', error)
        }
      }
    }

    if (!modelConfig) {
      throw new Error('未找到可用的管理员AI配置，请先在管理员设置中配置AI模型')
    }

    // 构建关系生成提示词
    const systemPrompt = this.buildRelationshipPrompt()
    const userContent = `声优A： ${seiyuuA.name} 的原始资料：
${seiyuuA.raw_profile_data}

声优B： ${seiyuuB.name} 的原始资料：
${seiyuuB.raw_profile_data}`

    // 构建消息列表
    const messages = [
      {
        role: 'system',
        content: systemPrompt
      },
      {
        role: 'user',
        content: userContent
      }
    ]

    // 调用AI模型
    let response: string
    switch (modelConfig.type) {
      case 'gemini':
        response = await this.callGeminiAPI(modelConfig as GeminiConfig, messages, {
          message: userContent,
          seiyuu_profile: '',
          conversation_history: [],
          temperature: 0.2, // 关系生成需要较低的温度以确保客观性
          max_tokens: modelConfig.max_tokens
        })
        break
      case 'openai':
        response = await this.callOpenAIAPI(modelConfig as OpenAIConfig, messages, {
          message: userContent,
          seiyuu_profile: '',
          conversation_history: [],
          temperature: 0.2,
          max_tokens: modelConfig.max_tokens
        })
        break
      default:
        throw new Error(`不支持的AI模型类型: ${modelConfig.type}`)
    }

    // 验证生成的内容不为空
    const relationshipDescription = response.trim()
    if (!relationshipDescription) {
      return `${seiyuuA.name} 和 ${seiyuuB.name} 同为声优行业的从业者。`
    }

    return relationshipDescription
  }

  /**
   * 构建声优关系生成提示词
   */
  private buildRelationshipPrompt(): string {
    return `请根据以下两位声优的资料分析他们之间可能存在的关系。

要求：
1. 基于提供的真实资料进行分析，不要添加虚构内容
2. 从中立的第三者视角描述关系，使用客观语言
3. 详细说明两人的具体关系背景、共同点或互动情况
4. 如果没有明显关系，描述他们作为同行的共同特点和专业领域

请直接返回关系描述文本，不要使用任何格式标记或代码块。`
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
