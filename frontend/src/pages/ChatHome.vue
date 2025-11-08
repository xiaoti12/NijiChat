<template>
  <div class="chat-page">
    <!-- 左侧对话列表 -->
    <ChatSidebar :conversations="conversations" :seiyuu-groups="seiyuuGroups" :dual-groups="dualGroups"
      :current-room-id="currentRoomId" :current-seiyuu-id="currentSeiyuuId" :collapsed="leftSidebarCollapsed"
      @select-conversation="handleSelectConversation" @select-seiyuu="handleSelectSeiyuu"
      @select-dual-conversation="handleSelectDualConversation" @new-conversation="handleNewConversation"
      @toggle="toggleLeftSidebar" />

    <!-- 中间聊天区域 -->
    <ChatInterface v-if="currentRoom" ref="chatInterfaceRef" :room="currentRoom" :messages="currentMessages"
      :seiyuu="currentSeiyuu" :use-real-AI="useRealAI" @send-message="handleSendMessage"
      @toggle-settings="toggleRightPanel" />

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
    <ChatRightPanel v-if="currentRoom" :room="currentRoom" :seiyuu="currentSeiyuu" :seiyuu-group="currentSeiyuuGroup"
      :collapsed="rightPanelCollapsed" @toggle="toggleRightPanel" @update-settings="handleUpdateSettings" />

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
import { getSeiyuuList } from '@/services/apiService'
import { aiService } from '@/services/aiService'
import type { Seiyuu } from '@/types'

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
const leftSidebarCollapsed = ref(false)
const rightPanelCollapsed = ref(false)
const loading = ref(false)
const showAIManager = ref(false)
const chatInterfaceRef = ref()
const currentSeiyuuId = ref<string | null>(null) // 当前选中的声优ID

// 计算属性
const conversations = computed(() => chatStore.conversations)
const currentRoomId = computed(() => chatStore.currentRoomId)
const currentRoom = computed(() => chatStore.currentRoom)
const currentMessages = computed(() => chatStore.currentMessages)

const currentSeiyuu = computed((): Seiyuu | null => {
  if (!currentRoom.value) return null
  // 1v1 对话获取对应的声优
  if (currentRoom.value.type === '1v1' && currentRoom.value.participants.length === 1) {
    return seiyuuStore.getSeiyuuById(currentRoom.value.participants[0]) || null
  }
  return null
})

// 声优分组列表
const seiyuuGroups = computed(() => chatStore.groupedConversations)

// 双人对话分组列表
const dualGroups = computed(() => chatStore.dualGroupedConversations)

// 当前声优的分组信息
const currentSeiyuuGroup = computed(() => {
  if (!currentSeiyuuId.value) return null
  return chatStore.getSeiyuuGroup(currentSeiyuuId.value)
})

// AI配置状态
const hasChatModels = computed(() => aiModelStore.chatModels.length > 0)

// 是否使用真实AI服务
const useRealAI = computed(() => {
  const selectedModel = configStore.config.selected_chat_model
  return hasChatModels.value && !!selectedModel && selectedModel.trim() !== ''
})

// 事件处理
function handleSelectConversation(roomId: string) {
  chatStore.setCurrentRoom(roomId)
  // 更新当前选中的声优ID
  const room = chatStore.getRoom(roomId)
  if (room && room.type === '1v1' && room.participants.length > 0) {
    currentSeiyuuId.value = room.participants[0]
  }
}

function handleSelectSeiyuu(seiyuuId: string) {
  // 选中声优时，设置当前声优ID
  currentSeiyuuId.value = seiyuuId

  // 获取该声优的所有会话
  const sessions = chatStore.getSessionsBySeiyuu(seiyuuId)

  // 如果有会话，切换到最新的一个
  if (sessions.length > 0) {
    chatStore.setCurrentRoom(sessions[0].id)
    console.log(`🔄 切换到声优的现有会话: ${sessions[0].session_name}`)
  } else {
    // 如果没有会话，创建一个新的
    const seiyuu = seiyuuStore.getSeiyuuById(seiyuuId)
    if (seiyuu) {
      const newRoom = chatStore.createNewSession(seiyuu.id, seiyuu.name, seiyuu.avatar_url)
      chatStore.setCurrentRoom(newRoom.id)
      console.log(`✨ 为声优创建第一个会话: ${newRoom.session_name}`)
    }
  }

  // 确保右侧面板展开，显示会话列表
  rightPanelCollapsed.value = false
}

