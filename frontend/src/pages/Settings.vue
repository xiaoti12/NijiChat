<template>
  <div class="settings-page">
    <div class="settings-container">
      <h1 class="settings-title">设置</h1>

      <div class="settings-section">
        <h2>AI模型配置</h2>
        <p class="section-desc">配置AI模型以开始对话</p>
        <!-- TODO: 添加AI模型配置表单 -->
        <div class="placeholder-box">
          AI模型配置功能即将推出
        </div>
      </div>

      <div class="settings-section">
        <h2>主题设置</h2>
        <div class="theme-options">
          <label>
            <input type="radio" value="light" v-model="theme" @change="updateTheme" />
            <span>浅色主题</span>
          </label>
          <label>
            <input type="radio" value="dark" v-model="theme" @change="updateTheme" />
            <span>深色主题</span>
          </label>
        </div>
      </div>

      <div class="settings-section">
        <h2>WebDAV同步</h2>
        <p class="section-desc">配置WebDAV服务器以同步聊天记录</p>
        <div class="placeholder-box">
          WebDAV同步功能即将推出
        </div>
      </div>

      <div class="settings-actions">
        <button @click="goBack" class="btn btn-secondary">返回</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '@/stores/configStore'

const router = useRouter()
const configStore = useConfigStore()

const theme = ref(configStore.config.theme)

function updateTheme() {
  configStore.setTheme(theme.value as 'light' | 'dark')
}

function goBack() {
  router.push({ name: 'Home' })
}
</script>

<style scoped>
.settings-page {
  width: 100%;
  min-height: 100vh;
  padding: var(--spacing-3xl);
  background: var(--bg-secondary);
}

.settings-container {
  max-width: 800px;
  margin: 0 auto;
}

.settings-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-3xl);
  color: var(--text-primary);
}

.settings-section {
  background: var(--bg-primary);
  padding: var(--spacing-2xl);
  border-radius: var(--radius-xl);
  margin-bottom: var(--spacing-xl);
}

.settings-section h2 {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
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

.theme-options {
  display: flex;
  gap: var(--spacing-lg);
}

.theme-options label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  cursor: pointer;
}

.settings-actions {
  display: flex;
  justify-content: center;
  margin-top: var(--spacing-2xl);
}
</style>
