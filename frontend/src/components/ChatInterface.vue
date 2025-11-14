<template>
  <div class="chat-interface">
    <!-- 聊天头部 -->
    <div class="chat-header">
      <!-- 移动端汉堡菜单按钮 -->
      <button class="menu-btn" @click="handleToggleSidebar">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
          stroke-linecap="round" stroke-linejoin="round">
          <line x1="3" y1="12" x2="21" y2="12"></line>
          <line x1="3" y1="6" x2="21" y2="6"></line>
          <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
      </button>

      <div class="header-avatar">
        <!-- 单人对话头像 -->
        <template v-if="!isDualConversation">
          <img v-if="currentSeiyuuAvatar" :src="currentSeiyuuAvatar" :alt="seiyuu?.name || '?'" class="avatar" />
          <div v-else class="avatar-placeholder">
            {{ seiyuu?.name?.charAt(0) || '?' }}
          </div>
        </template>

        <!-- 双人对话头像 -->
        <template v-else-if="dualSeiyuu">
          <div class="dual-avatar-group">
            <div class="dual-avatar dual-avatar-1">
              <img v-if="dualSeiyuu.initiator.avatar" :src="dualSeiyuu.initiator.avatar"
                :alt="dualSeiyuu.initiator.name" class="avatar" />
              <div v-else class="avatar-placeholder">
                {{ dualSeiyuu.initiator.name?.charAt(0) || '?' }}
              </div>
            </div>
            <div class="dual-avatar dual-avatar-2">
              <img v-if="dualSeiyuu.responder.avatar" :src="dualSeiyuu.responder.avatar"
                :alt="dualSeiyuu.responder.name" class="avatar" />
              <div v-else class="avatar-placeholder">
                {{ dualSeiyuu.responder.name?.charAt(0) || '?' }}
              </div>
            </div>
          </div>
        </template>
      </div>

      <div class="header-info">
        <h3 class="header-name">{{ headerTitle }}</h3>
        <p v-if="isDualConversation && room.dual_topic" class="header-topic">{{ room.dual_topic }}</p>
      </div>

      <div class="header-actions">
        <button class="action-btn" title="设置" @click="handleToggleSettings">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 消息容器 -->
    <div class="messages-container" ref="messagesContainer">
      <div class="messages-list">
        <!-- 欢迎消息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <!-- 单人对话欢迎消息 -->
          <template v-if="!isDualConversation">
            <div class="welcome-avatar">
              <img v-if="currentSeiyuuAvatar" :src="currentSeiyuuAvatar" :alt="seiyuu?.name || '?'" class="avatar" />
              <div v-else class="avatar-placeholder">
                {{ seiyuu?.name?.charAt(0) || '?' }}
              </div>
            </div>
            <div class="welcome-content">
              <h4>开始与 {{ seiyuu?.name }} 对话</h4>
              <p>你好！我是 {{ seiyuu?.name }}，很高兴和你聊天～有什么想说的吗？</p>
            </div>
          </template>

          <!-- 双人对话欢迎消息 -->
          <template v-else-if="dualSeiyuu">
            <div class="welcome-avatar">
              <div class="dual-avatar-group welcome-dual">
                <div class="dual-avatar dual-avatar-1">
                  <img v-if="dualSeiyuu.initiator.avatar" :src="dualSeiyuu.initiator.avatar"
                    :alt="dualSeiyuu.initiator.name" class="avatar" />
                  <div v-else class="avatar-placeholder">
                    {{ dualSeiyuu.initiator.name?.charAt(0) || '?' }}
                  </div>
                </div>
                <div class="dual-avatar dual-avatar-2">
                  <img v-if="dualSeiyuu.responder.avatar" :src="dualSeiyuu.responder.avatar"
                    :alt="dualSeiyuu.responder.name" class="avatar" />
                  <div v-else class="avatar-placeholder">
                    {{ dualSeiyuu.responder.name?.charAt(0) || '?' }}
                  </div>
                </div>
              </div>
            </div>
            <div class="welcome-content">
              <h4>{{ dualSeiyuu.initiator.name }} × {{ dualSeiyuu.responder.name }} 双人剧场</h4>
              <p>{{ room.dual_topic || '两位角色正准备开始对话，点击下方按钮开始剧场表演～' }}</p>
            </div>
          </template>
        </div>

        <!-- 消息列表 -->
        <div v-for="message in messages" :key="message.id" :class="[
          'message-wrapper',
          getMessageClass(message)
        ]">
          <!-- 声优消息 (左侧) 包括单人对话和双人对话的响应者 -->
          <div v-if="message.sender_id !== 'user-1' && getMessageClass(message) !== 'initiator-message'"
            class="message-row">
            <div class="message-avatar">
              <img v-if="getMessageSenderAvatar(message)" :src="getMessageSenderAvatar(message)"
                :alt="message.sender_name" class="avatar" />
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

          <!-- 发起者消息 (右侧) 包括用户消息和双人对话的发起者 -->
          <div v-else class="message-row user-row">
            <div class="message-avatar">
              <!-- 真实用户头像 -->
              <div v-if="message.sender_id === 'user-1'" class="user-avatar">
                <svg width="20" height="20" viewBox="0 0 16 16" fill="currentColor">
                  <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0z" />
                  <path
                    d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8-7a7 7 0 0 0-5.468 11.37C3.242 11.226 4.805 10 8 10s4.757 1.225 5.468 2.37A7 7 0 0 0 8 1z" />
                </svg>
              </div>
              <!-- 双人对话发起者头像 -->
              <template v-else>
                <img v-if="getMessageSenderAvatar(message)" :src="getMessageSenderAvatar(message)"
                  :alt="message.sender_name" class="avatar" />
                <div v-else class="avatar-placeholder">
                  {{ message.sender_name?.charAt(0) || '?' }}
                </div>
              </template>
            </div>

            <div class="message-content">
              <div class="message-header">
                <span class="message-time">{{ formatTime(message.timestamp) }}</span>
                <span class="sender-name">{{ message.sender_name }}</span>
              </div>
              <div class="user-bubble">
                {{ message.content }}
              </div>
            </div>
          </div>
        </div>

        <!-- 正在输入指示器 -->
        <div v-if="isTyping" class="typing-indicator">
          <!-- 双人对话的加载指示器 -->
          <div v-if="isDualConversation && nextSpeaker" class="message-row"
            :class="{ 'user-row': dualSeiyuu && nextSpeaker.id === dualSeiyuu.initiator.id }">
            <!-- 当是发起者时，头像在前（会被反转到右边） -->
            <template v-if="dualSeiyuu && nextSpeaker.id === dualSeiyuu.initiator.id">
              <div class="message-avatar">
                <img v-if="nextSpeaker.avatar" :src="nextSpeaker.avatar" :alt="nextSpeaker.name" class="avatar" />
                <div v-else class="avatar-placeholder">
                  {{ nextSpeaker.name?.charAt(0) || '?' }}
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
            </template>
            <!-- 当是响应者时，头像在前（正常排列在左边） -->
            <template v-else>
              <div class="message-avatar">
                <img v-if="nextSpeaker.avatar" :src="nextSpeaker.avatar" :alt="nextSpeaker.name" class="avatar" />
                <div v-else class="avatar-placeholder">
                  {{ nextSpeaker.name?.charAt(0) || '?' }}
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
            </template>
          </div>

          <!-- 单人对话的加载指示器 -->
          <div v-else class="message-row">
            <div class="message-avatar">
              <img v-if="currentSeiyuuAvatar" :src="currentSeiyuuAvatar" :alt="seiyuu?.name || '?'" class="avatar" />
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
    <div class="input-area" v-if="!isDualConversation">
      <div class="input-container">
        <div class="input-wrapper" :class="{ focused: inputFocused }">
          <textarea ref="messageInput" v-model="currentMessage" placeholder="输入消息..." class="message-input" rows="1"
            @keydown="handleKeyDown" @focus="inputFocused = true" @blur="inputFocused = false" @input="handleInput" />

          <button :disabled="!canSend" class="send-btn" @click="handleSend">
            <svg v-if="!isLoading" width="18" height="18" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M15.964.686a.5.5 0 0 0-.65-.65L.767 5.855H.766l-.452.18a.5.5 0 0 0-.082.887l.41.26.001.002 4.995 3.178 3.178 4.995.002.002.26.41a.5.5 0 0 0 .886-.083l6-15Zm-1.833 1.89L6.637 10.07l-.215-.338a.5.5 0 0 0-.154-.154l-.338-.215 7.494-7.494 1.178-.471-.47 1.178Z" />
            </svg>
            <div v-else class="loading-spinner"></div>
          </button>
        </div>

      </div>

      <!-- 输入提示 -->
      <div class="input-hint">
        按 <kbd>Enter</kbd> 发送，<kbd>Shift + Enter</kbd> 换行
      </div>
    </div>

    <!-- 双人对话控制区域 -->
    <div v-else class="dual-control-area">
      <button :disabled="isLoading" class="continue-btn" @click="handleContinueDual">
        <svg v-if="!isLoading" width="18" height="18" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="m11.596 8.697-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393z" />
        </svg>
        <div v-else class="loading-spinner"></div>
        <span>{{ isLoading ? '正在生成...' : '继续对话' }}</span>
      </button>

      <div class="dual-hint">
        点击按钮让角色们继续对话
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import type { Room, Message, Seiyuu } from '@/types'
import { useChatStore } from '@/stores/chatStore'
import { useSeiyuuStore } from '@/stores/seiyuuStore'

