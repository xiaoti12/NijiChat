<template>
  <div class="relationship-manager">
    <div class="manager-header">
      <h2>声优关系管理</h2>
      <button @click="showCreateForm = true" class="btn btn-primary">
        创建关系
      </button>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <div class="search-controls">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索声优名称..."
          class="input search-input"
        />
        <button @click="loadRelationships" class="btn btn-secondary">
          刷新
        </button>
      </div>
    </div>

    <!-- 关系列表 -->
    <div class="relationships-section">
      <div v-if="loading" class="loading">
        加载中...
      </div>

      <div v-else-if="relationships.length === 0" class="empty-state">
        <p>暂无声优关系记录</p>
      </div>

      <div v-else class="relationships-grid">
        <div
          v-for="relationship in filteredRelationships"
          :key="relationship.id"
          class="relationship-card"
        >
          <div class="relationship-header">
            <div class="seiyuu-pair">
              <span class="seiyuu-name">{{ relationship.seiyuu_name_a }}</span>
              <span class="relationship-arrow">↔</span>
              <span class="seiyuu-name">{{ relationship.seiyuu_name_b }}</span>
            </div>
            <div class="relationship-actions">
              <button
                @click="editRelationship(relationship)"
                class="btn btn-small btn-secondary"
              >
                编辑
              </button>
              <button
                @click="deleteRelationship(relationship.id)"
                class="btn btn-small btn-danger"
              >
                删除
              </button>
            </div>
          </div>

          <div class="relationship-description">
            {{ relationship.relationship_description }}
          </div>

          <div class="relationship-meta">
            <span class="meta-item">创建时间: {{ formatDate(relationship.created_at) }}</span>
            <span class="meta-item">更新时间: {{ formatDate(relationship.updated_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑关系表单模态框 -->
    <RelationshipForm
      v-if="showCreateForm || editingRelationship"
      :relationship="editingRelationship"
      @close="closeForm"
      @success="onFormSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { SeiyuuRelationship } from '@/types'
import {
  adminGetAllRelationships,
  adminDeleteRelationship
} from '@/services/apiService'
import RelationshipForm from './RelationshipForm.vue'

// 响应式数据
const loading = ref(false)
const relationships = ref<SeiyuuRelationship[]>([])
const searchQuery = ref('')
const showCreateForm = ref(false)
const editingRelationship = ref<SeiyuuRelationship | null>(null)

// 计算属性
const filteredRelationships = computed(() => {
  if (!searchQuery.value) {
    return relationships.value
  }

  const query = searchQuery.value.toLowerCase()
  return relationships.value.filter(rel =>
    rel.seiyuu_name_a?.toLowerCase().includes(query) ||
    rel.seiyuu_name_b?.toLowerCase().includes(query) ||
    rel.relationship_description.toLowerCase().includes(query)
  )
})

// 方法
async function loadRelationships() {
  loading.value = true
  try {
    const response = await adminGetAllRelationships()
    if (response.success) {
      relationships.value = response.data
    }
  } catch (error: any) {
    console.error('获取关系列表失败:', error.message)
  } finally {
    loading.value = false
  }
}

function editRelationship(relationship: SeiyuuRelationship) {
  editingRelationship.value = { ...relationship }
  showCreateForm.value = false
}

async function deleteRelationship(id: string) {
  if (!confirm('确定要删除这个关系吗？')) {
    return
  }

  try {
    const response = await adminDeleteRelationship(id)
    if (response.success) {
      await loadRelationships()
    }
  } catch (error: any) {
    console.error('删除关系失败:', error.message)
    alert(`删除失败：${error.message}`)
  }
}

function closeForm() {
  showCreateForm.value = false
  editingRelationship.value = null
}

function onFormSuccess() {
  closeForm()
  loadRelationships()
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 生命周期
onMounted(() => {
  loadRelationships()
})
</script>

<style scoped>
.relationship-manager {
  padding: var(--spacing-xl);
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xl);
}

.manager-header h2 {
  font-size: var(--font-size-xl);
  margin: 0;
}

.search-section {
  margin-bottom: var(--spacing-lg);
}

.search-controls {
  display: flex;
  gap: var(--spacing-md);
  align-items: center;
}

.search-input {
  flex: 1;
  max-width: 300px;
}

.loading, .empty-state {
  text-align: center;
  padding: var(--spacing-3xl);
  color: var(--text-muted);
}

.relationships-grid {
  display: grid;
  gap: var(--spacing-lg);
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
}

.relationship-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg);
}

.relationship-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-md);
}

.seiyuu-pair {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  flex: 1;
}

.seiyuu-name {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.relationship-arrow {
  color: var(--text-muted);
  font-size: var(--font-size-lg);
}

.relationship-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.relationship-description {
  background: var(--bg-tertiary);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  margin-bottom: var(--spacing-md);
  line-height: 1.5;
  color: var(--text-primary);
  border-left: 3px solid var(--color-primary);
}

.relationship-meta {
  display: flex;
  gap: var(--spacing-lg);
  font-size: var(--font-size-xs);
  color: var(--text-muted);
}

.meta-item {
  white-space: nowrap;
}

.btn-small {
  padding: var(--spacing-xs) var(--spacing-sm);
  font-size: var(--font-size-xs);
}

.btn-danger {
  background-color: #ef4444;
  color: white;
  border: 1px solid #dc2626;
}

.btn-danger:hover {
  background-color: #dc2626;
}

@media (max-width: 768px) {
  .relationships-grid {
    grid-template-columns: 1fr;
  }

  .relationship-header {
    flex-direction: column;
    gap: var(--spacing-md);
    align-items: flex-start;
  }

  .relationship-actions {
    align-self: flex-end;
  }

  .relationship-meta {
    flex-direction: column;
    gap: var(--spacing-xs);
  }
}
</style>