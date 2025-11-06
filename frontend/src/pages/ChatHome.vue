<template>
  <div class="chat-home">
    <!-- 顶部导航栏 -->
    <div class="chat-header">
      <h1 class="chat-title">NijiChat</h1>
      <div class="chat-actions">
        <button @click="goToSeiyuuLibrary" class="btn btn-primary">
          <span class="icon">+</span>
          选择声优
        </button>
      </div>
    </div>

    <!-- 主聊天区域 -->
    <div class="chat-main">
      <!-- 空状态 -->
      <div v-if="!selectedSeiyuu" class="empty-state">
        <div class="empty-icon">💬</div>
        <h2>欢迎使用 NijiChat</h2>
        <p>选择一位声优开始你的专属对话</p>
        <button @click="goToSeiyuuLibrary" class="btn btn-primary btn-large">
          浏览声优库
        </button>
      </div>

      <!-- 聊天界面 -->
      <div v-else class="chat-interface">
        <!-- 选中的声优信息 -->
        <div class="selected-seiyuu">
          <div class="seiyuu-avatar">
            <img :src="selectedSeiyuu.avatar_url || '/default-avatar.png'" :alt="selectedSeiyuu.name" />
          </div>
          <div class="seiyuu-info">
            <h3>{{ selectedSeiyuu.name }}</h3>
            <p>开始和 {{ selectedSeiyuu.name }} 对话吧！</p>
          </div>
        </div>

        <!-- 聊天消息区域 -->
        <div class="chat-messages">
          <div class="message-placeholder">
            <p>聊天功能正在开发中...</p>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="chat-input">
          <div class="input-wrapper">
            <input
              type="text"
              placeholder="输入消息..."
              class="message-input"
              disabled
            />
            <button class="send-btn" disabled>
              发送
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useSeiyuuStore } from '@/stores/seiyuuStore'
import type { Seiyuu } from '@/types'

const router = useRouter()
const route = useRoute()
const seiyuuStore = useSeiyuuStore()

const selectedSeiyuu = ref<Seiyuu | null>(null)

// 跳转到声优库
function goToSeiyuuLibrary() {
  router.push({ name: 'SeiyuuLibrary' })
}

// 从路由参数或查询参数中获取选中的声优
onMounted(() => {
  // 如果有 seiyuuId 查询参数，尝试获取声优信息
  const seiyuuId = route.query.seiyuuId as string
  if (seiyuuId) {
    selectedSeiyuu.value = seiyuuStore.getSeiyuuById(seiyuuId) || null
  }
})
</script>

<style scoped>
.chat-home {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
}

/* 顶部导航栏 */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-lg) var(--spacing-xl);
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
}

.chat-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.chat-actions .btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.icon {
  font-size: var(--font-size-lg);
  font-weight: bold;
}

/* 主聊天区域 */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: var(--spacing-3xl);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: var(--spacing-xl);
}

.empty-state h2 {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-md);
}

.empty-state p {
  color: var(--text-muted);
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-2xl);
}

.btn-large {
  padding: var(--spacing-md) var(--spacing-xl);
  font-size: var(--font-size-lg);
}

/* 聊天界面 */
.chat-interface {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 选中的声优信息 */
.selected-seiyuu {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg) var(--spacing-xl);
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
}

.seiyuu-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
}

.seiyuu-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.seiyuu-info h3 {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.seiyuu-info p {
  color: var(--text-muted);
  font-size: var(--font-size-sm);
}

/* 聊天消息区域 */
.chat-messages {
  flex: 1;
  padding: var(--spacing-xl);
  overflow-y: auto;
}

.message-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
}

/* 输入区域 */
.chat-input {
  padding: var(--spacing-lg) var(--spacing-xl);
  background: var(--bg-primary);
  border-top: 1px solid var(--border-color);
}

.input-wrapper {
  display: flex;
  gap: var(--spacing-md);
  align-items: center;
}

.message-input {
  flex: 1;
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: var(--font-size-md);
}

.message-input:focus {
  outline: none;
  border-color: var(--color-primary);
}

.send-btn {
  padding: var(--spacing-md) var(--spacing-lg);
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: var(--transition-normal);
}

.send-btn:hover:not(:disabled) {
  background: var(--color-primary-dark);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 按钮样式 */
.btn {
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: var(--transition-normal);
  text-decoration: none;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-dark);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--bg-quaternary);
}
</style>