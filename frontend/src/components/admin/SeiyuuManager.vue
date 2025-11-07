<template>
  <div class="seiyuu-manager">
    <!-- 头部操作栏 -->
    <div class="manager-header">
      <div class="header-left">
        <h2>声优管理</h2>
        <p class="header-desc">管理声优资料，包括创建、编辑、发布和状态管理</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm = true">
          新增声优
        </button>
        <button class="btn btn-secondary" @click="refreshList">
          刷新列表
        </button>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-section">
      <div class="filter-group">
        <label>状态筛选：</label>
        <select v-model="filterStatus" class="select">
          <option value="">全部</option>
          <option value="pending">待审核</option>
          <option value="active">已发布</option>
          <option value="inactive">已禁用</option>
        </select>
      </div>
      <div class="filter-group">
        <label>搜索：</label>
        <input
          v-model="searchKeyword"
          type="text"
          placeholder="搜索声优名称或标签"
          class="input"
        />
      </div>
    </div>

    <!-- 声优列表 -->
    <div class="seiyuu-list" v-if="!loading">
      <div class="list-header">
        <div class="batch-actions" v-if="selectedSeiyuu.length > 0">
          <span>已选择 {{ selectedSeiyuu.length }} 项</span>
          <button class="btn btn-sm btn-success" @click="batchPublish">
            批量发布
          </button>
          <button class="btn btn-sm btn-warning" @click="batchDisable">
            批量禁用
          </button>
        </div>
        <div class="list-stats" v-else>
          共 {{ filteredSeiyuuList.length }} 个声优
        </div>
      </div>

      <div class="seiyuu-table">
        <div class="table-header">
          <div class="th checkbox-col">
            <input
              type="checkbox"
              :checked="isAllSelected"
              @change="toggleSelectAll"
            />
          </div>
          <div class="th name-col">姓名</div>
          <div class="th status-col">状态</div>
          <div class="th tags-col">标签</div>
          <div class="th date-col">创建时间</div>
          <div class="th actions-col">操作</div>
        </div>

        <div
          v-for="seiyuu in filteredSeiyuuList"
          :key="seiyuu.id"
          class="table-row"
        >
          <div class="td checkbox-col">
            <input
              type="checkbox"
              :value="seiyuu.id"
              v-model="selectedSeiyuu"
            />
          </div>
          <div class="td name-col">
            <div class="seiyuu-info">
              <img
                v-if="seiyuu.avatar_url"
                :src="seiyuu.avatar_url"
                :alt="seiyuu.name"
                class="seiyuu-avatar"
              />
              <div class="seiyuu-name">{{ seiyuu.name }}</div>
            </div>
          </div>
          <div class="td status-col">
            <span
              class="status-badge"
              :class="seiyuu.status"
            >
              {{ getStatusText(seiyuu.status) }}
            </span>
          </div>
          <div class="td tags-col">
            <div class="tag-list">
              <span
                v-for="tag in seiyuu.tags.slice(0, 3)"
                :key="tag"
                class="tag"
              >
                {{ tag }}
              </span>
              <span v-if="seiyuu.tags.length > 3" class="tag-more">
                +{{ seiyuu.tags.length - 3 }}
              </span>
            </div>
          </div>
          <div class="td date-col">
            {{ formatDate(seiyuu.created_at) }}
          </div>
          <div class="td actions-col">
            <button
              class="btn btn-sm btn-secondary"
              @click="editSeiyuu(seiyuu)"
            >
              编辑
            </button>
            <button
              v-if="seiyuu.status === 'pending'"
              class="btn btn-sm btn-success"
              @click="publishSeiyuu(seiyuu.id)"
            >
              发布
            </button>
            <button
              v-if="seiyuu.status === 'active'"
              class="btn btn-sm btn-warning"
              @click="toggleSeiyuuStatus(seiyuu.id, 'inactive')"
            >
              禁用
            </button>
            <button
              v-if="seiyuu.status === 'inactive'"
              class="btn btn-sm btn-success"
              @click="toggleSeiyuuStatus(seiyuu.id, 'active')"
            >
              启用
            </button>
            <button
              class="btn btn-sm btn-danger"
              @click="deleteSeiyuu(seiyuu.id, seiyuu.name)"
            >
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-section">
      <div class="loading-spinner"></div>
      <p>加载声优列表中...</p>
    </div>

    <!-- 创建/编辑表单弹窗 -->
    <div v-if="showCreateForm || editingSeiyuu" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <SeiyuuForm
          :seiyuu="editingSeiyuu"
          @save="handleSave"
          @cancel="closeModal"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Seiyuu, SeiyuuStatus } from '@/types'
import {
  adminGetAllSeiyuu,
  adminDeleteSeiyuu,
  adminUpdateSeiyuu
} from '@/services/apiService'
import SeiyuuForm from './SeiyuuForm.vue'

// 响应式数据
const seiyuuList = ref<Seiyuu[]>([])
const loading = ref(true)
const selectedSeiyuu = ref<string[]>([])
const filterStatus = ref<string>('')
const searchKeyword = ref('')
const showCreateForm = ref(false)
const editingSeiyuu = ref<Seiyuu | null>(null)

// 计算属性
const filteredSeiyuuList = computed(() => {
  let filtered = seiyuuList.value

  // 状态筛选
  if (filterStatus.value) {
    filtered = filtered.filter(s => s.status === filterStatus.value)
  }

  // 搜索筛选
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(s =>
      s.name.toLowerCase().includes(keyword) ||
      s.tags.some(tag => tag.toLowerCase().includes(keyword))
    )
  }

  return filtered
})