// Props
const props = defineProps<{
  room: Room
  messages: Message[]
  seiyuu: Seiyuu | null
  useRealAI?: boolean
}>()

// Emits
const emit = defineEmits<{
  sendMessage: [content: string]
  toggleSettings: []
  toggleSidebar: []
}>()

// Store
const chatStore = useChatStore()
const seiyuuStore = useSeiyuuStore()

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

// 判断是否为双人对话
const isDualConversation = computed(() => {
  return props.room.type === 'dual_theater'
})

// 获取双人对话的声优信息（动态获取最新头像）
const dualSeiyuu = computed(() => {
  if (!isDualConversation.value) return null
  const dualInfo = chatStore.getCurrentDualSeiyuu()

  if (!dualInfo) return null

  // 动态获取最新头像
  const seiyuu1Latest = seiyuuStore.getSeiyuuById(dualInfo.seiyuu1.id)
  const seiyuu2Latest = seiyuuStore.getSeiyuuById(dualInfo.seiyuu2.id)

  return {
    ...dualInfo,
    seiyuu1: {
      ...dualInfo.seiyuu1,
      avatar: seiyuu1Latest?.avatar_url || dualInfo.seiyuu1.avatar
    },
    seiyuu2: {
      ...dualInfo.seiyuu2,
      avatar: seiyuu2Latest?.avatar_url || dualInfo.seiyuu2.avatar
    },
    initiator: {
      ...dualInfo.initiator,
      avatar: dualInfo.initiator.id === dualInfo.seiyuu1.id
        ? (seiyuu1Latest?.avatar_url || dualInfo.initiator.avatar)
        : (seiyuu2Latest?.avatar_url || dualInfo.initiator.avatar)
    },
    responder: {
      ...dualInfo.responder,
      avatar: dualInfo.responder.id === dualInfo.seiyuu1.id
        ? (seiyuu1Latest?.avatar_url || dualInfo.responder.avatar)
        : (seiyuu2Latest?.avatar_url || dualInfo.responder.avatar)
    }
  }
})

