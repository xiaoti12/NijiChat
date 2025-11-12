<template>
  <div class="dual-theater">
    <div class="theater-header">
      <div class="theater-title">
        <h1>双人剧场</h1>
        <p v-if="!selectedSeiyuu1 || !selectedSeiyuu2">选择两位声优开始对话</p>
        <p v-else class="theater-subtitle">
          {{ selectedSeiyuu1.name }} × {{ selectedSeiyuu2.name }}
        </p>
      </div>
      <button @click="goBack" class="btn btn-secondary">
        返回首页
      </button>
    </div>

    <!-- 声优选择阶段 -->
    <div v-if="!theaterStarted" class="setup-section">
      <div class="seiyuu-selection">
        <div class="selection-grid">
          <!-- 声优1选择 -->
          <div class="seiyuu-selector">
            <h3>发起聊天人</h3>
            <div class="selected-seiyuu" v-if="selectedSeiyuu1">
              <img v-if="selectedSeiyuu1.avatar_url" :src="selectedSeiyuu1.avatar_url" :alt="selectedSeiyuu1.name"
                class="seiyuu-avatar" />
              <div class="seiyuu-info">
                <div class="seiyuu-name">{{ selectedSeiyuu1.name }}</div>
                <div class="seiyuu-tags">
                  <span v-for="tag in selectedSeiyuu1.tags" :key="tag" class="tag">
                    {{ tag }}
                  </span>
                </div>
              </div>
              <button @click="selectedSeiyuu1 = null" class="change-btn">
                更换
              </button>
            </div>
            <div v-else class="seiyuu-picker">
              <SeiyuuPicker @select="onSelectSeiyuu1" />
            </div>
          </div>

          <!-- 关系显示 -->
          <div class="relationship-display">
            <div v-if="selectedSeiyuu1 && selectedSeiyuu2" class="relationship-content">
              <div class="relationship-header">
                <span class="relationship-icon">🤝</span>
                <h4>关系背景</h4>
              </div>
              <div v-if="relationship" class="relationship-description">
                {{ relationship.relationship_description }}
              </div>
              <div v-else-if="loadingRelationship" class="relationship-loading">
                正在加载关系信息...
              </div>
              <div v-else class="relationship-none">
                暂无关系记录，将使用默认关系
              </div>
            </div>
            <div v-else class="relationship-placeholder">
              <span class="relationship-icon">↔️</span>
              <p>选择两位声优后显示关系</p>
            </div>
          </div>

          <!-- 声优2选择 -->
          <div class="seiyuu-selector">
            <h3>回复聊天人</h3>
            <div class="selected-seiyuu" v-if="selectedSeiyuu2">
              <img v-if="selectedSeiyuu2.avatar_url" :src="selectedSeiyuu2.avatar_url" :alt="selectedSeiyuu2.name"
                class="seiyuu-avatar" />
              <div class="seiyuu-info">
                <div class="seiyuu-name">{{ selectedSeiyuu2.name }}</div>
                <div class="seiyuu-tags">
                  <span v-for="tag in selectedSeiyuu2.tags" :key="tag" class="tag">
                    {{ tag }}
                  </span>
                </div>
              </div>
              <button @click="selectedSeiyuu2 = null" class="change-btn">
                更换
              </button>
            </div>
            <div v-else class="seiyuu-picker">
              <SeiyuuPicker :excludeIds="selectedSeiyuu1 ? [selectedSeiyuu1.id] : []" @select="onSelectSeiyuu2" />
            </div>
          </div>
        </div>

        <!-- 话题设置 -->
        <div v-if="selectedSeiyuu1 && selectedSeiyuu2" class="topic-section">
          <div class="topic-input">
            <label class="topic-label">对话话题（可选）</label>
            <input v-model="currentTopic" type="text" placeholder="例如：聊聊最近的工作、讨论一部动漫作品等..." class="input topic-field" />
          </div>
          <button @click="startTheater" :disabled="!canStartTheater" class="btn btn-primary start-btn">
            创建剧场并开始对话
          </button>
        </div>
      </div>
    </div>

    <!-- 对话进行阶段已移除 - 对话功能已转移到ChatHome主页面 -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { Seiyuu, SeiyuuRelationship, Message } from '@/types'
