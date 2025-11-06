<template>
  <div class="chat-sidebar">
    <!-- 侧边栏头部 -->
    <div class="sidebar-header">
      <h1 class="app-title">NijiChat</h1>
      <p class="app-subtitle">与你的虚拟声优对话</p>
      <!-- 动画波浪背景 -->
      <div class="wave-background"></div>
    </div>

    <!-- 搜索和新建 -->
    <div class="sidebar-controls">
      <div class="search-box">
        <svg class="search-icon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"/>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索对话..."
          class="search-input"
        />
      </div>
      <button @click="handleNewConversation" class="new-conversation-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
        </svg>
      </button>
      <!-- 管理员入口 (开发模式) -->
      <button @click="goToAdmin" class="admin-btn" v-if="isDevelopmentMode" title="管理后台">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492zM5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0z"/>
          <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52l-.094-.319zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115l.094-.319z"/>
        </svg>
      </button>
    </div>

    <!-- 对话列表 -->
    <div class="conversations-container">
      <!-- 1v1 对话分组 -->
      <div v-if="filteredConversations.length > 0" class="conversation-section">
        <div class="section-header">
          <h3 class="section-title">对话</h3>
          <span class="section-count">{{ filteredConversations.length }}</span>
        </div>

        <div class="conversation-list">
          <div
            v-for="conversation in filteredConversations"
            :key="conversation.id"
            :class="[
              'conversation-item',
              { active: conversation.id === currentRoomId }
            ]"
            @click="handleSelectConversation(conversation.id)"
          >
            <div class="conversation-avatar">
              <img
                v-if="conversation.seiyuu?.avatar"
                :src="conversation.seiyuu.avatar"
                :alt="conversation.seiyuu.name"
                class="avatar"
              />
              <div v-else class="avatar-placeholder">
                {{ conversation.seiyuu?.name?.charAt(0) || '?' }}
              </div>
            </div>

            <div class="conversation-content">
              <div class="conversation-header">
                <h4 class="conversation-name">{{ conversation.seiyuu?.name || '未知声优' }}</h4>
                <span class="conversation-time">{{ formatTime(conversation.timestamp) }}</span>
              </div>
              <div class="conversation-preview">
                <p class="last-message">{{ conversation.lastMessage }}</p>
                <span v-if="conversation.unread > 0" class="conversation-badge">
                  {{ conversation.unread > 99 ? '99+' : conversation.unread }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 双声优剧场分组 -->
      <div v-if="dualTheaters.length > 0" class="conversation-section">
        <div class="section-header">
          <h3 class="section-title">双声优剧场</h3>
          <span class="section-count">{{ dualTheaters.length }}</span>
        </div>

        <div class="conversation-list">
          <div
            v-for="theater in dualTheaters"
            :key="theater.id"
            class="conversation-item theater-item"
            @click="handleSelectConversation(theater.id)"
          >
            <div class="theater-avatars">
              <img :src="theater.seiyuu1.avatar" :alt="theater.seiyuu1.name" class="avatar avatar-small" />
              <img :src="theater.seiyuu2.avatar" :alt="theater.seiyuu2.name" class="avatar avatar-small avatar-overlap" />
            </div>

            <div class="conversation-content">
              <div class="conversation-header">
                <h4 class="conversation-name">{{ theater.seiyuu1.name }} & {{ theater.seiyuu2.name }}</h4>
                <span :class="['theater-status', `status-${theater.status}`]">
                  {{ getStatusText(theater.status) }}
                </span>
              </div>
              <div class="conversation-preview">
                <p class="last-message">{{ theater.lastMessage }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 群组剧场分组 -->
      <div v-if="groupTheaters.length > 0" class="conversation-section">
        <div class="section-header">
          <h3 class="section-title">群组剧场</h3>
          <span class="section-count">{{ groupTheaters.length }}</span>
        </div>

        <div class="conversation-list">
          <div
            v-for="group in groupTheaters"
            :key="group.id"
            class="conversation-item group-item"
            @click="handleSelectConversation(group.id)"
          >
            <div class="group-avatar">
              <div class="avatar-placeholder group-placeholder">
                {{ group.name.charAt(0) }}
              </div>
            </div>

            <div class="conversation-content">
              <div class="conversation-header">
                <h4 class="conversation-name">{{ group.name }}</h4>
                <span class="participant-count">{{ group.participants.length }}人</span>
              </div>
              <div class="conversation-preview">
                <p class="last-message">{{ group.lastMessage }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredConversations.length === 0 && !searchQuery" class="empty-conversations">
        <div class="empty-icon">💭</div>
        <p class="empty-text">还没有对话呢</p>
        <button @click="handleNewConversation" class="btn btn-primary">开始聊天</button>
      </div>

      <!-- 搜索无结果 -->
      <div v-if="filteredConversations.length === 0 && searchQuery" class="empty-conversations">
        <div class="empty-icon">🔍</div>
        <p class="empty-text">没有找到匹配的对话</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Conversation, DualTheater, GroupTheater, TheaterStatus } from '@/types'

// Props
const props = defineProps<{
  conversations: Conversation[]
  currentRoomId: string | null
}>()

// Emits
const emit = defineEmits<{
  selectConversation: [roomId: string]
  newConversation: []
}>()

// 路由
const router = useRouter()

// 状态
const searchQuery = ref('')

// 开发模式检测
const isDevelopmentMode = computed(() => {
  return import.meta.env.DEV || import.meta.env.MODE === 'development'
})

// 模拟的双声优剧场数据
const dualTheaters = ref<DualTheater[]>([
  {
    id: 'dual-1',
    seiyuu1: {
      id: 'seiyuu-1',
      name: '花泽香菜',
      avatar: 'https://via.placeholder.com/32/FFB6C1/FFFFFF?text=花泽'
    },
    seiyuu2: {
      id: 'seiyuu-2',
      name: '钉宫理惠',
      avatar: 'https://via.placeholder.com/32/FFB6C1/FFFFFF?text=钉宫'
    },
    topic: '萝莉对决',
    status: 'active',
    lastMessage: '哼！我可不会输给你的！',
    timestamp: new Date(Date.now() - 1000 * 60 * 60).toISOString()
  }
])

// 模拟的群组剧场数据
const groupTheaters = ref<GroupTheater[]>([
  {
    id: 'group-1',
    group_id: 'group-1',
    name: '治愈系声优群',
    participants: [
      { id: 'seiyuu-3', name: '水树奈奈', avatar: 'https://via.placeholder.com/32/87CEEB/FFFFFF?text=水树' },
      { id: 'seiyuu-5', name: '堀江由衣', avatar: 'https://via.placeholder.com/32/98FB98/FFFFFF?text=堀江' },
      { id: 'seiyuu-7', name: '茅野爱衣', avatar: 'https://via.placeholder.com/32/E6E6FA/FFFFFF?text=茅野' }
    ],
    status: 'active',
    lastMessage: '大家一起唱首歌吧～',
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 3).toISOString()
  }
])

// 计算属性
const filteredConversations = computed(() => {
  if (!searchQuery.value) {
    return props.conversations
  }

  const query = searchQuery.value.toLowerCase()
  return props.conversations.filter(conv =>
    conv.seiyuu?.name.toLowerCase().includes(query) ||
    conv.lastMessage.toLowerCase().includes(query)
  )
})

// 方法
function handleSelectConversation(roomId: string) {
  emit('selectConversation', roomId)
}

function handleNewConversation() {
  emit('newConversation')
}

function goToAdmin() {
  router.push({ name: 'Admin' })
}

function formatTime(timestamp: string): string {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 1000 * 60) {
    return '刚刚'
  } else if (diff < 1000 * 60 * 60) {
    return `${Math.floor(diff / (1000 * 60))}分钟前`
  } else if (diff < 1000 * 60 * 60 * 24) {
    return `${Math.floor(diff / (1000 * 60 * 60))}小时前`
  } else if (diff < 1000 * 60 * 60 * 24 * 7) {
    return `${Math.floor(diff / (1000 * 60 * 60 * 24))}天前`
  } else {
    return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
  }
}

