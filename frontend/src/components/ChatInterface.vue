<template>
  <div class="chat-interface">
    <!-- 聊天头部 -->
    <div class="chat-header">
      <div class="header-avatar">
        <img
          v-if="seiyuu?.avatar_url"
          :src="seiyuu.avatar_url"
          :alt="seiyuu.name"
          class="avatar"
        />
        <div v-else class="avatar-placeholder">
          {{ seiyuu?.name?.charAt(0) || '?' }}
        </div>
      </div>

      <div class="header-info">
        <h3 class="header-name">{{ seiyuu?.name || '未知声优' }}</h3>
        <p class="header-status">在线 · 声优模拟对话</p>
      </div>

      <div class="header-actions">
        <button class="action-btn" title="设置">
          <svg width="18" height="18" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"/>
            <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319z"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- 消息容器 -->
    <div class="messages-container" ref="messagesContainer">
      <div class="messages-list">
        <!-- 欢迎消息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="welcome-avatar">
            <img
              v-if="seiyuu?.avatar_url"
              :src="seiyuu.avatar_url"
              :alt="seiyuu.name"
              class="avatar"
            />
            <div v-else class="avatar-placeholder">
              {{ seiyuu?.name?.charAt(0) || '?' }}
            </div>
          </div>
          <div class="welcome-content">
            <h4>开始与 {{ seiyuu?.name }} 对话</h4>
            <p>你好！我是 {{ seiyuu?.name }}，很高兴和你聊天～有什么想说的吗？</p>
          </div>
        </div>

        <!-- 消息列表 -->
        <div
          v-for="message in messages"
          :key="message.id"
          :class="[
            'message-wrapper',
            message.sender_id === 'user-1' ? 'user-message' : 'seiyuu-message'
          ]"
        >
          <!-- 声优消息 -->
          <div v-if="message.sender_id !== 'user-1'" class="message-row">
            <div class="message-avatar">
              <img
                v-if="message.sender_avatar"
                :src="message.sender_avatar"
                :alt="message.sender_name"
                class="avatar"
              />
              <div v-else class="avatar-placeholder">
                {{ message.sender_name?.charAt(0) || '?' }}
              </div>
            </div>

            <div class="message-content">
              <div class="message-header">
                <span class="sender-name">{{ message.sender_name }}</span>
                <span class="message-time">{{ formatTime(message.timestamp) }}</span>
              </div>
              <div class="seiyuu-bubble">
                {{ message.content }}
              </div>
            </div>
          </div>

          <!-- 用户消息 -->
          <div v-else class="message-row user-row">
            <div class="message-content">
              <div class="message-header">
                <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                <span class="sender-name">{{ message.sender_name }}</span>
              </div>
              <div class="user-bubble">
                {{ message.content }}
              </div>
            </div>

            <div class="message-avatar">
              <div class="user-avatar">
                <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0z"/>
                  <path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1z"/>
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- 正在输入指示器 -->
        <div v-if="isTyping" class="typing-indicator">
          <div class="message-row">
            <div class="message-avatar">
              <img
                v-if="seiyuu?.avatar_url"
                :src="seiyuu.avatar_url"
                :alt="seiyuu.name"
                class="avatar"
              />
              <div v-else class="avatar-placeholder">
                {{ seiyuu?.name?.charAt(0) || '?' }}
              </div>
            </div>

            <div class="message-content">
              <div class="typing-bubble">
                <div class="typing-dots">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <div class="input-container">
        <div class="input-wrapper" :class="{ focused: inputFocused }">
          <textarea
            ref="messageInput"
            v-model="currentMessage"
            placeholder="输入消息..."
            class="message-input"
            rows="1"
            @keydown="handleKeyDown"
            @focus="inputFocused = true"
            @blur="inputFocused = false"
            @input="handleInput"
          />

          <button
            :disabled="!canSend"
            class="send-btn"
            @click="handleSend"
          >
            <svg v-if="!isLoading" width="18" height="18" viewBox="0 0 16 16" fill="currentColor">
              <path d="M15.964.686a.5.5 0 0 0-.65-.65L.767 5.855H.766l-.452.18a.5.5 0 0 0-.082.887l.41.26.001.002 4.995 3.178 3.178 4.995.002.002.26.41a.5.5 0 0 0 .886-.083l6-15Zm-1.833 1.89L6.637 10.07l-.215-.338a.5.5 0 0 0-.154-.154l-.338-.215 7.494-7.494 1.178-.471-.47 1.178Z"/>
            </svg>
            <div v-else class="loading-spinner"></div>
          </button>
        </div>

        <!-- 快捷操作 -->
        <div class="input-actions">
          <button class="action-btn" title="表情">😊</button>
          <button class="action-btn" title="附件">📎</button>
        </div>
      </div>

      <!-- 输入提示 -->
      <div class="input-hint">
        按 <kbd>Enter</kbd> 发送，<kbd>Shift + Enter</kbd> 换行
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import type { Room, Message, Seiyuu } from '@/types'

// Props
const props = defineProps<{
  room: Room
  messages: Message[]
  seiyuu: Seiyuu | null
}>()

// Emits
const emit = defineEmits<{
  sendMessage: [content: string]
}>()

