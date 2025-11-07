/**
 * 聊天持久化服务
 * 使用 IndexedDB 存储聊天记录和房间信息
 */

import localforage from 'localforage'
import type { Message, Room } from '@/types'

// 配置存储实例
const messagesStore = localforage.createInstance({
  name: 'NijiChat',
  storeName: 'messages',
  version: 1.0,
  description: '聊天消息存储'
})

const roomsStore = localforage.createInstance({
  name: 'NijiChat',
  storeName: 'rooms',
  version: 1.0,
  description: '聊天房间存储'
})

export class ChatPersistenceService {
  // === 消息持久化 ===

  /**
   * 保存房间的所有消息
   */
  async saveRoomMessages(roomId: string, messages: Message[]): Promise<void> {
    try {
      await messagesStore.setItem(roomId, messages)
    } catch (error) {
      console.error(`保存房间 ${roomId} 消息失败:`, error)
      throw error
    }
  }

  /**
   * 获取房间的所有消息
   */
  async getRoomMessages(roomId: string): Promise<Message[]> {
    try {
      const messages = await messagesStore.getItem<Message[]>(roomId)
      return messages || []
    } catch (error) {
      console.error(`获取房间 ${roomId} 消息失败:`, error)
      return []
    }
  }

  /**
   * 添加单条消息到房间（增量保存）
   */
  async addMessageToRoom(roomId: string, message: Message): Promise<void> {
    try {
      const existingMessages = await this.getRoomMessages(roomId)
      existingMessages.push(message)
      await this.saveRoomMessages(roomId, existingMessages)
    } catch (error) {
      console.error(`添加消息到房间 ${roomId} 失败:`, error)
      throw error
    }
  }

  /**
   * 删除房间的消息
   */
  async deleteRoomMessages(roomId: string): Promise<void> {
    try {
      await messagesStore.removeItem(roomId)
    } catch (error) {
      console.error(`删除房间 ${roomId} 消息失败:`, error)
      throw error
    }
  }

  /**
   * 删除房间中的特定消息
   */
  async deleteMessageFromRoom(roomId: string, messageId: string): Promise<void> {
    try {
      const messages = await this.getRoomMessages(roomId)
      const updatedMessages = messages.filter(m => m.id !== messageId)
      await this.saveRoomMessages(roomId, updatedMessages)
    } catch (error) {
      console.error(`从房间 ${roomId} 删除消息 ${messageId} 失败:`, error)
      throw error
    }
  }

  /**
   * 获取所有房间的消息（用于初始化加载）
   */
  async getAllMessages(): Promise<Map<string, Message[]>> {
    try {
      const messagesMap = new Map<string, Message[]>()

      // 获取所有房间ID
      const roomIds = await messagesStore.keys()

      // 并行加载所有房间的消息
      const loadPromises = roomIds.map(async (roomId) => {
        const messages = await this.getRoomMessages(roomId)
        messagesMap.set(roomId, messages)
      })

      await Promise.all(loadPromises)
      return messagesMap
    } catch (error) {
      console.error('获取所有消息失败:', error)
      return new Map()
    }
  }

  // === 房间持久化 ===

  /**
   * 保存所有房间
   */
  async saveRooms(rooms: Room[]): Promise<void> {
    try {
      await roomsStore.setItem('rooms', rooms)
    } catch (error) {
      console.error('保存房间列表失败:', error)
      throw error
    }
  }

  /**
   * 获取所有房间
   */
  async getRooms(): Promise<Room[]> {
    try {
      const rooms = await roomsStore.getItem<Room[]>('rooms')
      return rooms || []
    } catch (error) {
      console.error('获取房间列表失败:', error)
      return []
    }
  }

  /**
   * 添加单个房间
   */
  async addRoom(room: Room): Promise<void> {
    try {
      const existingRooms = await this.getRooms()
      existingRooms.push(room)
      await this.saveRooms(existingRooms)
    } catch (error) {
      console.error(`添加房间 ${room.id} 失败:`, error)
      throw error
    }
  }

  /**
   * 更新房间信息
   */
  async updateRoom(roomId: string, updates: Partial<Room>): Promise<void> {
    try {
      const rooms = await this.getRooms()
      const roomIndex = rooms.findIndex(r => r.id === roomId)

      if (roomIndex !== -1) {
        rooms[roomIndex] = { ...rooms[roomIndex], ...updates }
        await this.saveRooms(rooms)
      }
    } catch (error) {
      console.error(`更新房间 ${roomId} 失败:`, error)
      throw error
    }
  }

  /**
   * 删除房间
   */
  async deleteRoom(roomId: string): Promise<void> {
    try {
      // 删除房间信息
      const rooms = await this.getRooms()
      const updatedRooms = rooms.filter(r => r.id !== roomId)
      await this.saveRooms(updatedRooms)

      // 删除房间的所有消息
      await this.deleteRoomMessages(roomId)
    } catch (error) {
      console.error(`删除房间 ${roomId} 失败:`, error)
      throw error
    }
  }

  // === 数据管理 ===

  /**
   * 清空所有聊天数据
   */
  async clearAllData(): Promise<void> {
    try {
      await Promise.all([
        messagesStore.clear(),
        roomsStore.clear()
      ])
      console.log('所有聊天数据已清空')
    } catch (error) {
      console.error('清空聊天数据失败:', error)
      throw error
    }
  }

  /**
   * 获取存储使用情况
   */
  async getStorageInfo(): Promise<{
    roomsCount: number
    totalMessages: number
    storageSize: number
  }> {
    try {
      const rooms = await this.getRooms()
      const allMessages = await this.getAllMessages()

      let totalMessages = 0
      for (const messages of allMessages.values()) {
        totalMessages += messages.length
      }

      // 估算存储大小（字节）
      const storageSize = JSON.stringify({
        rooms,
        messages: Object.fromEntries(allMessages)
      }).length

      return {
        roomsCount: rooms.length,
        totalMessages,
        storageSize
      }
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
   * 数据迁移和版本升级
   */
  async migrateData(): Promise<void> {
    try {
      // 检查是否存在旧版本数据需要迁移
      // 目前是首个版本，暂无迁移需求
      console.log('数据迁移检查完成')
    } catch (error) {
      console.error('数据迁移失败:', error)
      throw error
    }
  }
}

// 导出单例实例
export const chatPersistenceService = new ChatPersistenceService()