function getStatusText(status: TheaterStatus): string {
  const statusMap = {
    active: '进行中',
    completed: '已完成',
    paused: '已暂停'
  }
  return statusMap[status] || '未知'
}
</script>

<style scoped>
.chat-sidebar {
  width: 380px;
  height: 100vh;
  background: var(--bg-primary);
  border-right: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 侧边栏头部 */
.sidebar-header {
  padding: 20px;
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;
  overflow: hidden;
}

.wave-background {
  position: absolute;
  top: 0;
  left: -100px;
  width: 200%;
  height: 100%;
  background: url("data:image/svg+xml,%3Csvg width='100' height='100' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M0 50c25 0 25-25 50-25s25 25 50 25 25-25 50-25 25 25 50 25v50H0V50z' fill='%23ffffff' opacity='0.1'/%3E%3C/svg%3E") repeat-x;
  animation: wave 10s linear infinite;
}

@keyframes wave {
  0% { transform: translateX(0); }
  100% { transform: translateX(100px); }
}

.app-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 4px 0;
  position: relative;
  z-index: 1;
}

.app-subtitle {
  font-size: 13px;
  opacity: 0.9;
  margin: 0;
  position: relative;
  z-index: 1;
}

/* 搜索和控制区 */
.sidebar-controls {
  padding: 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  gap: 8px;
  align-items: center;
}

