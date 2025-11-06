/**
 * 声优数据状态管理
 * 管理声优列表和群组信息
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Seiyuu, SeiyuuGroup } from '@/types'

export const useSeiyuuStore = defineStore('seiyuu', () => {
  // 状态
  const seiyuuList = ref<Seiyuu[]>([])
  const groups = ref<SeiyuuGroup[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const activeSeiyuu = computed(() => {
    return seiyuuList.value.filter(s => s.status === 'active')
  })

  const seiyuuCount = computed(() => seiyuuList.value.length)
  const activeCount = computed(() => activeSeiyuu.value.length)

  // Actions

  // 设置声优列表
  function setSeiyuuList(list: Seiyuu[]) {
    seiyuuList.value = list
  }

  // 添加声优
  function addSeiyuu(seiyuu: Seiyuu) {
    seiyuuList.value.push(seiyuu)
  }

  // 更新声优
  function updateSeiyuu(id: string, updates: Partial<Seiyuu>) {
    const index = seiyuuList.value.findIndex(s => s.id === id)
    if (index !== -1) {
      seiyuuList.value[index] = {
        ...seiyuuList.value[index],
        ...updates,
        updated_at: new Date().toISOString()
      }
    }
  }

  // 删除声优
  function deleteSeiyuu(id: string) {
    const index = seiyuuList.value.findIndex(s => s.id === id)
    if (index !== -1) {
      seiyuuList.value.splice(index, 1)
    }
  }

  // 根据ID获取声优
  function getSeiyuuById(id: string): Seiyuu | undefined {
    return seiyuuList.value.find(s => s.id === id)
  }

  // 根据ID列表获取声优列表
  function getSeiyuuByIds(ids: string[]): Seiyuu[] {
    return seiyuuList.value.filter(s => ids.includes(s.id))
  }

  // 搜索声优
  function searchSeiyuu(keyword: string): Seiyuu[] {
    const lowerKeyword = keyword.toLowerCase()
    return seiyuuList.value.filter(s =>
      s.name.toLowerCase().includes(lowerKeyword) ||
      s.tags.some(tag => tag.toLowerCase().includes(lowerKeyword))
    )
  }

  // 按标签过滤声优
  function filterByTag(tag: string): Seiyuu[] {
    return seiyuuList.value.filter(s => s.tags.includes(tag))
  }

  // 获取所有标签
  const allTags = computed(() => {
    const tags = new Set<string>()
    seiyuuList.value.forEach(s => {
      s.tags.forEach(tag => tags.add(tag))
    })
    return Array.from(tags)
  })

  // === 群组管理 ===

  // 设置群组列表
  function setGroups(groupList: SeiyuuGroup[]) {
    groups.value = groupList
  }

  // 添加群组
  function addGroup(group: SeiyuuGroup) {
    groups.value.push(group)
  }

  // 更新群组
  function updateGroup(id: string, updates: Partial<SeiyuuGroup>) {
    const index = groups.value.findIndex(g => g.id === id)
    if (index !== -1) {
      groups.value[index] = {
        ...groups.value[index],
        ...updates,
        updated_at: new Date().toISOString()
      }
    }
  }

  // 删除群组
  function deleteGroup(id: string) {
    const index = groups.value.findIndex(g => g.id === id)
    if (index !== -1) {
      groups.value.splice(index, 1)
    }
  }

  // 根据ID获取群组
  function getGroupById(id: string): SeiyuuGroup | undefined {
    return groups.value.find(g => g.id === id)
  }

  // 获取群组成员
  function getGroupMembers(groupId: string): Seiyuu[] {
    const group = getGroupById(groupId)
    if (!group) return []
    return getSeiyuuByIds(group.member_ids)
  }

  // 设置加载状态
  function setLoading(value: boolean) {
    loading.value = value
  }

  // 设置错误信息
  function setError(message: string | null) {
    error.value = message
  }

  // 清空数据
  function clear() {
    seiyuuList.value = []
    groups.value = []
    error.value = null
  }

  return {
    // State
    seiyuuList,
    groups,
    loading,
    error,

    // Computed
    activeSeiyuu,
    seiyuuCount,
    activeCount,
    allTags,

    // Actions - 声优
    setSeiyuuList,
    addSeiyuu,
    updateSeiyuu,
    deleteSeiyuu,
    getSeiyuuById,
    getSeiyuuByIds,
    searchSeiyuu,
    filterByTag,

    // Actions - 群组
    setGroups,
    addGroup,
    updateGroup,
    deleteGroup,
    getGroupById,
    getGroupMembers,

    // Actions - 通用
    setLoading,
    setError,
    clear
  }
})
