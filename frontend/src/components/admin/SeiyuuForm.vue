<template>
  <div class="seiyuu-form">
    <div class="form-header">
      <h3>{{ isEditing ? '编辑声优' : '新增声优' }}</h3>
      <p class="form-desc">
        {{ isEditing ? '修改声优资料信息' : '创建新的声优资料，支持从萌娘百科获取数据' }}
      </p>
    </div>

    <div class="form-content">
      <!-- 基本信息 -->
      <div class="form-section">
        <h4>基本信息</h4>
        <div class="form-group">
          <label class="required">声优姓名</label>
          <input
            v-model="form.name"
            type="text"
            class="input"
            placeholder="请输入声优姓名"
            required
          />
        </div>

        <div class="form-group">
          <label>头像URL</label>
          <input
            v-model="form.avatar_url"
            type="url"
            class="input"
            placeholder="https://example.com/avatar.jpg"
          />
          <div v-if="form.avatar_url" class="avatar-preview">
            <img :src="form.avatar_url" :alt="form.name" />
          </div>
        </div>

        <div class="form-group">
          <label class="required">标签</label>
          <div class="tag-input-container">
            <div class="tag-list">
              <span
                v-for="(tag, index) in form.tags"
                :key="index"
                class="tag-item"
              >
                {{ tag }}
                <button
                  type="button"
                  @click="removeTag(index)"
                  class="tag-remove"
                >
                  ×
                </button>
              </span>
            </div>
            <div class="tag-input-group">
              <input
                v-model="newTag"
                type="text"
                class="input tag-input"
                placeholder="添加标签"
                @keyup.enter="addTag"
              />
              <button
                type="button"
                @click="addTag"
                class="btn btn-sm btn-secondary"
              >
                添加
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 萌娘百科数据获取 -->
      <div v-if="!isEditing" class="form-section">
        <h4>萌娘百科数据</h4>
        <p class="section-desc">可以从萌娘百科获取声优基础资料，然后通过AI处理为结构化内容</p>

        <div class="form-group">
          <label>萌娘百科页面名</label>
          <div class="moegirl-input-group">
            <input
              v-model="moegirlName"
              type="text"
              class="input"
              placeholder="请输入萌娘百科的声优页面名称"
            />
            <button
              type="button"
              @click="fetchMoegirlData"
              class="btn btn-secondary"
              :disabled="moegirlLoading || !moegirlName.trim()"
            >
              {{ moegirlLoading ? '获取中...' : '获取数据' }}
            </button>
          </div>
        </div>

        <div v-if="moegirlData" class="moegirl-result">
          <div class="result-header">
            <h5>获取到的原始数据</h5>
            <button
              type="button"
              @click="processWithAI"
              class="btn btn-primary"
              :disabled="aiProcessing"
            >
              {{ aiProcessing ? 'AI处理中...' : 'AI处理为Markdown' }}
            </button>
          </div>
          <textarea
            readonly
            :value="moegirlData.content"
            class="textarea moegirl-content"
            rows="8"
          ></textarea>
        </div>
      </div>

      <!-- Markdown资料 -->
      <div class="form-section">
        <h4>声优资料 (Markdown格式)</h4>
        <p class="section-desc">
          请使用Markdown格式编写声优详细资料，包括基本信息、作品经历、角色特点等
        </p>

        <div class="form-group">
          <label class="required">资料内容</label>
          <textarea
            v-model="form.profile_markdown"
            class="textarea profile-textarea"
            placeholder="请输入声优的详细资料，支持Markdown格式..."
            rows="20"
            required
          ></textarea>
          <div class="textarea-help">
            <p>支持Markdown语法：</p>
            <ul>
              <li># 标题，## 二级标题</li>
              <li>**粗体**，*斜体*</li>
              <li>- 列表项目</li>
              <li>[链接](URL)</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 状态设置 (编辑模式) -->
      <div v-if="isEditing" class="form-section">
        <h4>状态设置</h4>
        <div class="form-group">
          <label>发布状态</label>
          <select v-model="form.status" class="select">
            <option value="pending">待审核</option>
            <option value="active">已发布</option>
            <option value="inactive">已禁用</option>
          </select>
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
import type { Seiyuu, CreateSeiyuuRequest, UpdateSeiyuuRequest } from '@/types'
import {
  adminCreateSeiyuu,
  adminUpdateSeiyuu,
  getMoegirlData,
  processProfile
} from '@/services/apiService'

// Props和Emits
interface Props {
  seiyuu?: Seiyuu | null
}

const props = withDefaults(defineProps<Props>(), {
  seiyuu: null
})

const emit = defineEmits<{
  save: []
  cancel: []
}>()

// 响应式数据
const isEditing = computed(() => !!props.seiyuu)

const form = reactive({
  name: '',
  avatar_url: '',
  profile_markdown: '',
  tags: [] as string[],
  status: 'pending' as const
})

const newTag = ref('')
const moegirlName = ref('')
const moegirlData = ref<{ content: string; url: string } | null>(null)
const moegirlLoading = ref(false)
const aiProcessing = ref(false)
const saving = ref(false)

