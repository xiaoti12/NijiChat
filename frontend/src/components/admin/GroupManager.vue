<template>
  <div class="group-manager">
    <!-- 头部操作栏 -->
    <div class="manager-header">
      <div class="header-left">
        <h2>群组管理</h2>
        <p class="header-desc">管理声优群组，设置群聊模式和成员配置</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateForm = true">
          新建群组
        </button>
        <button class="btn btn-secondary" @click="refreshList">
          刷新列表
        </button>
      </div>
    </div>

    <!-- 群组列表 -->
    <div class="group-list" v-if="!loading">
      <div class="list-stats">
        共 {{ groupList.length }} 个群组
      </div>

      <div class="group-grid">
        <div
          v-for="group in groupList"
          :key="group.id"
          class="group-card"
        >
          <div class="card-header">
            <h3 class="group-name">{{ group.name }}</h3>
            <div class="group-actions">
              <button
                class="btn btn-sm btn-secondary"
                @click="editGroup(group)"
              >
                编辑
              </button>
              <button
                class="btn btn-sm btn-danger"
                @click="deleteGroup(group.id, group.name)"
              >
                删除
              </button>
            </div>
          </div>

          <div class="card-content">
            <p class="group-description">{{ group.description }}</p>

            <div class="group-meta">
              <div class="meta-item">
                <span class="meta-label">成员数量:</span>
                <span class="meta-value">{{ group.member_ids.length }} 人</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">讨论模式:</span>
                <span
                  class="meta-value mode-badge"
                  :class="{ active: group.is_discussion_mode }"
                >
                  {{ group.is_discussion_mode ? '开启' : '关闭' }}
                </span>
              </div>
              <div class="meta-item">
                <span class="meta-label">创建时间:</span>
                <span class="meta-value">{{ formatDate(group.created_at) }}</span>
              </div>
            </div>

            <div v-if="group.member_ids.length > 0" class="group-members">
              <h4>成员列表</h4>
              <div class="member-list">
                <div
                  v-for="memberId in group.member_ids"
                  :key="memberId"
                  class="member-item"
                >
                  <img
                    v-if="getMemberInfo(memberId)?.avatar_url"
                    :src="getMemberInfo(memberId)?.avatar_url"
                    :alt="getMemberInfo(memberId)?.name"
                    class="member-avatar"
                  />
                  <span class="member-name">
                    {{ getMemberInfo(memberId)?.name || '未知声优' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="groupList.length === 0" class="empty-state">
          <div class="empty-icon">📁</div>
          <h3>暂无群组</h3>
          <p>点击"新建群组"创建第一个声优群组</p>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else class="loading-section">
      <div class="loading-spinner"></div>
      <p>加载群组列表中...</p>
    </div>

    <!-- 创建/编辑表单弹窗 -->
    <div v-if="showCreateForm || editingGroup" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <GroupForm
          :group="editingGroup"
          :seiyuu-list="seiyuuList"
          @save="handleSave"
          @cancel="closeModal"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { SeiyuuGroup, Seiyuu } from '@/types'
import {
  adminGetAllGroups,
  adminDeleteGroup,
  adminGetAllSeiyuu
} from '@/services/apiService'
import GroupForm from './GroupForm.vue'

// 响应式数据
const groupList = ref<SeiyuuGroup[]>([])
const seiyuuList = ref<Seiyuu[]>([])
const loading = ref(true)
const showCreateForm = ref(false)
const editingGroup = ref<SeiyuuGroup | null>(null)

// 方法
async function loadData() {
  try {
    loading.value = true
    const [groupResponse, seiyuuResponse] = await Promise.all([
      adminGetAllGroups(),
      adminGetAllSeiyuu()
    ])

    if (groupResponse.success) {
      groupList.value = groupResponse.data
    }

    if (seiyuuResponse.success) {
      seiyuuList.value = seiyuuResponse.data
    }
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

function refreshList() {
  loadData()
}

function editGroup(group: SeiyuuGroup) {
  editingGroup.value = group
}

function closeModal() {
  showCreateForm.value = false
  editingGroup.value = null
}

async function handleSave() {
  closeModal()
  await refreshList()
}

async function deleteGroup(id: string, name: string) {
  if (!confirm(`确定要删除群组「${name}」吗？此操作不可撤销！`)) return

  try {
    const response = await adminDeleteGroup(id)
    if (response.success) {
      await refreshList()
    }
  } catch (error) {
    console.error('删除群组失败:', error)
  }
}

function getMemberInfo(memberId: string): Seiyuu | undefined {
  return seiyuuList.value.find(s => s.id === memberId)
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
  loadData()
})
</script>

<style scoped>
.group-manager {
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

/* 列表统计 */
.list-stats {
  margin-bottom: var(--spacing-lg);
  font-size: var(--font-size-sm);
  color: var(--text-muted);
}

/* 群组网格 */
.group-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: var(--spacing-xl);
}

/* 群组卡片 */
.group-card {
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

.group-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-lg);
}

.group-name {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
}

.group-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.card-content {
  color: var(--text-secondary);
}

.group-description {
  margin: 0 0 var(--spacing-lg) 0;
  line-height: 1.6;
  color: var(--text-secondary);
}

/* 元数据 */
.group-meta {
  margin-bottom: var(--spacing-lg);
}

.meta-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-color);
}

.meta-item:last-child {
  border-bottom: none;
}

.meta-label {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
}

.meta-value {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

.mode-badge {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  background: var(--bg-tertiary);
  color: var(--text-muted);
}

.mode-badge.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

/* 成员列表 */
.group-members h4 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
}

.member-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.member-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  background: var(--bg-primary);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.member-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.member-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: var(--spacing-3xl);
  color: var(--text-muted);
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: var(--spacing-lg);
}

.empty-state h3 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-lg);
  color: var(--text-secondary);
}

.empty-state p {
  margin: 0;
  font-size: var(--font-size-sm);
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
  max-width: 600px;
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

  .group-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .group-actions {
    align-self: flex-start;
  }

  .member-list {
    flex-direction: column;
  }
}
</style>