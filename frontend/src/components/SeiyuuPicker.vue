<template>
  <div class="seiyuu-picker">
    <!-- 搜索框 -->
    <div class="search-box">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索声优名称..."
        class="input search-input"
      />
    </div>

    <!-- 声优列表 -->
    <div class="seiyuu-list">
      <div v-if="loading" class="loading-state">
        加载中...
      </div>

      <div v-else-if="filteredSeiyuu.length === 0" class="empty-state">
        <p v-if="searchQuery">未找到匹配的声优</p>
        <p v-else>暂无可选择的声优</p>
      </div>

      <div
        v-for="seiyuu in filteredSeiyuu"
        :key="seiyuu.id"
        class="seiyuu-item"
        @click="selectSeiyuu(seiyuu)"
      >
        <img
          v-if="seiyuu.avatar_url"
          :src="seiyuu.avatar_url"
          :alt="seiyuu.name"
          class="seiyuu-avatar"
        />
        <div v-else class="seiyuu-avatar placeholder">
          {{ seiyuu.name.charAt(0) }}
        </div>

        <div class="seiyuu-info">
          <div class="seiyuu-name">{{ seiyuu.name }}</div>
          <div class="seiyuu-tags">
            <span
              v-for="tag in seiyuu.tags.slice(0, 3)"
              :key="tag"
              class="tag"
            >
              {{ tag }}
            </span>
            <span v-if="seiyuu.tags.length > 3" class="tag more">
              +{{ seiyuu.tags.length - 3 }}
            </span>
          </div>
        </div>

        <div class="select-indicator">
          <span class="select-icon">→</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Seiyuu } from '@/types'
import { getSeiyuuList } from '@/services/apiService'

// Props
interface Props {
  excludeIds?: string[]  // 排除的声优ID列表
}

const props = withDefaults(defineProps<Props>(), {
  excludeIds: () => []
})

// Emits
const emit = defineEmits<{
  select: [seiyuu: Seiyuu]
}>()

// 响应式数据
const seiyuuList = ref<Seiyuu[]>([])
const loading = ref(false)
const searchQuery = ref('')

// 计算属性
const filteredSeiyuu = computed(() => {
  let filtered = seiyuuList.value.filter(seiyuu => {
    // 过滤掉排除的声优
    if (props.excludeIds.includes(seiyuu.id)) {
      return false
    }

    // 只显示已发布的声优
    if (seiyuu.status !== 'active') {
      return false
    }

    return true
  })

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(seiyuu =>
      seiyuu.name.toLowerCase().includes(query) ||
      seiyuu.tags.some(tag => tag.toLowerCase().includes(query))
    )
  }

  return filtered
})

// 方法
async function loadSeiyuu() {
  loading.value = true
  try {
    const response = await getSeiyuuList()
    if (response.success) {
      seiyuuList.value = response.data
    }
  } catch (error: any) {
    console.error('获取声优列表失败:', error.message)
  } finally {
    loading.value = false
  }
}

function selectSeiyuu(seiyuu: Seiyuu) {
  emit('select', seiyuu)
}

// 生命周期
onMounted(() => {
  loadSeiyuu()
})
</script>

<style scoped>
.seiyuu-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.search-box {
  margin-bottom: var(--spacing-md);
}

.search-input {
  width: 100%;
  font-size: var(--font-size-sm);
}

.seiyuu-list {
  flex: 1;
  overflow-y: auto;
  max-height: 300px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: var(--spacing-xl);
  color: var(--text-muted);
}

.seiyuu-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.seiyuu-item:hover {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
}

.seiyuu-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.seiyuu-avatar.placeholder {
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-lg);
}

.seiyuu-info {
  flex: 1;
  min-width: 0;
}

.seiyuu-name {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.seiyuu-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.tag {
  background: var(--bg-secondary);
  color: var(--text-muted);
  padding: 2px var(--spacing-xs);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  white-space: nowrap;
}

.tag.more {
  background: var(--color-primary);
  color: white;
}

.select-indicator {
  color: var(--text-muted);
  font-size: var(--font-size-lg);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.seiyuu-item:hover .select-indicator {
  opacity: 1;
}

.select-icon {
  font-weight: var(--font-weight-bold);
}

/* 滚动条样式 */
.seiyuu-list::-webkit-scrollbar {
  width: 6px;
}

.seiyuu-list::-webkit-scrollbar-track {
  background: var(--bg-secondary);
  border-radius: 3px;
}

.seiyuu-list::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.seiyuu-list::-webkit-scrollbar-thumb:hover {
  background: var(--text-muted);
}
</style>