// 计算下一个发言者（用于加载动画）
const nextSpeaker = computed(() => {
  if (!isDualConversation.value || !dualSeiyuu.value) return null

  const { initiator, responder } = dualSeiyuu.value
  const lastMessage = props.messages[props.messages.length - 1]

  // 如果没有消息，下一个发言者是发起者
  if (!lastMessage) {
    return initiator
  }

  // 如果最后一条消息是发起者说的，下一个应该是响应者
  if (lastMessage.sender_id === initiator.id) {
    return responder
  }

  // 否则，下一个是发起者
  return initiator
})

// 头部显示名称
const headerTitle = computed(() => {
  if (isDualConversation.value && dualSeiyuu.value) {
    return `${dualSeiyuu.value.initiator.name} × ${dualSeiyuu.value.responder.name}`
  }
  return props.seiyuu?.name || '未知声优'
})

// 动态获取声优头像
const currentSeiyuuAvatar = computed(() => {
  if (!props.seiyuu?.id) return props.seiyuu?.avatar_url
  const latestSeiyuu = seiyuuStore.getSeiyuuById(props.seiyuu.id)
  return latestSeiyuu?.avatar_url || props.seiyuu.avatar_url
})

// 动态获取消息发送者头像
function getMessageSenderAvatar(message: Message): string | undefined {
  if (message.sender_id === 'user-1') return undefined

  // 从 seiyuuStore 获取最新头像
  const latestSeiyuu = seiyuuStore.getSeiyuuById(message.sender_id)
  return latestSeiyuu?.avatar_url || message.sender_avatar
}

// 方法
function getMessageClass(message: Message) {
  if (message.sender_id === 'user-1') {
    return 'user-message'
  }

  // 双人对话中，根据发起者确定消息位置
  if (isDualConversation.value && dualSeiyuu.value) {
    // 发起者的消息显示在右侧（类似用户消息）
    if (message.sender_id === dualSeiyuu.value.initiator.id) {
      return 'initiator-message'
    } else {
      return 'responder-message'
    }
  }

  return 'seiyuu-message'
}

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
    // 发生错误时隐藏加载状态
    isLoading.value = false
    isTyping.value = false
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

function handleToggleSettings() {
  emit('toggleSettings')
}

function handleToggleSidebar() {
  emit('toggleSidebar')
}

