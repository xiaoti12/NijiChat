<template>
  <div class="chat-page">
    <!-- 左侧对话列表 -->
    <ChatSidebar :conversations="conversations" :current-room-id="currentRoomId"
      @select-conversation="handleSelectConversation" @new-conversation="handleNewConversation" />

    <!-- 中间聊天区域 -->
    <ChatInterface v-if="currentRoom" :room="currentRoom" :messages="currentMessages" :seiyuu="currentSeiyuu"
      :use-real-AI="useRealAI" @send-message="handleSendMessage" />

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
              <path
                d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z" />
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
              <path
                d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z" />
            </svg>
            管理AI模型
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧设置面板 -->
    <ChatRightPanel v-if="currentRoom" :room="currentRoom" :seiyuu="currentSeiyuu" :collapsed="rightPanelCollapsed"
      @toggle="toggleRightPanel" @update-settings="handleUpdateSettings" />

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
import { useConfigStore } from '@/stores/configStore'
import { mockService } from '@/services/mockService'
import { aiService } from '@/services/aiService'
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
const configStore = useConfigStore()

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

// 是否使用真实AI服务
const useRealAI = computed(() => {
  const selectedModel = configStore.config.selected_chat_model
  return hasChatModels.value && !!selectedModel && selectedModel.trim() !== ''
})

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

    // 1. 先获取历史消息（在添加当前用户消息之前）
    const conversationHistory = chatStore.getRecentMessages(currentRoom.value.id, 10)

    // 2. 立即添加用户消息到存储（用户可以立即看到）
    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: 'user-1',
      sender_name: '用户',
      content
    })

    let aiReply: string

    // 判断是否使用真实AI服务
    const selectedModel = configStore.config.selected_chat_model

    console.log('🔍 AI服务检查:', {
      hasChatModels: hasChatModels.value,
      selectedModel: selectedModel,
      useRealAI: useRealAI.value
    })

    // 3. 异步生成AI回复
    if (useRealAI.value) {
      try {
        console.log('🤖 使用真实AI服务生成回复...')
        // 创建AI请求Promise
        const aiPromise = aiService.generateReply({
          message: content,
          seiyuu_profile: currentSeiyuu.value.profile_markdown || '',
          conversation_history: conversationHistory,
          model_id: selectedModel
        })

        // 创建超时Promise（30秒超时）
        const timeoutPromise = new Promise<never>((_, reject) => {
          setTimeout(() => reject(new Error('AI请求超时')), 30000)
        })

        // 使用Promise.race实现超时控制
        aiReply = await Promise.race([aiPromise, timeoutPromise])
      } catch (error) {
        console.warn('AI服务调用失败，使用mock回复:', error)
        // AI服务失败时降级到mock服务
        aiReply = await mockService.generateAIReply(
          currentSeiyuu.value.name,
          content,
          conversationHistory
        )
      }
    } else {
      console.log('🎭 使用演示模式生成回复...')
      // 使用mock服务
      aiReply = await mockService.generateAIReply(
        currentSeiyuu.value.name,
        content,
        conversationHistory
      )
    }

    // 4. 添加AI回复消息
    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: currentSeiyuu.value.id,
      sender_name: currentSeiyuu.value.name,
      sender_avatar: currentSeiyuu.value.avatar_url,
      content: aiReply
    })

  } catch (error) {
    console.error('发送消息失败:', error)

    // 添加友好的错误提示消息
    let errorMessage = '抱歉，我现在有点忙，稍后再回复你吧～'

    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: currentSeiyuu.value?.id || 'system',
      sender_name: currentSeiyuu.value?.name || '系统',
      sender_avatar: currentSeiyuu.value?.avatar_url,
      content: errorMessage
    })
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
    console.log('🚀 开始初始化应用数据...')

    // 1. 首先尝试从本地存储加载声优数据
    const hasStoredSeiyuu = seiyuuStore.loadFromStorage()

    // 2. 如果本地没有声优数据，从mock服务加载
    if (!hasStoredSeiyuu && mockService.isMockMode()) {
      console.log('🎭 从mock服务加载声优数据...')
      const response = await mockService.getSeiyuuList()
      if (response.success && response.data) {
        seiyuuStore.setSeiyuuList(response.data) // 这会自动保存到本地存储
      }
    }

    // 3. 从本地存储加载聊天历史数据
    await chatStore.loadFromStorage()

    // 4. 验证和修复数据完整性
    const validation = chatStore.validateDataIntegrity()
    if (!validation.isValid) {
      console.warn('⚠️ 发现数据完整性问题:', validation.issues)
      await chatStore.repairDataIntegrity()
    } else {
      console.log('✅ 数据完整性检查通过')
    }

    // 5. 数据完整性检查和初始化
    if (mockService.isMockMode()) {
      // 如果没有历史对话，创建初始对话
      if (chatStore.rooms.length === 0) {
        console.log('🎭 创建初始演示对话...')
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
      } else {
        console.log(`✅ 已加载 ${chatStore.rooms.length} 个历史对话`)

        // 验证历史对话的声优数据完整性
        const invalidRooms = chatStore.rooms.filter(room => {
          if (room.type === '1v1' && room.participants.length === 1) {
            const seiyuu = seiyuuStore.getSeiyuuById(room.participants[0])
            return !seiyuu
          }
          return false
        })

        if (invalidRooms.length > 0) {
          console.warn(`⚠️ 发现 ${invalidRooms.length} 个对话缺少声优数据，需要重新加载声优信息`)
          // 如果发现数据不完整，重新加载声优数据
          const response = await mockService.getSeiyuuList()
          if (response.success && response.data) {
            seiyuuStore.setSeiyuuList(response.data)
          }
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
  --gradient-primary: linear-gradient(135deg, #fead00 0%, #ff791b 100%);

  /* 主色 */
  --color-primary: #fead00;
  --color-primary-dark: #e6950d;
  --color-primary-light: rgba(254, 173, 0, 0.1);

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
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  font-size: var(--font-size-md);
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
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