import { getSeiyuuList, getRelationship } from '@/services/apiService'
import { useAIService } from '@/services/aiService'
import { useChatStore } from '@/stores/chatStore'
import SeiyuuPicker from '@/components/SeiyuuPicker.vue'

const router = useRouter()
const aiService = useAIService()
const chatStore = useChatStore()

// 响应式数据
const selectedSeiyuu1 = ref<Seiyuu | null>(null)
const selectedSeiyuu2 = ref<Seiyuu | null>(null)
const relationship = ref<SeiyuuRelationship | null>(null)
const loadingRelationship = ref(false)
const currentTopic = ref('')
const theaterStarted = ref(false)
const messages = ref<Message[]>([])
const generating = ref(false)
const messagesContainer = ref<HTMLElement>()
const currentSpeakerIndex = ref(0) // 0: 声优1, 1: 声优2

// 计算属性
const canStartTheater = computed(() => {
  return selectedSeiyuu1.value && selectedSeiyuu2.value && !loadingRelationship.value
})

const nextSpeaker = computed(() => {
  if (!selectedSeiyuu1.value || !selectedSeiyuu2.value) return null
  return currentSpeakerIndex.value === 0 ? selectedSeiyuu1.value : selectedSeiyuu2.value
})

const otherSpeaker = computed(() => {
  if (!selectedSeiyuu1.value || !selectedSeiyuu2.value) return null
  return currentSpeakerIndex.value === 0 ? selectedSeiyuu2.value : selectedSeiyuu1.value
})

// 方法
function goBack() {
  router.push({ name: 'Home' })
}

function onSelectSeiyuu1(seiyuu: Seiyuu) {
  selectedSeiyuu1.value = seiyuu
}

function onSelectSeiyuu2(seiyuu: Seiyuu) {
  selectedSeiyuu2.value = seiyuu
}

async function loadRelationship() {
  if (!selectedSeiyuu1.value || !selectedSeiyuu2.value) return

  loadingRelationship.value = true
  try {
    const response = await getRelationship(
      selectedSeiyuu1.value.id,
      selectedSeiyuu2.value.id
    )

    if (response.success) {
      relationship.value = Array.isArray(response.data) ? response.data[0] : response.data
    }
  } catch (error: any) {
    console.error('加载关系失败:', error.message)
  } finally {
    loadingRelationship.value = false
  }
}

function startTheater() {
  if (!canStartTheater.value || !selectedSeiyuu1.value || !selectedSeiyuu2.value) return

  try {
    console.log('🎭 开始创建双人剧场:', {
      seiyuu1: selectedSeiyuu1.value.name,
      seiyuu2: selectedSeiyuu2.value.name,
      topic: currentTopic.value
    })

    // 创建双人对话房间
    const relationshipDesc = relationship.value?.relationship_description ||
      `${selectedSeiyuu1.value.name}和${selectedSeiyuu2.value.name}是同行，彼此了解但不算特别熟悉的关系。`

    const seiyuu1Data = {
      id: selectedSeiyuu1.value.id,
      name: selectedSeiyuu1.value.name,
      avatar: selectedSeiyuu1.value.avatar_url
    }

    const seiyuu2Data = {
      id: selectedSeiyuu2.value.id,
      name: selectedSeiyuu2.value.name,
      avatar: selectedSeiyuu2.value.avatar_url
    }

    const sessionOptions = {
      topic: currentTopic.value || undefined,
      initiatorId: selectedSeiyuu1.value.id, // 第一个选择的声优作为发起者
      relationship: relationshipDesc
    }

    const newRoom = chatStore.createDualSession(
      seiyuu1Data,
      seiyuu2Data,
      sessionOptions
    )

    // 跳转到主页面，并选中新创建的双人对话房间
    router.push({ name: 'Home' })

    console.log('✅ 双人剧场已创建，跳转到主页面:', {
      roomId: newRoom.id,
      seiyuu1: selectedSeiyuu1.value.name,
      seiyuu2: selectedSeiyuu2.value.name,
      topic: currentTopic.value,
      relationship: relationshipDesc
    })
  } catch (error) {
    console.error('❌ 创建双人剧场失败:', error)
    alert('创建双人剧场失败，请检查控制台错误信息。')
  }
}

