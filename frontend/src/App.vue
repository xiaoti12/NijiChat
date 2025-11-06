<template>
  <div id="app" class="app-container">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useConfigStore } from '@/stores/configStore'

const configStore = useConfigStore()

onMounted(() => {
  // 应用主题
  if (configStore.config.theme) {
    document.documentElement.setAttribute('data-theme', configStore.config.theme)
  }

  // 其他初始化操作
  console.log('NijiChat 应用已启动')
})
</script>

<style>
.app-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>
