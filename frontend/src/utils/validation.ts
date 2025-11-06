/**
 * 数据验证工具函数
 */

/**
 * 验证邮箱格式
 */
export function isValidEmail(email: string): boolean {
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return regex.test(email)
}

/**
 * 验证URL格式
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 验证HTTP/HTTPS URL
 */
export function isValidHttpUrl(url: string): boolean {
  try {
    const urlObj = new URL(url)
    return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
  } catch {
    return false
  }
}

/**
 * 验证WebDAV URL格式
 */
export function isValidWebDAVUrl(url: string): boolean {
  if (!isValidHttpUrl(url)) return false
  // WebDAV通常以/remote.php/dav或/dav结尾
  return url.includes('/dav') || url.includes('/webdav')
}

/**
 * 验证用户名格式（3-20个字符，字母数字下划线）
 */
export function isValidUsername(username: string): boolean {
  const regex = /^[a-zA-Z0-9_]{3,20}$/
  return regex.test(username)
}

/**
 * 验证密码强度（至少8个字符）
 */
export function isValidPassword(password: string): boolean {
  return password.length >= 8
}

/**
 * 验证API密钥格式
 */
export function isValidApiKey(apiKey: string): boolean {
  // API密钥通常是32-128个字符的字母数字组合
  return apiKey.length >= 20 && apiKey.length <= 256
}

/**
 * 验证JSON格式
 */
export function isValidJSON(str: string): boolean {
  try {
    JSON.parse(str)
    return true
  } catch {
    return false
  }
}

/**
 * 验证非空字符串
 */
export function isNonEmptyString(str: string): boolean {
  return typeof str === 'string' && str.trim().length > 0
}

/**
 * 验证数字范围
 */
export function isInRange(value: number, min: number, max: number): boolean {
  return value >= min && value <= max
}

/**
 * 验证温度参数（0-2）
 */
export function isValidTemperature(temperature: number): boolean {
  return isInRange(temperature, 0, 2)
}

/**
 * 验证token数（1-100000）
 */
export function isValidMaxTokens(tokens: number): boolean {
  return Number.isInteger(tokens) && isInRange(tokens, 1, 100000)
}

/**
 * 验证声优标签（1-20个字符）
 */
export function isValidSeiyuuTag(tag: string): boolean {
  return isNonEmptyString(tag) && tag.length >= 1 && tag.length <= 20
}

/**
 * 清理XSS危险字符
 */
export function sanitizeHtml(str: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#x27;',
    '/': '&#x2F;'
  }
  const reg = /[&<>"'/]/g
  return str.replace(reg, (match) => map[match])
}

/**
 * 验证表单数据
 */
export interface ValidationRule {
  required?: boolean
  minLength?: number
  maxLength?: number
  pattern?: RegExp
  validator?: (value: any) => boolean
  message?: string
}

export interface ValidationResult {
  valid: boolean
  errors: Record<string, string>
}

export function validateForm(
  data: Record<string, any>,
  rules: Record<string, ValidationRule>
): ValidationResult {
  const errors: Record<string, string> = {}

  for (const [field, rule] of Object.entries(rules)) {
    const value = data[field]

    // 必填验证
    if (rule.required && !value) {
      errors[field] = rule.message || `${field}不能为空`
      continue
    }

    // 跳过空值的其他验证
    if (!value) continue

    // 最小长度验证
    if (rule.minLength && value.length < rule.minLength) {
      errors[field] = rule.message || `${field}长度不能少于${rule.minLength}个字符`
      continue
    }

    // 最大长度验证
    if (rule.maxLength && value.length > rule.maxLength) {
      errors[field] = rule.message || `${field}长度不能超过${rule.maxLength}个字符`
      continue
    }

    // 正则验证
    if (rule.pattern && !rule.pattern.test(value)) {
      errors[field] = rule.message || `${field}格式不正确`
      continue
    }

    // 自定义验证
    if (rule.validator && !rule.validator(value)) {
      errors[field] = rule.message || `${field}验证失败`
    }
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors
  }
}