// 响应式数据
const currentMessage = ref('')
const inputFocused = ref(false)
const isLoading = ref(false)
const isTyping = ref(false)
const messagesContainer = ref<HTMLElement>()
const messageInput = ref<HTMLTextAreaElement>()

// 计算属性
const canSend = computed(() => {
  return currentMessage.value.trim().length > 0 && !isLoading.value
})

// 方法
function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    handleSend()
  }
}

function handleInput() {
  // 自动调整输入框高度
  if (messageInput.value) {
    messageInput.value.style.height = 'auto'
    messageInput.value.style.height = messageInput.value.scrollHeight + 'px'
  }
}

async function handleSend() {
  if (!canSend.value) return

  const content = currentMessage.value.trim()
  currentMessage.value = ''

  // 重置输入框高度
  if (messageInput.value) {
    messageInput.value.style.height = 'auto'
  }

  // 显示正在输入状态
  isTyping.value = true
  isLoading.value = true

  try {
    await emit('sendMessage', content)
  } catch (error) {
    console.error('发送消息失败:', error)
  } finally {
    isLoading.value = false
    setTimeout(() => {
      isTyping.value = false
    }, 1000) // 1秒后隐藏输入状态
  }

  // 滚动到底部
  scrollToBottom()
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function formatTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 监听消息变化，自动滚动到底部
watch(
  () => props.messages.length,
  () => {
    scrollToBottom()
  }
)
</script>

<style scoped>
.chat-interface {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #fafafa;
}

/* 聊天头部 */
.chat-header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  padding: 16px 24px;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.header-avatar {
  margin-right: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(102, 126, 234, 0.1);
}

.avatar-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.header-info {
  flex: 1;
}

.header-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 2px 0;
}

.header-status {
  font-size: 13px;
  color: #10b981;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

/* 消息容器 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 0;
}

.messages-list {
  padding: 16px;
  min-height: 100%;
  display: flex;
  flex-direction: column;
}

/* 欢迎消息 */
.welcome-message {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  border: 1px solid #e5e7eb;
  margin-bottom: 16px;
}

.welcome-avatar .avatar {
  width: 56px;
  height: 56px;
}

.welcome-avatar .avatar-placeholder {
  width: 56px;
  height: 56px;
  font-size: 20px;
}

.welcome-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.welcome-content p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

/* 消息行 */
.message-wrapper {
  margin-bottom: 16px;
}

.message-row {
  display: flex;
  align-items: flex-end;
  gap: 12px;
}

.user-row {
  flex-direction: row-reverse;
}

.message-avatar .avatar {
  width: 32px;
  height: 32px;
}

.message-avatar .avatar-placeholder {
  width: 32px;
  height: 32px;
  font-size: 14px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.message-content {
  max-width: 70%;
  min-width: 0;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 11px;
  color: #9ca3af;
}

.user-message .message-header {
  justify-content: flex-end;
}

.sender-name {
  font-weight: 500;
}

/* 消息气泡 */
.seiyuu-bubble {
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
  border-radius: 16px 16px 16px 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.4;
}

.user-bubble {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-radius: 16px 16px 4px 16px;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.4;
}

/* 正在输入指示器 */
.typing-indicator {
  margin-bottom: 16px;
}

.typing-bubble {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px 16px 16px 4px;
  padding: 12px 16px;
  display: flex;
  align-items: center;
}

.typing-dots {
  display: flex;
  gap: 4px;
}

.typing-dots span {
  width: 6px;
  height: 6px;
  background: #9ca3af;
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-dots span:nth-child(1) { animation-delay: 0s; }
.typing-dots span:nth-child(2) { animation-delay: 0.2s; }
.typing-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.5;
  }
  30% {
    transform: translateY(-8px);
    opacity: 1;
  }
}

/* 输入区域 */
.input-area {
  background: white;
  border-top: 1px solid #e5e7eb;
  padding: 16px 24px;
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  margin-bottom: 8px;
}

.input-wrapper {
  flex: 1;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px;
  display: flex;
  align-items: flex-end;
  gap: 8px;
  transition: all 0.2s;
}

.input-wrapper.focused {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  background: white;
}

.message-input {
  flex: 1;
  border: none;
  background: transparent;
  resize: none;
  outline: none;
  font-size: 14px;
  line-height: 20px;
  min-height: 20px;
  max-height: 120px;
  font-family: inherit;
}

.message-input::placeholder {
  color: #9ca3af;
}

.send-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: all 0.2s;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.input-actions {
  display: flex;
  gap: 4px;
}

.input-actions .action-btn {
  width: 32px;
  height: 32px;
  font-size: 16px;
}

/* 输入提示 */
.input-hint {
  font-size: 11px;
  color: #9ca3af;
  text-align: center;
}

kbd {
  background: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 3px;
  padding: 1px 4px;
  font-size: 10px;
  font-family: monospace;
}

/* 响应式 */
@media (max-width: 768px) {
  .chat-header {
    padding: 12px 16px;
  }

  .messages-list {
    padding: 12px;
  }

  .message-content {
    max-width: 85%;
  }

  .input-area {
    padding: 12px 16px;
  }
}
</style>