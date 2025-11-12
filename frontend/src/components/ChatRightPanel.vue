<template>
  <div :class="['chat-right-panel', { collapsed }]">
    <!-- 面板头部 -->
    <div class="panel-header">
      <h3 class="panel-title">
        {{ getPanelTitle() }}
      </h3>
      <button @click="handleToggle" class="toggle-btn" :title="collapsed ? '展开面板' : '收起面板'">
        <svg :class="{ rotated: collapsed }" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z" />
        </svg>
      </button>
    </div>

    <!-- 面板内容 -->
    <div v-if="!collapsed" class="panel-content">
      <!-- 角色信息 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">{{ isDualConversation ? '双人剧场信息' : '角色信息' }}</h4>
        </div>

        <!-- 单人对话角色信息 -->
        <div v-if="!isDualConversation" class="seiyuu-profile">
          <div class="profile-avatar">
            <img v-if="currentSeiyuuAvatar" :src="currentSeiyuuAvatar" :alt="seiyuu.name" class="avatar" />
            <div v-else class="avatar-placeholder">
              {{ seiyuu?.name?.charAt(0) || '?' }}
            </div>
          </div>

          <div class="profile-info">
            <h5 class="profile-name">{{ seiyuu?.name || '未知声优' }}</h5>
            <div class="profile-tags">
              <span v-for="tag in seiyuu?.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
          </div>
        </div>

        <!-- 双人对话信息 -->
        <div v-else-if="dualSeiyuu" class="dual-profile">
          <div class="dual-info-item">
            <div class="dual-label">发起者</div>
            <div class="seiyuu-profile compact">
              <div class="profile-avatar">
                <img v-if="dualSeiyuu.initiator.avatar" :src="dualSeiyuu.initiator.avatar"
                     :alt="dualSeiyuu.initiator.name" class="avatar" />
                <div v-else class="avatar-placeholder">
                  {{ dualSeiyuu.initiator.name?.charAt(0) || '?' }}
                </div>
              </div>
              <div class="profile-info">
                <h5 class="profile-name">{{ dualSeiyuu.initiator.name }}</h5>
              </div>
            </div>
          </div>

          <div class="dual-info-item">
            <div class="dual-label">响应者</div>
            <div class="seiyuu-profile compact">
              <div class="profile-avatar">
                <img v-if="dualSeiyuu.responder.avatar" :src="dualSeiyuu.responder.avatar"
                     :alt="dualSeiyuu.responder.name" class="avatar" />
                <div v-else class="avatar-placeholder">
                  {{ dualSeiyuu.responder.name?.charAt(0) || '?' }}
                </div>
              </div>
              <div class="profile-info">
                <h5 class="profile-name">{{ dualSeiyuu.responder.name }}</h5>
              </div>
            </div>
          </div>

          <div class="dual-topic">
            <div class="dual-label">当前话题</div>
            <p class="topic-text">{{ room.dual_topic || '自由聊天' }}</p>
          </div>

        </div>
      </div>

      <!-- AI模型配置 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">AI模型配置</h4>
          <button @click="showAIManager = true" class="manage-btn" title="管理AI模型">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z" />
            </svg>
          </button>
        </div>

        <!-- 未配置AI模型时的提示 -->
        <div v-if="!hasAnyModels" class="no-models-tip">
          <div class="tip-icon">🤖</div>
          <p class="tip-text">还没有配置AI模型</p>
          <button @click="showAIManager = true" class="btn btn-primary btn-sm">
            立即配置
          </button>
        </div>

        <!-- 已配置AI模型时的设置 -->
        <div v-else class="model-settings">
          <!-- 模型选择（始终显示） -->
          <div class="setting-item model-selection">
            <label class="setting-label">
              <span>当前模型</span>
              <span class="setting-value">{{ selectedModel?.name || '未选择' }}</span>
            </label>
            <select v-model="selectedModelId" class="setting-select" @change="handleModelChange">
              <option value="">请选择AI模型</option>
              <option v-for="model in availableChatModels" :key="model.id" :value="model.id">
                {{ model.name }} ({{ model.model_name }})
              </option>
            </select>
          </div>

          <!-- 高级设置折叠触发器 -->
          <div class="advanced-settings-toggle" @click="handleToggleAdvancedSettings">
            <div class="toggle-content">
              <span class="toggle-text">
                {{ isAdvancedSettingsExpanded ? '收起高级设置' : '展开高级设置' }}
              </span>
              <svg class="toggle-icon" :class="{ expanded: isAdvancedSettingsExpanded }" width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
                <path fill-rule="evenodd" d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"/>
              </svg>
            </div>
          </div>

          <!-- 高级参数设置（可折叠） -->
          <div class="advanced-settings" :class="{ expanded: isAdvancedSettingsExpanded }">
            <div class="settings-list">
              <!-- 温度参数 -->
              <div class="setting-item">
                <label class="setting-label">
                  <span>创造性</span>
                  <span class="setting-value">{{ settings.temperature }}</span>
                </label>
                <input v-model.number="settings.temperature" type="range" min="0" max="1" step="0.1"
                  class="setting-slider" />
                <div class="slider-labels">
                  <span>保守</span>
                  <span>创新</span>
                </div>
              </div>

              <!-- 最大Token数 -->
              <div class="setting-item">
                <label class="setting-label">
                  <span>回复长度</span>
                  <span class="setting-value">{{ settings.maxTokens }}</span>
                </label>
                <input v-model.number="settings.maxTokens" type="range" min="50" max="1000" step="50"
                  class="setting-slider" />
                <div class="slider-labels">
                  <span>简短</span>
                  <span>详细</span>
                </div>
              </div>

              <!-- 上下文长度 -->
              <div class="setting-item">
                <label class="setting-label">
                  <span>记忆长度</span>
                  <span class="setting-value">{{ settings.contextLength }}条</span>
                </label>
                <input v-model.number="settings.contextLength" type="range" min="5" max="50" step="5"
                  class="setting-slider" />
                <div class="slider-labels">
                  <span>短期</span>
                  <span>长期</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 应用设置按钮 -->
          <button @click="handleApplySettings" class="apply-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M10.97 4.97a.235.235 0 0 0-.02.022L7.477 9.417 5.384 7.323a.75.75 0 0 0-1.06 1.061L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05z" />
            </svg>
            应用设置
          </button>
        </div>
      </div>

      <!-- 快捷操作 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">快捷操作</h4>
        </div>

        <!-- 1v1对话快捷操作 -->
        <div v-if="!isDualConversation" class="quick-actions">
          <button @click="handleNewSession" class="quick-action-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
            </svg>
            开启新对话
          </button>

          <button @click="handleRenameSession" class="quick-action-btn" v-if="sessionList.length > 0">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708l-3-3zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207l6.5-6.5zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.499.499 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11l.178-.178z" />
            </svg>
            重命名会话
          </button>

          <button @click="handleDeleteSession" class="quick-action-btn danger-btn" v-if="sessionList.length > 0">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z" />
              <path fill-rule="evenodd"
                d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z" />
            </svg>
            删除当前会话
          </button>
        </div>

        <!-- 双人对话快捷操作 -->
        <div v-else class="quick-actions">
          <button @click="handleNewDualSession" class="quick-action-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z" />
            </svg>
            开启新对话
          </button>

          <button @click="handleRenameDualSession" class="quick-action-btn" v-if="dualSessionList.length > 0">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708l-3-3zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207l6.5-6.5zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.499.499 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11l.178-.178z" />
            </svg>
            重命名对话
          </button>

          <button @click="handleDeleteDualSession" class="quick-action-btn danger-btn" v-if="dualSessionList.length > 0">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z" />
              <path fill-rule="evenodd"
                d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z" />
            </svg>
            删除当前对话
          </button>
        </div>
      </div>

      <!-- 对话历史 -->
      <div class="panel-section" v-if="(!isDualConversation && props.seiyuuGroup) || (isDualConversation && dualSeiyuu)">
        <div class="section-header">
          <h4 class="section-title">会话列表</h4>
          <span class="section-count">{{ isDualConversation ? dualSessionList.length : sessionList.length }}</span>
        </div>

        <!-- 1v1对话会话列表 -->
        <div v-if="!isDualConversation" class="history-list">
          <div v-for="session in sessionList" :key="session.id" :class="['history-item', { active: session.active }]"
            @click="handleSelectSession(session)">
            <div class="history-time">{{ session.time }}</div>
            <div class="history-session-name">{{ session.name }}</div>
            <div class="history-preview">{{ session.preview }}</div>
          </div>

          <div v-if="sessionList.length === 0" class="empty-history">
            <div class="empty-icon">📝</div>
            <p>暂无其他会话</p>
          </div>
        </div>

        <!-- 双人对话会话列表 -->
        <div v-else class="history-list">
          <div v-for="session in dualSessionList" :key="session.id" :class="['history-item', { active: session.active }]"
            @click="handleSelectDualSession(session)">
            <div class="history-time">{{ session.time }}</div>
            <div class="history-session-name">{{ session.name }}</div>
            <div class="history-topic" v-if="session.topic">{{ session.topic }}</div>
            <div class="history-preview">{{ session.preview }}</div>
          </div>

          <div v-if="dualSessionList.length === 0" class="empty-history">
            <div class="empty-icon">🎭</div>
            <p>暂无其他对话</p>
          </div>
        </div>
      </div>
    </div>

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

    <!-- 新建双人对话弹窗 -->
    <div v-if="showNewDualSessionDialog" class="dual-session-overlay" @click="handleCancelDualSession">
      <div class="dual-session-container" @click.stop>
        <div class="dual-session-header">
          <h3 class="dual-session-title">开启新对话</h3>
          <button @click="handleCancelDualSession" class="dual-session-close">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z"/>
            </svg>
          </button>
        </div>

        <div class="dual-session-content">
          <div v-if="dualSeiyuu" class="dual-session-seiyuu">
            <div class="seiyuu-info">
              <img v-if="dualSeiyuu.seiyuu1.avatar" :src="dualSeiyuu.seiyuu1.avatar" :alt="dualSeiyuu.seiyuu1.name" class="seiyuu-avatar" />
              <div v-else class="seiyuu-avatar-placeholder">{{ dualSeiyuu.seiyuu1.name.charAt(0) }}</div>
              <span class="seiyuu-name">{{ dualSeiyuu.seiyuu1.name }}</span>
            </div>
            <div class="seiyuu-connector">×</div>
            <div class="seiyuu-info">
              <img v-if="dualSeiyuu.seiyuu2.avatar" :src="dualSeiyuu.seiyuu2.avatar" :alt="dualSeiyuu.seiyuu2.name" class="seiyuu-avatar" />
              <div v-else class="seiyuu-avatar-placeholder">{{ dualSeiyuu.seiyuu2.name.charAt(0) }}</div>
              <span class="seiyuu-name">{{ dualSeiyuu.seiyuu2.name }}</span>
            </div>
          </div>

          <div class="topic-input-section">
            <label class="topic-label">对话主题</label>
            <input
              v-model="newDualTopic"
              type="text"
              class="topic-input"
              placeholder="例如：聊聊最近的工作、讨论一部动漫作品等..."
              @keyup.enter="handleCreateDualSession"
            />
            <p class="topic-hint">设置一个有趣的话题，让两位声优围绕这个主题开始对话</p>
          </div>
        </div>

        <div class="dual-session-footer">
          <button @click="handleCancelDualSession" class="btn btn-secondary">
            取消
          </button>
          <button @click="handleCreateDualSession" class="btn btn-primary" :disabled="!newDualTopic.trim()">
            创建对话
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useAIModelStore } from '@/stores/aiModelStore'
import { useConfigStore } from '@/stores/configStore'
import { useChatStore } from '@/stores/chatStore'
import { useSeiyuuStore } from '@/stores/seiyuuStore'
import type { Room, Seiyuu, SeiyuuConversationGroup } from '@/types'
import AIModelManager from './AIModelManager.vue'

