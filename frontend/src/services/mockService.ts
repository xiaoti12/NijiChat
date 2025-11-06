/**
 * Mock数据服务
 * 提供声优数据和AI回复模拟
 */

import type { Seiyuu, SeiyuuListResponse, Message } from '@/types'
import { generateId } from '@/utils/crypto'

// Mock声优数据
const mockSeiyuuList: Seiyuu[] = [
  {
    id: 'seiyuu-1',
    name: '花泽香菜',
    avatar_url: 'https://via.placeholder.com/64/FFB6C1/FFFFFF?text=花泽',
    profile_markdown: `# 花泽香菜

## 基本信息
- **生日**: 2月25日
- **出身地**: 东京都
- **声线特色**: 清澈甜美的声音，擅长少女角色

## 代表作品
- 《化物语》千石抚子
- 《Angel Beats!》天使
- 《我的青春恋爱物语果然有问题》雪之下雪乃

## 性格特点
温柔体贴，声音清新可人，经常会用敬语说话。喜欢和大家聊天，对粉丝很亲切。`,
    tags: ['温柔', '清纯', '少女音'],
    status: 'active',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z'
  },
  {
    id: 'seiyuu-2',
    name: '钉宫理惠',
    avatar_url: 'https://via.placeholder.com/64/FFB6C1/FFFFFF?text=钉宫',
    profile_markdown: `# 钉宫理惠

## 基本信息
- **生日**: 5月30日
- **出身地**: 大阪府
- **声线特色**: 著名的"钉宫病"，傲娇女王

## 代表作品
- 《龙与虎》逢坂大河
- 《灼眼的夏娜》夏娜
- 《零之使魔》露易丝

## 性格特点
傲娇属性MAX，经常会说"笨蛋"、"变态"等，但内心其实很温柔。是傲娇角色的代名词。`,
    tags: ['傲娇', '萝莉音', '经典'],
    status: 'active',
    created_at: '2024-01-02T00:00:00Z',
    updated_at: '2024-01-02T00:00:00Z'
  },
  {
    id: 'seiyuu-3',
    name: '水树奈奈',
    avatar_url: 'https://via.placeholder.com/64/87CEEB/FFFFFF?text=水树',
    profile_markdown: `# 水树奈奈

## 基本信息
- **生日**: 1月21日
- **出身地**: 爱媛县
- **声线特色**: 歌声与演技并重的全能声优

## 代表作品
- 《魔法少女奈叶》菲特·泰斯塔罗莎
- 《交响诗篇》塔尔荷
- 《Heartcatch光之美少女》花咲翼

## 性格特点
充满活力，既可以演绎坚强的战士，也能表现温柔的一面。声音有很强的感染力。`,
    tags: ['活力', '歌姬', '治愈'],
    status: 'active',
    created_at: '2024-01-03T00:00:00Z',
    updated_at: '2024-01-03T00:00:00Z'
  },
  {
    id: 'seiyuu-4',
    name: '田村由香里',
    avatar_url: 'https://via.placeholder.com/64/DDA0DD/FFFFFF?text=田村',
    profile_markdown: `# 田村由香里

## 基本信息
- **生日**: 2月27日
- **出身地**: 福冈县
- **声线特色**: 可爱与冷酷兼备，声域宽广

## 代表作品
- 《魔法少女奈叶》高町奈叶
- 《今天开始做魔王》渥利拉姆
- 《School Days》桂言叶

## 性格特点
很有个性，既能演可爱的魔法少女，也能演黑化的角色。说话直率，有点毒舌。`,
    tags: ['个性', '多面性', '实力派'],
    status: 'active',
    created_at: '2024-01-04T00:00:00Z',
    updated_at: '2024-01-04T00:00:00Z'
  },
  {
    id: 'seiyuu-5',
    name: '堀江由衣',
    avatar_url: 'https://via.placeholder.com/64/98FB98/FFFFFF?text=堀江',
    profile_markdown: `# 堀江由衣

## 基本信息
- **生日**: 9月20日
- **出身地**: 东京都
- **声线特色**: 温和甜美，被称为"堀江女神"

## 代表作品
- 《Love Live!》园田海未
- 《凉宫春日的忧郁》朝比奈实玖瑠
- 《潜行吧！奈亚子》奈亚子

## 性格特点
温柔体贴，声音温暖如春风。很会照顾人，给人很安心的感觉。经常用温柔的语调说话。`,
    tags: ['温柔', '治愈', '女神'],
    status: 'active',
    created_at: '2024-01-05T00:00:00Z',
    updated_at: '2024-01-05T00:00:00Z'
  },
  {
    id: 'seiyuu-6',
    name: '雨宫天',
    avatar_url: 'https://via.placeholder.com/64/F0E68C/FFFFFF?text=雨宫',
    profile_markdown: `# 雨宫天

## 基本信息
- **生日**: 8月28日
- **出身地**: 东京都
- **声线特色**: 清澈明亮，充满青春活力

## 代表作品
- 《在地下城寻求邂逅是否搞错了什么》赫斯缇雅
- 《Charlotte》友利奈绪
- 《刀剑神域》亚丝娜

## 性格特点
活泼开朗，声音很有亲和力。说话很有元气，总是充满正能量，能给人带来快乐。`,
    tags: ['活泼', '元气', '青春'],
    status: 'active',
    created_at: '2024-01-06T00:00:00Z',
    updated_at: '2024-01-06T00:00:00Z'
  },
  {
    id: 'seiyuu-7',
    name: '茅野爱衣',
    avatar_url: 'https://via.placeholder.com/64/E6E6FA/FFFFFF?text=茅野',
    profile_markdown: `# 茅野爱衣

## 基本信息
- **生日**: 9月13日
- **出身地**: 东京都
- **声线特色**: 温柔透明的声音，治愈系代表

## 代表作品
- 《我的青春恋爱物语果然有问题》由比滨结衣
- 《无职转生》希露菲叶特
- 《五等分的新娘》中野一花

## 性格特点
非常温柔，说话轻声细语，像微风一样温暖。很善解人意，总是为别人着想。`,
    tags: ['治愈', '温柔', '透明感'],
    status: 'active',
    created_at: '2024-01-07T00:00:00Z',
    updated_at: '2024-01-07T00:00:00Z'
  }
]

