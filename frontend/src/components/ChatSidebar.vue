<template>
  <div :class="['chat-sidebar', { collapsed: props.collapsed }]">
    <!-- 侧边栏头部 -->
    <div class="sidebar-header">
      <div v-if="!props.collapsed" class="header-content">
        <h1 class="app-title">NijiChat</h1>
        <p class="app-subtitle">与你的虚拟声优对话</p>
      </div>
      <div v-else class="header-collapsed">
        <div class="app-icon">N</div>
      </div>
      <!-- 动画波浪背景 -->
      <div class="wave-background"></div>
      <!-- 收起/展开按钮 -->
      <button class="toggle-btn" @click="handleToggle" :title="props.collapsed ? '展开侧边栏' : '收起侧边栏'">
        <svg :class="{ rotated: props.collapsed }" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z" />
        </svg>
      </button>
    </div>

    <!-- 搜索和新建 -->
    <div v-if="!props.collapsed" class="sidebar-controls">
      <div class="search-box">
        <svg class="search-icon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z" />
        </svg>
        <input v-model="searchQuery" type="text" placeholder="搜索对话..." class="search-input" />
      </div>
      <button @click="handleNewConversation" class="new-conversation-btn">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
        </svg>
      </button>
    </div>
    <!-- 收起状态下的新建按钮 -->
    <div v-else class="collapsed-controls">
      <button @click="handleNewConversation" class="new-conversation-btn-collapsed" title="开始新对话">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
        </svg>
      </button>
    </div>

    <!-- 功能模式选择 -->
    <div v-if="!props.collapsed" class="mode-selection">
      <div class="mode-header">
        <h4 class="mode-title">功能模式</h4>
      </div>
      <div class="mode-buttons">
        <button @click="handleModeSelection('chat')" :class="['mode-btn', { active: currentMode === 'chat' }]"
          title="1v1声优对话">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path
              d="M14 1a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H4.414A2 2 0 0 0 3 11.586l-2 2V2a1 1 0 0 1 1-1h12zM2 0a2 2 0 0 0-2 2v12.793a.5.5 0 0 0 .854.353l2.853-2.853A1 1 0 0 1 4.414 12H14a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2z" />
            <path
              d="M3 3.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5zM3 6a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9A.5.5 0 0 1 3 6zm0 2.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5z" />
          </svg>
          <span class="mode-label">普通对话</span>
        </button>

        <button @click="handleModeSelection('dual')" :class="['mode-btn', { active: currentMode === 'dual' }]"
          title="双人剧场对话">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path
              d="M6 12.5a.5.5 0 0 1 .5-.5h3a.5.5 0 0 1 0 1h-3a.5.5 0 0 1-.5-.5ZM3 8.062C3 6.76 4.235 5.765 5.53 5.886a26.58 26.58 0 0 0 4.94 0C11.765 5.765 13 6.76 13 8.062v1.157a.933.933 0 0 1-.765.935c-.845.147-2.34.346-4.235.346-1.895 0-3.39-.2-4.235-.346A.933.933 0 0 1 3 9.219V8.062Zm4.542-.827a.25.25 0 0 0-.217.068l-.92.9a24.767 24.767 0 0 1-1.871-.183.25.25 0 0 0-.068.495c.55.076 1.232.149 2.02.193a.25.25 0 0 0 .189-.071l.754-.736.847 1.71a.25.25 0 0 0 .404.062l.932-.97a25.286 25.286 0 0 0 1.922-.188.25.25 0 0 0-.068-.495c-.538.074-1.207.145-1.98.189a.25.25 0 0 0-.166.077l-.734.764-.617-1.44Z" />
            <path d="M8.5 4.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0Z" />
          </svg>
          <span class="mode-label">双人剧场</span>
        </button>

        <button @click="handleModeSelection('group')" :class="['mode-btn', { active: currentMode === 'group' }]"
          title="群组剧场对话">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M7 14s-1 0-1-1 1-4 5-4 5 3 5 4-1 1-1 1H7zm4-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
            <path fillRule="evenodd"
              d="M5.216 14A2.238 2.238 0 0 1 5 13c0-1.355.68-2.75 1.936-3.72A6.325 6.325 0 0 0 5 9c-4 0-5 3-5 4s1 1 1 1h4.216z" />
            <path d="M4.5 8a2.5 2.5 0 1 0 0-5 2.5 2.5 0 0 0 0 5z" />
          </svg>
          <span class="mode-label">群组剧场</span>
        </button>
      </div>
    </div>

    <!-- 对话列表 -->
    <div class="conversations-container">
      <!-- 声优分组列表 -->
      <div v-if="filteredSeiyuuGroups.length > 0" class="conversation-section">
        <div v-if="!props.collapsed" class="section-header">
          <h3 class="section-title">声优</h3>
          <span class="section-count">{{ filteredSeiyuuGroups.length }}</span>
        </div>

        <div :class="['conversation-list', { collapsed: props.collapsed }]">
          <div v-for="group in filteredSeiyuuGroups" :key="group.seiyuuId" :class="[
            'conversation-item',
            { active: group.seiyuuId === props.currentSeiyuuId },
            { collapsed: props.collapsed }
          ]" @click="handleSelectSeiyuu(group.seiyuuId)">
            <div class="conversation-avatar">
              <img v-if="group.seiyuuAvatar" :src="group.seiyuuAvatar" :alt="group.seiyuuName" class="avatar" />
              <div v-else class="avatar-placeholder">
                {{ group.seiyuuName.charAt(0) || '?' }}
              </div>
              <!-- 未读徽章 -->
              <span v-if="group.totalUnread > 0" class="conversation-badge">
                {{ group.totalUnread > 99 ? '99+' : group.totalUnread }}
              </span>
              <!-- 会话数量徽章 -->
              <span v-if="group.totalSessions > 1" class="session-count-badge">
                {{ group.totalSessions }}
              </span>
            </div>

            <div v-if="!props.collapsed" class="conversation-content">
              <div class="conversation-header">
                <h4 class="conversation-name">{{ group.seiyuuName }}</h4>
                <span class="conversation-time">{{ formatTime(group.lastActive) }}</span>
              </div>
              <div class="conversation-preview">
                <p class="last-message">{{ group.lastMessage?.content || '暂无消息' }}</p>
              </div>
              <div class="conversation-meta">
                <span class="session-info">{{ group.totalSessions }} 个对话</span>
              </div>
            </div>
            <!-- 收起状态下的悬浮提示 -->
            <div v-else class="conversation-tooltip">
              <div class="tooltip-content">
                <div class="tooltip-name">{{ group.seiyuuName }}</div>
                <div class="tooltip-preview">{{ group.lastMessage?.content || '暂无消息' }}</div>
                <div class="tooltip-meta">{{ group.totalSessions }} 个对话</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 双人对话分组列表 -->
      <div v-if="filteredDualGroups.length > 0" class="conversation-section">
        <div v-if="!props.collapsed" class="section-header">
          <h3 class="section-title">双人对话</h3>
          <span class="section-count">{{ filteredDualGroups.length }}</span>
        </div>

        <div :class="['conversation-list', { collapsed: props.collapsed }]">
          <div v-for="group in filteredDualGroups" :key="group.pairId" :class="[
            'conversation-item',
            'dual-conversation-item',
            { active: group.rooms.some(r => r.id === props.currentRoomId) },
            { collapsed: props.collapsed }
          ]" @click="handleSelectDualConversation(group.rooms[0]?.id)">
            <!-- 双人头像组合 -->
            <div class="dual-conversation-avatar">
              <div class="dual-avatar-container">
                <!-- 声优1头像 -->
                <img v-if="group.seiyuu1.avatar" :src="group.seiyuu1.avatar" :alt="group.seiyuu1.name"
                  class="dual-avatar dual-avatar-1" />
                <div v-else class="dual-avatar-placeholder dual-avatar-1">
                  {{ group.seiyuu1.name.charAt(0) || '?' }}
                </div>

                <!-- 声优2头像 -->
                <img v-if="group.seiyuu2.avatar" :src="group.seiyuu2.avatar" :alt="group.seiyuu2.name"
                  class="dual-avatar dual-avatar-2" />
                <div v-else class="dual-avatar-placeholder dual-avatar-2">
                  {{ group.seiyuu2.name.charAt(0) || '?' }}
                </div>
              </div>

              <!-- 未读徽章 -->
              <span v-if="group.totalUnread > 0" class="conversation-badge">
                {{ group.totalUnread > 99 ? '99+' : group.totalUnread }}
              </span>
              <!-- 会话数量徽章 -->
              <span v-if="group.totalSessions > 1" class="session-count-badge">
                {{ group.totalSessions }}
              </span>
            </div>

            <div v-if="!props.collapsed" class="conversation-content">
              <div class="conversation-header">
                <h4 class="conversation-name">{{ group.seiyuu1.name }} × {{ group.seiyuu2.name }}</h4>
                <span class="conversation-time">{{ formatTime(group.lastActive) }}</span>
              </div>
              <div class="conversation-preview">
                <p class="last-message">{{ group.lastMessage?.content || '暂无消息' }}</p>
              </div>
              <div class="conversation-meta">
                <span class="session-info">{{ group.totalSessions }} 个对话</span>
              </div>
            </div>

            <!-- 收起状态下的悬浮提示 -->
            <div v-else class="conversation-tooltip">
              <div class="tooltip-content">
                <div class="tooltip-name">{{ group.seiyuu1.name }} × {{ group.seiyuu2.name }}</div>
                <div class="tooltip-preview">{{ group.lastMessage?.content || '暂无消息' }}</div>
                <div class="tooltip-meta">{{ group.totalSessions }} 个对话</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredSeiyuuGroups.length === 0 && filteredDualGroups.length === 0 && !searchQuery"
        class="empty-conversations">
        <div class="empty-icon">💭</div>
        <p class="empty-text">还没有对话呢</p>
        <button @click="handleNewConversation" class="btn btn-primary">浏览声优库</button>
      </div>

      <!-- 搜索无结果 -->
      <div v-if="filteredSeiyuuGroups.length === 0 && searchQuery" class="empty-conversations">
        <div class="empty-icon">🔍</div>
        <p class="empty-text">没有找到匹配的声优</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Conversation, SeiyuuConversationGroup, DualConversationGroup, DualTheater, GroupTheater } from '@/types'