// Props
const props = defineProps<{
  room: Room
  seiyuu: Seiyuu | null
  seiyuuGroup: SeiyuuConversationGroup | null
  collapsed: boolean
}>()

// Emits
const emit = defineEmits<{
  toggle: []
  updateSettings: [settings: any]
}>()

// Store
const aiModelStore = useAIModelStore()
const configStore = useConfigStore()
const chatStore = useChatStore()
const seiyuuStore = useSeiyuuStore()

// 状态
const showAIManager = ref(false)
const isAdvancedSettingsExpanded = ref(false)
const selectedModelId = ref<string | null>(configStore.config.selected_chat_model || null)
const showNewDualSessionDialog = ref(false)
const newDualTopic = ref('')

// 设置
const settings = reactive({
  model_id: selectedModelId.value,
  temperature: 0.7,
  maxTokens: 200,
  contextLength: 10
})

// 计算属性
const selectedModel = computed(() => {
  return selectedModelId.value ? aiModelStore.getModel(selectedModelId.value) : null
})

// 动态获取声优头像
const currentSeiyuuAvatar = computed(() => {
  if (!props.seiyuu?.id) return props.seiyuu?.avatar_url
  const latestSeiyuu = seiyuuStore.getSeiyuuById(props.seiyuu.id)
  return latestSeiyuu?.avatar_url || props.seiyuu.avatar_url
})