function handleSelectDualConversation(roomId: string) {
  chatStore.setCurrentRoom(roomId)
  // 清空当前声优ID，因为这是双人对话
  currentSeiyuuId.value = null
}

function handleNewConversation() {
  router.push({ name: 'SeiyuuLibrary' })
}

async function handleSendMessage(content: string) {
  if (!currentRoom.value) return

  // 判断房间类型
  if (currentRoom.value.type === '1v1') {
    // 1v1对话需要当前声优
    if (!currentSeiyuu.value) return
    await handleSendMessage1v1(content)
  } else if (currentRoom.value.type === 'dual_theater') {
    // 双人对话
    await handleSendMessageDual(content)
  }
}

async function handleSendMessage1v1(content: string) {
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
        console.error('AI服务调用失败:', error)
        throw error // 重新抛出错误，让外层catch处理
      }
    } else {
      throw new Error('请先配置AI模型才能开始聊天')
    }

    // 4. 添加AI回复消息
    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: currentSeiyuu.value.id,
      sender_name: currentSeiyuu.value.name,
      sender_avatar: currentSeiyuu.value.avatar_url,
      content: aiReply
    })

    // AI回复完成后停止加载状态
    if (chatInterfaceRef.value) {
      chatInterfaceRef.value.stopTyping()
    }

  } catch (error) {
    console.error('发送消息失败:', error)

    // 根据错误类型提供不同的友好提示
    let errorMessage = '抱歉，我现在有点忙，稍后再回复你吧～'

    if (error instanceof Error) {
      if (error.message.includes('AI请求超时')) {
        errorMessage = '哎呀，思考太久了，请稍后再试吧～'
      } else if (error.message.includes('请先配置AI模型')) {
        errorMessage = '请先配置AI模型才能开始聊天哦～'
      } else if (error.message.includes('网络')) {
        errorMessage = '网络好像有点问题，请检查一下连接～'
      }
    }

    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: currentSeiyuu.value?.id || 'system',
      sender_name: currentSeiyuu.value?.name || '系统',
      sender_avatar: currentSeiyuu.value?.avatar_url,
      content: errorMessage
    })

    // 错误情况下也需要停止加载状态
    if (chatInterfaceRef.value) {
      chatInterfaceRef.value.stopTyping()
    }
  } finally {
    loading.value = false
  }
}

