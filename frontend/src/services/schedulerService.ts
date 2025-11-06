/**
 * 前端调度算法服务
 * 用于本地快速选择声优（降级方案）
 */

import type { Seiyuu } from '@/types'

/**
 * 前端调度服务类
 * 基于关键词匹配的简单调度算法
 */
export class SchedulerService {
  /**
   * 根据消息内容选择声优
   * @param message 用户消息
   * @param availableSeiyuu 可选声优列表
   * @returns 选中的声优和原因
   */
  selectSeiyuu(
    message: string,
    availableSeiyuu: Seiyuu[]
  ): { seiyuu: Seiyuu; reason: string; confidence: number } | null {
    if (availableSeiyuu.length === 0) {
      return null
    }

    // 如果只有一个声优，直接返回
    if (availableSeiyuu.length === 1) {
      return {
        seiyuu: availableSeiyuu[0],
        reason: '唯一可选声优',
        confidence: 1.0
      }
    }

    // 关键词匹配
    const scores = availableSeiyuu.map(seiyuu => {
      let score = 0
      const lowerMessage = message.toLowerCase()

      // 1. 名字匹配（最高权重）
      if (lowerMessage.includes(seiyuu.name.toLowerCase())) {
        score += 10
      }

      // 2. 标签匹配
      for (const tag of seiyuu.tags) {
        if (lowerMessage.includes(tag.toLowerCase())) {
          score += 5
        }
      }

      // 3. 资料内容匹配（简化版）
      const profileKeywords = this.extractKeywords(seiyuu.profile_markdown)
      for (const keyword of profileKeywords) {
        if (lowerMessage.includes(keyword.toLowerCase())) {
          score += 2
        }
      }

      return { seiyuu, score }
    })

    // 排序并选择得分最高的
    scores.sort((a, b) => b.score - a.score)
    const best = scores[0]

    // 如果所有得分都是0，随机选择
    if (best.score === 0) {
      const randomIndex = Math.floor(Math.random() * availableSeiyuu.length)
      return {
        seiyuu: availableSeiyuu[randomIndex],
        reason: '随机选择',
        confidence: 0.3
      }
    }

    // 计算置信度（0-1）
    const maxPossibleScore = 10 + seiyuu.tags.length * 5
    const confidence = Math.min(best.score / maxPossibleScore, 1.0)

    return {
      seiyuu: best.seiyuu,
      reason: this.buildReason(best.seiyuu, message),
      confidence
    }
  }

  /**
   * 从Markdown资料中提取关键词
   */
  private extractKeywords(markdown: string): string[] {
    // 简单实现：提取所有中文词语（2-4个字）
    const regex = /[\u4e00-\u9fa5]{2,4}/g
    const matches = markdown.match(regex) || []

    // 去重并返回前20个
    return Array.from(new Set(matches)).slice(0, 20)
  }

  /**
   * 构建选择原因
   */
  private buildReason(seiyuu: Seiyuu, message: string): string {
    const lowerMessage = message.toLowerCase()

    // 检查是否提到名字
    if (lowerMessage.includes(seiyuu.name.toLowerCase())) {
      return `用户提到了${seiyuu.name}的名字`
    }

    // 检查标签匹配
    const matchedTags = seiyuu.tags.filter(tag =>
      lowerMessage.includes(tag.toLowerCase())
    )
    if (matchedTags.length > 0) {
      return `话题与${seiyuu.name}的特点(${matchedTags.join('、')})相关`
    }

    return `${seiyuu.name}适合回复这个话题`
  }

  /**
   * 轮流选择声优（用于群聊公平发言）
   */
  selectInRoundRobin(
    availableSeiyuu: Seiyuu[],
    lastSpeakerId?: string
  ): Seiyuu | null {
    if (availableSeiyuu.length === 0) return null
    if (availableSeiyuu.length === 1) return availableSeiyuu[0]

    // 如果有上一个发言者，选择下一个
    if (lastSpeakerId) {
      const lastIndex = availableSeiyuu.findIndex(s => s.id === lastSpeakerId)
      if (lastIndex !== -1) {
        const nextIndex = (lastIndex + 1) % availableSeiyuu.length
        return availableSeiyuu[nextIndex]
      }
    }

    // 否则选择第一个
    return availableSeiyuu[0]
  }

  /**
   * 随机选择声优
   */
  selectRandom(availableSeiyuu: Seiyuu[]): Seiyuu | null {
    if (availableSeiyuu.length === 0) return null

    const randomIndex = Math.floor(Math.random() * availableSeiyuu.length)
    return availableSeiyuu[randomIndex]
  }
}

// 导出单例
export const schedulerService = new SchedulerService()

// 导出composable
export function useSchedulerService() {
  return schedulerService
}
