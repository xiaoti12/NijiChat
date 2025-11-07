/**
 * 聊天状态管理
 * 封装IndexedDB操作，管理聊天记录和会话
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message, Room, Conversation, ChatSettings } from '@/types'
import { generateId } from '@/utils/crypto'
import { chatPersistenceService } from '@/services/chatPersistenceService'

export const useChatStore = defineStore('chat', () => {
  // 状态
  const rooms = ref<Room[]>([])
  const messages = ref<Map<string, Message[]>>(new Map())
  const currentRoomId = ref<string | null>(null)
  const chatSettings = ref<ChatSettings>(loadChatSettings())
  const isLoaded = ref(false) // 标记是否已从本地存储加载数据

  // Computed
  const currentRoom = computed(() => {
    return rooms.value.find(r => r.id === currentRoomId.value)
  })

  const currentMessages = computed(() => {
    if (!currentRoomId.value) return []
    return messages.value.get(currentRoomId.value) || []
  })

  // 加载聊天设置
  function loadChatSettings(): ChatSettings {
    try {
      const stored = localStorage.getItem('nijichat_settings')
      if (stored) {
        return JSON.parse(stored)
      }
    } catch (error) {
      console.error('加载聊天设置失败:', error)
    }
    return {
      messageCount: 10,
      temperature: 0.7,
      maxTokens: 2000,
      autoScroll: true,
      showTimestamp: true
    }
  }

  // 保存聊天设置
  function saveChatSettings() {
    try {
      localStorage.setItem('nijichat_settings', JSON.stringify(chatSettings.value))
    } catch (error) {
      console.error('保存聊天设置失败:', error)
    }
  }

  // 更新聊天设置
  function updateChatSettings(updates: Partial<ChatSettings>) {
    chatSettings.value = { ...chatSettings.value, ...updates }
    saveChatSettings()
  }

  // === 数据初始化加载 ===

  /**
   * 从本地存储加载所有聊天数据
   */
  async function loadFromStorage(): Promise<void> {
    if (isLoaded.value) return // 避免重复加载

    try {
      console.log('🔄 正在从本地存储加载聊天数据...')

      // 并行加载房间和消息数据
      const [storedRooms, storedMessages] = await Promise.all([
        chatPersistenceService.getRooms(),
        chatPersistenceService.getAllMessages()
      ])

      // 更新状态
      rooms.value = storedRooms
      messages.value = storedMessages

      console.log(`✅ 成功加载 ${storedRooms.length} 个房间和 ${storedMessages.size} 个消息集合`)

      // 数据迁移
      await chatPersistenceService.migrateData()

      isLoaded.value = true
    } catch (error) {
      console.error('❌ 从本地存储加载数据失败:', error)
      // 加载失败时保持默认状态，不抛出异常
      isLoaded.value = true
    }
  }

  // === 房间管理 ===

  // 创建房间
  function createRoom(room: Omit<Room, 'id' | 'created_at'>): Room {
    const newRoom: Room = {
      id: generateId(),
      created_at: Date.now(),
      ...room
    }
    rooms.value.push(newRoom)
    messages.value.set(newRoom.id, [])

    // 持久化保存
    chatPersistenceService.addRoom(newRoom).catch(error => {
      console.error('保存新房间到本地存储失败:', error)
    })

    return newRoom
  }

  // 获取房间
  function getRoom(roomId: string): Room | undefined {
    return rooms.value.find(r => r.id === roomId)
  }

  // 删除房间
  function deleteRoom(roomId: string) {
    const index = rooms.value.findIndex(r => r.id === roomId)
    if (index !== -1) {
      rooms.value.splice(index, 1)
      messages.value.delete(roomId)
      if (currentRoomId.value === roomId) {
        currentRoomId.value = null
      }

      // 持久化删除
      chatPersistenceService.deleteRoom(roomId).catch(error => {
        console.error('从本地存储删除房间失败:', error)
      })
    }
  }

  // 设置当前房间
  function setCurrentRoom(roomId: string) {
    currentRoomId.value = roomId
    // 清空未读计数
    const room = getRoom(roomId)
    if (room) {
      room.unread_count = 0
    }
  }

  // 更新房间信息
  function updateRoom(roomId: string, updates: Partial<Room>) {
    const room = getRoom(roomId)
    if (room) {
      Object.assign(room, updates)

      // 持久化更新
      chatPersistenceService.updateRoom(roomId, updates).catch(error => {
        console.error('更新房间到本地存储失败:', error)
      })
    }
  }

  // === 消息管理 ===

  // 添加消息
  function addMessage(message: Omit<Message, 'id' | 'timestamp'>): Message {
    const newMessage: Message = {
      id: generateId(),
      timestamp: Date.now(),
      ...message
    }

    const roomMessages = messages.value.get(message.room_id) || []
    roomMessages.push(newMessage)
    messages.value.set(message.room_id, roomMessages)

    // 更新房间最后一条消息
    updateRoom(message.room_id, {
      last_message: newMessage,
      unread_count: currentRoomId.value === message.room_id ? 0 : (getRoom(message.room_id)?.unread_count || 0) + 1
    })

    // 持久化保存消息
    chatPersistenceService.addMessageToRoom(message.room_id, newMessage).catch(error => {
      console.error('保存消息到本地存储失败:', error)
    })

    return newMessage
  }

  // 获取房间消息列表
  function getRoomMessages(roomId: string, limit?: number): Message[] {
    const roomMessages = messages.value.get(roomId) || []
    if (limit) {
      return roomMessages.slice(-limit)
    }
    return roomMessages
  }

  // 获取最近的消息（用于AI上下文）
  function getRecentMessages(roomId: string, count: number): Message[] {
    const roomMessages = messages.value.get(roomId) || []
    return roomMessages.slice(-count)
  }

  // 删除消息
  function deleteMessage(roomId: string, messageId: string) {
    const roomMessages = messages.value.get(roomId)
    if (roomMessages) {
      const index = roomMessages.findIndex(m => m.id === messageId)
      if (index !== -1) {
        roomMessages.splice(index, 1)

        // 持久化删除
        chatPersistenceService.deleteMessageFromRoom(roomId, messageId).catch(error => {
          console.error('从本地存储删除消息失败:', error)
        })
      }
    }
  }

  // 清空房间消息
  function clearRoomMessages(roomId: string) {
    messages.value.set(roomId, [])
    updateRoom(roomId, { last_message: undefined })

    // 持久化清空
    chatPersistenceService.deleteRoomMessages(roomId).catch(error => {
      console.error('清空房间消息从本地存储失败:', error)
    })
  }

  // === 对话列表（用于侧边栏显示） ===
  const conversations = computed((): Conversation[] => {
    return rooms.value.map(room => ({
      id: room.id,
      type: room.type,
      seiyuu: room.participants.length === 1 ? {
        id: room.participants[0],
        name: room.name,
        avatar: room.avatar
      } : undefined,
      lastMessage: room.last_message?.content || '',
      timestamp: room.last_message ? new Date(room.last_message.timestamp).toISOString() : new Date(room.created_at).toISOString(),
      unread: room.unread_count
    }))
  })

  // 搜索对话
  function searchConversations(keyword: string): Conversation[] {
    const lowerKeyword = keyword.toLowerCase()
    return conversations.value.filter(conv =>
      conv.seiyuu?.name.toLowerCase().includes(lowerKeyword) ||
      conv.lastMessage.toLowerCase().includes(lowerKeyword)
    )
  }

  // === 数据管理 ===

  /**
   * 清空所有聊天数据
   */
  async function clearAllData(): Promise<void> {
    try {
      // 清空内存状态
      rooms.value = []
      messages.value.clear()
      currentRoomId.value = null

      // 清空本地存储
      await chatPersistenceService.clearAllData()

      console.log('✅ 所有聊天数据已清空')
    } catch (error) {
      console.error('❌ 清空聊天数据失败:', error)
      throw error
    }
  }

  /**
   * 获取存储使用情况
   */
  async function getStorageInfo() {
    try {
      return await chatPersistenceService.getStorageInfo()
    } catch (error) {
      console.error('获取存储信息失败:', error)
      return {
        roomsCount: 0,
        totalMessages: 0,
        storageSize: 0
      }
    }
  }

  /**
   * 验证数据完整性
   * 检查房间和消息数据是否一致
   */
  function validateDataIntegrity(): {
    isValid: boolean
    issues: string[]
    stats: {
      totalRooms: number
      totalMessages: number
      orphanedMessages: number
      emptyRooms: number
    }
  } {
    const issues: string[] = []
    let orphanedMessages = 0
    let emptyRooms = 0

    // 检查消息是否有对应的房间
    for (const [roomId, roomMessages] of messages.value.entries()) {
      const room = rooms.value.find(r => r.id === roomId)
      if (!room) {
        issues.push(`发现孤立消息集合: ${roomId} (${roomMessages.length} 条消息)`)
        orphanedMessages += roomMessages.length
      }
    }

    // 检查房间是否有消息
    for (const room of rooms.value) {
      const roomMessages = messages.value.get(room.id)
      if (!roomMessages || roomMessages.length === 0) {
        emptyRooms++
      }

      // 检查房间的最后一条消息是否存在
      if (room.last_message && roomMessages) {
        const lastMessageExists = roomMessages.some(m => m.id === room.last_message?.id)
        if (!lastMessageExists) {
          issues.push(`房间 ${room.name} 的最后一条消息引用无效`)
        }
      }
    }

    const totalMessages = Array.from(messages.value.values()).reduce(
      (sum, msgs) => sum + msgs.length, 0
    )

    return {
      isValid: issues.length === 0,
      issues,
      stats: {
        totalRooms: rooms.value.length,
        totalMessages,
        orphanedMessages,
        emptyRooms
      }
    }
  }

  /**
   * 修复数据完整性问题
   */
  async function repairDataIntegrity(): Promise<void> {
    const validation = validateDataIntegrity()

    if (validation.isValid) {
      console.log('✅ 数据完整性检查通过，无需修复')
      return
    }

    console.log('🔧 开始修复数据完整性问题...')

    // 清理孤立的消息集合
    for (const [roomId] of messages.value.entries()) {
      const room = rooms.value.find(r => r.id === roomId)
      if (!room) {
        console.log(`🗑️ 删除孤立消息集合: ${roomId}`)
        messages.value.delete(roomId)
        await chatPersistenceService.deleteRoomMessages(roomId)
      }
    }

    // 确保所有房间都有消息集合
    for (const room of rooms.value) {
      if (!messages.value.has(room.id)) {
        console.log(`📝 为房间创建空消息集合: ${room.name}`)
        messages.value.set(room.id, [])
      }
    }

    // 重新同步到持久化存储
    await chatPersistenceService.saveRooms(rooms.value)

    console.log('✅ 数据完整性修复完成')
  }

  return {
    // State
    rooms,
    messages,
    currentRoomId,
    chatSettings,
    isLoaded,

    // Computed
    currentRoom,
    currentMessages,
    conversations,

    // Actions - 初始化
    loadFromStorage,

    // Actions - 设置
    updateChatSettings,

    // Actions - 房间
    createRoom,
    getRoom,
    deleteRoom,
    setCurrentRoom,
    updateRoom,

    // Actions - 消息
    addMessage,
    getRoomMessages,
    getRecentMessages,
    deleteMessage,
    clearRoomMessages,

    // Actions - 搜索
    searchConversations,

    // Actions - 数据管理
    clearAllData,
    getStorageInfo,
    validateDataIntegrity,
    repairDataIntegrity
  }
})
