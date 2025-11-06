/**
 * 管理员状态管理
 * 管理管理员登录状态和权限
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

interface AdminUser {
  username: string
  token: string
  expires_at: number
}

const ADMIN_STORAGE_KEY = 'nijichat_admin_user'

export const useAdminStore = defineStore('admin', () => {
  const adminUser = ref<AdminUser | null>(loadAdminUser())

  // 从localStorage加载管理员信息
  function loadAdminUser(): AdminUser | null {
    try {
      const stored = localStorage.getItem(ADMIN_STORAGE_KEY)
      if (stored) {
        const user = JSON.parse(stored)
        // 检查是否过期
        if (user.expires_at > Date.now()) {
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
    adminUser.value = user
    saveAdminUser(user)
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
    // 检查token是否过期
    return adminUser.value.expires_at > Date.now()
  })

  // 获取当前用户名
  const username = computed(() => adminUser.value?.username || '')

  // 获取token
  const token = computed(() => adminUser.value?.token || '')

  // 检查token是否即将过期（1小时内）
  const isTokenExpiringSoon = computed(() => {
    if (!adminUser.value) return false
    const oneHour = 60 * 60 * 1000
    return adminUser.value.expires_at - Date.now() < oneHour
  })

  return {
    // State
    adminUser,

    // Computed
    isLoggedIn,
    username,
    token,
    isTokenExpiringSoon,

    // Actions
    login,
    logout
  }
})
