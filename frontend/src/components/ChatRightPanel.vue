<template>
  <div :class="['chat-right-panel', { collapsed }]">
    <!-- 面板头部 -->
    <div class="panel-header">
      <h3 class="panel-title">聊天设置</h3>
      <button @click="handleToggle" class="toggle-btn" :title="collapsed ? '展开面板' : '收起面板'">
        <svg :class="{ rotated: collapsed }" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"/>
        </svg>
      </button>
    </div>

    <!-- 面板内容 -->
    <div v-if="!collapsed" class="panel-content">
      <!-- 声优信息 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">角色信息</h4>
        </div>

        <div class="seiyuu-profile">
          <div class="profile-avatar">
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

          <div class="profile-info">
            <h5 class="profile-name">{{ seiyuu?.name || '未知声优' }}</h5>
            <div class="profile-tags">
              <span v-for="tag in seiyuu?.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- AI模型配置 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">AI模型配置</h4>
          <button @click="showAIManager = true" class="manage-btn" title="管理AI模型">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
              <path d="M9.405 1.05c-.413-1.4-2.397-1.4-2.81 0l-.1.34a1.464 1.464 0 0 1-2.105.872l-.31-.17c-1.283-.698-2.686.705-1.987 1.987l.169.311c.446.82.023 1.841-.872 2.105l-.34.1c-1.4.413-1.4 2.397 0 2.81l.34.1a1.464 1.464 0 0 1 .872 2.105l-.17.31c-.698 1.283.705 2.686 1.987 1.987l.311-.169a1.464 1.464 0 0 1 2.105.872l.1.34c.413 1.4 2.397 1.4 2.81 0l.1-.34a1.464 1.464 0 0 1 2.105-.872l.31.17c1.283.698 2.686-.705 1.987-1.987l-.169-.311a1.464 1.464 0 0 1 .872-2.105l.34-.1c1.4-.413 1.4-2.397 0-2.81l-.34-.1a1.464 1.464 0 0 1-.872-2.105l.17-.31c.698-1.283-.705-2.686-1.987-1.987l-.311.169a1.464 1.464 0 0 1-2.105-.872l-.1-.34zM8 10.93a2.929 2.929 0 1 1 0-5.86 2.929 2.929 0 0 1 0 5.858z"/>
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
        <div v-else class="settings-list">
          <!-- 模型选择 -->
          <div class="setting-item">
            <label class="setting-label">
              <span>当前模型</span>
              <span class="setting-value">{{ selectedModel?.name || '未选择' }}</span>
            </label>
            <select v-model="selectedModelId" class="setting-select">
              <option value="">请选择AI模型</option>
              <option
                v-for="model in availableChatModels"
                :key="model.id"
                :value="model.id"
              >
                {{ model.name }} ({{ model.model_name }})
              </option>
            </select>
          </div>

          <!-- 温度参数 -->
          <div class="setting-item">
            <label class="setting-label">
              <span>创造性</span>
              <span class="setting-value">{{ settings.temperature }}</span>
            </label>
            <input
              v-model.number="settings.temperature"
              type="range"
              min="0"
              max="1"
              step="0.1"
              class="setting-slider"
            />
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
            <input
              v-model.number="settings.maxTokens"
              type="range"
              min="50"
              max="1000"
              step="50"
              class="setting-slider"
            />
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
            <input
              v-model.number="settings.contextLength"
              type="range"
              min="5"
              max="50"
              step="5"
              class="setting-slider"
            />
            <div class="slider-labels">
              <span>短期</span>
              <span>长期</span>
            </div>
          </div>
        </div>

        <!-- 应用设置按钮 -->
        <button @click="handleApplySettings" class="apply-btn">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M10.97 4.97a.235.235 0 0 0-.02.022L7.477 9.417 5.384 7.323a.75.75 0 0 0-1.06 1.061L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-1.071-1.05z"/>
          </svg>
          应用设置
        </button>
      </div>

      <!-- 快捷操作 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">快捷操作</h4>
        </div>

        <div class="quick-actions">
          <button @click="handleResetChat" class="quick-action-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M11.534 7h3.932a.25.25 0 0 1 .192.41l-1.966 2.36a.25.25 0 0 1-.384 0l-1.966-2.36a.25.25 0 0 1 .192-.41zm-11 2h3.932a.25.25 0 0 0 .192-.41L2.692 6.23a.25.25 0 0 0-.384 0L.342 8.59A.25.25 0 0 0 .534 9z"/>
              <path fill-rule="evenodd" d="M8 3c-1.552 0-2.94.707-3.857 1.818a.5.5 0 1 1-.771-.636A6.002 6.002 0 0 1 13.917 7H12.9A5.002 5.002 0 0 0 8 3zM3.1 9a5.002 5.002 0 0 0 8.757 2.182.5.5 0 1 1 .771.636A6.002 6.002 0 0 1 2.083 9H3.1z"/>
            </svg>
            重置对话
          </button>

          <button @click="handleExportChat" class="quick-action-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"/>
              <path d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708l3-3z"/>
            </svg>
            导出聊天记录
          </button>

          <button @click="handleChangeModel" class="quick-action-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M1 2.828c.885-.37 2.154-.769 3.388-.893 1.33-.134 2.458.063 3.112.752v9.746c-.935-.53-2.12-.603-3.213-.493-1.18.12-2.37.461-3.287.811V2.828zm7.5-.141c.654-.689 1.782-.886 3.112-.752 1.234.124 2.503.523 3.388.893v9.923c-.918-.35-2.107-.692-3.287-.81-1.094-.111-2.278-.039-3.213.492V2.687zM8 1.783C7.015.936 5.587.81 4.287.94c-1.514.153-3.042.672-3.994 1.105A.5.5 0 0 0 0 2.5v11a.5.5 0 0 0 .707.455c.882-.4 2.303-.881 3.68-1.02 1.409-.142 2.59.087 3.223.877a.5.5 0 0 0 .78 0c.633-.79 1.814-1.019 3.222-.877 1.378.139 2.8.62 3.681 1.02A.5.5 0 0 0 16 13.5v-11a.5.5 0 0 0-.293-.455c-.952-.433-2.48-.952-3.994-1.105C10.413.809 8.985.936 8 1.783z"/>
            </svg>
            切换AI模型
          </button>
        </div>
      </div>

      <!-- 对话历史 -->
      <div class="panel-section">
        <div class="section-header">
          <h4 class="section-title">对话历史</h4>
          <span class="section-count">{{ historyItems.length }}</span>
        </div>

        <div class="history-list">
          <div
            v-for="item in historyItems"
            :key="item.id"
            :class="['history-item', { active: item.active }]"
            @click="handleSelectHistory(item)"
          >
            <div class="history-time">{{ item.time }}</div>
            <div class="history-preview">{{ item.preview }}</div>
          </div>

          <div v-if="historyItems.length === 0" class="empty-history">
            <div class="empty-icon">📝</div>
            <p>暂无历史记录</p>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useAIModelStore } from '@/stores/aiModelStore'
