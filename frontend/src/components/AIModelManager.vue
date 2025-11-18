<template>
  <div class="ai-model-manager">
    <!-- 头部 -->
    <div class="manager-header">
      <h3 class="manager-title">AI 模型配置</h3>
      <div class="header-actions">
        <button @click="showTestCodeForm = true" class="test-code-btn">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M2.5 8a5.5 5.5 0 0 1 8.25-4.764.5.5 0 0 0 .5-.866A6.5 6.5 0 1 0 14.5 8a.5.5 0 0 0-1 0 5.5 5.5 0 1 1-11 0z"/>
            <path d="M15.354 3.354a.5.5 0 0 0-.708-.708L8 9.293 5.354 6.646a.5.5 0 1 0-.708.708l3 3a.5.5 0 0 0 .708 0l7-7z"/>
          </svg>
          使用体验码
        </button>
        <button @click="showAddForm = true" class="add-btn" :disabled="showAddForm">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path
              d="M8 2a.5.5 0 0 1 .5.5v5h5a.5.5 0 0 1 0 1h-5v5a.5.5 0 0 1-1 0v-5h-5a.5.5 0 0 1 0-1h5v-5A.5.5 0 0 1 8 2z" />
          </svg>
          添加模型
        </button>
      </div>
    </div>

    <!-- 模型列表 -->
    <div class="models-list">
      <div v-if="aiModelStore.models.length === 0 && !showAddForm" class="empty-state">
        <div class="empty-icon">🤖</div>
        <p class="empty-text">还没有配置AI模型</p>
        <button @click="showAddForm = true" class="btn btn-primary">
          配置第一个模型
        </button>
      </div>

      <!-- 已配置的模型 -->
      <div v-for="model in aiModelStore.models" :key="model.id" class="model-item">
        <div class="model-info">
          <div class="model-header">
            <div class="model-name-section">
              <span class="model-type">{{ getModelTypeLabel(model.type) }}</span>
              <h4 class="model-name">{{ model.name }}</h4>
            </div>
          </div>
          <div class="model-details">
            <span class="model-detail">{{ model.model_name }}</span>
            <span class="model-detail">最大Token: {{ model.max_tokens }}</span>
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
          <h4>{{ editingModel ? '编辑模型配置' : '添加AI模型' }}</h4>
          <button @click="closeForm" class="close-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
            </svg>
          </button>
        </div>

        <form @submit.prevent="handleSubmit" class="form-content">
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
              <input v-model.number="formData.max_tokens" type="number" class="form-input" placeholder="8000" />
            </div>
          </div>


          <div class="form-actions">
            <button type="button" @click="closeForm" class="btn btn-secondary">
              取消
            </button>
            <button type="submit" class="btn btn-primary" :disabled="!isFormValid">
              {{ editingModel ? '保存修改' : '添加模型' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 体验码输入弹窗 -->
    <div v-if="showTestCodeForm" class="model-form-overlay" @click="closeTestCodeForm">
      <div class="model-form test-code-form" @click.stop>
        <div class="form-header">
          <h4>使用体验码获取配置</h4>
          <button @click="closeTestCodeForm" class="close-btn">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path
                d="M2.146 2.854a.5.5 0 1 1 .708-.708L8 7.293l5.146-5.147a.5.5 0 0 1 .708.708L8.707 8l5.147 5.146a.5.5 0 0 1-.708.708L8 8.707l-5.146 5.147a.5.5 0 0 1-.708-.708L7.293 8 2.146 2.854Z" />
            </svg>
          </button>
        </div>

        <div class="form-content">
          <div class="form-group">
            <label class="form-label">请输入体验码</label>
            <input
              v-model="testCode"
              type="text"
              class="form-input"
              placeholder=""
              @keyup.enter="submitTestCode"
            />
            <p class="form-hint">输入体验码后，系统将自动获取预配置的AI模型参数，让您快速体验功能</p>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeTestCodeForm" class="btn btn-secondary">
              取消
            </button>
            <button
              @click="submitTestCode"
              class="btn btn-primary"
              :disabled="!testCode.trim() || isSubmittingTestCode"
            >
              {{ isSubmittingTestCode ? '获取中...' : '获取配置' }}
            </button>
          </div>
        </div>
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
import { useAIModelStore } from '@/stores/aiModelStore'
import { aiService } from '@/services/aiService'
import { getModelConfigByTestCode } from '@/services/apiService'
import type { AIModelConfig, AIModelType, AIModelTestResult } from '@/types'

const aiModelStore = useAIModelStore()

// 状态
const showAddForm = ref(false)
const editingModel = ref<AIModelConfig | null>(null)
const testingModel = ref<string | null>(null)
const testResult = ref<AIModelTestResult | null>(null)

// 测试码相关状态
const showTestCodeForm = ref(false)
const testCode = ref('')
const isSubmittingTestCode = ref(false)

// 表单数据
const formData = reactive<Omit<AIModelConfig, 'id'>>({
  name: '',
  type: '' as AIModelType,
  api_key: '',
  api_endpoint: '',
  model_name: '',
  is_lightweight: false, // 始终为false，用户不可配置
  max_tokens: 2000,
  temperature: 0.7 // 使用默认值，用户不可配置
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
    'openai': '例如: gpt-3.5-turbo, gpt-4, gpt-4o',
    'gemini': '例如: gemini-pro, gemini-pro-vision'
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
  formData.temperature = 0.7
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
      aiModelStore.updateModel(editingModel.value.id, formData)
      console.log('模型配置已更新')
    } else {
      // 添加新模型
      const newModel = aiModelStore.addModel({
        ...formData,
        api_endpoint: formData.api_endpoint || getDefaultEndpoint(formData.type)
      })
      console.log('新模型配置已添加:', newModel)
    }

    closeForm()
  } catch (error) {
    console.error('保存模型配置失败:', error)
  }
}

function deleteModel(model: AIModelConfig) {
  if (confirm(`确定要删除 "${model.name}" 配置吗？此操作不可恢复。`)) {
    aiModelStore.deleteModel(model.id)
    console.log('模型配置已删除')
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
      if (testResult.value && testResult.value === testResult.value) {
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

// 测试码相关方法
function closeTestCodeForm() {
  showTestCodeForm.value = false
  testCode.value = ''
  isSubmittingTestCode.value = false
}

async function submitTestCode() {
  if (!testCode.value.trim() || isSubmittingTestCode.value) return

  isSubmittingTestCode.value = true
  testResult.value = null

  try {
    // 调用后端API获取配置
    const response = await getModelConfigByTestCode(testCode.value.trim())

    if (response.success && response.data) {
      // 将配置保存到store
      const newModel = aiModelStore.addModel({
        name: response.data.name,
        type: response.data.type as AIModelType,
        api_key: response.data.api_key,
        api_endpoint: response.data.api_endpoint,
        model_name: response.data.model_name,
        is_lightweight: response.data.is_lightweight,
        max_tokens: response.data.max_tokens,
        temperature: response.data.temperature
      })

      // 显示成功提示
      testResult.value = {
        success: true,
        message: `成功添加模型配置: ${response.data.name}`
      }

      // 关闭测试码弹窗
      closeTestCodeForm()

      // 自动关闭成功提示
      setTimeout(() => {
        testResult.value = null
      }, 3000)

      console.log('通过测试码添加的模型:', newModel)
    } else {
      throw new Error(response.error || '获取配置失败')
    }
  } catch (error: any) {
    testResult.value = {
      success: false,
      message: error.response?.data?.error || error.message || '获取配置失败，请检查测试码是否正确'
    }

    // 自动关闭错误提示
    setTimeout(() => {
      testResult.value = null
    }, 5000)
  } finally {
    isSubmittingTestCode.value = false
  }
}
</script>

<style scoped>
.ai-model-manager {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  padding: 16px;
}

/* 头部 */
.manager-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.manager-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.test-code-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.test-code-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.add-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
}

.add-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 模型列表 */
.models-list {
  max-height: 400px;
  overflow-y: auto;
}

.empty-state {
  text-align: center;
  padding: 32px 16px;
  color: #6b7280;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 12px;
  opacity: 0.6;
}

.empty-text {
  font-size: 14px;
  margin-bottom: 16px;
}

/* 模型项 */
.model-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.2s;
}

.model-item:hover {
  border-color: #fead00;
  background: rgba(254, 173, 0, 0.05);
}

.model-info {
  flex: 1;
  min-width: 0;
}

.model-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.model-name-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}


.model-type {
  font-size: 11px;
  padding: 2px 6px;
  background: #fead00;
  color: white;
  border-radius: 4px;
  font-weight: 500;
}

.model-badge.lightweight {
  font-size: 11px;
  padding: 2px 6px;
  background: #10b981;
  color: white;
  border-radius: 4px;
  font-weight: 500;
}

.model-details {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #6b7280;
}

.model-detail {
  background: #e5e7eb;
  padding: 2px 6px;
  border-radius: 4px;
}

.model-actions {
  display: flex;
  gap: 6px;
  margin-left: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 4px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  color: #6b7280;
}

.action-btn:hover {
  border-color: #fead00;
  color: #fead00;
}

.test-btn:hover {
  background: rgba(34, 197, 94, 0.1);
  border-color: #22c55e;
  color: #22c55e;
}

.edit-btn:hover {
  background: rgba(245, 158, 11, 0.1);
  border-color: #f59e0b;
  color: #f59e0b;
}

.delete-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
  color: #ef4444;
}