// 以下方法已不再使用，对话功能已转移到主页面
async function nextDialogue() {
  // 不再使用 - 对话功能已转移到ChatHome页面
  console.log('nextDialogue方法已废弃，对话功能已转移到ChatHome页面')
}

function resetTheater() {
  // 不再使用 - 对话功能已转移到ChatHome页面
  theaterStarted.value = false
  messages.value = []
  currentSpeakerIndex.value = 0
  currentTopic.value = ''
}

function getSpeakerAvatar(senderId: string): string | undefined {
  if (senderId === selectedSeiyuu1.value?.id) {
    return selectedSeiyuu1.value.avatar_url
  } else if (senderId === selectedSeiyuu2.value?.id) {
    return selectedSeiyuu2.value.avatar_url
  }
  return undefined
}

function formatTime(timestamp: number): string {
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 监听器
watch([selectedSeiyuu1, selectedSeiyuu2], () => {
  if (selectedSeiyuu1.value && selectedSeiyuu2.value) {
    loadRelationship()
  } else {
    relationship.value = null
  }
})

// 生命周期
onMounted(() => {
  // 页面加载完成
})
</script>

<style scoped>
.dual-theater {
  width: 100%;
  min-height: 100vh;
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
}

.theater-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xl) var(--spacing-2xl);
  background: var(--bg-primary);
  border-bottom: 1px solid var(--border-color);
}

.theater-title h1 {
  font-size: var(--font-size-2xl);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.theater-subtitle {
  font-size: var(--font-size-md);
  color: var(--text-muted);
  margin: 0;
}

.setup-section {
  flex: 1;
  padding: var(--spacing-2xl);
}

.selection-grid {
  display: grid;
  grid-template-columns: minmax(300px, 1fr) minmax(380px, 500px) minmax(300px, 1fr);
  gap: var(--spacing-md);
  max-width: 1200px;
  margin: 0 auto;
  align-items: start;
}

.seiyuu-selector {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  min-height: 260px;
  min-width: 280px;
  max-width: 350px;
  margin: 0 auto;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
}

.seiyuu-selector h3 {
  text-align: center;
  margin-bottom: var(--spacing-lg);
  color: var(--text-primary);
}

.selected-seiyuu {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
}

.seiyuu-avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.seiyuu-info {
  text-align: center;
}

.seiyuu-name {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.seiyuu-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  justify-content: center;
}

.tag {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
}

.change-btn {
  padding: var(--spacing-xs) var(--spacing-md);
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  color: var(--text-primary);
}

.change-btn:hover {
  background: var(--bg-tertiary);
}

.relationship-display {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  min-width: 380px;
  min-height: 260px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.relationship-content,
.relationship-placeholder {
  text-align: center;
}

.relationship-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  justify-content: center;
  margin-bottom: var(--spacing-md);
}

.relationship-icon {
  font-size: var(--font-size-2xl);
}

.relationship-header h4 {
  margin: 0;
  color: var(--text-primary);
}

.relationship-description {
  background: var(--bg-secondary);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  border-left: 3px solid var(--color-primary);
  line-height: 1.6;
  color: var(--text-primary);
  height: 180px;
  overflow-y: auto;
  width: 100%;
  margin: 0;
  box-sizing: border-box;
  flex: 1;
}

.relationship-loading,
.relationship-none {
  color: var(--text-muted);
  font-style: italic;
}

.topic-section {
  margin-top: var(--spacing-2xl);
  text-align: center;
}

.topic-input {
  margin-bottom: var(--spacing-lg);
}

.topic-label {
  display: block;
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.topic-field {
  width: 100%;
  max-width: 500px;
}

.start-btn {
  padding: var(--spacing-md) var(--spacing-2xl);
  font-size: var(--font-size-lg);
}

.theater-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: var(--spacing-xl);
}

.theater-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--bg-primary);
  padding: var(--spacing-lg);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
}