.search-box {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: #9ca3af;
  z-index: 1;
}

.search-input {
  width: 100%;
  padding: 10px 12px 10px 36px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: #f9fafb;
  transition: all 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.new-conversation-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #667eea;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.new-conversation-btn:hover {
  background: #5568d3;
  transform: translateY(-1px);
}

/* 管理员按钮 */
.admin-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #f59e0b;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.admin-btn:hover {
  background: #d97706;
  transform: translateY(-1px);
}

/* 对话容器 */
.conversations-container {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

/* 分组标题 */
.conversation-section {
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #4b5563;
  margin: 0;
}

.section-count {
  font-size: 11px;
  color: #9ca3af;
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

/* 对话列表 */
.conversation-list {
  padding: 0 8px;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 0 20px 20px 0;
  margin-bottom: 2px;
  margin-right: 8px;
}

.conversation-item:hover {
  background: rgba(102, 126, 234, 0.05);
}

.conversation-item.active {
  background: linear-gradient(90deg,
    rgba(102, 126, 234, 0.1) 0%,
    rgba(102, 126, 234, 0.05) 100%);
  border-right: 3px solid #667eea;
}

/* 头像 */
.conversation-avatar {
  margin-right: 12px;
  flex-shrink: 0;
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.avatar-placeholder {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
}

/* 小尺寸头像（双人剧场） */
.avatar-small {
  width: 28px;
  height: 28px;
}

.avatar-overlap {
  margin-left: -12px;
  border-color: white;
}

.theater-avatars {
  display: flex;
  margin-right: 12px;
  flex-shrink: 0;
}

/* 群组头像 */
.group-avatar {
  margin-right: 12px;
  flex-shrink: 0;
}

.group-placeholder {
  background: linear-gradient(135deg, #10b981, #059669);
}

/* 对话内容 */
.conversation-content {
  flex: 1;
  min-width: 0;
}

.conversation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.conversation-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  margin-right: 8px;
}

.conversation-time {
  font-size: 11px;
  color: #9ca3af;
  white-space: nowrap;
}

.conversation-preview {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.last-message {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  margin-right: 8px;
}

/* 徽章 */
.conversation-badge {
  background: #667eea;
  color: white;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

/* 状态标签 */
.theater-status {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
}

.status-active {
  background: #dcfce7;
  color: #166534;
}

.status-completed {
  background: #f3f4f6;
  color: #6b7280;
}

.status-paused {
  background: #fef3c7;
  color: #92400e;
}

.participant-count {
  font-size: 11px;
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  padding: 2px 6px;
  border-radius: 10px;
}

/* 空状态 */
.empty-conversations {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-text {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 16px 0;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-size: 13px;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* 响应式 */
@media (max-width: 768px) {
  .chat-sidebar {
    width: 100%;
    height: 200px;
    border-right: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  }

  .sidebar-header {
    padding: 15px;
  }

  .app-title {
    font-size: 20px;
  }
}
</style>