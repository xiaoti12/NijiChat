/**
 * 日期处理工具函数
 */

/**
 * 格式化时间戳为可读格式
 * @param timestamp 时间戳（毫秒）
 * @param format 格式：'full' | 'date' | 'time' | 'relative'
 */
export function formatTimestamp(timestamp: number, format: 'full' | 'date' | 'time' | 'relative' = 'relative'): string {
  const date = new Date(timestamp)
  const now = new Date()

  if (format === 'relative') {
    return formatRelativeTime(timestamp)
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  switch (format) {
    case 'full':
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    case 'date':
      return `${year}-${month}-${day}`
    case 'time':
      return `${hours}:${minutes}`
    default:
      return date.toLocaleString('zh-CN')
  }
}

/**
 * 格式化相对时间（如"刚刚"、"5分钟前"）
 * @param timestamp 时间戳（毫秒）
 */
export function formatRelativeTime(timestamp: number): string {
  const now = Date.now()
  const diff = now - timestamp

  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const weeks = Math.floor(days / 7)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)

  if (seconds < 60) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  if (weeks < 4) return `${weeks}周前`
  if (months < 12) return `${months}个月前`
  return `${years}年前`
}

/**
 * 格式化消息时间（智能显示）
 * 今天：显示时间
 * 昨天：显示"昨天"
 * 本周：显示星期
 * 更早：显示日期
 */
export function formatMessageTime(timestamp: number): string {
  const date = new Date(timestamp)
  const now = new Date()

  const isToday = date.toDateString() === now.toDateString()
  const isYesterday = new Date(now.setDate(now.getDate() - 1)).toDateString() === date.toDateString()

  if (isToday) {
    return formatTimestamp(timestamp, 'time')
  }

  if (isYesterday) {
    return '昨天'
  }

  const daysDiff = Math.floor((Date.now() - timestamp) / (1000 * 60 * 60 * 24))

  if (daysDiff < 7) {
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
    return weekdays[date.getDay()]
  }

  return formatTimestamp(timestamp, 'date')
}

/**
 * 判断是否是今天
 */
export function isToday(timestamp: number): boolean {
  const date = new Date(timestamp)
  const today = new Date()
  return date.toDateString() === today.toDateString()
}

/**
 * 判断是否是本周
 */
export function isThisWeek(timestamp: number): boolean {
  const date = new Date(timestamp)
  const now = new Date()
  const weekStart = new Date(now.setDate(now.getDate() - now.getDay()))
  return date >= weekStart
}

/**
 * 获取时间戳（毫秒）
 */
export function getTimestamp(): number {
  return Date.now()
}

/**
 * 解析ISO日期字符串为时间戳
 */
export function parseISODate(isoString: string): number {
  return new Date(isoString).getTime()
}
