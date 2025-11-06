<template>
  <div class="chat-page">
    <!-- 左侧对话列表 -->
    <ChatSidebar
      :conversations="conversations"
      :current-room-id="currentRoomId"
      @select-conversation="handleSelectConversation"
      @new-conversation="handleNewConversation"
    />

    <!-- 中间聊天区域 -->
    <ChatInterface
      v-if="currentRoom"
      :room="currentRoom"
      :messages="currentMessages"
      :seiyuu="currentSeiyuu"
      @send-message="handleSendMessage"
    />

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">💬</div>
      <h2 class="empty-title">欢迎使用 NijiChat</h2>

      <!-- AI未配置提示 -->
      <div v-if="!hasChatModels" class="ai-setup-notice">
        <div class="notice-icon">🤖</div>
        <p class="notice-text">要开始聊天，请先配置AI模型</p>
        <div class="action-buttons">
          <button @click="showAIManager = true" class="btn btn-primary btn-large">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" class="btn-icon">
              <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
            </svg>
            配置AI模型
          </button>
          <button @click="handleNewConversation" class="btn btn-secondary btn-large">
            浏览声优库
          </button>
        </div>
      </div>

      <!-- AI已配置，可以开始聊天 -->
      <div v-else class="ready-state">
        <p class="empty-text">选择一个对话开始聊天，或创建新的对话</p>
        <div class="action-buttons">
          <button @click="handleNewConversation" class="btn btn-primary btn-large">
            开始新对话
          </button>
          <button @click="showAIManager = true" class="btn btn-secondary btn-large">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" class="btn-icon">
              <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
            </svg>
            管理AI模型
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧设置面板 -->
    <ChatRightPanel
      v-if="currentRoom"
      :room="currentRoom"
      :seiyuu="currentSeiyuu"
      :collapsed="rightPanelCollapsed"
      @toggle="toggleRightPanel"
      @update-settings="handleUpdateSettings"
    />

    <!-- AI模型管理器弹窗 -->
    <div v-if="showAIManager" class="ai-manager-overlay" @click="showAIManager = false">
      <div class="ai-manager-container" @click.stop>
        <AIModelManager />
        <div class="manager-footer">
          <button @click="showAIManager = false" class="btn btn-secondary">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useChatStore } from '@/stores/chatStore'
import { useSeiyuuStore } from '@/stores/seiyuuStore'
import { useAIModelStore } from '@/stores/aiModelStore'
import { mockService } from '@/services/mockService'
import type { Room, Message, Seiyuu } from '@/types'

// 导入组件
import ChatSidebar from '@/components/ChatSidebar.vue'
import ChatInterface from '@/components/ChatInterface.vue'
import ChatRightPanel from '@/components/ChatRightPanel.vue'
import AIModelManager from '@/components/AIModelManager.vue'

const router = useRouter()
const route = useRoute()
const chatStore = useChatStore()
const seiyuuStore = useSeiyuuStore()
const aiModelStore = useAIModelStore()

// 状态
const rightPanelCollapsed = ref(false)
const loading = ref(false)
const showAIManager = ref(false)

// 计算属性
const conversations = computed(() => chatStore.conversations)
const currentRoomId = computed(() => chatStore.currentRoomId)
const currentRoom = computed(() => chatStore.currentRoom)
const currentMessages = computed(() => chatStore.currentMessages)

const currentSeiyuu = computed(() => {
  if (!currentRoom.value) return null
  // 1v1 对话获取对应的声优
  if (currentRoom.value.type === '1v1' && currentRoom.value.participants.length === 1) {
    return seiyuuStore.getSeiyuuById(currentRoom.value.participants[0])
  }
  return null
})

// AI配置状态
const hasAIModels = computed(() => aiModelStore.models.length > 0)
const hasChatModels = computed(() => aiModelStore.chatModels.length > 0)

// 事件处理
function handleSelectConversation(roomId: string) {
  chatStore.setCurrentRoom(roomId)
}

function handleNewConversation() {
  router.push({ name: 'SeiyuuLibrary' })
}

async function handleSendMessage(content: string) {
  if (!currentRoom.value || !currentSeiyuu.value) return

  try {
    loading.value = true

    // 添加用户消息
    const userMessage = chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: 'user-1',
      sender_name: '用户',
      content
    })

    // 生成AI回复
    const aiReply = await mockService.generateAIReply(
      currentSeiyuu.value.name,
      content,
      chatStore.getRecentMessages(currentRoom.value.id, 10)
    )

    // 添加AI回复消息
    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: currentSeiyuu.value.id,
      sender_name: currentSeiyuu.value.name,
      sender_avatar: currentSeiyuu.value.avatar_url,
      content: aiReply
    })

  } catch (error) {
    console.error('发送消息失败:', error)
  } finally {
    loading.value = false
  }
}

