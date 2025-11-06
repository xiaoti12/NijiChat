<template>
  <div class="admin-page">
    <div class="admin-container">
      <h1 class="admin-title">管理后台</h1>

      <div v-if="!adminStore.isLoggedIn" class="login-section">
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

      <div v-else class="admin-content">
        <div class="admin-header">
          <p>欢迎，{{ adminStore.username }}</p>
          <button @click="handleLogout" class="btn btn-secondary">退出登录</button>
        </div>

        <div class="admin-sections">
          <div class="admin-section">
            <h2>声优管理</h2>
            <p class="section-desc">添加、编辑和发布声优资料</p>
            <div class="placeholder-box">
              声优管理功能即将推出
            </div>
          </div>

          <div class="admin-section">
            <h2>群组管理</h2>
            <p class="section-desc">创建和管理声优群组</p>
            <div class="placeholder-box">
              群组管理功能即将推出
            </div>
          </div>
        </div>
      </div>

      <div class="admin-actions">
        <button @click="goBack" class="btn btn-secondary">返回首页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/adminStore'
import { adminLogin } from '@/services/apiService'

const router = useRouter()
const adminStore = useAdminStore()

const loginForm = reactive({
  username: '',
  password: ''
})

const logging = ref(false)
const loginError = ref<string | null>(null)

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

.admin-sections {
  display: grid;
  gap: var(--spacing-xl);
}

.admin-section {
  background: var(--bg-primary);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-xl);
}

.admin-section h2 {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-sm);
}

.section-desc {
  color: var(--text-muted);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-lg);
}

.placeholder-box {
  padding: var(--spacing-xl);
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  text-align: center;
  color: var(--text-muted);
}

.admin-actions {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-2xl);
}
</style>
