<template>
  <div class="group-form">
    <div class="form-header">
      <h3>{{ isEditing ? '编辑群组' : '新建群组' }}</h3>
      <p class="form-desc">
        {{ isEditing ? '修改群组配置和成员设置' : '创建新的声优群组，设置成员和讨论模式' }}
      </p>
    </div>

    <div class="form-content">
      <!-- 基本信息 -->
      <div class="form-section">
        <h4>基本信息</h4>

        <div class="form-group">
          <label class="required">群组名称</label>
          <input
            v-model="form.name"
            type="text"
            class="input"
            placeholder="请输入群组名称"
            required
          />
        </div>

        <div class="form-group">
          <label class="required">群组描述</label>
          <textarea
            v-model="form.description"
            class="textarea"
            placeholder="请描述这个群组的主题、规则和特色..."
            rows="4"
            required
          ></textarea>
        </div>

        <div class="form-group">
          <label>
            <input
              v-model="form.is_discussion_mode"
              type="checkbox"
              class="checkbox"
            />
            启用群讨论模式
          </label>
          <div class="help-text">
            开启后，多个声优可以在同一个对话中轮流发言，创造更丰富的互动体验
          </div>
        </div>
      </div>

      <!-- 成员选择 -->
      <div class="form-section">
        <h4>成员管理</h4>
        <p class="section-desc">选择要加入此群组的声优成员</p>

        <div class="member-selection">
          <div class="search-box">
            <input
              v-model="memberSearchKeyword"
              type="text"
              class="input"
              placeholder="搜索声优..."
            />
          </div>

          <div class="member-grid">
            <div
              v-for="seiyuu in filteredSeiyuuList"
              :key="seiyuu.id"
              class="member-option"
              :class="{ selected: form.member_ids.includes(seiyuu.id) }"
              @click="toggleMember(seiyuu.id)"
            >
              <div class="member-info">
                <img
                  v-if="seiyuu.avatar_url"
                  :src="seiyuu.avatar_url"
                  :alt="seiyuu.name"
                  class="member-avatar"
                />
                <div class="member-details">
                  <div class="member-name">{{ seiyuu.name }}</div>
                  <div class="member-status">
                    <span
                      class="status-badge"
                      :class="seiyuu.status"
                    >
                      {{ getStatusText(seiyuu.status) }}
                    </span>
                    <div class="member-tags">
                      <span
                        v-for="tag in seiyuu.tags.slice(0, 2)"
                        :key="tag"
                        class="tag"
                      >
                        {{ tag }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
              <div class="selection-indicator">
                <div v-if="form.member_ids.includes(seiyuu.id)" class="checkmark">✓</div>
              </div>
            </div>
          </div>

          <div class="selected-summary">
            已选择 {{ form.member_ids.length }} 名声优
          </div>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="form-footer">
      <button
        type="button"
        @click="$emit('cancel')"
        class="btn btn-secondary"
      >
        取消
      </button>
      <button
        type="button"
        @click="handleSave"
        class="btn btn-primary"
        :disabled="saving || !isFormValid"
      >
        {{ saving ? '保存中...' : '保存' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import type { SeiyuuGroup, Seiyuu, SeiyuuStatus } from '@/types'
import { adminCreateGroup, adminUpdateGroup } from '@/services/apiService'

// Props和Emits
interface Props {
  group?: SeiyuuGroup | null
  seiyuuList: Seiyuu[]
}

const props = withDefaults(defineProps<Props>(), {
  group: null
})

const emit = defineEmits<{
  save: []
  cancel: []
}>()

// 响应式数据
const isEditing = computed(() => !!props.group)

const form = reactive({
  name: '',
  description: '',
  member_ids: [] as string[],
  is_discussion_mode: false
})

const memberSearchKeyword = ref('')
const saving = ref(false)

// 计算属性
const isFormValid = computed(() => {
  return form.name.trim() && form.description.trim() && form.member_ids.length > 0
})

const filteredSeiyuuList = computed(() => {
  let filtered = props.seiyuuList.filter(s => s.status === 'active' || s.status === 'pending')

  if (memberSearchKeyword.value.trim()) {
    const keyword = memberSearchKeyword.value.toLowerCase()
    filtered = filtered.filter(s =>
      s.name.toLowerCase().includes(keyword) ||
      s.tags.some(tag => tag.toLowerCase().includes(keyword))
    )
  }

  return filtered
})

// 监听props变化，初始化表单
watch(() => props.group, (group) => {
  if (group) {
    form.name = group.name
    form.description = group.description
    form.member_ids = [...group.member_ids]
    form.is_discussion_mode = group.is_discussion_mode
  } else {
    // 重置表单
    form.name = ''
    form.description = ''
    form.member_ids = []
    form.is_discussion_mode = false
  }
}, { immediate: true })

// 方法
function toggleMember(seiyuuId: string) {
  const index = form.member_ids.indexOf(seiyuuId)
  if (index >= 0) {
    form.member_ids.splice(index, 1)
  } else {
    form.member_ids.push(seiyuuId)
  }
}

function getStatusText(status: SeiyuuStatus): string {
  const statusMap = {
    pending: '待审核',
    active: '已发布',
    inactive: '已禁用'
  }
  return statusMap[status] || status
}

async function handleSave() {
  if (!isFormValid.value || saving.value) return

  try {
    saving.value = true

    if (isEditing.value && props.group) {
      // 编辑模式
      const updateData = {
        name: form.name,
        description: form.description,
        member_ids: form.member_ids,
        is_discussion_mode: form.is_discussion_mode
      }

      const response = await adminUpdateGroup(props.group.id, updateData)
      if (response.success) {
        emit('save')
      }
    } else {
      // 创建模式
      const createData = {
        name: form.name,
        description: form.description,
        member_ids: form.member_ids,
        is_discussion_mode: form.is_discussion_mode
      }

      const response = await adminCreateGroup(createData)
      if (response.success) {
        emit('save')
      }
    }
  } catch (error) {
    console.error('保存群组失败:', error)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.group-form {
  width: 100%;
  max-width: 600px;
}

/* 表单头部 */
.form-header {
  padding: var(--spacing-xl);
  border-bottom: 1px solid var(--border-color);
}

.form-header h3 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
}

.form-desc {
  margin: 0;
  color: var(--text-muted);
  font-size: var(--font-size-sm);
}

/* 表单内容 */
.form-content {
  padding: var(--spacing-xl);
  max-height: 60vh;
  overflow-y: auto;
}

.form-section {
  margin-bottom: var(--spacing-2xl);
}

.form-section:last-child {
  margin-bottom: 0;
}

.form-section h4 {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.section-desc {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: var(--font-size-sm);
  color: var(--text-muted);
  line-height: 1.5;
}

/* 表单组 */
.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.form-group label.required::after {
  content: ' *';
  color: var(--color-error);
}

.form-group label input[type="checkbox"] {
  margin-right: var(--spacing-sm);
}

.help-text {
  margin-top: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--text-muted);
  line-height: 1.4;
}

/* 成员选择 */
.member-selection {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
}

.search-box {
  margin-bottom: var(--spacing-lg);
}

.member-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-sm);
  max-height: 300px;
  overflow-y: auto;
  margin-bottom: var(--spacing-lg);
}

.member-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-primary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.member-option:hover {
  border-color: #667eea;
  background: var(--bg-secondary);
}

.member-option.selected {
  border-color: #667eea;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

.member-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex: 1;
}

.member-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  background: var(--bg-tertiary);
}

.member-details {
  flex: 1;
}

.member-name {
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--spacing-xs);
}

.member-status {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.status-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.status-badge.pending {
  background: #fef3cd;
  color: #856404;
}

.status-badge.active {
  background: #d1ecf1;
  color: #0c5460;
}

.status-badge.inactive {
  background: #f8d7da;
  color: #721c24;
}

.member-tags {
  display: flex;
  gap: var(--spacing-xs);
}

.tag {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
}

.selection-indicator {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkmark {
  width: 20px;
  height: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-bold);
}

.selected-summary {
  text-align: center;
  padding: var(--spacing-md);
  background: var(--bg-primary);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  color: var(--text-muted);
}

/* 表单底部 */
.form-footer {
  padding: var(--spacing-xl);
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-md);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .group-form {
    max-width: 100%;
  }

  .member-info {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }

  .member-status {
    flex-wrap: wrap;
  }

  .form-footer {
    flex-direction: column;
  }
}
</style>