// AI回复模板
const replyTemplates = {
  '花泽香菜': [
    '谢谢你的话呢～我会继续努力的！',
    '诶？是这样吗？真的很开心～',
    '能和你聊天真好呢，感觉很温暖～',
    '请多多指教哦！我会加油的！',
    '今天也要一起快乐地聊天呢～'
  ],
  '钉宫理惠': [
    '哼！才、才不是为了你呢！',
    '笨蛋！说什么傻话呢！',
    '虽然不讨厌啦...但是不要误会了！',
    '真是的！你这个笨蛋！',
    '哼～既然你这么说了，就勉强接受吧！'
  ],
  '水树奈奈': [
    '好的！我会全力以赴的！',
    '真的吗？太棒了！一起加油吧！',
    '唱歌的时候最开心了～',
    '谢谢你的支持！这给了我很大的力量！',
    '今天也要充满活力地度过呢！'
  ],
  '田村由香里': [
    '呵呵，有趣呢～',
    '是吗？那我可要好好考虑一下了～',
    '你的想法很独特呢～我喜欢～',
    '嗯～不愧是你呢～',
    '说得对，就是要这样有个性！'
  ],
  '堀江由衣': [
    '真的吗？你真温柔呢～',
    '和你说话总是让我很安心～',
    '请不要勉强自己哦～',
    '今天辛苦了呢，要好好休息～',
    '能帮到你的话我也很开心～'
  ],
  '雨宫天': [
    '哇！太棒了！超级开心的！',
    '是的是的！我也这么觉得！',
    '今天也很有精神呢！',
    '一起努力吧！我会加油的！',
    '真的吗？好期待啊～！'
  ],
  '茅野爱衣': [
    '是这样呢...谢谢你告诉我～',
    '嗯嗯...我明白了～',
    '你说得很对呢...真的很温柔～',
    '能听到你的声音...真好呢～',
    '请不要太勉强自己哦～'
  ]
}

// 默认回复模板（当找不到特定声优时使用）
const defaultReplies = [
  '是这样呢～',
  '谢谢你～',
  '我明白了～',
  '真的吗？',
  '和你聊天很开心呢～'
]

/**
 * Mock声优列表API
 */
export function getMockSeiyuuList(): Promise<SeiyuuListResponse> {
  return new Promise((resolve) => {
    // 模拟网络延迟
    setTimeout(() => {
      resolve({
        success: true,
        message: '获取成功',
        data: mockSeiyuuList
      })
    }, 300)
  })
}

