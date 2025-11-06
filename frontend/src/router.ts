/**
 * Vue Router 路由配置
 */

import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAdminStore } from '@/stores/adminStore'

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/ChatHome.vue'),
    meta: {
      title: 'NijiChat - 聊天'
    }
  },
  {
    path: '/seiyuu',
    name: 'SeiyuuLibrary',
    component: () => import('@/pages/SeiyuuLibrary.vue'),
    meta: {
      title: 'NijiChat - 声优库'
    }
  },
  {
    path: '/chat/:seiyuuId',
    name: 'Chat1v1',
    component: () => import('@/pages/ChatPage.vue'),
    meta: {
      title: '聊天'
    }
  },
  {
    path: '/dual-theater',
    name: 'DualTheater',
    component: () => import('@/pages/DualTheater.vue'),
    meta: {
      title: '双声优剧场'
    }
  },
  {
    path: '/group-theater/:groupId?',
    name: 'GroupTheater',
    component: () => import('@/pages/GroupTheater.vue'),
    meta: {
      title: '群组剧场'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/pages/Settings.vue'),
    meta: {
      title: '设置'
    }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/pages/Admin.vue'),
    meta: {
      title: '管理后台',
      requiresAuth: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFound.vue'),
    meta: {
      title: '页面不存在'
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  const title = to.meta.title as string
  if (title) {
    document.title = title
  }

  // 检查是否需要管理员权限
  if (to.meta.requiresAuth) {
    // 开发模式下跳过登录校验
    if (import.meta.env.DEV || import.meta.env.MODE === 'development') {
      console.log('开发模式：跳过管理员权限校验')
      next()
      return
    }

    const adminStore = useAdminStore()
    if (!adminStore.isLoggedIn) {
      // 未登录，跳转到首页
      console.warn('需要管理员权限')
      next({ name: 'Home' })
      return
    }
  }

  next()
})

// 全局后置钩子
router.afterEach((to, from) => {
  // 可以在这里添加页面访问统计等
  console.log(`路由跳转: ${from.path} -> ${to.path}`)
})

export default router