import { useSeiyuuStore } from '@/stores/seiyuuStore'

// Props
const props = defineProps<{
  conversations: Conversation[]
  currentRoomId: string | null
  currentSeiyuuId: string | null
  seiyuuGroups: SeiyuuConversationGroup[]
  dualGroups: DualConversationGroup[]
  collapsed?: boolean
}>()

// Emits
const emit = defineEmits<{
  selectConversation: [roomId: string]
  selectSeiyuu: [seiyuuId: string]
  selectDualConversation: [roomId: string]
  newConversation: []
  toggle: []
}>()

// 路由
const router = useRouter()

// Store
const seiyuuStore = useSeiyuuStore()

// 状态
const searchQuery = ref('')
const currentMode = ref('chat') // 当前选中的模式：chat, dual, group


// 双声优剧场数据（从后端API获取）
const dualTheaters = ref<DualTheater[]>([])

// 群组剧场数据（从后端API获取）
const groupTheaters = ref<GroupTheater[]>([])

// 动态获取声优头像
function getSeiyuuLatestAvatar(seiyuuId: string): string | undefined {
  const latestSeiyuu = seiyuuStore.getSeiyuuById(seiyuuId)
  return latestSeiyuu?.avatar_url
}

// 动态更新声优分组的头像
const filteredSeiyuuGroups = computed(() => {
  let groups = props.seiyuuGroups

  // 更新每个分组的头像为最新头像
  groups = groups.map(group => ({
    ...group,
    seiyuuAvatar: getSeiyuuLatestAvatar(group.seiyuuId) || group.seiyuuAvatar
  }))

  if (!searchQuery.value) {
    return groups
  }

  const query = searchQuery.value.toLowerCase()
  return groups.filter(group =>
    group.seiyuuName.toLowerCase().includes(query)
  )
})

