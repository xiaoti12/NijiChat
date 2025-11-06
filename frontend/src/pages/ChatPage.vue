<template>
  <div class="chat-page">
    <!-- TODO: 集成实际的聊天组件 -->
    <div class="chat-placeholder">
      <h2>聊天页面</h2>
      <p>声优ID: {{ seiyuuId }}</p>
      <p>此页面将集成完整的聊天功能</p>
      <button @click="goBack" class="btn btn-secondary">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const seiyuuId = route.params.seiyuuId as string

function goBack() {
  router.push({ name: 'Home' })
}
</script>

<style scoped>
.chat-page {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
}

.chat-placeholder {
  text-align: center;
  padding: var(--spacing-3xl);
}

.chat-placeholder h2 {
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-lg);
}

.chat-placeholder p {
  color: var(--text-muted);
  margin-bottom: var(--spacing-md);
}
</style>