const availableChatModels = computed(() => {
  return aiModelStore.chatModels
})

const hasAnyModels = computed(() => {
  return aiModelStore.models.length > 0
})

// 判断是否为双人对话
const isDualConversation = computed(() => {
  return props.room.type === 'dual_theater'
})

// 获取双人对话的声优信息
const dualSeiyuu = computed(() => {
  if (!isDualConversation.value) return null
  return chatStore.getCurrentDualSeiyuu()
})

// 历史记录 - 获取当前声优的所有会话
const sessionList = computed(() => {
  if (!props.seiyuu) return []

  return chatStore.getSessionsBySeiyuu(props.seiyuu.id).map(room => ({
    id: room.id,
    name: chatStore.getSessionDisplayName(room),
    time: new Date(room.created_at).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    }),
    preview: room.last_message?.content || '暂无消息',
    active: props.room.id === room.id,
    session_name: room.session_name
  }))
})

// 双人对话会话列表
const dualSessionList = computed(() => {
  if (!isDualConversation.value || !dualSeiyuu.value) return []

  const { seiyuu1, seiyuu2 } = dualSeiyuu.value
  const sessions = chatStore.getDualSessions(seiyuu1.id, seiyuu2.id)

  return sessions.map(room => ({
    id: room.id,
    name: room.session_name || room.dual_topic || '双人对话',
    topic: room.dual_topic,
    time: new Date(room.created_at).toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    }),
    preview: room.last_message?.content || '暂无消息',
    active: props.room.id === room.id
  }))
})