// 动态更新双人对话分组的头像
const filteredDualGroups = computed(() => {
  return props.dualGroups.map(group => ({
    ...group,
    seiyuu1: {
      ...group.seiyuu1,
      avatar: getSeiyuuLatestAvatar(group.seiyuu1.id) || group.seiyuu1.avatar
    },
    seiyuu2: {
      ...group.seiyuu2,
      avatar: getSeiyuuLatestAvatar(group.seiyuu2.id) || group.seiyuu2.avatar
    }
  }))
})

// 方法
function handleSelectConversation(roomId: string) {
  emit('selectConversation', roomId)
}

function handleSelectSeiyuu(seiyuuId: string) {
  emit('selectSeiyuu', seiyuuId)
}

function handleSelectDualConversation(roomId: string) {
  emit('selectDualConversation', roomId)
}

function handleNewConversation() {
  emit('newConversation')
}


function handleToggle() {
  emit('toggle')
}

function handleModeSelection(mode: string) {
  currentMode.value = mode

  switch (mode) {
    case 'chat':
      // 普通对话模式，跳转到声优库选择
      handleNewConversation()
      break
    case 'dual':
      // 双声优剧场模式
      router.push({ name: 'DualTheater' })
      break
    case 'group':
      // 群组剧场模式
      router.push({ name: 'GroupTheater' })
      break
  }
}

