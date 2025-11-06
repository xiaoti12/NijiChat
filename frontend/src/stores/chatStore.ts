/**
 * 聊天状态管理
 * 封装IndexedDB操作，管理聊天记录和会话
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message, Room, Conversation, ChatSettings } from '@/types'
import { generateId } from '@/utils/crypto'

export const useChatStore = defineStore('chat', () => {
  // 状态
  const rooms = ref<Room[]>([])
  const messages = ref<Map<string, Message[]>>(new Map())
  const currentRoomId = ref<string | null>(null)
  const chatSettings = ref<ChatSettings>(loadChatSettings())

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
      }
    }
  }

  // 清空房间消息
  function clearRoomMessages(roomId: string) {
    messages.value.set(roomId, [])
    updateRoom(roomId, { last_message: undefined })
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

  // === 持久化（后续集成IndexedDB） ===

  // TODO: 实现IndexedDB持久化
  // - 消息存储和检索
  // - 房间信息存储
  // - 批量操作优化

  // TODO: 实现WebDAV同步
  // - 数据备份
  // - 跨设备同步

  return {
    // State
    rooms,
    messages,
    currentRoomId,
    chatSettings,

    // Computed
    currentRoom,
    currentMessages,
    conversations,

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
    searchConversations
  }
})