// 监听selectedModelId变化，自动更新settings
watch(selectedModelId, (newValue) => {
  settings.model_id = newValue
})

// 方法
function getPanelTitle() {
  if (isDualConversation.value && dualSeiyuu.value) {
    return '双人剧场设置'
  }
  return props.seiyuuGroup?.seiyuuName ? `${props.seiyuuGroup.seiyuuName} - 会话管理` : '聊天设置'
}

function handleToggle() {
  emit('toggle')
}

function handleToggleAdvancedSettings() {
  isAdvancedSettingsExpanded.value = !isAdvancedSettingsExpanded.value
}

function handleModelChange() {
  // 模型选择变化时自动保存配置
  if (selectedModelId.value) {
    configStore.setSelectedChatModel(selectedModelId.value)
    console.log('✅ 已选择AI模型:', selectedModelId.value)
  } else {
    // 取消选择时清除配置
    configStore.updateConfig({ selected_chat_model: null })
    console.log('⚠️ 未选择AI模型，已清除配置')
  }
}

function handleApplySettings() {
  // 更新模型ID到设置对象
  settings.model_id = selectedModelId.value

  emit('updateSettings', { ...settings })
  console.log('设置已应用:', settings)
}

function handleNewSession() {
  if (!props.seiyuu) return

  if (confirm('确定要开启新对话吗？将创建一个全新的对话会话。')) {
    const newRoom = chatStore.createNewSession(
      props.seiyuu.id,
      props.seiyuu.name,
      props.seiyuu.avatar_url
    )

    // 切换到新会话
    chatStore.switchToSession(newRoom.id)
    console.log('✅ 已创建新对话会话')
  }
}

