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
      <p class="empty-text">选择一个对话开始聊天，或创建新的对话 </p>
      <button @click="handleNewConversation" class="btn btn-primary btn-large">
        开始新对话
      </button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useChatStore } from '@/stores/chatStore'
import { useSeiyuuStore } from '@/stores/seiyuuStore'
import { mockService } from '@/services/mockService'
import type { Room, Message, Seiyuu } from '@/types'

// 导入组件（稍后创建）
import ChatSidebar from '@/components/ChatSidebar.vue'
import ChatInterface from '@/components/ChatInterface.vue'
import ChatRightPanel from '@/components/ChatRightPanel.vue'

const router = useRouter()
const route = useRoute()
const chatStore = useChatStore()
const seiyuuStore = useSeiyuuStore()

// 状态
const rightPanelCollapsed = ref(false)
const loading = ref(false)

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

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-page {
    flex-direction: column;
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