async function handleSendMessageDual(content: string) {
  if (!currentRoom.value || currentRoom.value.type !== 'dual_theater') return

  try {
    loading.value = true

    // 双人对话不需要用户输入，直接生成下一句对话
    const conversationHistory = chatStore.getRecentMessages(currentRoom.value.id, 15)
    const selectedModel = configStore.config.selected_chat_model

    if (!useRealAI.value) {
      throw new Error('请先配置AI模型才能开始双人对话')
    }

    // 获取双人对话的声优信息
    const dualInfo = chatStore.getCurrentDualSeiyuu()
    if (!dualInfo) {
      throw new Error('无法获取双人对话信息')
    }

    const { initiator, responder } = dualInfo

    // 确定下一个发言者（简单轮流）
    let nextSpeaker = initiator
    let otherSpeaker = responder

    if (conversationHistory.length > 0) {
      const lastMessage = conversationHistory[conversationHistory.length - 1]
      // 如果最后一条消息是发起者说的，下一个应该是响应者
      if (lastMessage.sender_id === initiator.id) {
        nextSpeaker = responder
        otherSpeaker = initiator
      }
    }

    console.log('🎭 双人对话AI回复生成:', { nextSpeaker, otherSpeaker })

    // 生成双人对话AI回复
    const aiReply = await aiService.generateDualReply({
      responder_profile: nextSpeaker.profile_markdown || '',
      initiator_profile: otherSpeaker.profile_markdown || '',
      relationship_description: currentRoom.value.dual_relationship || '两人是朋友关系',
      conversation_history: conversationHistory,
      current_topic: currentRoom.value.dual_topic,
      model_id: selectedModel
    })

    // 添加AI回复消息
    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: nextSpeaker.id,
      sender_name: nextSpeaker.name,
      sender_avatar: nextSpeaker.avatar,
      content: aiReply
    })

    // AI回复完成后停止加载状态
    if (chatInterfaceRef.value) {
      chatInterfaceRef.value.stopTyping()
    }

  } catch (error) {
    console.error('双人对话生成失败:', error)

    let errorMessage = '双人对话遇到了一些问题...'
    if (error instanceof Error) {
      if (error.message.includes('请先配置AI模型')) {
        errorMessage = '请先配置AI模型才能开始双人对话～'
      }
    }

    chatStore.addMessage({
      room_id: currentRoom.value.id,
      sender_id: 'system',
      sender_name: '系统',
      content: errorMessage
    })

    if (chatInterfaceRef.value) {
      chatInterfaceRef.value.stopTyping()
    }
  } finally {
    loading.value = false
  }
}

function toggleLeftSidebar() {
  leftSidebarCollapsed.value = !leftSidebarCollapsed.value
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

    // 2. 如果本地没有声优数据，从后端API加载
    if (!hasStoredSeiyuu) {
      console.log('🌐 从后端API加载声优数据...')
      try {
        const response = await getSeiyuuList()
        if (response.success && response.data) {
          seiyuuStore.setSeiyuuList(response.data) // 这会自动保存到本地存储
        } else {
          console.warn('⚠️ 从后端加载声优数据失败:', response.message)
        }
      } catch (error) {
        console.error('❌ 从后端加载声优数据出错:', error)
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

    // 5. 数据迁移：为现有房间分配会话ID
    chatStore.migrateExistingRoomsToSessions()


    // 处理路由参数
    const seiyuuId = route.query.seiyuuId as string
    if (seiyuuId) {
      const seiyuu = seiyuuStore.getSeiyuuById(seiyuuId)
      if (seiyuu) {
        // 检查该声优是否已有会话
        const existingSessions = chatStore.getSessionsBySeiyuu(seiyuu.id)

        if (existingSessions.length > 0) {
          // 如果有现有会话，切换到最新的一个
          chatStore.switchToSession(existingSessions[0].id)
          currentSeiyuuId.value = seiyuu.id
          console.log(`✅ 切换到声优 ${seiyuu.name} 的现有会话`)
        } else {
          // 如果没有现有会话，创建新会话
          const newRoom = chatStore.createNewSession(
            seiyuu.id,
            seiyuu.name,
            seiyuu.avatar_url
          )
          chatStore.setCurrentRoom(newRoom.id)
          currentSeiyuuId.value = seiyuu.id
          console.log(`✅ 为声优 ${seiyuu.name} 创建第一个会话`)
        }
      }
    } else if (chatStore.rooms.length > 0) {
      // 没有路由参数但有现有房间，恢复到最后一次的会话状态
      const lastRoom = chatStore.rooms.reduce((prev, current) =>
        current.created_at > prev.created_at ? current : prev
      )

      if (lastRoom.type === '1v1' && lastRoom.participants.length > 0) {
        chatStore.setCurrentRoom(lastRoom.id)
        currentSeiyuuId.value = lastRoom.participants[0]
        console.log(`✅ 恢复到上次会话: ${lastRoom.session_name || '默认对话'}`)
      }
    } else {
      // 没有任何房间数据，保持空状态
      console.log('ℹ️ 暂无对话数据，等待用户开始聊天')
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