function handleRenameSession() {
  if (!props.room) return

  const currentName = props.room.session_name || chatStore.getSessionDisplayName(props.room)
  const newName = prompt('请输入新的会话名称:', currentName)

  if (newName && newName.trim() && newName !== currentName) {
    chatStore.updateSessionName(props.room.id, newName.trim())
    console.log('✅ 会话已重命名为:', newName)
  }
}

function handleDeleteSession() {
  if (!props.room) return

  const confirmDelete = confirm(`确定要删除会话"${props.room.session_name || '当前对话'}"吗？这将清除该会话的所有聊天记录。`)

  if (confirmDelete) {
    const seiyuuId = props.seiyuu?.id
    chatStore.deleteSession(props.room.id)

    // 删除后，如果该声优还有其他会话，切换到最新的一个
    if (seiyuuId) {
      const sessions = chatStore.getSessionsBySeiyuu(seiyuuId)
      if (sessions.length > 0) {
        chatStore.switchToSession(sessions[0].id)
      }
    }

    console.log('✅ 会话已删除')
  }
}

function handleSelectSession(session: any) {
  // 切换到选中的会话
  chatStore.switchToSession(session.id)
  console.log('切换到会话:', session.name)
}

// === 双人对话相关方法 ===

function handleNewDualSession() {
  if (!dualSeiyuu.value) return

  showNewDualSessionDialog.value = true
  newDualTopic.value = ''
}

function handleCreateDualSession() {
  if (!dualSeiyuu.value) return

  const { seiyuu1, seiyuu2 } = dualSeiyuu.value
  const topic = newDualTopic.value.trim()

  if (!topic) {
    alert('请输入对话主题')
    return
  }

  // 创建新的双人对话会话
  const newRoom = chatStore.createDualSession(
    { id: seiyuu1.id, name: seiyuu1.name, avatar: seiyuu1.avatar },
    { id: seiyuu2.id, name: seiyuu2.name, avatar: seiyuu2.avatar },
    {
      topic: topic,
      initiatorId: seiyuu1.id
    }
  )

  // 切换到新会话
  chatStore.switchToSession(newRoom.id)

  // 关闭对话框
  showNewDualSessionDialog.value = false
  newDualTopic.value = ''

  console.log('✅ 已创建新双人对话会话:', topic)
}

function handleCancelDualSession() {
  showNewDualSessionDialog.value = false
  newDualTopic.value = ''
}

function handleSelectDualSession(session: any) {
  // 切换到选中的双人对话会话
  chatStore.switchToSession(session.id)
  console.log('切换到双人对话:', session.name)
}

function handleRenameDualSession() {
  if (!props.room) return

  const currentName = props.room.session_name || props.room.dual_topic || '双人对话'
  const newName = prompt('请输入新的对话名称:', currentName)

  if (newName && newName.trim() && newName !== currentName) {
    chatStore.updateSessionName(props.room.id, newName.trim())
    console.log('✅ 双人对话已重命名为:', newName)
  }
}

function handleDeleteDualSession() {
  if (!props.room || !dualSeiyuu.value) return

  const sessionName = props.room.session_name || props.room.dual_topic || '双人对话'
  const confirmDelete = confirm(`确定要删除对话"${sessionName}"吗？这将清除该对话的所有聊天记录。`)

  if (confirmDelete) {
    const { seiyuu1, seiyuu2 } = dualSeiyuu.value
    chatStore.deleteSession(props.room.id)

    // 删除后，如果该声优组合还有其他会话，切换到最新的一个
    const sessions = chatStore.getDualSessions(seiyuu1.id, seiyuu2.id)
    if (sessions.length > 0) {
      chatStore.switchToSession(sessions[0].id)
    }

    console.log('✅ 双人对话已删除')
  }
}
</script>

<style scoped>
.chat-right-panel {
  width: 350px;
  background: var(--bg-primary);
  border-left: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  height: 100vh;
  transition: all 0.3s ease;
  overflow: hidden;
}

.chat-right-panel.collapsed {
  width: 0;
  border-left: none;
}