function formatTime(timestamp: string | number): string {
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
  transition: width 0.3s ease;
}

.chat-sidebar.collapsed {
  width: 80px;
}

/* 侧边栏头部 */
.sidebar-header {
  padding: 20px;
  text-align: center;
  background: linear-gradient(135deg, #fead00 0%, #ff791b 100%);
  color: white;
  position: relative;
  overflow: hidden;
  transition: padding 0.3s ease;
}

.chat-sidebar.collapsed .sidebar-header {
  padding: 15px 8px;
}

.header-content {
  position: relative;
  z-index: 1;
}

.header-collapsed {
  position: relative;
  z-index: 1;
}

.app-icon {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  margin: 0 auto;
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
  0% {
    transform: translateX(0);
  }

  100% {
    transform: translateX(100px);
  }
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

/* Toggle按钮 */
.toggle-btn {
  position: absolute;
  top: 50%;
  right: -16px;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border: 2px solid #e5e7eb;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow:
    0 4px 8px rgba(0, 0, 0, 0.12),
    0 2px 4px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  z-index: 10;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.chat-sidebar.collapsed .toggle-btn {
  right: -12px;
  width: 28px;
  height: 28px;
  top: 20px;
  transform: translateY(0);
}

.toggle-btn:hover {
  background: linear-gradient(135deg, #fead00 0%, #ff791b 100%);
  border-color: #fead00;
  color: white;
  transform: translateY(-50%) scale(1.15);
  box-shadow:
    0 6px 12px rgba(254, 173, 0, 0.25),
    0 3px 6px rgba(254, 173, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.toggle-btn:active {
  transform: translateY(-50%) scale(1.05);
  transition: all 0.1s ease;
}

.toggle-btn svg {
  width: 14px;
  height: 14px;
  transition: all 0.4s cubic-bezier(0.68, -0.55, 0.265, 1.55);
  filter: drop-shadow(0 1px 1px rgba(0, 0, 0, 0.1));
}

.toggle-btn svg.rotated {
  transform: rotate(180deg);
}

.toggle-btn:hover svg {
  filter: drop-shadow(0 2px 2px rgba(0, 0, 0, 0.2));
  transform: translateX(1px);
}

.toggle-btn:hover svg.rotated {
  transform: translateX(1px) rotate(180deg);
}

/* 呼吸灯效果 */
.toggle-btn::before {
  content: '';
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  background: inherit;
  opacity: 0;
  z-index: -1;
  transition: opacity 0.4s ease;
}

.toggle-btn:hover::before {
  opacity: 0.1;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 0.1;
  }

  50% {
    transform: scale(1.1);
    opacity: 0.05;
  }

  100% {
    transform: scale(1);
    opacity: 0.1;
  }
}

/* 涟漪效果 */
.toggle-btn::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(254, 173, 0, 0.3);
  transform: translate(-50%, -50%);
  transition: width 0.6s ease, height 0.6s ease, opacity 0.6s ease;
  opacity: 0;
}

.toggle-btn:active::after {
  width: 100%;
  height: 100%;
  opacity: 0;
  transition: width 0.6s ease, height 0.6s ease, opacity 0.3s ease;
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
  border-color: #fead00;
  background: white;
  box-shadow: 0 0 0 3px rgba(254, 173, 0, 0.1);
}

.new-conversation-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: #fead00;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.new-conversation-btn:hover {
  background: #e6950d;
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

.conversation-list.collapsed {
  padding: 0 4px;
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
  position: relative;
}

.conversation-item.collapsed {
  padding: 8px;
  margin-right: 0;
  border-radius: 12px;
  flex-direction: column;
  justify-content: center;
}

.conversation-item:hover {
  background: rgba(254, 173, 0, 0.05);
}

.conversation-item.active {
  background: linear-gradient(90deg,
      rgba(254, 173, 0, 0.1) 0%,
      rgba(254, 173, 0, 0.05) 100%);
  border-right: 3px solid #fead00;
}

/* 头像 */
.conversation-avatar {
  margin-right: 12px;
  flex-shrink: 0;
  position: relative;
}

.conversation-item.collapsed .conversation-avatar {
  margin-right: 0;
}

/* 双人对话头像 */
.dual-conversation-avatar {
  margin-right: 12px;
  flex-shrink: 0;
  position: relative;
}

.conversation-item.collapsed .dual-conversation-avatar {
  margin-right: 0;
}

.dual-avatar-container {
  position: relative;
  width: 52px;
  height: 44px;
}

.conversation-item.collapsed .dual-avatar-container {
  width: 42px;
  height: 36px;
}

.dual-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.9);
  position: absolute;
}

.conversation-item.collapsed .dual-avatar {
  width: 26px;
  height: 26px;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.dual-avatar-1 {
  left: 0;
  top: 0;
  z-index: 2;
}

.dual-avatar-2 {
  right: 0;
  bottom: 0;
  z-index: 1;
}

.dual-avatar-placeholder {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #fead00, #ff791b);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 12px;
  position: absolute;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.conversation-item.collapsed .dual-avatar-placeholder {
  width: 26px;
  height: 26px;
  font-size: 10px;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.dual-avatar-placeholder.dual-avatar-1 {
  left: 0;
  top: 0;
  z-index: 2;
}

.dual-avatar-placeholder.dual-avatar-2 {
  right: 0;
  bottom: 0;
  z-index: 1;
  background: linear-gradient(135deg, #22c55e, #16a34a);
}

.avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.conversation-item.collapsed .avatar {
  width: 36px;
  height: 36px;
  border: 2px solid rgba(255, 255, 255, 0.9);
}

.avatar-placeholder {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #fead00, #ff791b);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
}

.conversation-item.collapsed .avatar-placeholder {
  width: 36px;
  height: 36px;
  font-size: 14px;
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
  background: #fead00;
  color: white;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
  position: absolute;
  top: -2px;
  right: -2px;
  border: 2px solid white;
}

.conversation-item.collapsed .conversation-badge {
  top: -4px;
  right: -4px;
  font-size: 10px;
  padding: 1px 4px;
  min-width: 16px;
}

/* 会话数量徽章 */
.session-count-badge {
  background: rgba(254, 173, 0, 0.2);
  color: #fead00;
  font-size: 10px;
  font-weight: 600;
  padding: 1px 4px;
  border-radius: 8px;
  border: 1px solid #fead00;
  position: absolute;
  bottom: -2px;
  right: -2px;
  border: 1px solid white;
}

.conversation-item.collapsed .session-count-badge {
  bottom: -4px;
  right: -4px;
  font-size: 9px;
  padding: 1px 3px;
}

/* 对话元信息 */
.conversation-meta {
  display: flex;
  align-items: center;
  margin-top: 4px;
}

.session-info {
  font-size: 11px;
  color: #9ca3af;
  background: rgba(156, 163, 175, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
}

/* 悬浮提示 */
.conversation-tooltip {
  position: absolute;
  left: 100%;
  top: 50%;
  transform: translateY(-50%);
  margin-left: 8px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s ease;
  z-index: 1000;
  min-width: 200px;
  pointer-events: none;
}

.conversation-item:hover .conversation-tooltip {
  opacity: 1;
  visibility: visible;
  transform: translateY(-50%) translateX(4px);
}

.tooltip-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tooltip-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.tooltip-preview {
  font-size: 12px;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.tooltip-meta {
  font-size: 11px;
  color: #9ca3af;
  background: rgba(156, 163, 175, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  align-self: flex-start;
  font-weight: 500;
}

/* 功能模式选择 */
.mode-selection {
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid #e5e7eb;
}

.mode-header {
  margin-bottom: 8px;
}

.mode-title {
  font-size: 12px;
  font-weight: 600;
  color: #4b5563;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.mode-buttons {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mode-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #6b7280;
  font-size: 13px;
  font-weight: 500;
  text-align: left;
  width: 100%;
}

.mode-btn:hover {
  background: rgba(254, 173, 0, 0.05);
  color: #374151;
  border-color: rgba(254, 173, 0, 0.2);
}

.mode-btn.active {
  background: rgba(254, 173, 0, 0.1);
  color: #fead00;
  border-color: rgba(254, 173, 0, 0.3);
  font-weight: 600;
}

.mode-btn.active svg {
  color: #fead00;
}

.mode-btn svg {
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.mode-label {
  white-space: nowrap;
}

/* 收起状态下的控制按钮 */
.collapsed-controls {
  padding: 8px 12px;
  display: flex;
  justify-content: center;
}

.new-conversation-btn-collapsed {
  width: 56px;
  height: 56px;
  background: #fead00;
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.new-conversation-btn-collapsed:hover {
  background: #e6950d;
  transform: scale(1.05);
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
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  font-size: 13px;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
}

/* 响应式 */
@media (max-width: 768px) {
  .chat-sidebar {
    width: 100%;
    height: 200px;
    border-right: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  }

  .chat-sidebar.collapsed {
    height: 60px;
    width: 100%;
  }

  /* 侧边栏头部优化 - 简化显示 */
  .sidebar-header {
    padding: 12px 16px;
  }

  .chat-sidebar.collapsed .sidebar-header {
    padding: 10px 16px;
  }

  .app-title {
    font-size: 20px;
    margin-bottom: 0;
  }

  /* 移动端隐藏副标题 */
  .app-subtitle {
    display: none;
  }

  /* Toggle按钮优化 */
  .toggle-btn {
    width: 32px;
    height: 32px;
    right: -14px;
    top: 50%;
  }

  .chat-sidebar.collapsed .toggle-btn {
    width: 28px;
    height: 28px;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
  }

  .app-icon {
    width: 36px;
    height: 36px;
    font-size: 18px;
  }

  /* 搜索和控制区优化 */
  .sidebar-controls {
    padding: 10px 16px;
  }

  /* 模式选择区域 - 横向滑动布局 */
  .mode-selection {
    padding: 8px 12px;
  }

  .mode-header {
    margin-bottom: 6px;
  }

  .mode-buttons {
    flex-direction: row;
    gap: 6px;
    display: flex;
    justify-content: space-between;
  }

  .mode-btn {
    flex: 1;
    min-width: 0;
    padding: 8px 10px;
    font-size: 13px;
    min-height: 44px;
    /* 移动端最佳触控区域 */
    gap: 6px;
    justify-content: center;
  }

  .mode-btn svg {
    width: 16px;
    height: 16px;
  }

  .mode-label {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .collapsed-controls {
    padding: 8px 16px;
  }

  .new-conversation-btn-collapsed {
    width: 44px;
    height: 44px;
  }
}

@media (max-width: 480px) {

  /* 超小屏幕优化 */
  .sidebar-header {
    padding: 10px 12px;
  }

  .app-title {
    font-size: 18px;
  }

  .toggle-btn {
    width: 30px;
    height: 30px;
    right: -12px;
  }

  .chat-sidebar.collapsed .toggle-btn {
    width: 26px;
    height: 26px;
    right: 10px;
  }

  .app-icon {
    width: 32px;
    height: 32px;
    font-size: 16px;
  }

  /* 搜索控制区优化 */
  .sidebar-controls {
    padding: 8px 12px;
  }

  .search-input {
    font-size: 13px;
    padding: 8px 10px 8px 32px;
  }

  /* 模式选择按钮优化 */
  .mode-selection {
    padding: 6px 10px;
  }

  .mode-btn {
    flex: 1;
    min-width: 0;
    padding: 8px 6px;
    font-size: 12px;
    min-height: 44px;
    gap: 5px;
    justify-content: center;
  }

  .mode-btn svg {
    width: 15px;
    height: 15px;
  }

  .collapsed-controls {
    padding: 6px 12px;
  }

  .new-conversation-btn-collapsed {
    width: 44px;
    height: 44px;
  }
}
</style>