function toggleRightPanel() {
  rightPanelCollapsed.value = !rightPanelCollapsed.value
}

function handleUpdateSettings(settings: any) {
  // 处理设置更新
  console.log('更新设置:', settings)
}

// 初始化
onMounted(async () => {
  try {
    // 加载mock数据
    if (mockService.isMockMode()) {
      const response = await mockService.getSeiyuuList()
      if (response.success && response.data) {
        seiyuuStore.setSeiyuuList(response.data)
      }

      // 创建初始对话
      const initialConversations = mockService.generateInitialConversations()
      for (const conv of initialConversations) {
        const seiyuu = seiyuuStore.getSeiyuuById(conv.seiyuuId)
        if (seiyuu) {
          const room = chatStore.createRoom({
            type: '1v1',
            name: seiyuu.name,
            avatar: seiyuu.avatar_url,
            participants: [seiyuu.id],
            unread_count: 0
          })

          // 添加最后一条消息
          chatStore.addMessage({
            room_id: room.id,
            sender_id: seiyuu.id,
            sender_name: seiyuu.name,
            sender_avatar: seiyuu.avatar_url,
            content: conv.lastMessage
          })
        }
      }
    }

    // 处理路由参数
    const seiyuuId = route.query.seiyuuId as string
    if (seiyuuId) {
      // 查找或创建该声优的对话
      const existingRoom = chatStore.rooms.find(room =>
        room.type === '1v1' && room.participants.includes(seiyuuId)
      )

      if (existingRoom) {
        chatStore.setCurrentRoom(existingRoom.id)
      } else {
        // 创建新对话
        const seiyuu = seiyuuStore.getSeiyuuById(seiyuuId)
        if (seiyuu) {
          const room = chatStore.createRoom({
            type: '1v1',
            name: seiyuu.name,
            avatar: seiyuu.avatar_url,
            participants: [seiyuu.id],
            unread_count: 0
          })
          chatStore.setCurrentRoom(room.id)
        }
      }
    }
  } catch (error) {
    console.error('初始化失败:', error)
  }
})
</script>

<style scoped>
/* CSS 变量定义 */
:root {
  /* 品牌渐变色 */
  --gradient-primary: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

  /* 主色 */
  --color-primary: #667eea;
  --color-primary-dark: #5568d3;
  --color-primary-light: rgba(102, 126, 234, 0.1);

  /* 文本颜色 */
  --text-primary: #1f2937;
  --text-secondary: #374151;
  --text-muted: #6b7280;

  /* 背景颜色 */
  --bg-primary: #ffffff;
  --bg-secondary: #fafafa;
  --bg-tertiary: #f9fafb;

  /* 边框颜色 */
  --border-light: #e5e7eb;

  /* 间距 */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 12px;
  --spacing-lg: 16px;
  --spacing-xl: 20px;
  --spacing-2xl: 24px;
  --spacing-3xl: 32px;

  /* 圆角 */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;

  /* 字体 */
  --font-size-sm: 13px;
  --font-size-md: 14px;
  --font-size-lg: 16px;
  --font-size-xl: 18px;
  --font-size-2xl: 20px;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;

  /* 过渡 */
  --transition-normal: all 0.2s ease;
}

/* 三栏布局 */
.chat-page {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-secondary);
}

/* 空状态样式 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 48px;
  background: var(--bg-primary);
}

.empty-icon {
  width: 64px;
  height: 64px;
  background: #f3f4f6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  font-size: 24px;
  color: #9ca3af;
}

.empty-title {
  font-size: 18px;
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-text {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 24px;
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-size: var(--font-size-md);
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-large {
  padding: 10px 24px;
  font-size: var(--font-size-lg);
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #e5e7eb;
  border-color: #9ca3af;
}

.btn-icon {
  margin-right: 6px;
}

/* AI配置提示样式 */
.ai-setup-notice {
  margin-top: 16px;
  padding: 24px;
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 12px;
  text-align: center;
}

.notice-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.notice-text {
  font-size: 14px;
  color: #92400e;
  margin-bottom: 20px;
  font-weight: 500;
}

.ready-state {
  margin-top: 16px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

/* AI模型管理器弹窗 */
.ai-manager-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.ai-manager-container {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 700px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.manager-footer {
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  background: #fafafa;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-page {
    flex-direction: column;
  }

  .action-buttons {
    flex-direction: column;
    align-items: stretch;
  }

  .btn-large {
    width: 100%;
  }

  .ai-manager-container {
    max-width: 95%;
    max-height: 90vh;
  }

  .empty-state {
    padding: 32px 20px;
  }

  .ai-setup-notice {
    padding: 20px 16px;
  }
}

@media (max-width: 1400px) {
  /* 右侧面板浮动 */
  .chat-page :deep(.chat-right-panel) {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
  }
}
</style>