function handleContinueDual() {
  if (!isDualConversation.value) return

  isTyping.value = true
  isLoading.value = true

  try {
    // 双人对话使用空字符串触发继续对话
    emit('sendMessage', '')
  } catch (error) {
    console.error('继续双人对话失败:', error)
  } finally {
    // 无论成功或失败，都重置加载状态，但 isTyping 的重置由父组件控制
    isLoading.value = false
    // isTyping.value 不在此处重置，等待父组件调用 stopTyping
  }
}

// 暴露方法给父组件，用于控制加载状态
function stopTyping() {
  isTyping.value = false
  isLoading.value = false
}

// 暴露给父组件的方法
defineExpose({
  stopTyping
})

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

/* 移动端汉堡菜单按钮 */
.menu-btn {
  display: none;
  /* 默认隐藏,仅在移动端显示 */
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: all 0.2s;
  margin-right: 8px;
  padding: 0;
}

.menu-btn:hover {
  background: #f3f4f6;
  color: #374151;
}

.menu-btn:active {
  transform: scale(0.95);
}

.header-avatar {
  margin-right: 12px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(254, 173, 0, 0.1);
}

.avatar-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #fead00, #ff791b);
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
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 2px 0;
}

.header-topic {
  font-size: 12px;
  color: #6b7280;
  margin: 0;
  font-style: italic;
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

/* 双人头像组合 */
.dual-avatar-group {
  position: relative;
  width: 56px;
  height: 48px;
  display: flex;
  align-items: center;
}

.dual-avatar-group.welcome-dual {
  width: 64px;
  height: 56px;
}

.dual-avatar {
  position: absolute;
  border: 2px solid white;
  border-radius: 50%;
  transition: transform 0.2s ease;
}

.dual-avatar:hover {
  transform: scale(1.05);
}

.dual-avatar-1 {
  left: 0;
  z-index: 2;
  box-shadow: 0 2px 8px rgba(254, 173, 0, 0.2);
}

.dual-avatar-2 {
  right: 0;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.2);
}

.dual-avatar .avatar,
.dual-avatar .avatar-placeholder {
  width: 32px;
  height: 32px;
  font-size: 12px;
}

.welcome-dual .dual-avatar .avatar,
.welcome-dual .dual-avatar .avatar-placeholder {
  width: 40px;
  height: 40px;
  font-size: 16px;
}

/* 双人头像组合在聊天头部的样式 */
.chat-header .dual-avatar-group {
  width: 60px;
  height: 48px;
}

.chat-header .dual-avatar .avatar,
.chat-header .dual-avatar .avatar-placeholder {
  width: 36px;
  height: 36px;
  font-size: 14px;
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

/* 双人对话消息样式 */
.message-wrapper.initiator-message {
  /* 发起者消息显示在右侧，类似用户消息 */
}

.message-wrapper.initiator-message .user-bubble {
  /* 发起者消息使用稍微不同的渐变，表示是剧场发起者 */
  background: linear-gradient(135deg, #fead00, #ff9800);
  box-shadow: 0 2px 8px rgba(254, 173, 0, 0.35);
}

.message-wrapper.responder-message {
  /* 响应者消息显示在左侧，类似声优消息 */
}

.message-wrapper.responder-message .seiyuu-bubble {
  /* 响应者消息使用稍微不同的边框，表示是剧场响应者 */
  border-left: 3px solid #4f46e5;
  background: white;
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
  background: linear-gradient(135deg, #fead00, #ff791b);
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
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  border-radius: 16px 16px 4px 16px;
  box-shadow: 0 2px 8px rgba(254, 173, 0, 0.3);
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

.typing-dots span:nth-child(1) {
  animation-delay: 0s;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {

  0%,
  60%,
  100% {
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
  padding: 8px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
  min-height: 44px;
}

.input-wrapper.focused {
  border-color: #fead00;
  box-shadow: 0 0 0 3px rgba(254, 173, 0, 0.1);
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
  min-height: 24px;
  max-height: 120px;
  font-family: inherit;
  padding: 2px 0;
}

.message-input::placeholder {
  color: #9ca3af;
}

.send-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: linear-gradient(135deg, #fead00, #ff791b);
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
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
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
  to {
    transform: rotate(360deg);
  }
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

/* 双人对话控制区域 */
.dual-control-area {
  background: white;
  border-top: 1px solid #e5e7eb;
  padding: 20px 24px;
  text-align: center;
}

.continue-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 8px;
}

.continue-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
}

.continue-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.dual-hint {
  font-size: 12px;
  color: #9ca3af;
}

/* 响应式 */
@media (max-width: 768px) {

  /* 移动端显示汉堡菜单按钮 */
  .menu-btn {
    display: flex;
  }

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