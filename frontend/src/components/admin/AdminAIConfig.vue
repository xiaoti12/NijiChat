<template>
  <div class="admin-ai-config">
    <!-- 头部 -->
    <div class="manager-header">
      <div class="header-left">
        <h3 class="manager-title">管理员AI配置</h3>
        <p class="manager-desc">管理员专用AI模型配置，用于声优资料处理等管理功能</p>
      </div>
      <button @click="showAddForm = true" class="add-btn" :disabled="showAddForm">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M8 2a.5.5 0 0 1 .5.5v5h5a.5.5 0 0 1 0 1h-5v5a.5.5 0 0 1-1 0v-5h-5a.5.5 0 0 1 0-1h5v-5A.5.5 0 0 1 8 2z" />
        </svg>
        添加配置
      </button>
    </div>

    <!-- 模型列表 -->
    <div class="models-list">
      <div v-if="adminStore.adminAIModels.length === 0 && !showAddForm" class="empty-state">
        <div class="empty-icon">⚙️</div>
        <p class="empty-text">还没有配置管理员AI模型</p>
        <p class="empty-desc">配置AI模型后可以在前端直接处理声优资料，提升管理效率</p>
        <button @click="showAddForm = true" class="btn btn-primary">
          配置第一个模型
        </button>
      </div>

      <!-- 已配置的模型 -->
      <div v-for="model in adminStore.adminAIModels" :key="model.id" class="model-item">
        <div class="model-info">
          <div class="model-header">
            <div class="model-name-section">
              <span class="model-type">{{ getModelTypeLabel(model.type) }}</span>
              <h4 class="model-name">{{ model.name }}</h4>
              <span class="admin-badge">管理员专用</span>
            </div>
          </div>
          <div class="model-details">
            <span class="model-detail">{{ model.model_name }}</span>
            <span class="model-detail">最大Token: {{ model.max_tokens }}</span>
            <span class="model-detail">温度: {{ model.temperature }}</span>
          </div>
        </div>
        <div class="model-actions">
          <button @click="testModel(model)" class="action-btn test-btn" :disabled="testingModel === model.id">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z" />
              <path
                d="M6.271 5.055a.5.5 0 0 1 .52.444L7 7.105V13a.5.5 0 0 1-1 0V7.5a.5.5 0 0 1 .004-.085L6.271 5.055z" />
            </svg>
            {{ testingModel === model.id ? '测试中...' : '测试' }}
          </button>
          <button @click="editModel(model)" class="action-btn edit-btn">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708L9.707 9.707a.5.5 0 0 1-.708 0L6 6.707a.5.5 0 0 1 0-.708l6.146-6.146z" />
            </svg>
            编辑
          </button>
          <button @click="deleteModel(model)" class="action-btn delete-btn">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z" />
              <path fill-rule="evenodd"
                d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z" />
            </svg>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 添加/编辑表单 -->
    <div v-if="showAddForm || editingModel" class="model-form-overlay" @click="closeForm">
      <div class="model-form" @click.stop>
        <div class="form-header">
          <h4>{{ editingModel ? '编辑管理员AI配置' : '添加管理员AI配置' }}</h4>
          <button @click="closeForm" class="close-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="form-content">
          <div class="form-notice">
            <div class="notice-icon">ℹ️</div>
            <div class="notice-text">
              <strong>提示：</strong>此配置仅供管理员使用，用于声优资料AI处理等功能，与用户聊天AI配置完全独立。
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">配置名称 *</label>
              <input v-model="formData.name" type="text" class="form-input" placeholder="为这个配置起个名字" required />
            </div>
            <div class="form-group">
              <label class="form-label">模型类型 *</label>
              <select v-model="formData.type" class="form-select" required>
                <option value="">请选择模型类型</option>
                <option value="openai">OpenAI</option>
                <option value="gemini">Gemini</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">API Key *</label>
              <input v-model="formData.api_key" type="password" class="form-input" placeholder="输入API密钥" required />
            </div>
            <div class="form-group">
              <label class="form-label">API 端点</label>
              <input v-model="formData.api_endpoint" type="url" class="form-input"
                :placeholder="getDefaultEndpoint(formData.type)" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">模型名称 *</label>
              <input v-model="formData.model_name" type="text" class="form-input"
                :placeholder="getModelPlaceholder(formData.type)" required />
            </div>
            <div class="form-group">
              <label class="form-label">最大Token数</label>
              <input v-model.number="formData.max_tokens" type="number" class="form-input" placeholder="64000" />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label class="form-label">温度设置 ({{ formData.temperature }})</label>
              <input v-model.number="formData.temperature" type="range" class="form-slider" min="0" max="1"
                step="0.1" />
              <div class="slider-labels">
                <span>精确 (0)</span>
                <span>平衡 (0.5)</span>
                <span>创意 (1)</span>
              </div>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeForm" class="btn btn-secondary">
              取消
            </button>
            <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
              {{ editingModel ? '保存修改' : '添加配置' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 测试结果提示 -->
    <div v-if="testResult" class="test-result" :class="testResult.success ? 'success' : 'error'">
      <div class="result-icon">
        {{ testResult.success ? '✅' : '❌' }}
      </div>
      <div class="result-content">
        <div class="result-message">{{ testResult.message }}</div>
        <div v-if="testResult.latency" class="result-latency">
          响应时间: {{ testResult.latency }}ms
        </div>
      </div>
      <button @click="testResult = null" class="result-close">
        <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor">
          <path
            d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useAdminStore } from '@/stores/adminStore'
import { aiService } from '@/services/aiService'
import type { AIModelConfig, AIModelType, AIModelTestResult } from '@/types'

const adminStore = useAdminStore()

// 状态
const showAddForm = ref(false)
const editingModel = ref<AIModelConfig | null>(null)
const testingModel = ref<string | null>(null)
const testResult = ref<AIModelTestResult | null>(null)

// 表单数据
const formData = reactive<Omit<AIModelConfig, 'id'>>({
  name: '',
  type: '' as AIModelType,
  api_key: '',
  api_endpoint: '',
  model_name: '',
  is_lightweight: false, // 管理员模型通常不是轻量级
  max_tokens: 2000,
  temperature: 0.3 // 管理员任务通常需要更精确的输出
})

// 计算属性
const isFormValid = computed(() => {
  return formData.name &&
    formData.type &&
    formData.api_key &&
    formData.model_name
})

// 方法
function getModelTypeLabel(type: AIModelType): string {
  const labels = {
    'openai': 'OpenAI',
    'gemini': 'Gemini'
  }
  return labels[type] || type
}

function getDefaultEndpoint(type: AIModelType): string {
  const endpoints = {
    'openai': 'https://api.openai.com/v1',
    'gemini': 'https://generativelanguage.googleapis.com/v1beta'
  }
  return endpoints[type] || ''
}

function getModelPlaceholder(type: AIModelType): string {
  const placeholders = {
    'openai': '例如: gpt-4o-mini, gpt-4o, gpt-3.5-turbo',
    'gemini': '例如: gemini-pro, gemini-1.5-flash'
  }
  return placeholders[type] || '请输入模型名称'
}

function resetForm() {
  formData.name = ''
  formData.type = '' as AIModelType
  formData.api_key = ''
  formData.api_endpoint = ''
  formData.model_name = ''
  formData.is_lightweight = false
  formData.max_tokens = 2000
  formData.temperature = 0.3
}

function closeForm() {
  showAddForm.value = false
  editingModel.value = null
  resetForm()
}

function editModel(model: AIModelConfig) {
  editingModel.value = model
  Object.assign(formData, model)
  showAddForm.value = false
}

function handleSubmit() {
  if (!isFormValid.value) return

  try {
    if (editingModel.value) {
      // 更新现有模型
      adminStore.updateAdminAIModel(editingModel.value.id, formData)
      console.log('管理员AI配置已更新')
    } else {
      // 添加新模型
      const newModel = adminStore.addAdminAIModel({
        ...formData,
        api_endpoint: formData.api_endpoint || getDefaultEndpoint(formData.type)
      })
      console.log('新管理员AI配置已添加:', newModel)
    }

    closeForm()
  } catch (error) {
    console.error('保存管理员AI配置失败:', error)
  }
}

function deleteModel(model: AIModelConfig) {
  if (confirm(`确定要删除 "${model.name}" 配置吗？此操作不可恢复。`)) {
    adminStore.deleteAdminAIModel(model.id)
    console.log('管理员AI配置已删除')
  }
}

async function testModel(model: AIModelConfig) {
  testingModel.value = model.id
  testResult.value = null

  try {
    const startTime = Date.now()

    // 调用真实的AI服务测试
    const result = await aiService.testModel(model)
    const latency = Date.now() - startTime

    testResult.value = {
      success: result.success,
      message: result.message,
      latency: result.success ? latency : undefined
    }

    // 自动关闭测试结果提示
    setTimeout(() => {
      if (testResult.value) {
        testResult.value = null
      }
    }, 5000)

  } catch (error: any) {
    testResult.value = {
      success: false,
      message: `测试失败: ${error.message || error}`,
    }

    // 自动关闭错误提示
    setTimeout(() => {
      testResult.value = null
    }, 8000)
  } finally {
    testingModel.value = null
  }
}
</script>

<style scoped>
.admin-ai-config {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  padding: 0;
}

/* 头部 */
.manager-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.header-left {
  flex: 1;
}

.manager-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.manager-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
  line-height: 1.5;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.add-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.add-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 模型列表 */
.models-list {
  max-height: 500px;
  overflow-y: auto;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #6b7280;
  background: #f9fafb;
  border-radius: 12px;
  border: 2px dashed #d1d5db;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.empty-text {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #374151;
}

.empty-desc {
  font-size: 14px;
  margin-bottom: 24px;
  line-height: 1.5;
}

/* 模型项 */
.model-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  margin-bottom: 16px;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.model-item:hover {
  border-color: #667eea;
  background: linear-gradient(135deg, #f0f4ff 0%, #e6f2ff 100%);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

.model-info {
  flex: 1;
  min-width: 0;
}

.model-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.model-name-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.model-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.model-type {
  font-size: 11px;
  padding: 3px 8px;
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  border-radius: 6px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.admin-badge {
  font-size: 11px;
  padding: 3px 8px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border-radius: 6px;
  font-weight: 500;
}

.model-details {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #6b7280;
  flex-wrap: wrap;
}

.model-detail {
  background: rgba(255, 255, 255, 0.8);
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
  backdrop-filter: blur(10px);
}

.model-actions {
  display: flex;
  gap: 8px;
  margin-left: 16px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  color: #6b7280;
  backdrop-filter: blur(10px);
}

.action-btn:hover {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.test-btn:hover {
  border-color: #22c55e;
  color: #22c55e;
}

.edit-btn:hover {
  border-color: #f59e0b;
  color: #f59e0b;
}

.delete-btn:hover {
  border-color: #ef4444;
  color: #ef4444;
}

/* 表单样式 */
.model-form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.model-form {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 700px;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-bottom: 1px solid #e5e7eb;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 16px 16px 0 0;
}

.form-header h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.close-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: all 0.2s;
}

.close-btn:hover {
  background: white;
  color: #374151;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-content {
  padding: 28px;
}

.form-notice {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  margin-bottom: 24px;
}

.notice-icon {
  font-size: 16px;
}

.notice-text {
  font-size: 13px;
  color: #1e40af;
  line-height: 1.5;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.form-input,
.form-select {
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.2s;
  background: white;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-slider {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #e5e7eb;
  outline: none;
  -webkit-appearance: none;
  cursor: pointer;
  margin-bottom: 8px;
}

.form-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #9ca3af;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

/* 测试结果 */
.test-result {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-radius: 12px;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
  z-index: 1001;
  max-width: 350px;
  backdrop-filter: blur(10px);
}

.test-result.success {
  background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
  border: 1px solid #22c55e;
}

.test-result.error {
  background: linear-gradient(135deg, #fef2f2 0%, #fee2e2 100%);
  border: 1px solid #ef4444;
}

.result-icon {
  font-size: 18px;
}

.result-content {
  flex: 1;
  min-width: 0;
}

.result-message {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.result-latency {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

.result-close {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #6b7280;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.result-close:hover {
  background: rgba(0, 0, 0, 0.1);
}

/* 滚动条 */
.models-list::-webkit-scrollbar,
.model-form::-webkit-scrollbar {
  width: 6px;
}

.models-list::-webkit-scrollbar-track,
.model-form::-webkit-scrollbar-track {
  background: transparent;
}

.models-list::-webkit-scrollbar-thumb,
.model-form::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 3px;
}

.models-list::-webkit-scrollbar-thumb:hover,
.model-form::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

/* 响应式 */
@media (max-width: 768px) {
  .manager-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .add-btn {
    align-self: flex-start;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .model-item {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .model-actions {
    margin-left: 0;
    justify-content: space-between;
  }

  .test-result {
    bottom: 10px;
    right: 10px;
    left: 10px;
    max-width: none;
  }
}
</style>