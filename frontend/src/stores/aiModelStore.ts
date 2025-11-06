/**
 * AI模型配置状态管理
 * 管理AI模型配置（敏感数据本地存储）
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AIModelConfig } from '@/types'
import { encrypt, decrypt, generateId } from '@/utils/crypto'

const MODELS_STORAGE_KEY = 'nijichat_ai_models'

export const useAIModelStore = defineStore('aiModel', () => {
  const models = ref<AIModelConfig[]>(loadModels())

  // 从localStorage加载模型配置（解密）
  function loadModels(): AIModelConfig[] {
    try {
      const stored = localStorage.getItem(MODELS_STORAGE_KEY)
      if (stored) {
        const encrypted = JSON.parse(stored)
        const decrypted = decrypt(encrypted)
        return JSON.parse(decrypted)
      }
    } catch (error) {
      console.error('加载AI模型配置失败:', error)
    }
    return []
  }

  // 保存模型配置到localStorage（加密）
  function saveModels() {
    try {
      const json = JSON.stringify(models.value)
      const encrypted = encrypt(json)
      localStorage.setItem(MODELS_STORAGE_KEY, JSON.stringify(encrypted))
    } catch (error) {
      console.error('保存AI模型配置失败:', error)
    }
  }

  // 添加模型配置
  function addModel(model: Omit<AIModelConfig, 'id'>): AIModelConfig {
    const newModel: AIModelConfig = {
      id: generateId(),
      ...model
    }
    models.value.push(newModel)
    saveModels()
    return newModel
  }

  // 更新模型配置
  function updateModel(id: string, updates: Partial<AIModelConfig>) {
    const index = models.value.findIndex(m => m.id === id)
    if (index !== -1) {
      models.value[index] = { ...models.value[index], ...updates }
      saveModels()
    }
  }

  // 删除模型配置
  function deleteModel(id: string) {
    const index = models.value.findIndex(m => m.id === id)
    if (index !== -1) {
      models.value.splice(index, 1)
      saveModels()
    }
  }

  // 获取模型配置
  function getModel(id: string): AIModelConfig | undefined {
    return models.value.find(m => m.id === id)
  }

  // 获取聊天模型列表（非轻量级）
  const chatModels = computed(() => {
    return models.value.filter(m => !m.is_lightweight)
  })

  // 获取轻量级模型列表
  const lightweightModels = computed(() => {
    return models.value.filter(m => m.is_lightweight)
  })

  // 按类型获取模型
  function getModelsByType(type: AIModelConfig['type']): AIModelConfig[] {
    return models.value.filter(m => m.type === type)
  }

  // 检查是否有可用的聊天模型
  const hasChatModel = computed(() => chatModels.value.length > 0)

  // 检查是否有可用的轻量级模型
  const hasLightModel = computed(() => lightweightModels.value.length > 0)

  return {
    models,
    chatModels,
    lightweightModels,
    hasChatModel,
    hasLightModel,
    addModel,
    updateModel,
    deleteModel,
    getModel,
    getModelsByType
  }
})