import type { Room, Seiyuu } from '@/types'
import AIModelManager from './AIModelManager.vue'

// Props
const props = defineProps<{
  room: Room
  seiyuu: Seiyuu | null
  collapsed: boolean
}>()

// Emits
const emit = defineEmits<{
  toggle: []
  updateSettings: [settings: any]
}>()

// Store
const aiModelStore = useAIModelStore()

// 状态
const showAIManager = ref(false)
const selectedModelId = ref<string | null>(null)

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

const availableChatModels = computed(() => {
  return aiModelStore.chatModels
})

const hasAnyModels = computed(() => {
  return aiModelStore.models.length > 0
})

// 历史记录
const historyItems = ref([
  {
    id: 1,
    time: '14:30',
    preview: '今天天气真好呢～',
    active: false
  },
  {
    id: 2,
    time: '14:25',
    preview: '你好，初次见面！',
    active: true
  }
])

// 方法
function handleToggle() {
  emit('toggle')
}

function handleApplySettings() {
  emit('updateSettings', { ...settings })
  // 显示应用成功提示
  console.log('设置已应用:', settings)
}

function handleResetChat() {
  if (confirm('确定要重置当前对话吗？这将清除所有聊天记录。')) {
    // 实现重置逻辑
    console.log('重置对话')
  }
}

function handleExportChat() {
  // 实现导出逻辑
  console.log('导出聊天记录')
}

function handleChangeModel() {
  // 实现模型切换逻辑
  console.log('切换AI模型')
}

function handleSelectHistory(item: any) {
  // 切换历史记录激活状态
  historyItems.value.forEach(h => h.active = false)
  item.active = true

  // 实现历史记录选择逻辑
  console.log('选择历史记录:', item)
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

/* 设置列表 */
.settings-list {
  margin-bottom: 16px;
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

/* 历史记录 */
.history-list {
  max-height: 200px;
  overflow-y: auto;
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
.panel-content::-webkit-scrollbar,
.history-list::-webkit-scrollbar {
  width: 4px;
}

.panel-content::-webkit-scrollbar-track,
.history-list::-webkit-scrollbar-track {
  background: transparent;
}

.panel-content::-webkit-scrollbar-thumb,
.history-list::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

.panel-content::-webkit-scrollbar-thumb:hover,
.history-list::-webkit-scrollbar-thumb:hover {
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
  .chat-right-panel {
    display: none;
  }

  .ai-manager-container {
    max-width: 95%;
    max-height: 90vh;
  }
}
</style>