<template>
  <div class="modal-overlay" @click="$emit('close')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ isEdit ? '编辑关系' : '创建关系' }}</h3>
        <button @click="$emit('close')" class="close-btn">&times;</button>
      </div>

      <form @submit.prevent="handleSubmit" class="relationship-form">
        <!-- 声优选择 -->
        <div v-if="!isEdit" class="form-group">
          <label class="form-label">声优A</label>
          <select v-model="form.seiyuu_id_a" class="input" required>
            <option value="">请选择声优</option>
            <option
              v-for="seiyuu in availableSeiyuuA"
              :key="seiyuu.id"
              :value="seiyuu.id"
            >
              {{ seiyuu.name }}
            </option>
          </select>
        </div>

        <div v-if="!isEdit" class="form-group">
          <label class="form-label">声优B</label>
          <select v-model="form.seiyuu_id_b" class="input" required>
            <option value="">请选择声优</option>
            <option
              v-for="seiyuu in availableSeiyuuB"
              :key="seiyuu.id"
              :value="seiyuu.id"
            >
              {{ seiyuu.name }}
            </option>
          </select>
        </div>

        <!-- 编辑模式显示声优信息 -->
        <div v-if="isEdit" class="selected-seiyuu">
          <div class="seiyuu-info">
            <span class="seiyuu-name">{{ props.relationship?.seiyuu_name_a }}</span>
            <span class="relationship-arrow">↔</span>
            <span class="seiyuu-name">{{ props.relationship?.seiyuu_name_b }}</span>
          </div>
        </div>

        <!-- AI生成按钮 -->
        <div v-if="!isEdit && canGenerate" class="ai-generate-section">
          <button
            type="button"
            @click="generateRelationship"
            :disabled="generating"
            class="btn btn-secondary btn-full"
          >
            {{ generating ? 'AI生成中...' : 'AI智能生成关系描述' }}
          </button>
        </div>

        <!-- 关系描述 -->
        <div class="form-group">
          <label class="form-label">关系描述</label>
          <textarea
            v-model="form.relationship_description"
            class="input textarea"
            rows="6"
            placeholder="请描述两个声优之间的关系..."
            required
          ></textarea>
        </div>

        <div class="form-actions">
          <button type="button" @click="$emit('close')" class="btn btn-secondary">
            取消
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="btn btn-primary"
          >
            {{ submitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import type { Seiyuu, SeiyuuRelationship } from '@/types'
import {
  adminGetAllSeiyuu,
  adminCreateRelationship,
  adminUpdateRelationship,
  adminGenerateRelationship
} from '@/services/apiService'

// Props
interface Props {
  relationship?: SeiyuuRelationship | null
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits<{
  close: []
  success: []
}>()

// 响应式数据
const seiyuuList = ref<Seiyuu[]>([])
const submitting = ref(false)
const generating = ref(false)

const form = reactive({
  seiyuu_id_a: '',
  seiyuu_id_b: '',
  relationship_description: ''
})

// 计算属性
const isEdit = computed(() => !!props.relationship)

const availableSeiyuuA = computed(() => {
  return seiyuuList.value.filter(s => s.status === 'active')
})

const availableSeiyuuB = computed(() => {
  return seiyuuList.value.filter(s =>
    s.status === 'active' && s.id !== form.seiyuu_id_a
  )
})

const canGenerate = computed(() => {
  return form.seiyuu_id_a && form.seiyuu_id_b && !isEdit.value
})

// 方法
async function loadSeiyuu() {
  try {
    const response = await adminGetAllSeiyuu()
    if (response.success) {
      seiyuuList.value = response.data
    }
  } catch (error: any) {
    console.error('获取声优列表失败:', error.message)
  }
}

async function generateRelationship() {
  if (!form.seiyuu_id_a || !form.seiyuu_id_b) {
    return
  }

  generating.value = true
  try {
    const response = await adminGenerateRelationship({
      seiyuu_id_a: form.seiyuu_id_a,
      seiyuu_id_b: form.seiyuu_id_b
    })

    if (response.success && response.data?.relationship_description) {
      form.relationship_description = response.data.relationship_description
    } else {
      alert('AI生成失败，请手动填写关系描述')
    }
  } catch (error: any) {
    console.error('AI生成关系失败:', error.message)
    alert(`AI生成失败：${error.message}`)
  } finally {
    generating.value = false
  }
}

async function handleSubmit() {
  submitting.value = true

  try {
    let response
    if (isEdit.value && props.relationship) {
      // 编辑模式
      response = await adminUpdateRelationship(props.relationship.id, {
        relationship_description: form.relationship_description
      })
    } else {
      // 创建模式
      response = await adminCreateRelationship({
        seiyuu_id_a: form.seiyuu_id_a,
        seiyuu_id_b: form.seiyuu_id_b,
        relationship_description: form.relationship_description
      })
    }

    if (response.success) {
      emit('success')
    }
  } catch (error: any) {
    console.error('保存关系失败:', error.message)
    alert(`保存失败：${error.message}`)
  } finally {
    submitting.value = false
  }
}

function initForm() {
  if (props.relationship) {
    form.seiyuu_id_a = props.relationship.seiyuu_id_a
    form.seiyuu_id_b = props.relationship.seiyuu_id_b
    form.relationship_description = props.relationship.relationship_description
  } else {
    form.seiyuu_id_a = ''
    form.seiyuu_id_b = ''
    form.relationship_description = ''
  }
}

// 生命周期
onMounted(() => {
  loadSeiyuu()
  initForm()
})
</script>

<style scoped>
.modal-overlay {
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
}

.modal-content {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xl);
  border-bottom: 1px solid var(--border-color);
}

.modal-header h3 {
  margin: 0;
  font-size: var(--font-size-lg);
}

.close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-2xl);
  cursor: pointer;
  color: var(--text-muted);
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: var(--text-primary);
}

.relationship-form {
  padding: var(--spacing-xl);
}

.form-group {
  margin-bottom: var(--spacing-lg);
}

.form-label {
  display: block;
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
}

.textarea {
  resize: vertical;
  min-height: 120px;
}

.selected-seiyuu {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border-left: 3px solid var(--color-primary);
}

.seiyuu-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.seiyuu-name {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.relationship-arrow {
  color: var(--text-muted);
  font-size: var(--font-size-lg);
}

.ai-generate-section {
  margin-bottom: var(--spacing-lg);
}

.btn-full {
  width: 100%;
}

.form-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
  margin-top: var(--spacing-xl);
  padding-top: var(--spacing-lg);
  border-top: 1px solid var(--border-color);
}

@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: var(--spacing-md);
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .form-actions .btn {
    width: 100%;
  }
}
</style>