/* 面板头部 */
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  background: #fafafa;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.toggle-btn {
  width: 28px;
  height: 28px;
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

.toggle-btn:hover {
  background: #e5e7eb;
}

.toggle-btn svg.rotated {
  transform: rotate(180deg);
}

/* 面板内容 */
.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
}

/* 面板分组 */
.panel-section {
  margin-bottom: 24px;
}

.panel-section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
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

/* 声优信息 */
.seiyuu-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.profile-avatar .avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.profile-avatar .avatar-placeholder {
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

.profile-info {
  flex: 1;
  min-width: 0;
}

.profile-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 6px 0;
}

.profile-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag {
  font-size: 10px;
  padding: 2px 6px;
  background: rgba(254, 173, 0, 0.1);
  color: #fead00;
  border-radius: 4px;
  font-weight: 500;
}

/* 双人对话信息样式 */
.dual-profile {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dual-info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dual-label {
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.seiyuu-profile.compact {
  padding: 8px 12px;
  gap: 8px;
}

.seiyuu-profile.compact .profile-avatar .avatar,
.seiyuu-profile.compact .profile-avatar .avatar-placeholder {
  width: 32px;
  height: 32px;
  font-size: 14px;
}

.seiyuu-profile.compact .profile-name {
  font-size: 13px;
  margin: 0;
}

.dual-topic {
  padding: 8px 12px;
  background: #f9fafb;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.topic-text {
  font-size: 13px;
  color: #374151;
  margin: 0;
  line-height: 1.4;
}

/* 设置列表 */
.settings-list {
  margin-bottom: 16px;
}

/* 模型设置容器 */
.model-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 高级设置折叠触发器 */
.advanced-settings-toggle {
  cursor: pointer;
  padding: 8px 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: all 0.2s;
}

.advanced-settings-toggle:hover {
  border-color: #fead00;
  background: rgba(254, 173, 0, 0.05);
}

.toggle-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.toggle-text {
  font-size: 13px;
  font-weight: 500;
  color: #4b5563;
}

.toggle-icon {
  transition: transform 0.3s ease;
  color: #6b7280;
}

.toggle-icon.expanded {
  transform: rotate(180deg);
}

/* 高级设置折叠容器 */
.advanced-settings {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease-out, opacity 0.2s ease-out;
  opacity: 0;
}

.advanced-settings.expanded {
  max-height: 400px;
  opacity: 1;
  transition: max-height 0.3s ease-in, opacity 0.2s ease-in;
}

/* 高级设置内部列表 */
.advanced-settings .settings-list {
  margin-bottom: 0;
}

.advanced-settings .setting-item {
  margin-bottom: 12px;
}

.advanced-settings .setting-item:last-child {
  margin-bottom: 0;
}

/* 模型选择特殊样式 */
.model-selection {
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 16px;
  margin-bottom: 8px;
}

.model-selection .setting-item {
  margin-bottom: 0;
}

.setting-item {
  margin-bottom: 16px;
}

.setting-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  font-weight: 500;
  color: #4b5563;
  margin-bottom: 8px;
}

.setting-value {
  color: #fead00;
  font-weight: 600;
}

.setting-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.setting-select:hover {
  border-color: #fead00;
}

.setting-select:focus {
  outline: none;
  border-color: #fead00;
  box-shadow: 0 0 0 3px rgba(254, 173, 0, 0.1);
}

.setting-slider {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #e5e7eb;
  outline: none;
  -webkit-appearance: none;
  cursor: pointer;
}

.setting-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #fead00;
  cursor: pointer;
  transition: all 0.2s;
}

.setting-slider::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 0 0 4px rgba(254, 173, 0, 0.2);
}

.slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #9ca3af;
  margin-top: 4px;
}

/* 应用按钮 */
.apply-btn {
  width: 100%;
  padding: 10px;
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 8px;
}

.apply-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
}

/* 快捷操作 */
.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  color: #4b5563;
  cursor: pointer;
  transition: all 0.2s;
  width: 100%;
}

.quick-action-btn:hover {
  border-color: #fead00;
  background: rgba(254, 173, 0, 0.05);
  color: #fead00;
}

