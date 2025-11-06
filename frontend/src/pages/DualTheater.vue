<template>
  <div class="dual-theater">
    <div class="theater-placeholder">
      <h2>双声优剧场</h2>
      <p>选择两位声优，让他们进行对话</p>
      <button @click="goBack" class="btn btn-secondary">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

function goBack() {
  router.push({ name: 'Home' })
}
</script>

<style scoped>
.dual-theater {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-secondary);
}

.theater-placeholder {
  text-align: center;
  padding: var(--spacing-3xl);
}

.theater-placeholder h2 {
  font-size: var(--font-size-2xl);
  margin-bottom: var(--spacing-lg);
}

.theater-placeholder p {
  color: var(--text-muted);
  margin-bottom: var(--spacing-md);
}
</style>
