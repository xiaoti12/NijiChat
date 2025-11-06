/**
 * Vue应用入口文件
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 导入全局样式
import '@/styles/global.css'

// 创建Vue应用实例
const app = createApp(App)

// 使用Pinia状态管理
const pinia = createPinia()
app.use(pinia)

// 使用Vue Router
app.use(router)

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误:', err, info)
  // TODO: 可以添加错误上报服务
}

// 全局警告处理（仅开发环境）
app.config.warnHandler = (msg, instance, trace) => {
  console.warn('全局警告:', msg, trace)
}

// 挂载应用
app.mount('#app')

// 开发环境输出信息
if (import.meta.env.DEV) {
  console.log('🎨 NijiChat 开发模式')
  console.log('📦 Vue版本:', app.version)
}
