/**
 * 加密解密工具函数
 * 用于敏感数据的本地加密存储
 */

/**
 * 简单的XOR加密（用于本地存储）
 * 注意：这不是强加密，仅用于防止明文存储
 */
function xorEncrypt(text: string, key: string): string {
  let result = ''
  for (let i = 0; i < text.length; i++) {
    result += String.fromCharCode(text.charCodeAt(i) ^ key.charCodeAt(i % key.length))
  }
  return result
}

/**
 * Base64编码
 */
function base64Encode(str: string): string {
  return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, (match, p1) => {
    return String.fromCharCode(parseInt(p1, 16))
  }))
}

/**
 * Base64解码
 */
function base64Decode(str: string): string {
  return decodeURIComponent(Array.prototype.map.call(atob(str), (c) => {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
  }).join(''))
}

/**
 * 加密文本
 * @param text 要加密的文本
 * @param password 密码（可选，默认使用设备指纹）
 */
export function encrypt(text: string, password?: string): string {
  const key = password || getDeviceFingerprint()
  const encrypted = xorEncrypt(text, key)
  return base64Encode(encrypted)
}

/**
 * 解密文本
 * @param encryptedText 加密的文本
 * @param password 密码（可选，默认使用设备指纹）
 */
export function decrypt(encryptedText: string, password?: string): string {
  try {
    const key = password || getDeviceFingerprint()
    const decoded = base64Decode(encryptedText)
    return xorEncrypt(decoded, key)
  } catch (error) {
    console.error('解密失败:', error)
    return ''
  }
}

/**
 * 获取设备指纹（简单实现）
 * 用于生成设备唯一标识
 */
export function getDeviceFingerprint(): string {
  // 使用浏览器特征生成简单的设备指纹
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (ctx) {
    ctx.textBaseline = 'top'
    ctx.font = '14px Arial'
    ctx.fillText('fingerprint', 2, 2)
  }

  const fingerprint = [
    navigator.userAgent,
    navigator.language,
    screen.colorDepth,
    screen.width + 'x' + screen.height,
    new Date().getTimezoneOffset(),
    canvas.toDataURL()
  ].join('###')

  return hashCode(fingerprint)
}

/**
 * 简单的哈希函数
 */
function hashCode(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash
  }
  return hash.toString(36)
}

/**
 * 生成随机ID
 */
export function generateId(): string {
  return Date.now().toString(36) + Math.random().toString(36).substr(2)
}

/**
 * 生成UUID v4
 */
export function generateUUID(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

/**
 * 验证密码强度
 * @returns 强度等级 0-4
 */
export function checkPasswordStrength(password: string): number {
  let strength = 0

  if (password.length >= 8) strength++
  if (password.length >= 12) strength++
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++
  if (/\d/.test(password)) strength++
  if (/[^a-zA-Z0-9]/.test(password)) strength++

  return Math.min(strength, 4)
}