/**
 * 根据ID获取声优信息
 */
export function getMockSeiyuuById(id: string): Seiyuu | null {
  return mockSeiyuuList.find(s => s.id === id) || null
}

/**
 * 生成AI回复
 */
export function generateMockAIReply(seiyuuName: string, userMessage: string, history: Message[] = []): Promise<string> {
  return new Promise((resolve) => {
    // 模拟AI思考时间
    const delay = Math.random() * 1000 + 500 // 0.5-1.5秒

    setTimeout(() => {
      // 获取该声优的回复模板
      const templates = replyTemplates[seiyuuName] || defaultReplies

      // 基于用户消息和历史选择合适的回复
      let reply: string

      // 简单的关键词匹配逻辑
      if (userMessage.includes('你好') || userMessage.includes('初次见面')) {
        if (seiyuuName === '钉宫理惠') {
          reply = '哼！你好什么的...才不是第一次见面呢！'
        } else if (seiyuuName === '花泽香菜') {
          reply = '你好～初次见面，请多多指教呢！'
        } else {
          reply = '你好！很高兴认识你～'
        }
      } else if (userMessage.includes('再见') || userMessage.includes('拜拜')) {
        if (seiyuuName === '钉宫理惠') {
          reply = '哼！要走了吗...那就...再见吧！'
        } else if (seiyuuName === '茅野爱衣') {
          reply = '再见呢...要保重身体哦～'
        } else {
          reply = '再见～下次再聊呢！'
        }
      } else if (userMessage.includes('喜欢') || userMessage.includes('爱')) {
        if (seiyuuName === '钉宫理惠') {
          reply = '喜、喜欢什么的...才没有呢！笨蛋！'
        } else if (seiyuuName === '雨宫天') {
          reply = '诶？！真的吗？我也很喜欢和你聊天呢！'
        } else {
          reply = '谢谢你～听到这样的话很开心呢～'
        }
      } else {
        // 随机选择回复模板
        const randomIndex = Math.floor(Math.random() * templates.length)
        reply = templates[randomIndex]
      }

      resolve(reply)
    }, delay)
  })
}

/**
 * 搜索声优
 */
export function searchMockSeiyuu(keyword: string): Seiyuu[] {
  if (!keyword.trim()) return mockSeiyuuList

  const lowerKeyword = keyword.toLowerCase()
  return mockSeiyuuList.filter(seiyuu =>
    seiyuu.name.toLowerCase().includes(lowerKeyword) ||
    seiyuu.tags.some(tag => tag.toLowerCase().includes(lowerKeyword))
  )
}

/**
 * 获取热门声优（用于推荐）
 */
export function getPopularSeiyuu(): Seiyuu[] {
  // 返回前5个声优作为热门推荐
  return mockSeiyuuList.slice(0, 5)
}

/**
 * 检查是否为Mock模式
 */
export function isMockMode(): boolean {
  // 可以通过环境变量或localStorage控制
  return localStorage.getItem('nijichat_mock_mode') === 'true' || import.meta.env.DEV
}

/**
 * 设置Mock模式
 */
export function setMockMode(enabled: boolean): void {
  localStorage.setItem('nijichat_mock_mode', enabled.toString())
}

/**
 * 生成初始对话数据（用于演示）
 */
export function generateInitialConversations(): Array<{
  seiyuuId: string
  lastMessage: string
  timestamp: number
}> {
  return [
    {
      seiyuuId: 'seiyuu-1',
      lastMessage: '今天也要一起快乐地聊天呢～',
      timestamp: Date.now() - 1000 * 60 * 30 // 30分钟前
    },
    {
      seiyuuId: 'seiyuu-2',
      lastMessage: '哼！才不是想你了呢！',
      timestamp: Date.now() - 1000 * 60 * 60 * 2 // 2小时前
    },
    {
      seiyuuId: 'seiyuu-3',
      lastMessage: '一起加油吧！',
      timestamp: Date.now() - 1000 * 60 * 60 * 24 // 1天前
    }
  ]
}

// 导出默认的MockService实例
export const mockService = {
  getSeiyuuList: getMockSeiyuuList,
  getSeiyuuById: getMockSeiyuuById,
  generateAIReply: generateMockAIReply,
  searchSeiyuu: searchMockSeiyuu,
  getPopularSeiyuu,
  isMockMode,
  setMockMode,
  generateInitialConversations
}

export default mockService