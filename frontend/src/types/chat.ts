/**
 * 聊天相关类型定义
 */

// 消息类型
export interface Message {
  id: string                    // 消息ID
  room_id: string               // 所属房间ID
  sender_id: string             // 发送者ID ('user-1' 或 声优ID)
  sender_name: string           // 发送者名称
  sender_avatar?: string        // 发送者头像
  content: string               // 消息内容
  timestamp: number             // 时间戳
  type?: MessageType            // 消息类型
}

// 消息类型枚举
export type MessageType = 'text' | 'image' | 'system'

// 对话房间
export interface Room {
  id: string                    // 房间ID
  type: RoomType                // 房间类型
  name: string                  // 房间名称
  avatar?: string               // 房间头像
  participants: string[]        // 参与者ID列表
  last_message?: Message        // 最后一条消息
  unread_count: number          // 未读消息数
  created_at: number            // 创建时间
  session_id?: string          // 会话ID，用于区分同一声优的不同对话会话
  session_name?: string        // 会话名称，用于用户识别不同会话
}

// 房间类型
export type RoomType = '1v1' | 'dual_theater' | 'group_theater'

// 对话（用于侧边栏）
export interface Conversation {
  id: string
  type: RoomType
  seiyuu?: {
    id: string
    name: string
    avatar?: string
  }
  lastMessage: string
  timestamp: string
  unread: number
  session_id?: string          // 会话ID，用于区分同一声优的不同对话会话
  session_name?: string        // 会话名称
}

// 双声优剧场
export interface DualTheater {
  id: string
  seiyuu1: {
    id: string
    name: string
    avatar?: string
  }
  seiyuu2: {
    id: string
    name: string
    avatar?: string
  }
  topic?: string                // 剧场主题
  status: TheaterStatus
  lastMessage: string
  timestamp: string
}

// 群组剧场
export interface GroupTheater {
  id: string
  group_id: string              // 关联的声优群组ID
  name: string
  participants: Array<{
    id: string
    name: string
    avatar?: string
  }>
  status: TheaterStatus
  lastMessage: string
  timestamp: string
}

// 剧场状态
export type TheaterStatus = 'active' | 'completed' | 'paused'

// 剧场消息（扩展Message）
export interface TheaterMessage extends Message {
  seiyuu_id: string             // 发言的声优ID
  selection_reason?: string     // AI选择该声优的原因
  confidence?: number           // 选择置信度
}

// 按声优分组的对话
export interface SeiyuuConversationGroup {
  seiyuuId: string              // 声优ID
  seiyuuName: string            // 声优名称
  seiyuuAvatar?: string         // 声优头像
  rooms: Room[]                 // 该声优的所有房间
  totalSessions: number          // 会话总数
  totalUnread: number           // 未读消息总数
  lastMessage?: Message          // 最后一条消息
  lastActive?: number            // 最后活跃时间
}

// 聊天设置
export interface ChatSettings {
  messageCount: number          // 消息数量
  temperature: number           // 温度参数
  maxTokens: number            // 最大token数
  autoScroll: boolean          // 自动滚动
  showTimestamp: boolean       // 显示时间戳
}
