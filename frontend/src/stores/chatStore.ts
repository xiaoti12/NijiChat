/**
 * 聊天状态管理
 * 封装IndexedDB操作，管理聊天记录和会话
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message, Room, Conversation, SeiyuuConversationGroup, DualConversationGroup, ChatSettings } from '@/types'
import { generateId } from '@/utils/crypto'
import { chatPersistenceService } from '@/services/chatPersistenceService'
import { useSeiyuuStore } from './seiyuuStore'

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
      unread: room.unread_count,
      session_id: room.session_id,
      session_name: room.session_name
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

  // === 分组对话（按声优） ===
  const groupedConversations = computed((): SeiyuuConversationGroup[] => {
    const groups = new Map<string, SeiyuuConversationGroup>()

    // 只处理1v1类型的房间
    const rooms1v1 = rooms.value.filter(room => room.type === '1v1')

    rooms1v1.forEach(room => {
      const seiyuuId = room.participants[0]

      if (!groups.has(seiyuuId)) {
        groups.set(seiyuuId, {
          seiyuuId,
          seiyuuName: room.name,
          seiyuuAvatar: room.avatar,
          rooms: [],
          totalSessions: 0,
          totalUnread: 0,
          lastActive: 0
        })
      }

      const group = groups.get(seiyuuId)!
      group.rooms.push(room)
      group.totalSessions++
      group.totalUnread += room.unread_count

      // 更新最后活跃时间
      if (room.last_message?.timestamp) {
        group.lastActive = Math.max(group.lastActive, room.last_message.timestamp)
      }
      if (room.created_at) {
        group.lastActive = Math.max(group.lastActive, room.created_at)
      }

      // 设置最后一条消息（取最新的）
      if (!group.lastMessage ||
          (room.last_message && group.lastMessage.timestamp < room.last_message.timestamp)) {
        group.lastMessage = room.last_message
      }
    })

    // 为每个组设置最后一条消息（如果没有的话，取房间创建时间）
    groups.forEach(group => {
      if (!group.lastMessage && group.rooms.length > 0) {
        const latestRoom = group.rooms.reduce((prev, current) =>
          current.created_at > prev.created_at ? current : prev
        )
        group.lastActive = latestRoom.created_at
      }
    })

    // 转换为数组并按最后活跃时间排序
    return Array.from(groups.values()).sort((a, b) => b.lastActive - a.lastActive)
  })

  // 搜索分组对话
  function searchGroupedConversations(keyword: string): SeiyuuConversationGroup[] {
    if (!keyword) return groupedConversations.value

    const lowerKeyword = keyword.toLowerCase()
    return groupedConversations.value.filter(group =>
      group.seiyuuName.toLowerCase().includes(lowerKeyword)
    )
  }

  // 获取指定声优的分组信息
  function getSeiyuuGroup(seiyuuId: string): SeiyuuConversationGroup | undefined {
    return groupedConversations.value.find(group => group.seiyuuId === seiyuuId)
  }

  // === 双人对话分组 ===
  const dualGroupedConversations = computed((): DualConversationGroup[] => {
    const groups = new Map<string, DualConversationGroup>()

    // 只处理双人剧场类型的房间
    const dualRooms = rooms.value.filter(room => room.type === 'dual_theater' && room.participants.length === 2)

    dualRooms.forEach(room => {
      const [seiyuu1Id, seiyuu2Id] = room.participants
      // 创建声优对ID，保持一致的顺序 (较小的ID在前)
      const pairId = [seiyuu1Id, seiyuu2Id].sort().join('_')

      if (!groups.has(pairId)) {
        // 从room信息中获取声优基本信息
        const seiyuu1Info = getSeiyuuInfoFromRoom(room, seiyuu1Id)
        const seiyuu2Info = getSeiyuuInfoFromRoom(room, seiyuu2Id)

        groups.set(pairId, {
          pairId,
          seiyuu1: seiyuu1Info,
          seiyuu2: seiyuu2Info,
          rooms: [],
          totalSessions: 0,
          totalUnread: 0,
          lastActive: 0
        })
      }

      const group = groups.get(pairId)!
      group.rooms.push(room)
      group.totalSessions++
      group.totalUnread += room.unread_count

      // 更新最后活跃时间
      if (room.last_message?.timestamp) {
        group.lastActive = Math.max(group.lastActive, room.last_message.timestamp)
      }
      if (room.created_at) {
        group.lastActive = Math.max(group.lastActive, room.created_at)
      }

      // 设置最后一条消息（取最新的）
      if (!group.lastMessage ||
          (room.last_message && group.lastMessage.timestamp < room.last_message.timestamp)) {
        group.lastMessage = room.last_message
      }
    })

    // 为每个组设置最后一条消息（如果没有的话，取房间创建时间）
    groups.forEach(group => {
      if (!group.lastMessage && group.rooms.length > 0) {
        const latestRoom = group.rooms.reduce((prev, current) =>
          current.created_at > prev.created_at ? current : prev
        )
        group.lastActive = latestRoom.created_at
      }
    })

    // 转换为数组并按最后活跃时间排序
    return Array.from(groups.values()).sort((a, b) => b.lastActive - a.lastActive)
  })

  // 辅助函数：从房间信息中提取声优信息
  function getSeiyuuInfoFromRoom(room: Room, seiyuuId: string) {
    const seiyuuStore = useSeiyuuStore()

    // 从声优store获取详细信息
    const seiyuu = seiyuuStore.getSeiyuuById(seiyuuId)

    if (seiyuu) {
      return {
        id: seiyuu.id,
        name: seiyuu.name,
        avatar: seiyuu.avatar_url
      }
    }

    // 如果在声优store中找不到，返回基本信息作为备选
    return {
      id: seiyuuId,
      name: '未知声优',
      avatar: undefined
    }
  }

  // 获取双人对话分组信息
  function getDualGroup(pairId: string): DualConversationGroup | undefined {
    return dualGroupedConversations.value.find(group => group.pairId === pairId)
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

  // === 会话管理 ===

  /**
   * 为声优创建新对话会话
   */
  function createNewSession(seiyuuId: string, seiyuuName: string, seiyuuAvatar?: string): Room {
    const sessionId = generateId()
    const now = Date.now()

    // 生成默认会话名称
    const sessionCount = rooms.value.filter(room =>
      room.type === '1v1' &&
      room.participants.includes(seiyuuId)
    ).length

    const sessionName = `对话 ${sessionCount + 1}`

    const newRoom: Room = {
      id: generateId(),
      type: '1v1',
      name: seiyuuName,
      avatar: seiyuuAvatar,
      participants: [seiyuuId],
      unread_count: 0,
      created_at: now,
      session_id: sessionId,
      session_name: sessionName
    }

    rooms.value.push(newRoom)
    messages.value.set(newRoom.id, [])

    // 持久化保存
    chatPersistenceService.addRoom(newRoom).catch(error => {
      console.error('保存新会话房间到本地存储失败:', error)
    })

    console.log(`✅ 为声优 ${seiyuuName} 创建新会话: ${sessionName}`)
    return newRoom
  }

  /**
   * 获取声优的所有会话
   */
  function getSessionsBySeiyuu(seiyuuId: string): Room[] {
    return rooms.value.filter(room =>
      room.type === '1v1' &&
      room.participants.includes(seiyuuId)
    ).sort((a, b) => b.created_at - a.created_at)
  }

  /**
   * 切换到指定会话
   */
  function switchToSession(roomId: string) {
    const room = getRoom(roomId)
    if (room) {
      setCurrentRoom(roomId)
      console.log(`🔄 切换到会话: ${room.session_name || room.name}`)
    }
  }

  /**
   * 更新会话名称
   */
  function updateSessionName(roomId: string, newName: string) {
    const room = getRoom(roomId)
    if (room) {
      room.session_name = newName
      updateRoom(roomId, { session_name: newName })
      console.log(`📝 会话名称已更新: ${newName}`)
    }
  }

  /**
   * 删除会话
   */
  function deleteSession(roomId: string) {
    const room = getRoom(roomId)
    if (room) {
      console.log(`🗑️ 删除会话: ${room.session_name || room.name}`)
      deleteRoom(roomId)
    }
  }

  /**
   * 获取会话的友好显示名称
   */
  function getSessionDisplayName(room: Room): string {
    return room.session_name || `对话 ${room.created_at}`
  }

  /**
   * 数据迁移：为现有房间分配session_id
   */
  function migrateExistingRoomsToSessions() {
    const roomsWithoutSession = rooms.value.filter(room =>
      !room.session_id && room.type === '1v1'
    )

    if (roomsWithoutSession.length > 0) {
      console.log(`🔄 正在为 ${roomsWithoutSession.length} 个现有房间分配会话ID...`)

      roomsWithoutSession.forEach(room => {
        room.session_id = generateId()
        room.session_name = '默认对话'

        // 持久化更新
        chatPersistenceService.updateRoom(room.id, {
          session_id: room.session_id,
          session_name: room.session_name
        }).catch(error => {
          console.error(`更新房间 ${room.id} 的会话信息失败:`, error)
        })
      })

      console.log(`✅ 已为 ${roomsWithoutSession.length} 个现有房间分配会话ID`)
    }
  }

  // === 双人对话管理 ===

  /**
   * 创建双人对话会话
   */
  function createDualSession(
    seiyuu1: { id: string; name: string; avatar?: string },
    seiyuu2: { id: string; name: string; avatar?: string },
    options?: {
      topic?: string
      initiatorId?: string
      relationship?: string
    }
  ): Room {
    const sessionId = generateId()
    const now = Date.now()

    // 生成会话名称
    const sessionCount = rooms.value.filter(room =>
      room.type === 'dual_theater' &&
      room.participants.includes(seiyuu1.id) &&
      room.participants.includes(seiyuu2.id)
    ).length

    const sessionName = options?.topic || `对话 ${sessionCount + 1}`
    const roomName = `${seiyuu1.name} × ${seiyuu2.name}`

    const newRoom: Room = {
      id: generateId(),
      type: 'dual_theater',
      name: roomName,
      participants: [seiyuu1.id, seiyuu2.id],
      unread_count: 0,
      created_at: now,
      session_id: sessionId,
      session_name: sessionName,
      // 双人对话专用字段
      dual_topic: options?.topic,
      dual_initiator_id: options?.initiatorId || seiyuu1.id,
      dual_relationship: options?.relationship
    }

    rooms.value.push(newRoom)
    messages.value.set(newRoom.id, [])

    // 持久化保存
    chatPersistenceService.addRoom(newRoom).catch(error => {
      console.error('保存双人对话房间到本地存储失败:', error)
    })

    console.log(`✅ 创建双人对话会话: ${roomName} - ${sessionName}`)
    return newRoom
  }

  /**
   * 获取双人对话会话列表
   */
  function getDualSessions(seiyuu1Id: string, seiyuu2Id: string): Room[] {
    return rooms.value.filter(room =>
      room.type === 'dual_theater' &&
      room.participants.includes(seiyuu1Id) &&
      room.participants.includes(seiyuu2Id)
    ).sort((a, b) => b.created_at - a.created_at)
  }

  /**
   * 查找或创建双人对话会话
   */
  function findOrCreateDualSession(
    seiyuu1: { id: string; name: string; avatar?: string },
    seiyuu2: { id: string; name: string; avatar?: string },
    options?: {
      topic?: string
      initiatorId?: string
      relationship?: string
    }
  ): Room {
    // 查找现有会话
    const existingSessions = getDualSessions(seiyuu1.id, seiyuu2.id)

    if (existingSessions.length > 0) {
      // 如果有现有会话，返回最新的一个
      return existingSessions[0]
    } else {
      // 创建新会话
      return createDualSession(seiyuu1, seiyuu2, options)
    }
  }

  /**
   * 切换到双人对话会话
   */
  function switchToDualSession(
    seiyuu1: { id: string; name: string; avatar?: string },
    seiyuu2: { id: string; name: string; avatar?: string },
    options?: {
      topic?: string
      initiatorId?: string
      relationship?: string
    }
  ): Room {
    const room = findOrCreateDualSession(seiyuu1, seiyuu2, options)
    setCurrentRoom(room.id)
    return room
  }

  /**
   * 获取当前房间的声优信息 (用于双人对话)
   */
  function getCurrentDualSeiyuu(): { seiyuu1: any; seiyuu2: any; initiator: any; responder: any } | null {
    const room = currentRoom.value
    if (!room || room.type !== 'dual_theater' || room.participants.length !== 2) {
      return null
    }

    const [seiyuu1Id, seiyuu2Id] = room.participants
    const initiatorId = room.dual_initiator_id || seiyuu1Id

    const seiyuuStore = useSeiyuuStore()

    // 从声优store获取详细信息
    const seiyuu1Data = seiyuuStore.getSeiyuuById(seiyuu1Id)
    const seiyuu2Data = seiyuuStore.getSeiyuuById(seiyuu2Id)

    const seiyuu1 = {
      id: seiyuu1Id,
      name: seiyuu1Data?.name || '未知声优',
      avatar: seiyuu1Data?.avatar_url
    }

    const seiyuu2 = {
      id: seiyuu2Id,
      name: seiyuu2Data?.name || '未知声优',
      avatar: seiyuu2Data?.avatar_url
    }

    return {
      seiyuu1,
      seiyuu2,
      initiator: initiatorId === seiyuu1Id ? seiyuu1 : seiyuu2,
      responder: initiatorId === seiyuu1Id ? seiyuu2 : seiyuu1
    }
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
    groupedConversations,
    dualGroupedConversations,

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
    searchGroupedConversations,
    getSeiyuuGroup,

    // Actions - 数据管理
    clearAllData,
    getStorageInfo,
    validateDataIntegrity,
    repairDataIntegrity,

    // Actions - 会话管理
    createNewSession,
    getSessionsBySeiyuu,
    switchToSession,
    updateSessionName,
    deleteSession,
    getSessionDisplayName,
    migrateExistingRoomsToSessions,

    // Actions - 双人对话
    createDualSession,
    getDualSessions,
    findOrCreateDualSession,
    switchToDualSession,
    getCurrentDualSeiyuu,
    getDualGroup
  }
})
