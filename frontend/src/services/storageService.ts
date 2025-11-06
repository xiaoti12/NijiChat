/**
 * 本地存储服务
 * 封装IndexedDB操作
 */

import localforage from 'localforage'
import type { Message, Room } from '@/types'

// 配置localforage
const chatDB = localforage.createInstance({
  name: 'NijiChat',
  storeName: 'chat_records',
  description: '聊天记录数据库'
})

const roomDB = localforage.createInstance({
  name: 'NijiChat',
  storeName: 'rooms',
  description: '房间信息数据库'
})

/**
 * 本地存储服务类
 */
export class StorageService {
  // ========== 消息存储 ==========

  /**
   * 保存消息
   */
  async saveMessage(message: Message): Promise<void> {
    const key = `${message.room_id}_${message.id}`
    await chatDB.setItem(key, message)
  }

  /**
   * 批量保存消息
   */
  async saveMessages(messages: Message[]): Promise<void> {
    const promises = messages.map(msg => this.saveMessage(msg))
    await Promise.all(promises)
  }

  /**
   * 获取房间所有消息
   */
  async getRoomMessages(roomId: string): Promise<Message[]> {
    const messages: Message[] = []
    await chatDB.iterate<Message, void>((value, key) => {
      if (key.startsWith(`${roomId}_`)) {
        messages.push(value)
      }
    })
    // 按时间戳排序
    return messages.sort((a, b) => a.timestamp - b.timestamp)
  }

  /**
   * 删除消息
   */
  async deleteMessage(roomId: string, messageId: string): Promise<void> {
    const key = `${roomId}_${messageId}`
    await chatDB.removeItem(key)
  }

  /**
   * 清空房间消息
   */
  async clearRoomMessages(roomId: string): Promise<void> {
    const keysToDelete: string[] = []
    await chatDB.iterate<Message, void>((value, key) => {
      if (key.startsWith(`${roomId}_`)) {
        keysToDelete.push(key)
      }
    })
    const promises = keysToDelete.map(key => chatDB.removeItem(key))
    await Promise.all(promises)
  }

  // ========== 房间存储 ==========

  /**
   * 保存房间信息
   */
  async saveRoom(room: Room): Promise<void> {
    await roomDB.setItem(room.id, room)
  }

  /**
   * 获取所有房间
   */
  async getAllRooms(): Promise<Room[]> {
    const rooms: Room[] = []
    await roomDB.iterate<Room, void>((value) => {
      rooms.push(value)
    })
    // 按创建时间倒序排序
    return rooms.sort((a, b) => b.created_at - a.created_at)
  }

  /**
   * 获取单个房间
   */
  async getRoom(roomId: string): Promise<Room | null> {
    return await roomDB.getItem<Room>(roomId)
  }

  /**
   * 删除房间
   */
  async deleteRoom(roomId: string): Promise<void> {
    await roomDB.removeItem(roomId)
    // 同时删除该房间的所有消息
    await this.clearRoomMessages(roomId)
  }

  /**
   * 更新房间信息
   */
  async updateRoom(roomId: string, updates: Partial<Room>): Promise<void> {
    const room = await this.getRoom(roomId)
    if (room) {
      const updatedRoom = { ...room, ...updates }
      await this.saveRoom(updatedRoom)
    }
  }

  // ========== 数据库维护 ==========

  /**
   * 清空所有数据
   */
  async clearAll(): Promise<void> {
    await chatDB.clear()
    await roomDB.clear()
  }

  /**
   * 获取数据库大小（估算）
   */
  async getDatabaseSize(): Promise<{ messages: number; rooms: number }> {
    const messageCount = await chatDB.length()
    const roomCount = await roomDB.length()
    return {
      messages: messageCount,
      rooms: roomCount
    }
  }

  /**
   * 导出数据（用于备份）
   */
  async exportData(): Promise<{ messages: Message[]; rooms: Room[] }> {
    const messages: Message[] = []
    const rooms: Room[] = []

    await chatDB.iterate<Message, void>((value) => {
      messages.push(value)
    })

    await roomDB.iterate<Room, void>((value) => {
      rooms.push(value)
    })

    return { messages, rooms }
  }

  /**
   * 导入数据（用于恢复）
   */
  async importData(data: { messages: Message[]; rooms: Room[] }): Promise<void> {
    // 清空现有数据
    await this.clearAll()

    // 导入新数据
    await this.saveMessages(data.messages)
    const roomPromises = data.rooms.map(room => this.saveRoom(room))
    await Promise.all(roomPromises)
  }
}

// 导出单例
export const storageService = new StorageService()

// 导出composable
export function useStorageService() {
  return storageService
}
