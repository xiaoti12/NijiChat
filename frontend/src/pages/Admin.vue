<template>
  <div class="admin-page">
    <div class="admin-container">
      <h1 class="admin-title">管理后台</h1>

      <div v-if="!adminStore.isLoggedIn && !isDevelopmentMode" class="login-section">
        <h2>管理员登录</h2>
        <form @submit.prevent="handleLogin" class="login-form">
          <input
            v-model="loginForm.username"
            type="text"
            placeholder="用户名"
            class="input"
            required
          />
          <input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            class="input"
            required
          />
          <button type="submit" class="btn btn-primary" :disabled="logging">
            {{ logging ? '登录中...' : '登录' }}
          </button>
          <p v-if="loginError" class="error-text">{{ loginError }}</p>
        </form>
      </div>

      <div v-if="adminStore.isLoggedIn || isDevelopmentMode" class="admin-content">
        <div class="admin-header">
          <div class="admin-user-info">
            <p>欢迎，{{ adminStore.username || (isDevelopmentMode ? '开发者' : '') }}</p>
            <p v-if="isDevelopmentMode" class="dev-mode-hint">当前为开发模式，无需登录验证</p>
          </div>
          <button v-if="!isDevelopmentMode" @click="handleLogout" class="btn btn-secondary">退出登录</button>
        </div>

        <!-- Tab 导航 -->
        <div class="admin-tabs">
          <button
            class="tab-button"
            :class="{ active: activeTab === 'seiyuu' }"
            @click="activeTab = 'seiyuu'"
          >
            声优管理
          </button>
          <button
            class="tab-button"
            :class="{ active: activeTab === 'groups' }"
            @click="activeTab = 'groups'"
          >
            群组管理
          </button>
        </div>

        <!-- Tab 内容 -->
        <div class="tab-content">
          <SeiyuuManager v-if="activeTab === 'seiyuu'" />
          <GroupManager v-if="activeTab === 'groups'" />
        </div>
      </div>

      <div class="admin-actions">
        <button @click="goBack" class="btn btn-secondary">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/adminStore'
import { adminLogin } from '@/services/apiService'
import SeiyuuManager from '@/components/admin/SeiyuuManager.vue'
import GroupManager from '@/components/admin/GroupManager.vue'

const router = useRouter()
const adminStore = useAdminStore()

const loginForm = reactive({
  username: '',
  password: ''
})

const logging = ref(false)
const loginError = ref<string | null>(null)
const activeTab = ref('seiyuu')

// 开发模式检测
const isDevelopmentMode = computed(() => {
  return import.meta.env.DEV || import.meta.env.MODE === 'development'
})

async function handleLogin() {
  logging.value = true
  loginError.value = null

  try {
    const response = await adminLogin(loginForm)
    if (response.success && response.data) {
      adminStore.login(
        response.data.username,
        response.data.token,
        response.data.expires_at
      )
    } else {
      loginError.value = response.message || '登录失败'
    }
  } catch (err: any) {
    loginError.value = err.message || '网络错误'
  } finally {
    logging.value = false
  }
}

function handleLogout() {
  adminStore.logout()
}

function goBack() {
  router.push({ name: 'Home' })
}
</script>

<style scoped>
.admin-page {
  width: 100%;
  min-height: 100vh;
  padding: var(--spacing-3xl);
  background: var(--bg-secondary);
}

.admin-container {
  max-width: 900px;
  margin: 0 auto;
}

.admin-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-3xl);
  text-align: center;
  color: var(--text-primary);
}

.login-section {
  background: var(--bg-primary);
  padding: var(--spacing-3xl);
  border-radius: var(--radius-xl);
  max-width: 400px;
  margin: 0 auto;
}

.login-section h2 {
  font-size: var(--font-size-xl);
  margin-bottom: var(--spacing-xl);
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.error-text {
  color: var(--color-error);
  font-size: var(--font-size-sm);
  text-align: center;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xl);
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-xl);
}

.admin-user-info p {
  margin: 0;
  font-weight: var(--font-weight-medium);
}

.dev-mode-hint {
  font-size: var(--font-size-xs);
  color: #f59e0b;
  font-weight: var(--font-weight-normal) !important;
  margin-top: var(--spacing-xs) !important;
}

/* Tab 导航样式 */
.admin-tabs {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xl);
  background: var(--bg-primary);
  padding: var(--spacing-md);
  border-radius: var(--radius-lg);
}

.tab-button {
  flex: 1;
  padding: var(--spacing-md) var(--spacing-lg);
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--text-muted);
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-button:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.tab-button.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

/* Tab 内容区域 */
.tab-content {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  min-height: 600px;
}

.admin-actions {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-2xl);
}
</style>
