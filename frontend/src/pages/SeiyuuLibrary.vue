<template>
  <div class="seiyuu-library">
    <div class="library-header">
      <div class="header-nav">
        <button @click="goBack" class="btn btn-secondary">
          ← 返回
        </button>
      </div>
      <h1 class="library-title">声优库</h1>
      <p class="library-subtitle">选择你喜欢的声优开始对话</p>
    </div>

    <div class="library-content">
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <p>{{ error }}</p>
        <button @click="loadSeiyuu" class="btn btn-primary">重试</button>
      </div>

      <div v-else class="seiyuu-grid">
        <div v-for="seiyuu in activeSeiyuu" :key="seiyuu.id" class="seiyuu-card" @click="openChat(seiyuu.id)">
          <div class="seiyuu-avatar">
            <img :src="seiyuu.avatar_url || '/default-avatar.png'" :alt="seiyuu.name" />
          </div>
          <div class="seiyuu-info">
            <h3 class="seiyuu-name">{{ seiyuu.name }}</h3>
            <div class="seiyuu-tags">
              <span v-for="tag in seiyuu.tags.slice(0, 3)" :key="tag" class="tag">{{ tag }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSeiyuuStore } from '@/stores/seiyuuStore'
import { getSeiyuuList } from '@/services/apiService'

const router = useRouter()
const seiyuuStore = useSeiyuuStore()

const loading = ref(false)
const error = ref<string | null>(null)

const activeSeiyuu = computed(() => seiyuuStore.activeSeiyuu)

async function loadSeiyuu() {
  loading.value = true
  error.value = null

  try {
    const response = await getSeiyuuList()
    if (response.success && response.data) {
      seiyuuStore.setSeiyuuList(response.data)
    } else {
      error.value = response.message || '加载失败'
    }
  } catch (err: any) {
    error.value = err.message || '网络错误'
  } finally {
    loading.value = false
  }
}

function openChat(seiyuuId: string) {
  router.push({ name: 'Home', query: { seiyuuId } })
}

function goBack() {
  router.push({ name: 'Home' })
}

onMounted(() => {
  loadSeiyuu()
})
</script>

<style scoped>
.seiyuu-library {
  width: 100%;
  height: 100vh;
  padding: var(--spacing-3xl);
  background: var(--bg-secondary);
  overflow-y: auto;
}

.library-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
  position: relative;
}

.header-nav {
  position: absolute;
  left: 0;
  top: 0;
}

.library-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: var(--spacing-md);
}

.library-subtitle {
  color: var(--text-muted);
  font-size: var(--font-size-lg);
}

.seiyuu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-xl);
  max-width: 1400px;
  margin: 0 auto;
}

.seiyuu-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  cursor: pointer;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-light);
}

.seiyuu-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-normal);
}

.seiyuu-avatar {
  width: 100%;
  aspect-ratio: 1;
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: var(--spacing-lg);
}

.seiyuu-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.seiyuu-name {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
}

.seiyuu-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
}

.tag {
  padding: 4px 8px;
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-size: var(--font-size-xs);
  border-radius: var(--radius-sm);
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: var(--spacing-lg);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--color-primary-light);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.btn {
  padding: var(--spacing-sm) var(--spacing-md);
  border: none;
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: var(--transition-normal);
  text-decoration: none;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-dark);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-quaternary);
}
</style>