.quick-action-btn.danger-btn:hover {
  border-color: #ef4444;
  background: rgba(239, 68, 68, 0.05);
  color: #ef4444;
}

/* 历史记录 */
.history-list {
  overflow-y: visible;
}

.history-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: all 0.2s;
  cursor: pointer;
  margin-bottom: 8px;
}

.history-item:hover {
  border-color: #fead00;
  background: rgba(254, 173, 0, 0.05);
}

.history-item.active {
  border-color: #fead00;
  background: rgba(254, 173, 0, 0.1);
}

.history-time {
  font-size: 11px;
  color: #9ca3af;
  font-weight: 500;
}

.history-session-name {
  font-size: 12px;
  color: #374151;
  font-weight: 600;
  margin: 2px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-preview {
  font-size: 13px;
  color: #4b5563;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-history {
  text-align: center;
  padding: 20px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 8px;
  opacity: 0.5;
}

.empty-history p {
  font-size: 13px;
  margin: 0;
}

/* 滚动条样式 */
.panel-content::-webkit-scrollbar {
  width: 4px;
}

.panel-content::-webkit-scrollbar-track {
  background: transparent;
}

.panel-content::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

.panel-content::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

/* 管理按钮 */
.manage-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  transition: all 0.2s;
}

.manage-btn:hover {
  background: #f3f4f6;
  color: #fead00;
}

/* 无模型提示 */
.no-models-tip {
  text-align: center;
  padding: 20px 16px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.no-models-tip .btn {
  margin: 0 auto;
}

.tip-icon {
  font-size: 24px;
  margin-bottom: 8px;
  opacity: 0.6;
}

.tip-text {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
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

/* 双人对话特有样式 */
.history-topic {
  font-size: 11px;
  color: #fead00;
  font-weight: 500;
  margin: 2px 0;
  padding: 2px 6px;
  background: rgba(254, 173, 0, 0.1);
  border-radius: 4px;
  display: inline-block;
}

/* 新建双人对话弹窗 */
.dual-session-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.dual-session-container {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.dual-session-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.dual-session-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.dual-session-close {
  width: 32px;
  height: 32px;
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

.dual-session-close:hover {
  background: #f3f4f6;
  color: #374151;
}

.dual-session-content {
  padding: 24px;
  flex: 1;
}

.dual-session-seiyuu {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.seiyuu-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.seiyuu-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.seiyuu-avatar-placeholder {
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

.seiyuu-name {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.seiyuu-connector {
  font-size: 18px;
  font-weight: 600;
  color: #9ca3af;
}

.topic-input-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.topic-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.topic-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  color: #374151;
  transition: all 0.2s;
}

.topic-input:focus {
  outline: none;
  border-color: #fead00;
  box-shadow: 0 0 0 3px rgba(254, 173, 0, 0.1);
}

.topic-hint {
  font-size: 12px;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

.dual-session-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  background: #fafafa;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 80px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
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

/* 响应式 */
@media (max-width: 1400px) {
  .chat-right-panel {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
  }
}

@media (max-width: 768px) {
  /* 移动端：面板固定定位，从右侧滑入 */
  .chat-right-panel {
    position: fixed;
    right: 0;
    top: 0;
    bottom: 0;
    width: 85vw;
    max-width: 350px;
    z-index: 999;
    box-shadow: -4px 0 16px rgba(0, 0, 0, 0.2);
    transform: translateX(0);
    transition: transform 0.3s ease-in-out;
    animation: slideInRight 0.3s ease-in-out;
  }

  @keyframes slideInRight {
    from {
      transform: translateX(100%);
    }
    to {
      transform: translateX(0);
    }
  }

  /* 收起状态：滑出屏幕 */
  .chat-right-panel.collapsed {
    transform: translateX(100%);
  }

  .ai-manager-container,
  .dual-session-container {
    max-width: 95%;
    max-height: 90vh;
  }

  .dual-session-content {
    padding: 20px;
  }

  .dual-session-seiyuu {
    flex-direction: column;
    gap: 12px;
  }

  .seiyuu-connector {
    transform: rotate(90deg);
  }
}
</style>