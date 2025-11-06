/**
 * 用户配置状态管理
 * 封装localStorage操作
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UserConfig } from '@/types'

const CONFIG_STORAGE_KEY = 'nijichat_user_config'

const defaultConfig: UserConfig = {
  theme: 'light',
  language: 'zh-CN',
  chat_bubble_style: 'default'
}

export const useConfigStore = defineStore('config', () => {
  const config = ref<UserConfig>(loadConfig())

  // 从localStorage加载配置
  function loadConfig(): UserConfig {
    try {
      const stored = localStorage.getItem(CONFIG_STORAGE_KEY)
      if (stored) {
        return { ...defaultConfig, ...JSON.parse(stored) }
      }
    } catch (error) {
      console.error('加载配置失败:', error)
    }
    return defaultConfig
  }

  // 保存配置到localStorage
  function saveConfig() {
    try {
      localStorage.setItem(CONFIG_STORAGE_KEY, JSON.stringify(config.value))
    } catch (error) {
      console.error('保存配置失败:', error)
    }
  }

  // 更新配置
  function updateConfig(updates: Partial<UserConfig>) {
    config.value = { ...config.value, ...updates }
    saveConfig()
  }

  // 设置主题
  function setTheme(theme: 'light' | 'dark') {
    updateConfig({ theme })
    document.documentElement.setAttribute('data-theme', theme)
  }

  // 设置语言
  function setLanguage(language: 'zh-CN' | 'en-US') {
    updateConfig({ language })
  }

  // 设置选中的聊天模型
  function setSelectedChatModel(modelId: string) {
    updateConfig({ selected_chat_model: modelId })
  }

  // 设置选中的轻量级模型
  function setSelectedLightModel(modelId: string) {
    updateConfig({ selected_light_model: modelId })
  }

  // 重置配置
  function resetConfig() {
    config.value = defaultConfig
    saveConfig()
  }

  // 初始化时应用主题
  if (config.value.theme) {
    document.documentElement.setAttribute('data-theme', config.value.theme)
  }

  return {
    config,
    updateConfig,
    setTheme,
    setLanguage,
    setSelectedChatModel,
    setSelectedLightModel,
    resetConfig
  }
})