const isAllSelected = computed(() => {
  return filteredSeiyuuList.value.length > 0 &&
         selectedSeiyuu.value.length === filteredSeiyuuList.value.length
})

// 方法
async function loadSeiyuuList() {
  try {
    loading.value = true
    const response = await adminGetAllSeiyuu()
    if (response.success) {
      seiyuuList.value = response.data
    }
  } catch (error) {
    console.error('加载声优列表失败:', error)
  } finally {
    loading.value = false
  }
}

function refreshList() {
  selectedSeiyuu.value = []
  loadSeiyuuList()
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedSeiyuu.value = []
  } else {
    selectedSeiyuu.value = filteredSeiyuuList.value.map(s => s.id)
  }
}

function editSeiyuu(seiyuu: Seiyuu) {
  editingSeiyuu.value = seiyuu
}

function closeModal() {
  showCreateForm.value = false
  editingSeiyuu.value = null
}

async function handleSave() {
  closeModal()
  await refreshList()
}

async function publishSeiyuu(id: string) {
  if (!confirm('确定要发布这个声优吗？')) return

  try {
    const response = await adminUpdateSeiyuu(id, { status: 'active' })
    if (response.success) {
      await refreshList()
    }
  } catch (error) {
    console.error('发布声优失败:', error)
  }
}

async function toggleSeiyuuStatus(id: string, status: SeiyuuStatus) {
  const action = status === 'active' ? '启用' : '禁用'
  if (!confirm(`确定要${action}这个声优吗？`)) return

  try {
    const response = await adminUpdateSeiyuu(id, { status })
    if (response.success) {
      await refreshList()
    }
  } catch (error) {
    console.error(`${action}声优失败:`, error)
  }
}

async function deleteSeiyuu(id: string, name: string) {
  if (!confirm(`确定要删除声优「${name}」吗？此操作不可撤销！`)) return

  try {
    const response = await adminDeleteSeiyuu(id)
    if (response.success) {
      await refreshList()
    }
  } catch (error) {
    console.error('删除声优失败:', error)
  }
}

async function batchPublish() {
  if (!confirm(`确定要批量发布 ${selectedSeiyuu.value.length} 个声优吗？`)) return

  try {
    await Promise.all(
      selectedSeiyuu.value.map(id => adminUpdateSeiyuu(id, { status: 'active' }))
    )
    selectedSeiyuu.value = []
    await refreshList()
  } catch (error) {
    console.error('批量发布失败:', error)
  }
}

async function batchDisable() {
  if (!confirm(`确定要批量禁用 ${selectedSeiyuu.value.length} 个声优吗？`)) return

  try {
    await Promise.all(
      selectedSeiyuu.value.map(id => adminUpdateSeiyuu(id, { status: 'inactive' }))
    )
    selectedSeiyuu.value = []
    await refreshList()
  } catch (error) {
    console.error('批量禁用失败:', error)
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

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 生命周期
onMounted(() => {
  loadSeiyuuList()
})
</script>

<style scoped>
.seiyuu-manager {
  padding: var(--spacing-xl);
  height: 100%;
}

/* 头部 */
.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
}

.header-left h2 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
}

.header-desc {
  margin: 0;
  color: var(--text-muted);
  font-size: var(--font-size-sm);
}

.header-actions {
  display: flex;
  gap: var(--spacing-md);
}

/* 筛选区域 */
.filter-section {
  display: flex;
  gap: var(--spacing-xl);
  margin-bottom: var(--spacing-xl);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.filter-group label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  white-space: nowrap;
}

.filter-group .select,
.filter-group .input {
  min-width: 150px;
}

/* 列表区域 */
.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.batch-actions span {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
}

.list-stats {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
}

/* 表格样式 */
.seiyuu-table {
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.table-header,
.table-row {
  display: grid;
  grid-template-columns: 40px 1fr 100px 200px 120px 200px;
  gap: var(--spacing-md);
  align-items: center;
}

.table-header {
  background: var(--bg-tertiary);
  padding: var(--spacing-md);
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  border-bottom: 1px solid var(--border-color);
}

.table-row {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
  transition: background-color 0.2s ease;
}

.table-row:hover {
  background: var(--bg-secondary);
}

.table-row:last-child {
  border-bottom: none;
}

.th,
.td {
  overflow: hidden;
}

.checkbox-col {
  text-align: center;
}

/* 声优信息 */
.seiyuu-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.seiyuu-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  background: var(--bg-tertiary);
}

.seiyuu-name {
  font-weight: var(--font-weight-medium);
}

/* 状态徽章 */
.status-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  text-align: center;
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

/* 标签 */
.tag-list {
  display: flex;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
}

.tag {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  white-space: nowrap;
}

.tag-more {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  color: var(--text-muted);
}

/* 操作按钮 */
.actions-col {
  display: flex;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
}

/* 加载状态 */
.loading-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3xl);
  color: var(--text-muted);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--bg-tertiary);
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: var(--spacing-lg);
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 弹窗 */
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
  max-width: 800px;
  max-height: 90vh;
  overflow: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .manager-header {
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .filter-section {
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .table-header,
  .table-row {
    grid-template-columns: 1fr;
    gap: var(--spacing-sm);
  }

  .th,
  .td {
    padding: var(--spacing-sm);
  }
}
</style>