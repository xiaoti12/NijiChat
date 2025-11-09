/**
 * 管理员状态管理
 * 管理管理员登录状态、权限和专用AI配置
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AIModelConfig } from '@/types'
import { encrypt, decrypt, generateId } from '@/utils/crypto'

interface AdminUser {
  username: string
  token: string
  expires_at: number
}

const ADMIN_STORAGE_KEY = 'nijichat_admin_user'
const ADMIN_AI_MODELS_STORAGE_KEY = 'nijichat_admin_ai_models'

export const useAdminStore = defineStore('admin', () => {
  const adminUser = ref<AdminUser | null>(loadAdminUser())
  const adminAIModels = ref<AIModelConfig[]>(loadAdminAIModels())

  // 从localStorage加载管理员AI模型配置（解密）
  function loadAdminAIModels(): AIModelConfig[] {
    try {
      const stored = localStorage.getItem(ADMIN_AI_MODELS_STORAGE_KEY)
      if (stored) {
        const encrypted = JSON.parse(stored)
        const decrypted = decrypt(encrypted)
        return JSON.parse(decrypted)
      }
    } catch (error) {
      console.error('加载管理员AI模型配置失败:', error)
    }
    return []
  }

  // 保存管理员AI模型配置到localStorage（加密）
  function saveAdminAIModels() {
    try {
      const json = JSON.stringify(adminAIModels.value)
      const encrypted = encrypt(json)
      localStorage.setItem(ADMIN_AI_MODELS_STORAGE_KEY, JSON.stringify(encrypted))
    } catch (error) {
      console.error('保存管理员AI模型配置失败:', error)
    }
  }

  // 从localStorage加载管理员信息
  function loadAdminUser(): AdminUser | null {
    try {
      const stored = localStorage.getItem(ADMIN_STORAGE_KEY)
      if (stored) {
        const user = JSON.parse(stored)
        // 检查是否过期（后端返回秒级时间戳，需要转换为毫秒级）
        const expiresAtMs = user.expires_at * 1000
        if (expiresAtMs > Date.now()) {
          return user
        } else {
          // 已过期，清除
          localStorage.removeItem(ADMIN_STORAGE_KEY)
        }
      }
    } catch (error) {
      console.error('加载管理员信息失败:', error)
    }
    return null
  }

  // 保存管理员信息
  function saveAdminUser(user: AdminUser) {
    try {
      localStorage.setItem(ADMIN_STORAGE_KEY, JSON.stringify(user))
      // 同时保存token到单独的key（给request.ts使用）
      localStorage.setItem('admin_token', user.token)
    } catch (error) {
      console.error('保存管理员信息失败:', error)
    }
  }

  // 登录
  function login(username: string, token: string, expiresAt: number) {
    const user: AdminUser = {
      username,
      token,
      expires_at: expiresAt
    }

    // 先更新状态
    adminUser.value = user

    // 然后保存到localStorage（即使保存失败也不影响状态）
    try {
      saveAdminUser(user)
    } catch (error) {
      console.warn('adminStore: localStorage保存失败，但登录状态仍然有效', error)
    }
  }

  // 登出
  function logout() {
    adminUser.value = null
    localStorage.removeItem(ADMIN_STORAGE_KEY)
    localStorage.removeItem('admin_token')
  }

  // 检查是否已登录
  const isLoggedIn = computed(() => {
    if (!adminUser.value) return false
    // 检查token是否过期（后端返回秒级时间戳，需要转换为毫秒级）
    const expiresAtMs = adminUser.value.expires_at * 1000
    return expiresAtMs > Date.now()
  })

  // 获取当前用户名
  const username = computed(() => adminUser.value?.username || '')

  // 获取token
  const token = computed(() => adminUser.value?.token || '')

  // 检查token是否即将过期（1小时内）
  const isTokenExpiringSoon = computed(() => {
    if (!adminUser.value) return false
    const oneHour = 60 * 60 * 1000
    // 后端返回秒级时间戳，需要转换为毫秒级
    const expiresAtMs = adminUser.value.expires_at * 1000
    return expiresAtMs - Date.now() < oneHour
  })

  // ========== 管理员AI配置管理 ==========

  // 添加管理员AI模型配置
  function addAdminAIModel(model: Omit<AIModelConfig, 'id'>): AIModelConfig {
    const newModel: AIModelConfig = {
      id: generateId(),
      ...model
    }
    adminAIModels.value.push(newModel)
    saveAdminAIModels()
    return newModel
  }

  // 更新管理员AI模型配置
  function updateAdminAIModel(id: string, updates: Partial<AIModelConfig>) {
    const index = adminAIModels.value.findIndex(m => m.id === id)
    if (index !== -1) {
      adminAIModels.value[index] = { ...adminAIModels.value[index], ...updates }
      saveAdminAIModels()
    }
  }

  // 删除管理员AI模型配置
  function deleteAdminAIModel(id: string) {
    const index = adminAIModels.value.findIndex(m => m.id === id)
    if (index !== -1) {
      adminAIModels.value.splice(index, 1)
      saveAdminAIModels()
    }
  }

  // 获取管理员AI模型配置
  function getAdminAIModel(id: string): AIModelConfig | undefined {
    return adminAIModels.value.find(m => m.id === id)
  }

  // 获取管理员AI聊天模型列表（非轻量级）
  const adminChatModels = computed(() => {
    return adminAIModels.value.filter(m => !m.is_lightweight)
  })

  // 获取管理员AI轻量级模型列表
  const adminLightweightModels = computed(() => {
    return adminAIModels.value.filter(m => m.is_lightweight)
  })

  // 按类型获取管理员AI模型
  function getAdminAIModelsByType(type: AIModelConfig['type']): AIModelConfig[] {
    return adminAIModels.value.filter(m => m.type === type)
  }

  // 检查是否有可用的管理员AI聊天模型
  const hasAdminChatModel = computed(() => adminChatModels.value.length > 0)

  // 检查是否有可用的管理员AI轻量级模型
  const hasAdminLightModel = computed(() => adminLightweightModels.value.length > 0)

  // 检查是否有任何管理员AI配置
  const hasAdminAIModel = computed(() => adminAIModels.value.length > 0)

  return {
    // State
    adminUser,
    adminAIModels,

    // Computed
    isLoggedIn,
    username,
    token,
    isTokenExpiringSoon,
    adminChatModels,
    adminLightweightModels,
    hasAdminChatModel,
    hasAdminLightModel,
    hasAdminAIModel,

    // Actions
    login,
    logout,
    addAdminAIModel,
    updateAdminAIModel,
    deleteAdminAIModel,
    getAdminAIModel,
    getAdminAIModelsByType
  }
})
