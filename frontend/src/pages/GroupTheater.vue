<template>
  <div class="group-theater">
    <div class="theater-placeholder">
      <h2>群组剧场</h2>
      <p v-if="groupId">群组ID: {{ groupId }}</p>
      <p>多位声优智能对话场景</p>
      <button @click="goBack" class="btn btn-secondary">返回</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const groupId = route.params.groupId as string | undefined

function goBack() {
  router.push({ name: 'Home' })
}
</script>

<style scoped>
.group-theater {
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