/* 表单覆盖层 */
.model-form-overlay {
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

.model-form {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.form-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.close-btn {
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

.close-btn:hover {
  background: #f3f4f6;
}

.form-content {
  padding: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 20px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-label {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-hint {
  font-size: 12px;
  color: #6b7280;
  margin-top: 6px;
  margin-bottom: 0;
}

.test-code-form {
  max-width: 450px;
}

.form-input,
.form-select {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s;
  background: white;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #fead00;
  box-shadow: 0 0 0 3px rgba(254, 173, 0, 0.1);
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
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #fead00;
  cursor: pointer;
  transition: all 0.2s;
}

.slider-labels {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #9ca3af;
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  cursor: pointer;
}

.checkbox-label {
  font-size: 13px;
  color: #374151;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  padding: 10px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #fead00, #ff791b);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(254, 173, 0, 0.4);
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
}

/* 测试结果 */
.test-result {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  max-width: 300px;
}

.test-result.success {
  background: #f0fdf4;
  border: 1px solid #22c55e;
}

.test-result.error {
  background: #fef2f2;
  border: 1px solid #ef4444;
}

.result-icon {
  font-size: 16px;
}

.result-content {
  flex: 1;
  min-width: 0;
}

.result-message {
  font-size: 13px;
  font-weight: 500;
  color: #1f2937;
}

.result-latency {
  font-size: 11px;
  color: #6b7280;
  margin-top: 2px;
}

.result-close {
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #6b7280;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.result-close:hover {
  background: rgba(0, 0, 0, 0.1);
}

/* 滚动条 */
.models-list::-webkit-scrollbar,
.model-form::-webkit-scrollbar {
  width: 4px;
}

.models-list::-webkit-scrollbar-track,
.model-form::-webkit-scrollbar-track {
  background: transparent;
}

.models-list::-webkit-scrollbar-thumb,
.model-form::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

/* 响应式 */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .model-item {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .model-actions {
    margin-left: 0;
    justify-content: space-between;
  }
}
</style>