.current-topic {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.topic-text {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.control-buttons {
  display: flex;
  gap: var(--spacing-md);
}

.chat-container {
  flex: 1;
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  padding: var(--spacing-xl);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.message-wrapper {
  display: flex;
}

.message-wrapper.message-seiyuu1 {
  justify-content: flex-start;
}

.message-wrapper.message-seiyuu2 {
  justify-content: flex-end;
}

.message-content {
  max-width: 70%;
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  position: relative;
}

.message-seiyuu1 .message-content {
  background: linear-gradient(135deg, #e3f2fd 0%, #f3e5f5 100%);
  border-bottom-left-radius: var(--radius-sm);
}

.message-seiyuu2 .message-content {
  background: linear-gradient(135deg, #f0f4c3 0%, #dcedc8 100%);
  border-bottom-right-radius: var(--radius-sm);
}

.message-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.message-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.message-name {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  font-size: var(--font-size-sm);
}

.message-time {
  color: var(--text-muted);
  font-size: var(--font-size-xs);
  margin-left: auto;
}

.message-text {
  color: var(--text-primary);
  line-height: 1.5;
}

.generating-indicator {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  color: var(--text-muted);
}

.typing-animation {
  display: flex;
  gap: var(--spacing-xs);
}

.typing-dot {
  width: 6px;
  height: 6px;
  background: var(--color-primary);
  border-radius: 50%;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-dot:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {

  0%,
  80%,
  100% {
    transform: scale(0);
  }

  40% {
    transform: scale(1);
  }
}

.generating-text {
  font-style: italic;
}

/* 关系描述滚动条美化 */
.relationship-description::-webkit-scrollbar {
  width: 6px;
}

.relationship-description::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

.relationship-description::-webkit-scrollbar-thumb {
  background: var(--color-primary);
  border-radius: 3px;
  opacity: 0.7;
}

.relationship-description::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary-dark);
  opacity: 1;
}

/* Firefox 滚动条样式 */
.relationship-description {
  scrollbar-width: thin;
  scrollbar-color: var(--color-primary) rgba(0, 0, 0, 0.05);
}

@media (max-width: 1200px) {
  .selection-grid {
    grid-template-columns: 1.5fr minmax(280px, 1fr) 1.5fr;
    gap: var(--spacing-lg);
    max-width: 1200px;
  }

  .seiyuu-selector {
    min-width: 240px;
  }

  .relationship-display {
    min-width: 280px;
  }
}

@media (max-width: 1024px) {
  .selection-grid {
    grid-template-columns: 1fr;
    gap: var(--spacing-lg);
    max-width: 600px;
  }

  .relationship-display {
    order: -1;
    min-width: auto;
    max-width: 500px;
    margin: 0 auto;
  }

  .seiyuu-selector {
    min-width: auto;
    max-width: 500px;
    margin: 0 auto;
  }
}

@media (max-width: 768px) {
  .setup-section {
    padding: var(--spacing-lg);
  }

  .selection-grid {
    gap: var(--spacing-md);
  }

  .seiyuu-selector,
  .relationship-display {
    padding: var(--spacing-lg);
    min-height: auto;
  }

  .relationship-description {
    max-height: 120px;
  }

  .theater-header {
    flex-direction: column;
    gap: var(--spacing-md);
    text-align: center;
    padding: var(--spacing-lg);
  }

  .theater-controls {
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .control-buttons {
    width: 100%;
    justify-content: center;
  }

  .message-content {
    max-width: 85%;
  }
}
</style>