// 计算属性
const isFormValid = computed(() => {
  return form.name.trim() && form.profile_markdown.trim() && form.tags.length > 0
})

// 监听props变化，初始化表单
watch(() => props.seiyuu, (seiyuu) => {
  if (seiyuu) {
    form.name = seiyuu.name
    form.avatar_url = seiyuu.avatar_url || ''
    form.profile_markdown = seiyuu.profile_markdown
    form.tags = [...seiyuu.tags]
    form.status = seiyuu.status
  } else {
    // 重置表单
    form.name = ''
    form.avatar_url = ''
    form.profile_markdown = ''
    form.tags = []
    form.status = 'pending'
  }
}, { immediate: true })

// 方法
function addTag() {
  const tag = newTag.value.trim()
  if (tag && !form.tags.includes(tag)) {
    form.tags.push(tag)
    newTag.value = ''
  }
}

function removeTag(index: number) {
  form.tags.splice(index, 1)
}

async function fetchMoegirlData() {
  if (!moegirlName.value.trim()) return

  try {
    moegirlLoading.value = true
    const response = await getMoegirlData(moegirlName.value.trim())
    if (response.success) {
      moegirlData.value = response.data
      // 自动设置声优名称（如果还没有设置）
      if (!form.name) {
        form.name = moegirlName.value.trim()
      }
    }
  } catch (error) {
    console.error('获取萌娘百科数据失败:', error)
  } finally {
    moegirlLoading.value = false
  }
}

async function processWithAI() {
  if (!moegirlData.value) return

  try {
    aiProcessing.value = true
    const response = await processProfile({
      raw_data: moegirlData.value.content,
      seiyuu_name: form.name || moegirlName.value.trim()
    })

    if (response.success && response.data) {
      // 设置AI处理后的内容
      form.profile_markdown = response.data.markdown

      // 如果有建议的标签，添加到标签列表
      if (response.data.suggested_tags) {
        response.data.suggested_tags.forEach(tag => {
          if (!form.tags.includes(tag)) {
            form.tags.push(tag)
          }
        })
      }
    }
  } catch (error) {
    console.error('AI处理失败:', error)
  } finally {
    aiProcessing.value = false
  }
}

async function handleSave() {
  if (!isFormValid.value || saving.value) return

  try {
    saving.value = true

    if (isEditing.value && props.seiyuu) {
      // 编辑模式
      const updateData: UpdateSeiyuuRequest = {
        name: form.name,
        avatar_url: form.avatar_url || undefined,
        profile_markdown: form.profile_markdown,
        tags: form.tags,
        status: form.status
      }

      const response = await adminUpdateSeiyuu(props.seiyuu.id, updateData)
      if (response.success) {
        emit('save')
      }
    } else {
      // 创建模式
      const createData: CreateSeiyuuRequest = {
        name: form.name,
        avatar_url: form.avatar_url || undefined,
        profile_markdown: form.profile_markdown,
        tags: form.tags
      }

      const response = await adminCreateSeiyuu(createData)
      if (response.success) {
        emit('save')
      }
    }
  } catch (error) {
    console.error('保存声优失败:', error)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.seiyuu-form {
  width: 100%;
  max-width: 800px;
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

/* 头像预览 */
.avatar-preview {
  margin-top: var(--spacing-sm);
}

.avatar-preview img {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--border-color);
}

/* 标签输入 */
.tag-input-container {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  background: var(--bg-primary);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-md);
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-sm);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-sm);
}

.tag-remove {
  background: none;
  border: none;
  color: white;
  font-size: var(--font-size-md);
  cursor: pointer;
  padding: 0;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.tag-remove:hover {
  background: rgba(255, 255, 255, 0.2);
}

.tag-input-group {
  display: flex;
  gap: var(--spacing-sm);
}

.tag-input {
  flex: 1;
}

/* 萌娘百科数据 */
.moegirl-input-group {
  display: flex;
  gap: var(--spacing-sm);
}

.moegirl-result {
  margin-top: var(--spacing-lg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.result-header h5 {
  margin: 0;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
}

.moegirl-content {
  background: var(--bg-primary);
  font-family: monospace;
  font-size: var(--font-size-sm);
}

/* 资料文本域 */
.profile-textarea {
  font-family: monospace;
  font-size: var(--font-size-sm);
  line-height: 1.6;
  resize: vertical;
}

.textarea-help {
  margin-top: var(--spacing-sm);
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  color: var(--text-muted);
}

.textarea-help p {
  margin: 0 0 var(--spacing-sm) 0;
  font-weight: var(--font-weight-medium);
}

.textarea-help ul {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.textarea-help li {
  margin-bottom: var(--spacing-xs);
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
  .seiyuu-form {
    max-width: 100%;
  }

  .moegirl-input-group,
  .tag-input-group {
    flex-direction: column;
  }

  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }

  .form-footer {
    flex-direction: column;
  }
}
</style>